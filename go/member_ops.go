package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/router-for-me/CLIProxyAPI/v7/sdk/pluginapi"
)

func normalizeMemberOperationRequest(req pluginapi.ManagementRequest) (memberOperationRequest, pluginapi.ManagementResponse, bool) {
	var body memberOperationRequest
	if len(req.Body) == 0 {
		resp, _ := jsonResponse(http.StatusBadRequest, map[string]any{"error": "request body is required"})
		return body, resp, false
	}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		resp, _ := jsonResponse(http.StatusBadRequest, map[string]any{"error": err.Error()})
		return body, resp, false
	}
	body.KeyID = strings.TrimSpace(body.KeyID)
	body.RuleID = strings.TrimSpace(body.RuleID)
	body.Member = strings.TrimSpace(body.Member)
	body.Operation = strings.TrimSpace(body.Operation)
	body.Reason = strings.TrimSpace(body.Reason)
	body.Secondary = strings.TrimSpace(body.Secondary)
	body.PreviewToken = strings.TrimSpace(body.PreviewToken)
	if body.KeyID == "" || body.RuleID == "" || body.Operation == "" {
		resp, _ := jsonResponse(http.StatusBadRequest, map[string]any{"error": "key_id, rule_id, and operation are required"})
		return body, resp, false
	}
	needsMember := !strings.EqualFold(body.Operation, "pool-drain") && !strings.EqualFold(body.Operation, "pool-resume") && !strings.EqualFold(body.Operation, "canary-split") && !strings.EqualFold(body.Operation, "rebalance-by-health") && !strings.EqualFold(body.Operation, "restore-default-weights") && !strings.EqualFold(body.Operation, "shift-provider-traffic")
	if needsMember && body.Member == "" {
		resp, _ := jsonResponse(http.StatusBadRequest, map[string]any{"error": "member is required for this operation"})
		return body, resp, false
	}
	return body, pluginapi.ManagementResponse{}, true
}

func findMemberOperationRulePreview(body memberOperationRequest) (ruleConfig, []weightedRoute, []weightedRoute, pluginapi.ManagementResponse, bool) {
	for i := range gatewayState.config.KeyPolicies {
		if candidateKeyID(gatewayState.config.KeyPolicies[i]) != body.KeyID {
			continue
		}
		for j := range gatewayState.config.KeyPolicies[i].Rules {
			rule := gatewayState.config.KeyPolicies[i].Rules[j]
			if strings.TrimSpace(rule.ID) != body.RuleID {
				continue
			}
			members := cloneWeightedRoutes(rule.Actions.WeightedRoutes)
			if rule.Actions.RoutePool != nil && len(rule.Actions.RoutePool.Members) > 0 {
				members = cloneWeightedRoutes(rule.Actions.RoutePool.Members)
			}
			beforeMembers := cloneWeightedRoutes(members)
			return rule, members, beforeMembers, pluginapi.ManagementResponse{}, true
		}
		resp, _ := jsonResponse(http.StatusNotFound, map[string]any{"error": "rule not found"})
		return ruleConfig{}, nil, nil, resp, false
	}
	resp, _ := jsonResponse(http.StatusNotFound, map[string]any{"error": "policy not found"})
	return ruleConfig{}, nil, nil, resp, false
}

func routeMemberPreview(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	body, resp, ok := normalizeMemberOperationRequest(req)
	if !ok {
		return resp, nil
	}
	gatewayState.mu.Lock()
	defer gatewayState.mu.Unlock()
	rule, members, beforeMembers, respLookup, found := findMemberOperationRulePreview(body)
	_ = rule
	if !found {
		return respLookup, nil
	}
	updated, errMsg := previewApplyMemberOperationToRoutes(members, body)
	if errMsg != "" {
		status := http.StatusNotFound
		if strings.Contains(errMsg, "required") || strings.Contains(errMsg, "unsupported") {
			status = http.StatusBadRequest
		}
		return jsonResponse(status, map[string]any{"error": errMsg})
	}
	if !updated {
		return jsonResponse(http.StatusNotFound, map[string]any{"error": "member not found"})
	}
	issuedAt := time.Now()
	preview := buildMemberOperationPreview(body, beforeMembers, members, issuedAt)
	gatewayState.prunePreviewTokensLocked(issuedAt)
	gatewayState.previewTokens[preview.PreviewToken] = previewTokenRecord{KeyID: body.KeyID, RuleID: body.RuleID, Operation: body.Operation, Target: preview.TargetMember, Secondary: body.Secondary, BeforeState: preview.BeforeState, AfterState: preview.AfterState, Token: preview.PreviewToken, IssuedAt: issuedAt}
	return jsonResponse(http.StatusOK, preview)
}

func buildMemberOperationResult(body memberOperationRequest, beforeMembers, members []weightedRoute) memberOperationResult {
	target := body.Member
	if strings.TrimSpace(target) == "" && len(members) > 0 {
		target = members[0].Model
		if strings.TrimSpace(target) == "" {
			target = members[0].Provider + "/" + members[0].Suffix
		}
	}
	beforeState := summarizeWeightedRoutes(beforeMembers)
	afterState := summarizeWeightedRoutes(members)
	return memberOperationResult{
		OK:           true,
		KeyID:        body.KeyID,
		RuleID:       body.RuleID,
		Operation:    body.Operation,
		Reason:       firstNonEmpty(body.Reason, body.Operation),
		TargetMember: target,
		Secondary:    body.Secondary,
		BeforeState:  beforeState,
		AfterState:   afterState,
		Diff:         diffWeightedRoutes(beforeMembers, members),
		Members:      members,
	}
}

func buildMemberOperationPreview(body memberOperationRequest, beforeMembers, members []weightedRoute, issuedAt time.Time) memberOperationPreview {
	result := buildMemberOperationResult(body, beforeMembers, members)
	previewToken := signPreviewToken(body.KeyID, body.RuleID, body.Operation, result.TargetMember, body.Secondary, result.BeforeState, result.AfterState, issuedAt)
	return memberOperationPreview{
		memberOperationResult: result,
		PreviewToken:          previewToken,
	}
}

func previewApplyMemberOperationToRoutes(members []weightedRoute, body memberOperationRequest) (bool, string) {
	updated := false
	op := strings.ToLower(body.Operation)
	if op == "pool-drain" || op == "pool-resume" {
		for k := range members {
			if op == "pool-drain" {
				value := true
				members[k].Enabled = &value
				members[k].Status = "drain"
				members[k].Reason = firstNonEmpty(body.Reason, "pool-drain")
			} else {
				value := true
				members[k].Enabled = &value
				members[k].Status = "active"
				members[k].Reason = firstNonEmpty(body.Reason, "pool-resume")
			}
		}
		updated = len(members) > 0
	}
	if op == "canary-split" {
		if body.Member == "" || body.Secondary == "" {
			return false, "member and secondary are required for canary-split"
		}
		primaryIndex := -1
		secondaryIndex := -1
		for k := range members {
			label := strings.TrimSpace(members[k].Model)
			if label == "" {
				label = strings.TrimSpace(members[k].Provider + "/" + members[k].Suffix)
			}
			if strings.EqualFold(label, body.Member) {
				primaryIndex = k
			}
			if strings.EqualFold(label, body.Secondary) {
				secondaryIndex = k
			}
		}
		if primaryIndex < 0 || secondaryIndex < 0 {
			return false, "primary or secondary member not found"
		}
		primaryWeight := body.PrimaryWeight
		if primaryWeight <= 0 {
			primaryWeight = 90
		}
		canaryWeight := body.CanaryWeight
		if canaryWeight <= 0 {
			canaryWeight = 10
		}
		members[primaryIndex].Weight = primaryWeight
		members[primaryIndex].Status = "active"
		members[primaryIndex].Reason = firstNonEmpty(body.Reason, "canary-primary")
		members[secondaryIndex].Weight = canaryWeight
		members[secondaryIndex].Status = "active"
		members[secondaryIndex].Reason = firstNonEmpty(body.Reason, "canary-secondary")
		updated = true
	}
	if op == "restore-default-weights" {
		for k := range members {
			members[k].Weight = 100
			members[k].Reason = firstNonEmpty(body.Reason, "restore-default-weights")
		}
		updated = len(members) > 0
	}
	if op == "shift-provider-traffic" {
		provider := strings.ToLower(strings.TrimSpace(body.Secondary))
		percent := clampInt(body.CanaryWeight, 1, 100)
		matched := 0
		for k := range members {
			memberProvider := strings.ToLower(strings.TrimSpace(firstNonEmpty(members[k].Provider, providerFromModel(members[k].Model))))
			if provider != "" && memberProvider == provider {
				members[k].Weight = percent
				members[k].Reason = firstNonEmpty(body.Reason, "shift-provider-traffic")
				matched++
			} else {
				members[k].Weight = clampInt(100-percent, 1, 10000)
			}
		}
		updated = matched > 0
	}
	if op == "rebalance-by-health" {
		totalHealth := 0
		for k := range members {
			health := members[k].Health
			if health <= 0 {
				health = 1
			}
			totalHealth += health
		}
		if totalHealth > 0 {
			for k := range members {
				health := members[k].Health
				if health <= 0 {
					health = 1
				}
				members[k].Weight = clampInt((health*100)/totalHealth, 1, 10000)
				members[k].Reason = firstNonEmpty(body.Reason, "rebalance-by-health")
			}
			updated = len(members) > 0
		}
	}
	if !updated {
		for k := range members {
			label := strings.TrimSpace(members[k].Model)
			if label == "" {
				label = strings.TrimSpace(members[k].Provider + "/" + members[k].Suffix)
			}
			if !strings.EqualFold(label, body.Member) {
				continue
			}
			switch op {
			case "active":
				value := true
				members[k].Enabled = &value
				members[k].Status = "active"
				if body.Reason != "" {
					members[k].Reason = body.Reason
				}
			case "drain":
				value := true
				members[k].Enabled = &value
				members[k].Status = "drain"
				members[k].Reason = firstNonEmpty(body.Reason, "manual-drain")
			case "offline":
				value := false
				members[k].Enabled = &value
				members[k].Status = "offline"
				members[k].Reason = firstNonEmpty(body.Reason, "manual-offline")
			case "cap-up":
				members[k].TrafficCap = clampInt(members[k].TrafficCap+maxInt(body.Delta, 10), 0, 100)
				if body.Reason != "" {
					members[k].Reason = body.Reason
				}
			case "cap-down":
				delta := body.Delta
				if delta <= 0 {
					delta = 10
				}
				members[k].TrafficCap = clampInt(members[k].TrafficCap-delta, 0, 100)
				if body.Reason != "" {
					members[k].Reason = body.Reason
				}
			case "weight-up":
				members[k].Weight = clampInt(maxInt(members[k].Weight, 1)+maxInt(body.Delta, 1), 1, 10000)
				if body.Reason != "" {
					members[k].Reason = body.Reason
				}
			case "weight-down":
				delta := body.Delta
				if delta <= 0 {
					delta = 1
				}
				members[k].Weight = clampInt(maxInt(members[k].Weight, 1)-delta, 1, 10000)
				if body.Reason != "" {
					members[k].Reason = body.Reason
				}
			case "priority-up":
				members[k].Priority = clampInt(members[k].Priority+maxInt(body.Delta, 10), 0, 10000)
				if body.Reason != "" {
					members[k].Reason = body.Reason
				}
			case "priority-down":
				delta := body.Delta
				if delta <= 0 {
					delta = 10
				}
				members[k].Priority = clampInt(members[k].Priority-delta, 0, 10000)
				if body.Reason != "" {
					members[k].Reason = body.Reason
				}
			case "health-up":
				members[k].Health = clampInt(members[k].Health+maxInt(body.Delta, 10), 0, 100)
				if body.Reason != "" {
					members[k].Reason = body.Reason
				}
			case "health-down":
				delta := body.Delta
				if delta <= 0 {
					delta = 10
				}
				members[k].Health = clampInt(members[k].Health-delta, 0, 100)
				if body.Reason != "" {
					members[k].Reason = body.Reason
				}
			default:
				return false, "unsupported member operation"
			}
			updated = true
			break
		}
	}
	if !updated {
		return false, "member not found"
	}
	return true, ""
}

func validatePreviewTokenLocked(body memberOperationRequest) (*previewTokenRecord, error) {
	if strings.TrimSpace(body.PreviewToken) == "" {
		return nil, fmt.Errorf("preview token is required; preview the operation first")
	}
	record, ok := gatewayState.previewTokens[strings.TrimSpace(body.PreviewToken)]
	if !ok {
		return nil, fmt.Errorf("preview token not found")
	}
	if time.Since(record.IssuedAt) > 10*time.Minute {
		delete(gatewayState.previewTokens, strings.TrimSpace(body.PreviewToken))
		return nil, fmt.Errorf("preview token expired")
	}
	if record.KeyID != body.KeyID || record.RuleID != body.RuleID || !strings.EqualFold(record.Operation, body.Operation) {
		return nil, fmt.Errorf("preview token does not match requested operation")
	}
	if strings.TrimSpace(body.Secondary) != strings.TrimSpace(record.Secondary) {
		return nil, fmt.Errorf("preview token does not match requested secondary target")
	}
	if strings.TrimSpace(body.Member) != "" && !strings.EqualFold(strings.TrimSpace(body.Member), strings.TrimSpace(record.Target)) {
		return nil, fmt.Errorf("preview token does not match requested member")
	}
	expected := signPreviewToken(record.KeyID, record.RuleID, record.Operation, record.Target, record.Secondary, record.BeforeState, record.AfterState, record.IssuedAt)
	if record.Token != expected || strings.TrimSpace(body.PreviewToken) != expected {
		return nil, fmt.Errorf("preview token signature mismatch")
	}
	return &record, nil
}

func findMemberOperationRuleApply(body memberOperationRequest) (*keyPolicyConfig, *ruleConfig, []weightedRoute, []weightedRoute, pluginapi.ManagementResponse, bool) {
	for i := range gatewayState.config.KeyPolicies {
		if candidateKeyID(gatewayState.config.KeyPolicies[i]) != body.KeyID {
			continue
		}
		for j := range gatewayState.config.KeyPolicies[i].Rules {
			rule := &gatewayState.config.KeyPolicies[i].Rules[j]
			if strings.TrimSpace(rule.ID) != body.RuleID {
				continue
			}
			members := cloneWeightedRoutes(rule.Actions.WeightedRoutes)
			if rule.Actions.RoutePool != nil && len(rule.Actions.RoutePool.Members) > 0 {
				members = cloneWeightedRoutes(rule.Actions.RoutePool.Members)
			}
			beforeMembers := cloneWeightedRoutes(members)
			return &gatewayState.config.KeyPolicies[i], rule, members, beforeMembers, pluginapi.ManagementResponse{}, true
		}
		resp, _ := jsonResponse(http.StatusNotFound, map[string]any{"error": "rule not found"})
		return nil, nil, nil, nil, resp, false
	}
	resp, _ := jsonResponse(http.StatusNotFound, map[string]any{"error": "policy not found"})
	return nil, nil, nil, nil, resp, false
}

func routeMemberOperation(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	body, resp, ok := normalizeMemberOperationRequest(req)
	if !ok {
		return resp, nil
	}
	gatewayState.mu.Lock()
	tokenRecord, errToken := validatePreviewTokenLocked(body)
	if errToken != nil {
		gatewayState.mu.Unlock()
		return jsonResponse(http.StatusForbidden, map[string]any{"error": errToken.Error()})
	}
	defer gatewayState.mu.Unlock()
	policy, rule, members, beforeMembers, respLookup, found := findMemberOperationRuleApply(body)
	if !found {
		return respLookup, nil
	}
	if tokenRecord != nil && tokenRecord.BeforeState != summarizeWeightedRoutes(beforeMembers) {
		delete(gatewayState.previewTokens, strings.TrimSpace(body.PreviewToken))
		return jsonResponse(http.StatusConflict, map[string]any{"error": "preview token is stale; preview the operation again"})
	}
	updated, errMsg := previewApplyMemberOperationToRoutes(members, body)
	if errMsg != "" {
		status := http.StatusNotFound
		if strings.Contains(errMsg, "required") || strings.Contains(errMsg, "unsupported") {
			status = http.StatusBadRequest
		}
		return jsonResponse(status, map[string]any{"error": errMsg})
	}
	if !updated {
		return jsonResponse(http.StatusNotFound, map[string]any{"error": "member not found"})
	}
	result := buildMemberOperationResult(body, beforeMembers, members)
	if tokenRecord != nil && tokenRecord.AfterState != result.AfterState {
		delete(gatewayState.previewTokens, strings.TrimSpace(body.PreviewToken))
		return jsonResponse(http.StatusConflict, map[string]any{"error": "preview token result mismatch; preview the operation again"})
	}
	if rule.Actions.RoutePool != nil && len(rule.Actions.RoutePool.Members) > 0 {
		rule.Actions.RoutePool.Members = members
	} else {
		rule.Actions.WeightedRoutes = members
	}
	*policy = normalizeKeyPolicies([]keyPolicyConfig{*policy})[0]
	if tokenRecord != nil {
		delete(gatewayState.previewTokens, strings.TrimSpace(body.PreviewToken))
	}
	preview := memberOperationPreview{memberOperationResult: result}
	gatewayState.appendAuditLocked(auditEntry{Time: time.Now(), PolicyID: policy.KeyID, PolicyName: policy.DisplayName, Decision: "operator", RuleID: rule.ID, Reason: result.Reason, FinalModel: result.TargetMember, Provider: providerFromModel(result.TargetMember), EventType: "operator", OperatorAction: body.Operation, TargetMember: result.TargetMember, Secondary: body.Secondary, BeforeState: result.BeforeState, AfterState: result.AfterState, Diff: result.Diff})
	if err := gatewayState.savePersistentStateLocked(); err != nil {
		return persistenceErrorResponse(err)
	}
	return jsonResponse(http.StatusOK, map[string]any{"ok": true, "members": members, "preview": preview, "result": result})
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func clampInt(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func cloneWeightedRoutes(items []weightedRoute) []weightedRoute {
	if len(items) == 0 {
		return nil
	}
	out := make([]weightedRoute, len(items))
	copy(out, items)
	return out
}

func diffWeightedRoutes(before, after []weightedRoute) []memberChangeSummary {
	labels := map[string]weightedRoute{}
	for _, item := range before {
		label := strings.TrimSpace(item.Model)
		if label == "" {
			label = strings.TrimSpace(item.Provider + "/" + item.Suffix)
		}
		labels[strings.ToLower(label)] = item
	}
	changes := make([]memberChangeSummary, 0)
	seen := make(map[string]struct{})
	for _, item := range after {
		label := strings.TrimSpace(item.Model)
		if label == "" {
			label = strings.TrimSpace(item.Provider + "/" + item.Suffix)
		}
		key := strings.ToLower(label)
		beforeItem, ok := labels[key]
		if ok {
			seen[key] = struct{}{}
			if beforeItem != item {
				changes = append(changes, memberChangeSummary{Member: label, Before: beforeItem, After: item})
			}
			continue
		}
		changes = append(changes, memberChangeSummary{Member: label, Before: weightedRoute{}, After: item})
	}
	for _, item := range before {
		label := strings.TrimSpace(item.Model)
		if label == "" {
			label = strings.TrimSpace(item.Provider + "/" + item.Suffix)
		}
		key := strings.ToLower(label)
		if _, ok := seen[key]; ok {
			continue
		}
		changes = append(changes, memberChangeSummary{Member: label, Before: item, After: weightedRoute{}})
	}
	return changes
}

func summarizeWeightedRoutes(items []weightedRoute) string {
	if len(items) == 0 {
		return "[]"
	}
	parts := make([]string, 0, len(items))
	for _, item := range items {
		label := strings.TrimSpace(item.Model)
		if label == "" {
			label = strings.TrimSpace(item.Provider + "/" + item.Suffix)
		}
		parts = append(parts, fmt.Sprintf("%s[w=%d,p=%d,h=%d,cap=%d,status=%s]", label, item.Weight, item.Priority, item.Health, item.TrafficCap, firstNonEmpty(item.Status, "active")))
	}
	return strings.Join(parts, "; ")
}

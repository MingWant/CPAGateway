package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/router-for-me/CLIProxyAPI/v7/sdk/pluginapi"
)

func routeKeys(_ pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	keys := gatewayState.listKeys()
	return jsonResponse(http.StatusOK, map[string]any{"keys": keys})
}

func routePolicies(_ pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	gatewayState.mu.RLock()
	defer gatewayState.mu.RUnlock()
	return jsonResponse(http.StatusOK, sanitizedStoredPolicy(storedPolicy{Version: 1, DefaultPolicy: gatewayState.config.Default, KeyPolicies: gatewayState.config.KeyPolicies}))
}

func routePutPolicies(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	var body storedPolicy
	if len(req.Body) == 0 {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "request body is required"})
	}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	gatewayState.mu.Lock()
	gatewayState.config.Default = normalizePolicy(body.DefaultPolicy)
	gatewayState.config.KeyPolicies = normalizeKeyPolicies(gatewayState.preservePolicySecretsLocked(body.KeyPolicies))
	gatewayState.mu.Unlock()
	return jsonResponse(http.StatusOK, map[string]any{"ok": true})
}

func routeExportPolicies(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	name := strings.TrimSpace(req.Query.Get("name"))
	if name == "" {
		name = "gateway-policy-bundle"
	}
	description := strings.TrimSpace(req.Query.Get("description"))
	tags := compactUnique(strings.Split(strings.TrimSpace(req.Query.Get("tags")), ","))
	gatewayState.mu.RLock()
	defer gatewayState.mu.RUnlock()
	bundle := policyBundle{
		Version:       1,
		Name:          name,
		Description:   description,
		Tags:          tags,
		DefaultPolicy: gatewayState.config.Default,
		KeyPolicies:   append([]keyPolicyConfig(nil), gatewayState.config.KeyPolicies...),
		ExportedAt:    time.Now(),
	}
	return jsonResponse(http.StatusOK, sanitizedPolicyBundle(bundle))
}

func routeImportPolicies(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	var body policyBundle
	if len(req.Body) == 0 {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "request body is required"})
	}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	mode := strings.ToLower(strings.TrimSpace(req.Query.Get("mode")))
	if mode == "" {
		mode = "merge"
	}
	gatewayState.mu.Lock()
	defer gatewayState.mu.Unlock()
	incoming := normalizeKeyPolicies(gatewayState.preservePolicySecretsLocked(body.KeyPolicies))
	if mode == "replace" {
		gatewayState.config.Default = normalizePolicy(body.DefaultPolicy)
		gatewayState.config.KeyPolicies = incoming
		return jsonResponse(http.StatusOK, map[string]any{"ok": true, "mode": mode, "imported": len(incoming)})
	}
	merged := append([]keyPolicyConfig(nil), gatewayState.config.KeyPolicies...)
	seen := map[string]int{}
	for i, item := range merged {
		seen[candidateKeyID(item)] = i
	}
	imported := 0
	updated := 0
	for _, item := range incoming {
		if idx, exists := seen[candidateKeyID(item)]; exists {
			merged[idx] = item
			updated++
			continue
		}
		merged = append(merged, item)
		seen[candidateKeyID(item)] = len(merged) - 1
		imported++
	}
	gatewayState.config.KeyPolicies = normalizeKeyPolicies(merged)
	if policyHasContent(body.DefaultPolicy) {
		gatewayState.config.Default = normalizePolicy(body.DefaultPolicy)
	}
	return jsonResponse(http.StatusOK, map[string]any{"ok": true, "mode": mode, "imported": imported, "updated": updated})
}

func routeClonePolicy(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	keyID := strings.TrimSpace(req.Query.Get("key_id"))
	if keyID == "" {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "key_id query is required"})
	}
	nameSuffix := strings.TrimSpace(req.Query.Get("name_suffix"))
	if nameSuffix == "" {
		nameSuffix = " Copy"
	}
	gatewayState.mu.Lock()
	defer gatewayState.mu.Unlock()
	for _, item := range gatewayState.config.KeyPolicies {
		if candidateKeyID(item) != keyID {
			continue
		}
		clone := item
		clone.KeyID = stableKeyID(item.DisplayName + nameSuffix + time.Now().Format(time.RFC3339Nano))
		clone.DisplayName = strings.TrimSpace(item.DisplayName + nameSuffix)
		clone.MatchAPIKey = ""
		gatewayState.config.KeyPolicies = append(gatewayState.config.KeyPolicies, normalizeKeyPolicies([]keyPolicyConfig{clone})[0])
		return jsonResponse(http.StatusCreated, map[string]any{"ok": true, "key_id": clone.KeyID})
	}
	return jsonResponse(http.StatusNotFound, map[string]any{"error": "policy not found"})
}

func routeAddPolicy(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	var body keyPolicyConfig
	if len(req.Body) == 0 {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "request body is required"})
	}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	if strings.TrimSpace(body.DisplayName) == "" {
		body.DisplayName = "New Policy"
	}
	if strings.TrimSpace(body.KeyID) == "" {
		seed := body.MatchAPIKey
		if strings.TrimSpace(seed) == "" {
			seed = body.DisplayName + time.Now().Format(time.RFC3339Nano)
		}
		body.KeyID = stableKeyID(seed)
	}
	gatewayState.mu.Lock()
	defer gatewayState.mu.Unlock()
	for _, item := range gatewayState.config.KeyPolicies {
		if strings.TrimSpace(body.MatchAPIKey) != "" && strings.TrimSpace(item.MatchAPIKey) == strings.TrimSpace(body.MatchAPIKey) {
			return jsonResponse(http.StatusConflict, map[string]any{"error": "policy already exists for key"})
		}
		if candidateKeyID(item) == candidateKeyID(body) {
			return jsonResponse(http.StatusConflict, map[string]any{"error": "policy id already exists"})
		}
	}
	gatewayState.config.KeyPolicies = append(gatewayState.config.KeyPolicies, normalizeKeyPolicies([]keyPolicyConfig{body})[0])
	return jsonResponse(http.StatusCreated, map[string]any{"ok": true})
}

func routeDeletePolicy(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	keyID := strings.TrimSpace(req.Query.Get("key_id"))
	if keyID == "" {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "key_id query is required"})
	}
	gatewayState.mu.Lock()
	defer gatewayState.mu.Unlock()
	out := make([]keyPolicyConfig, 0, len(gatewayState.config.KeyPolicies))
	removed := false
	for _, item := range gatewayState.config.KeyPolicies {
		if candidateKeyID(item) == keyID {
			removed = true
			continue
		}
		out = append(out, item)
	}
	if !removed {
		return jsonResponse(http.StatusNotFound, map[string]any{"error": "policy not found"})
	}
	gatewayState.config.KeyPolicies = out
	delete(gatewayState.usage, keyID)
	delete(gatewayState.requestWindow, keyID)
	return jsonResponse(http.StatusOK, map[string]any{"ok": true})
}

func routePatchPolicy(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	keyID := strings.TrimSpace(req.Query.Get("key_id"))
	if keyID == "" {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "key_id query is required"})
	}
	var body keyPolicyConfig
	if len(req.Body) == 0 {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "request body is required"})
	}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	gatewayState.mu.Lock()
	defer gatewayState.mu.Unlock()
	for i := range gatewayState.config.KeyPolicies {
		if candidateKeyID(gatewayState.config.KeyPolicies[i]) != keyID {
			continue
		}
		body.KeyID = keyID
		if strings.TrimSpace(body.MatchAPIKey) == "" {
			body.MatchAPIKey = gatewayState.config.KeyPolicies[i].MatchAPIKey
		}
		gatewayState.config.KeyPolicies[i] = normalizeKeyPolicies([]keyPolicyConfig{body})[0]
		return jsonResponse(http.StatusOK, map[string]any{"ok": true})
	}
	return jsonResponse(http.StatusNotFound, map[string]any{"error": "policy not found"})
}

func routeAddRule(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	keyID := strings.TrimSpace(req.Query.Get("key_id"))
	if keyID == "" {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "key_id query is required"})
	}
	var body ruleConfig
	if len(req.Body) == 0 {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "request body is required"})
	}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	if strings.TrimSpace(body.ID) == "" {
		body.ID = "rule-" + time.Now().Format("20060102150405")
	}
	gatewayState.mu.Lock()
	defer gatewayState.mu.Unlock()
	for i := range gatewayState.config.KeyPolicies {
		if candidateKeyID(gatewayState.config.KeyPolicies[i]) != keyID {
			continue
		}
		gatewayState.config.KeyPolicies[i].Rules = append(gatewayState.config.KeyPolicies[i].Rules, body)
		gatewayState.config.KeyPolicies[i] = normalizeKeyPolicies([]keyPolicyConfig{gatewayState.config.KeyPolicies[i]})[0]
		return jsonResponse(http.StatusCreated, map[string]any{"ok": true})
	}
	return jsonResponse(http.StatusNotFound, map[string]any{"error": "policy not found"})
}

func routePatchRule(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	keyID := strings.TrimSpace(req.Query.Get("key_id"))
	ruleID := strings.TrimSpace(req.Query.Get("rule_id"))
	if keyID == "" || ruleID == "" {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "key_id and rule_id query are required"})
	}
	var body ruleConfig
	if len(req.Body) == 0 {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "request body is required"})
	}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	gatewayState.mu.Lock()
	defer gatewayState.mu.Unlock()
	for i := range gatewayState.config.KeyPolicies {
		if candidateKeyID(gatewayState.config.KeyPolicies[i]) != keyID {
			continue
		}
		for j := range gatewayState.config.KeyPolicies[i].Rules {
			if strings.TrimSpace(gatewayState.config.KeyPolicies[i].Rules[j].ID) != ruleID {
				continue
			}
			body.ID = ruleID
			gatewayState.config.KeyPolicies[i].Rules[j] = body
			gatewayState.config.KeyPolicies[i] = normalizeKeyPolicies([]keyPolicyConfig{gatewayState.config.KeyPolicies[i]})[0]
			return jsonResponse(http.StatusOK, map[string]any{"ok": true})
		}
	}
	return jsonResponse(http.StatusNotFound, map[string]any{"error": "rule not found"})
}

func routeDeleteRule(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	keyID := strings.TrimSpace(req.Query.Get("key_id"))
	ruleID := strings.TrimSpace(req.Query.Get("rule_id"))
	if keyID == "" || ruleID == "" {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "key_id and rule_id query are required"})
	}
	gatewayState.mu.Lock()
	defer gatewayState.mu.Unlock()
	for i := range gatewayState.config.KeyPolicies {
		if candidateKeyID(gatewayState.config.KeyPolicies[i]) != keyID {
			continue
		}
		out := make([]ruleConfig, 0, len(gatewayState.config.KeyPolicies[i].Rules))
		removed := false
		for _, rule := range gatewayState.config.KeyPolicies[i].Rules {
			if strings.TrimSpace(rule.ID) == ruleID {
				removed = true
				continue
			}
			out = append(out, rule)
		}
		if !removed {
			return jsonResponse(http.StatusNotFound, map[string]any{"error": "rule not found"})
		}
		gatewayState.config.KeyPolicies[i].Rules = out
		gatewayState.config.KeyPolicies[i] = normalizeKeyPolicies([]keyPolicyConfig{gatewayState.config.KeyPolicies[i]})[0]
		return jsonResponse(http.StatusOK, map[string]any{"ok": true})
	}
	return jsonResponse(http.StatusNotFound, map[string]any{"error": "policy not found"})
}

func routeResetUsage(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	keyID := strings.TrimSpace(req.Query.Get("key_id"))
	if keyID == "" {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "key_id query is required"})
	}
	gatewayState.mu.Lock()
	defer gatewayState.mu.Unlock()
	delete(gatewayState.usage, keyID)
	delete(gatewayState.requestWindow, keyID)
	return jsonResponse(http.StatusOK, map[string]any{"ok": true})
}

func routeUI(_ pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	return pluginapi.ManagementResponse{
		StatusCode: http.StatusOK,
		Headers:    http.Header{"Content-Type": []string{"text/html; charset=utf-8"}},
		Body:       []byte(gatewayUIHTML()),
	}, nil
}

func routeUsage(_ pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	return jsonResponse(http.StatusOK, map[string]any{
		"usage":                 gatewayState.listUsage(),
		"member_hits":           gatewayState.memberHitsSnapshot(),
		"rule_hits":             gatewayState.ruleHitsSnapshot(),
		"stage_hits":            gatewayState.stageHitsSnapshot(),
		"member_hits_last_5m":   gatewayState.memberHitsWindowSnapshot(5 * time.Minute),
		"rule_hits_last_5m":     gatewayState.ruleHitsWindowSnapshot(5 * time.Minute),
		"stage_hits_last_5m":    gatewayState.stageHitsWindowSnapshot(5 * time.Minute),
		"member_hits_last_hour": gatewayState.memberHitsWindowSnapshot(time.Hour),
		"rule_hits_last_hour":   gatewayState.ruleHitsWindowSnapshot(time.Hour),
		"stage_hits_last_hour":  gatewayState.stageHitsWindowSnapshot(time.Hour),
		"member_hits_last_24h":  gatewayState.memberHitsWindowSnapshot(24 * time.Hour),
		"rule_hits_last_24h":    gatewayState.ruleHitsWindowSnapshot(24 * time.Hour),
		"stage_hits_last_24h":   gatewayState.stageHitsWindowSnapshot(24 * time.Hour),
	})
}

func routeAudit(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	limit := 50
	if raw := strings.TrimSpace(req.Query.Get("limit")); raw != "" {
		if parsed, err := fmt.Sscanf(raw, "%d", &limit); parsed == 0 || err != nil {
			limit = 50
		}
	}
	filters := map[string]string{
		"decision":   strings.TrimSpace(req.Query.Get("decision")),
		"key":        strings.TrimSpace(req.Query.Get("key")),
		"rule":       strings.TrimSpace(req.Query.Get("rule")),
		"reason":     strings.TrimSpace(req.Query.Get("reason")),
		"policy":     strings.TrimSpace(req.Query.Get("policy")),
		"model":      strings.TrimSpace(req.Query.Get("model")),
		"provider":   strings.TrimSpace(req.Query.Get("provider")),
		"event_type": strings.TrimSpace(req.Query.Get("event_type")),
		"operator":   strings.TrimSpace(req.Query.Get("operator")),
		"member":     strings.TrimSpace(req.Query.Get("member")),
		"from":       strings.TrimSpace(req.Query.Get("from")),
		"to":         strings.TrimSpace(req.Query.Get("to")),
	}
	return jsonResponse(http.StatusOK, map[string]any{"items": gatewayState.listAudit(limit, filters)})
}

func routeAuditDetail(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	timeValue := strings.TrimSpace(req.Query.Get("time"))
	ruleID := strings.TrimSpace(req.Query.Get("rule"))
	reason := strings.TrimSpace(req.Query.Get("reason"))
	gatewayState.mu.RLock()
	defer gatewayState.mu.RUnlock()
	for i := len(gatewayState.auditLog) - 1; i >= 0; i-- {
		entry := gatewayState.auditLog[i]
		if timeValue != "" && entry.Time.Format(time.RFC3339Nano) != timeValue {
			continue
		}
		if ruleID != "" && !strings.EqualFold(strings.TrimSpace(entry.RuleID), ruleID) {
			continue
		}
		if reason != "" && !strings.EqualFold(strings.TrimSpace(entry.Reason), reason) {
			continue
		}
		return jsonResponse(http.StatusOK, entry)
	}
	return jsonResponse(http.StatusNotFound, map[string]any{"error": "audit entry not found"})
}

func routeAuditSummary(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	filters := map[string]string{
		"decision":   strings.TrimSpace(req.Query.Get("decision")),
		"key":        strings.TrimSpace(req.Query.Get("key")),
		"rule":       strings.TrimSpace(req.Query.Get("rule")),
		"reason":     strings.TrimSpace(req.Query.Get("reason")),
		"policy":     strings.TrimSpace(req.Query.Get("policy")),
		"model":      strings.TrimSpace(req.Query.Get("model")),
		"provider":   strings.TrimSpace(req.Query.Get("provider")),
		"event_type": strings.TrimSpace(req.Query.Get("event_type")),
		"operator":   strings.TrimSpace(req.Query.Get("operator")),
		"member":     strings.TrimSpace(req.Query.Get("member")),
		"from":       strings.TrimSpace(req.Query.Get("from")),
		"to":         strings.TrimSpace(req.Query.Get("to")),
	}
	return jsonResponse(http.StatusOK, gatewayState.auditSummary(filters))
}

func routeTemplates(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	search := strings.ToLower(strings.TrimSpace(req.Query.Get("search")))
	category := strings.ToLower(strings.TrimSpace(req.Query.Get("category")))
	scenario := strings.ToLower(strings.TrimSpace(req.Query.Get("scenario")))
	maturity := strings.ToLower(strings.TrimSpace(req.Query.Get("maturity")))
	tag := strings.ToLower(strings.TrimSpace(req.Query.Get("tag")))
	gatewayState.mu.RLock()
	defer gatewayState.mu.RUnlock()
	items := make([]ruleTemplate, 0, len(gatewayState.templates))
	for _, item := range gatewayState.templates {
		if search != "" {
			haystack := strings.ToLower(item.Name + " " + item.Description + " " + item.Scenario + " " + strings.Join(item.Tags, " "))
			if !strings.Contains(haystack, search) {
				continue
			}
		}
		if category != "" && strings.ToLower(strings.TrimSpace(item.Category)) != category {
			continue
		}
		if scenario != "" && strings.ToLower(strings.TrimSpace(item.Scenario)) != scenario {
			continue
		}
		if maturity != "" && strings.ToLower(strings.TrimSpace(item.Maturity)) != maturity {
			continue
		}
		if tag != "" && !containsFold(item.Tags, tag) {
			continue
		}
		items = append(items, item)
	}
	return jsonResponse(http.StatusOK, map[string]any{"items": items})
}

func routeAddTemplate(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	var body ruleTemplate
	if len(req.Body) == 0 {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "request body is required"})
	}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	if strings.TrimSpace(body.ID) == "" {
		body.ID = "template-" + time.Now().Format("20060102150405")
	}
	if strings.TrimSpace(body.Name) == "" {
		body.Name = body.ID
	}
	gatewayState.mu.Lock()
	defer gatewayState.mu.Unlock()
	for _, item := range gatewayState.templates {
		if strings.EqualFold(strings.TrimSpace(item.ID), strings.TrimSpace(body.ID)) {
			return jsonResponse(http.StatusConflict, map[string]any{"error": "template id already exists"})
		}
	}
	gatewayState.templates = append(gatewayState.templates, body)
	return jsonResponse(http.StatusCreated, map[string]any{"ok": true})
}

func routeDeleteTemplate(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	templateID := strings.TrimSpace(req.Query.Get("template_id"))
	if templateID == "" {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "template_id query is required"})
	}
	gatewayState.mu.Lock()
	defer gatewayState.mu.Unlock()
	out := make([]ruleTemplate, 0, len(gatewayState.templates))
	removed := false
	for _, item := range gatewayState.templates {
		if strings.EqualFold(strings.TrimSpace(item.ID), templateID) {
			removed = true
			continue
		}
		out = append(out, item)
	}
	if !removed {
		return jsonResponse(http.StatusNotFound, map[string]any{"error": "template not found"})
	}
	gatewayState.templates = out
	return jsonResponse(http.StatusOK, map[string]any{"ok": true})
}

func routeCloneTemplate(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	templateID := strings.TrimSpace(req.Query.Get("template_id"))
	if templateID == "" {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "template_id query is required"})
	}
	gatewayState.mu.Lock()
	defer gatewayState.mu.Unlock()
	for _, item := range gatewayState.templates {
		if !strings.EqualFold(strings.TrimSpace(item.ID), templateID) {
			continue
		}
		clone := item
		clone.ID = clone.ID + "-copy"
		clone.Name = clone.Name + " Copy"
		gatewayState.templates = append(gatewayState.templates, clone)
		return jsonResponse(http.StatusCreated, map[string]any{"ok": true, "id": clone.ID})
	}
	return jsonResponse(http.StatusNotFound, map[string]any{"error": "template not found"})
}

func routePatchTemplate(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	templateID := strings.TrimSpace(req.Query.Get("template_id"))
	if templateID == "" {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "template_id query is required"})
	}
	var body ruleTemplate
	if len(req.Body) == 0 {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "request body is required"})
	}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	gatewayState.mu.Lock()
	defer gatewayState.mu.Unlock()
	for i := range gatewayState.templates {
		if !strings.EqualFold(strings.TrimSpace(gatewayState.templates[i].ID), templateID) {
			continue
		}
		body.ID = templateID
		if strings.TrimSpace(body.Name) == "" {
			body.Name = templateID
		}
		gatewayState.templates[i] = body
		return jsonResponse(http.StatusOK, map[string]any{"ok": true})
	}
	return jsonResponse(http.StatusNotFound, map[string]any{"error": "template not found"})
}

func routeExportTemplates(_ pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	gatewayState.mu.RLock()
	defer gatewayState.mu.RUnlock()
	return jsonResponse(http.StatusOK, map[string]any{"items": append([]ruleTemplate(nil), gatewayState.templates...)})
}

func routeImportTemplates(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	var body struct {
		Items []ruleTemplate `json:"items"`
	}
	if len(req.Body) == 0 {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "request body is required"})
	}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	gatewayState.mu.Lock()
	defer gatewayState.mu.Unlock()
	seen := make(map[string]struct{}, len(gatewayState.templates))
	for _, item := range gatewayState.templates {
		seen[strings.ToLower(strings.TrimSpace(item.ID))] = struct{}{}
	}
	imported := 0
	for _, item := range body.Items {
		if strings.TrimSpace(item.ID) == "" {
			item.ID = "template-" + time.Now().Format("20060102150405")
		}
		key := strings.ToLower(strings.TrimSpace(item.ID))
		if _, exists := seen[key]; exists {
			continue
		}
		if strings.TrimSpace(item.Name) == "" {
			item.Name = item.ID
		}
		gatewayState.templates = append(gatewayState.templates, item)
		seen[key] = struct{}{}
		imported++
	}
	return jsonResponse(http.StatusOK, map[string]any{"ok": true, "imported": imported})
}

func routeDryRun(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	var body pluginapi.RequestInterceptRequest
	if len(req.Body) == 0 {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": "request body is required"})
	}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return jsonResponse(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	result := gatewayState.evaluateDryRun(body)
	return jsonResponse(http.StatusOK, result)
}

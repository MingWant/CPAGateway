package main

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/router-for-me/CLIProxyAPI/v7/sdk/pluginapi"
)

func newPluginState() *pluginState {
	return &pluginState{
		usage:           make(map[string]*usageCounter),
		requestWindow:   make(map[string][]time.Time),
		auditLog:        make([]auditEntry, 0, 128),
		templates:       builtInRuleTemplates(),
		memberHitCounts: make(map[string]int),
		ruleHitCounts:   make(map[string]int),
		stageHitCounts:  make(map[string]int),
		memberHitTimes:  make(map[string][]time.Time),
		ruleHitTimes:    make(map[string][]time.Time),
		stageHitTimes:   make(map[string][]time.Time),
		previewTokens:   make(map[string]previewTokenRecord),
	}
}

func (s *pluginState) apply(req pluginapi.RequestInterceptRequest, afterAuth bool) pluginapi.RequestInterceptResponse {
	result := s.evaluate(req, afterAuth)
	finalModel := strings.TrimSpace(result.FinalModel)
	if finalModel == "" {
		finalModel = strings.TrimSpace(req.Model)
	}
	s.recordMemberHit(finalModel)
	for _, ruleID := range result.MatchedRules {
		s.recordRuleHit(ruleID)
	}
	for _, trace := range result.StageTrace {
		s.recordStageHit(trace.Stage)
	}
	s.appendAudit(auditEntry{Time: time.Now(), PolicyID: stringMetadata(map[string]any{"v": result.Response.Headers.Get("X-Gateway-Policy-Id")}, "v"), PolicyName: result.Response.Headers.Get("X-Gateway-Policy-Name"), Decision: result.Decision, RuleID: result.RuleID, Reason: result.Reason, RequestedModel: req.Model, FinalModel: finalModel, Mirrors: stringListHeader(result.Response.Headers, "X-Gateway-Mirror-Models"), Path: stringMetadata(req.Metadata, "request_path"), APIKey: maskKey(stringMetadata(req.Metadata, "access.api_key")), Provider: providerFromModel(finalModel)})
	return result.Response
}

func (s *pluginState) evaluate(req pluginapi.RequestInterceptRequest, afterAuth bool) dryRunResult {
	now := time.Now()
	policy, keyID, displayName, maskedKey := s.lookupPolicy(req)
	if !policy.Enabled {
		resp := withGatewayMetadata(pluginapi.RequestInterceptResponse{Headers: cloneHeader(req.Headers), Body: cloneBytes(req.Body)}, map[string]string{"gateway.policy_id": keyID, "gateway.policy_name": displayName, "gateway.decision": "pass", "gateway.reason": "policy_disabled"})
		return dryRunResult{Decision: "pass", Reason: "policy_disabled", FinalModel: req.Model, Response: resp}
	}
	if reject := s.enforceLimits(keyID, displayName, maskedKey, policy.Limits, now, afterAuth); reject != nil {
		resp := withGatewayMetadata(*reject, map[string]string{"gateway.policy_id": keyID, "gateway.policy_name": displayName, "gateway.decision": "reject", "gateway.reason": reject.RejectCode})
		return dryRunResult{Decision: "reject", Reason: reject.RejectCode, FinalModel: req.Model, Response: resp}
	}
	result := applyRules(req, policy, afterAuth, now)
	result.Response = withGatewayMetadata(result.Response, map[string]string{"gateway.policy_id": keyID, "gateway.policy_name": displayName, "gateway.decision": result.Decision, "gateway.reason": result.Reason})
	return result
}

func (s *pluginState) lookupPolicy(req pluginapi.RequestInterceptRequest) (policy policyConfig, keyID, displayName, maskedKey string) {
	apiKey := strings.TrimSpace(stringMetadata(req.Metadata, "access.api_key"))
	requestedKeyID := firstNonEmpty(stringMetadata(req.Metadata, "access.key_id"), stringMetadata(req.Metadata, "gateway.key_id"))
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, candidate := range s.config.KeyPolicies {
		candidateID := candidateKeyID(candidate)
		if requestedKeyID != "" && candidateID == requestedKeyID {
			return policyConfig{Enabled: candidate.Enabled, Limits: candidate.Limits, Rules: cloneRuleConfigs(candidate.Rules), StagePolicy: cloneMap(candidate.StagePolicy)}, candidateID, candidate.DisplayName, maskKey(candidate.MatchAPIKey)
		}
		if strings.TrimSpace(candidate.MatchAPIKey) == "" || apiKey == "" || candidate.MatchAPIKey != apiKey {
			continue
		}
		return policyConfig{Enabled: candidate.Enabled, Limits: candidate.Limits, Rules: cloneRuleConfigs(candidate.Rules), StagePolicy: cloneMap(candidate.StagePolicy)}, candidateID, candidate.DisplayName, maskKey(candidate.MatchAPIKey)
	}
	defaultKeyID := stableKeyID(apiKey)
	if requestedKeyID != "" {
		defaultKeyID = requestedKeyID
	}
	return clonePolicyConfig(s.config.Default), defaultKeyID, "default", maskKey(apiKey)
}

func (s *pluginState) enforceLimits(keyID, displayName, maskedKey string, limits limitConfig, now time.Time, afterAuth bool) *pluginapi.RequestInterceptResponse {
	if keyID == "" {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	entry := s.usage[keyID]
	if entry == nil {
		entry = &usageCounter{DisplayName: displayName, MaskedKey: maskedKey}
		s.usage[keyID] = entry
	}
	entry.DisplayName = displayName
	entry.MaskedKey = maskedKey
	entry.LastSeenAt = now
	today := now.Format("2006-01-02")
	if entry.DayBucket != today {
		entry.RequestsToday = 0
		entry.DayBucket = today
	}
	if !withinAbsoluteWindow(limits, now) || !withinSchedules(limits.Schedules, now) {
		return &pluginapi.RequestInterceptResponse{Reject: true, RejectStatusCode: http.StatusForbidden, RejectMessage: "gateway schedule rejected request", RejectCode: "gateway_schedule_denied"}
	}
	window := pruneRecent(s.requestWindow[keyID], now.Add(-time.Minute))
	s.requestWindow[keyID] = window
	if !afterAuth {
		if limits.RequestsPerMin > 0 && len(window) >= limits.RequestsPerMin {
			return &pluginapi.RequestInterceptResponse{Reject: true, RejectStatusCode: http.StatusTooManyRequests, RejectMessage: "gateway rate limit exceeded", RejectCode: "gateway_rate_limit_exceeded"}
		}
		if limits.RequestsPerDay > 0 && entry.RequestsToday >= limits.RequestsPerDay {
			return &pluginapi.RequestInterceptResponse{Reject: true, RejectStatusCode: http.StatusForbidden, RejectMessage: "gateway daily quota exceeded", RejectCode: "gateway_quota_exceeded"}
		}
		if limits.MaxInflight > 0 && entry.Inflight >= limits.MaxInflight {
			return &pluginapi.RequestInterceptResponse{Reject: true, RejectStatusCode: http.StatusTooManyRequests, RejectMessage: "gateway concurrency limit exceeded", RejectCode: "gateway_concurrency_exceeded"}
		}
		entry.RequestsToday++
		entry.RequestsMinute = len(window) + 1
		entry.Inflight++
		s.requestWindow[keyID] = append(window, now)
		return nil
	}
	if entry.Inflight > 0 {
		entry.Inflight--
	}
	entry.RequestsMinute = len(window)
	return nil
}

func (s *pluginState) listKeys() []map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]map[string]any, 0, len(s.config.KeyPolicies))
	for _, item := range s.config.KeyPolicies {
		out = append(out, map[string]any{"key_id": candidateKeyID(item), "display_name": item.DisplayName, "masked_key": maskKey(item.MatchAPIKey), "enabled": item.Enabled})
	}
	return out
}

func (s *pluginState) listUsage() []usageEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]usageEntry, 0, len(s.usage))
	for keyID, entry := range s.usage {
		window := pruneRecent(append([]time.Time(nil), s.requestWindow[keyID]...), time.Now().Add(-time.Minute))
		out = append(out, usageEntry{KeyID: keyID, DisplayName: entry.DisplayName, MaskedKey: entry.MaskedKey, RequestsToday: entry.RequestsToday, RequestsMinute: len(window), Inflight: entry.Inflight, LastSeenAt: entry.LastSeenAt})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].KeyID < out[j].KeyID })
	return out
}

func (s *pluginState) listAudit(limit int, filters map[string]string) []auditEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if limit <= 0 || limit > len(s.auditLog) {
		limit = len(s.auditLog)
	}
	out := make([]auditEntry, 0, limit)
	for i := len(s.auditLog) - 1; i >= 0 && len(out) < limit; i-- {
		entry := s.auditLog[i]
		if !auditEntryMatches(entry, filters) {
			continue
		}
		out = append(out, entry)
	}
	return out
}

func (s *pluginState) auditSummary(filters map[string]string) auditSummary {
	items := s.listAudit(len(s.auditLog), filters)
	summary := auditSummary{TotalByDecision: map[string]int{}, TotalByReason: map[string]int{}, TotalByRule: map[string]int{}, TotalByPolicy: map[string]int{}, TotalByModel: map[string]int{}, TotalByProvider: map[string]int{}, Timeline: make([]auditBucket, 0)}
	timeline := map[string]int{}
	for _, item := range items {
		summary.TotalByDecision[item.Decision]++
		if strings.TrimSpace(item.Reason) != "" {
			summary.TotalByReason[item.Reason]++
		}
		if strings.TrimSpace(item.RuleID) != "" {
			summary.TotalByRule[item.RuleID]++
		}
		if strings.TrimSpace(item.PolicyName) != "" {
			summary.TotalByPolicy[item.PolicyName]++
		}
		if strings.TrimSpace(item.FinalModel) != "" {
			summary.TotalByModel[item.FinalModel]++
		}
		if strings.TrimSpace(item.Provider) != "" {
			summary.TotalByProvider[item.Provider]++
		}
		bucket := item.Time.Format("2006-01-02 15:04")
		timeline[bucket]++
	}
	keys := make([]string, 0, len(timeline))
	for key := range timeline {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		summary.Timeline = append(summary.Timeline, auditBucket{Window: key, Count: timeline[key]})
	}
	return summary
}

func builtInRuleTemplates() []ruleTemplate {
	return []ruleTemplate{
		{ID: "route-openai", Name: "Route To OpenAI Model", Category: "routing", Scenario: "model-migration", Maturity: "stable", Tags: []string{"migration", "openai", "rewrite"}, Description: "Routes a matched model directly to another upstream model.", Rule: ruleConfig{ID: "route-template", Enabled: true, Priority: 10, OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{RouteToModel: "openai/gpt-5.4"}}},
		{ID: "weighted-split", Name: "Weighted Split Route", Category: "routing", Scenario: "traffic-split", Maturity: "beta", Tags: []string{"weighted", "ab-test", "routing"}, Description: "Splits matched traffic across multiple upstream models using deterministic weights.", Rule: ruleConfig{ID: "weighted-template", Enabled: true, Priority: 15, OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{ShardBy: "api_key", WeightedRoutes: []weightedRoute{{Model: "openai/gpt-5.4", Weight: 80}, {Model: "codex/gpt-5.4", Weight: 20}}}}},
		{ID: "fallback-mini", Name: "Fallback To Mini Chain", Category: "fallback", Scenario: "cost-control", Maturity: "stable", Tags: []string{"fallback", "cost", "resilience"}, Description: "Falls back through smaller models when the primary route is not preferred.", Rule: ruleConfig{ID: "fallback-template", Enabled: true, Priority: 20, OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{FallbackModels: []string{"openai/gpt-5.4-mini", "openai/gpt-4.1-mini"}}}},
		{ID: "mirror-safety", Name: "Mirror For Audit", Category: "routing", Scenario: "shadow-release", Maturity: "beta", Tags: []string{"mirror", "shadow", "audit"}, Description: "Tags a request with mirror targets for audit or future shadow routing.", Rule: ruleConfig{ID: "mirror-template", Enabled: true, Priority: 25, OnMatch: "continue", Match: matchConfig{Paths: []string{"/v1/responses"}}, Actions: actionConfig{MirrorModels: []string{"openai/gpt-5.4-mini"}, TagMetadata: map[string]string{"mirror.mode": "shadow"}}}},
		{ID: "deny-provider", Name: "Deny Provider", Category: "security", Scenario: "provider-guardrail", Maturity: "stable", Tags: []string{"security", "deny", "provider"}, Description: "Blocks requests that resolve to a disallowed provider.", Rule: ruleConfig{ID: "deny-template", Enabled: true, Priority: 30, OnMatch: "stop", Match: matchConfig{Providers: []string{"claude"}}, Actions: actionConfig{Deny: &denyConfig{StatusCode: 403, Message: "provider denied", Code: "gateway_provider_denied"}}}},
	}
}

func (s *pluginState) appendAudit(entry auditEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.appendAuditLocked(entry)
}

func (s *pluginState) appendAuditLocked(entry auditEntry) {
	s.auditLog = append(s.auditLog, entry)
	if len(s.auditLog) > 200 {
		s.auditLog = append([]auditEntry(nil), s.auditLog[len(s.auditLog)-200:]...)
	}
}

func (s *pluginState) prunePreviewTokensLocked(now time.Time) {
	for token, record := range s.previewTokens {
		if now.Sub(record.IssuedAt) > 10*time.Minute {
			delete(s.previewTokens, token)
		}
	}
}

func (s *pluginState) recordMemberHit(model string) {
	model = strings.ToLower(strings.TrimSpace(model))
	if model == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.memberHitCounts == nil {
		s.memberHitCounts = make(map[string]int)
	}
	if s.memberHitTimes == nil {
		s.memberHitTimes = make(map[string][]time.Time)
	}
	s.memberHitCounts[model]++
	now := time.Now()
	s.memberHitTimes[model] = append(pruneRecent(s.memberHitTimes[model], now.Add(-24*time.Hour)), now)
}

func (s *pluginState) memberHitCount(model string) int {
	model = strings.ToLower(strings.TrimSpace(model))
	if model == "" {
		return 0
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.memberHitCounts[model]
}

func (s *pluginState) memberHitsSnapshot() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]int, len(s.memberHitCounts))
	for key, value := range s.memberHitCounts {
		out[key] = value
	}
	return out
}

func (s *pluginState) recordRuleHit(ruleID string) {
	ruleID = strings.TrimSpace(ruleID)
	if ruleID == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ruleHitCounts == nil {
		s.ruleHitCounts = make(map[string]int)
	}
	if s.ruleHitTimes == nil {
		s.ruleHitTimes = make(map[string][]time.Time)
	}
	s.ruleHitCounts[ruleID]++
	now := time.Now()
	s.ruleHitTimes[ruleID] = append(pruneRecent(s.ruleHitTimes[ruleID], now.Add(-24*time.Hour)), now)
}

func (s *pluginState) recordStageHit(stage string) {
	stage = strings.TrimSpace(stage)
	if stage == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stageHitCounts == nil {
		s.stageHitCounts = make(map[string]int)
	}
	if s.stageHitTimes == nil {
		s.stageHitTimes = make(map[string][]time.Time)
	}
	s.stageHitCounts[stage]++
	now := time.Now()
	s.stageHitTimes[stage] = append(pruneRecent(s.stageHitTimes[stage], now.Add(-24*time.Hour)), now)
}

func (s *pluginState) ruleHitsSnapshot() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]int, len(s.ruleHitCounts))
	for key, value := range s.ruleHitCounts {
		out[key] = value
	}
	return out
}

func (s *pluginState) stageHitsSnapshot() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]int, len(s.stageHitCounts))
	for key, value := range s.stageHitCounts {
		out[key] = value
	}
	return out
}

func (s *pluginState) memberHitsWindowSnapshot(window time.Duration) map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return windowedHitSnapshotLocked(s.memberHitTimes, window, time.Now())
}

func (s *pluginState) ruleHitsWindowSnapshot(window time.Duration) map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return windowedHitSnapshotLocked(s.ruleHitTimes, window, time.Now())
}

func (s *pluginState) stageHitsWindowSnapshot(window time.Duration) map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return windowedHitSnapshotLocked(s.stageHitTimes, window, time.Now())
}

func windowedHitSnapshotLocked(source map[string][]time.Time, window time.Duration, now time.Time) map[string]int {
	out := make(map[string]int, len(source))
	cutoff := now.Add(-window)
	for key, items := range source {
		count := 0
		for _, item := range items {
			if item.Before(cutoff) {
				continue
			}
			count++
		}
		if count > 0 {
			out[key] = count
		}
	}
	return out
}

func auditEntryMatches(entry auditEntry, filters map[string]string) bool {
	if filters == nil {
		return true
	}
	if value := strings.TrimSpace(filters["decision"]); value != "" && !strings.EqualFold(entry.Decision, value) {
		return false
	}
	if value := strings.TrimSpace(filters["key"]); value != "" && !strings.Contains(strings.ToLower(entry.APIKey), strings.ToLower(value)) {
		return false
	}
	if value := strings.TrimSpace(filters["rule"]); value != "" && !strings.Contains(strings.ToLower(entry.RuleID), strings.ToLower(value)) {
		return false
	}
	if value := strings.TrimSpace(filters["reason"]); value != "" && !strings.Contains(strings.ToLower(entry.Reason), strings.ToLower(value)) {
		return false
	}
	if value := strings.TrimSpace(filters["policy"]); value != "" && !strings.Contains(strings.ToLower(entry.PolicyName), strings.ToLower(value)) {
		return false
	}
	if value := strings.TrimSpace(filters["model"]); value != "" {
		modelHit := strings.Contains(strings.ToLower(entry.RequestedModel), strings.ToLower(value)) || strings.Contains(strings.ToLower(entry.FinalModel), strings.ToLower(value))
		if !modelHit {
			return false
		}
	}
	if value := strings.TrimSpace(filters["provider"]); value != "" && !strings.Contains(strings.ToLower(entry.Provider), strings.ToLower(value)) {
		return false
	}
	if value := strings.TrimSpace(filters["event_type"]); value != "" && !strings.Contains(strings.ToLower(entry.EventType), strings.ToLower(value)) {
		return false
	}
	if value := strings.TrimSpace(filters["operator"]); value != "" && !strings.Contains(strings.ToLower(entry.OperatorAction), strings.ToLower(value)) {
		return false
	}
	if value := strings.TrimSpace(filters["member"]); value != "" && !strings.Contains(strings.ToLower(entry.TargetMember), strings.ToLower(value)) {
		return false
	}
	if value := strings.TrimSpace(filters["from"]); value != "" {
		if from, err := time.Parse(time.RFC3339, value); err == nil && entry.Time.Before(from) {
			return false
		}
	}
	if value := strings.TrimSpace(filters["to"]); value != "" {
		if to, err := time.Parse(time.RFC3339, value); err == nil && entry.Time.After(to) {
			return false
		}
	}
	return true
}

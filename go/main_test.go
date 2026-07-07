package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/router-for-me/CLIProxyAPI/v7/sdk/pluginapi"
)

func TestApplyRulesRewriteModelAndResponsesPayload(t *testing.T) {
	policy := policyConfig{
		Enabled: true,
		Rules: []ruleConfig{{
			Enabled: true,
			OnMatch: "stop",
			Match:   matchConfig{PathKinds: []string{"/responses"}, Models: []string{"gpt-5.5"}},
			Actions: actionConfig{RewriteModel: "gpt-5.4", RewriteEndpoint: &endpointRewrite{TargetModel: "gpt-5.4"}},
		}},
	}
	resp := applyRules(pluginapi.RequestInterceptRequest{
		SourceFormat: "openai-response",
		Model:        "gpt-5.5",
		Headers:      http.Header{},
		Body:         []byte(`{"model":"gpt-5.5","messages":[{"role":"user","content":"hi"}]}`),
		Metadata:     map[string]any{"request_path": "/v1/responses"},
	}, policy, false, time.Date(2026, 7, 3, 10, 0, 0, 0, time.Local))
	if string(resp.Response.Body) == "" || !contains(string(resp.Response.Body), `"model":"gpt-5.4"`) || !contains(string(resp.Response.Body), `"input"`) {
		t.Fatalf("rewritten body = %s", resp.Response.Body)
	}
}

func TestApplyRulesRejectsDisallowedProvider(t *testing.T) {
	policy := policyConfig{
		Enabled: true,
		Rules:   []ruleConfig{{Enabled: true, OnMatch: "stop", Actions: actionConfig{AllowOnlyProviders: []string{"openai", "codex"}}}},
	}
	resp := applyRules(pluginapi.RequestInterceptRequest{Model: "claude/gpt-5.5", Headers: http.Header{}, Body: []byte(`{"model":"claude/gpt-5.5"}`)}, policy, false, time.Now())
	if !resp.Response.Reject || resp.Response.RejectCode != "gateway_provider_denied" {
		t.Fatalf("response = %#v, want provider deny", resp)
	}
}

func TestPluginStateEnforcesPerMinuteLimit(t *testing.T) {
	state := newPluginState()
	state.config = pluginConfig{Default: policyConfig{Enabled: true, Limits: limitConfig{RequestsPerMin: 1}}}
	req := pluginapi.RequestInterceptRequest{Headers: http.Header{}, Body: []byte(`{"model":"gpt-5.5"}`), Metadata: map[string]any{"access.api_key": "abc"}}
	first := state.apply(req, false)
	if first.Reject {
		t.Fatalf("first request rejected unexpectedly: %#v", first)
	}
	second := state.apply(req, false)
	if !second.Reject || second.RejectCode != "gateway_rate_limit_exceeded" {
		t.Fatalf("second request = %#v, want rate limit reject", second)
	}
}

func TestRouteDryRunDoesNotMutateUsage(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = pluginConfig{Default: policyConfig{Enabled: true, Limits: limitConfig{RequestsPerDay: 1, RequestsPerMin: 1, MaxInflight: 1}}}
	body, err := json.Marshal(pluginapi.RequestInterceptRequest{Model: "gpt-5.5", Headers: http.Header{}, Metadata: map[string]any{"access.api_key": "abc"}})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	for i := 0; i < 2; i++ {
		resp, err := routeDryRun(pluginapi.ManagementRequest{Body: body})
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("routeDryRun(%d) = %#v, %v", i, resp, err)
		}
	}
	if _, ok := gatewayState.usage[stableKeyID("abc")]; ok {
		t.Fatalf("dry-run created usage entry: %#v", gatewayState.usage[stableKeyID("abc")])
	}
	if got := len(gatewayState.requestWindow[stableKeyID("abc")]); got != 0 {
		t.Fatalf("dry-run request window len = %d, want 0", got)
	}
}

func TestPluginStateEnforcesTokenQuotaAfterUsage(t *testing.T) {
	state := newPluginState()
	state.config = pluginConfig{Default: policyConfig{Enabled: true, Limits: limitConfig{TokensPerDay: 10}}}
	req := pluginapi.RequestInterceptRequest{Headers: http.Header{}, Body: []byte(`{"model":"gpt-5.5"}`), Metadata: map[string]any{"access.api_key": "abc"}}
	first := state.apply(req, false)
	if first.Reject {
		t.Fatalf("first request rejected unexpectedly: %#v", first)
	}
	state.recordUsage(pluginapi.UsageRecord{APIKey: "abc", Detail: pluginapi.UsageDetail{TotalTokens: 11}})
	usage := state.usage[stableKeyID("abc")]
	if usage == nil || usage.Inflight != 0 || usage.TokensToday != 11 {
		t.Fatalf("usage after record = %#v, want inflight 0 and 11 tokens", usage)
	}
	second := state.apply(req, false)
	if !second.Reject || second.RejectCode != "gateway_token_quota_exceeded" {
		t.Fatalf("second request = %#v, want token quota reject", second)
	}
}

func TestApplyRulesReturnsDecisionAndRouteToModel(t *testing.T) {
	policy := policyConfig{Enabled: true, Rules: []ruleConfig{{ID: "route-1", Enabled: true, OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{RouteToModel: "openai/gpt-5.4"}}}}
	result := applyRules(pluginapi.RequestInterceptRequest{Model: "gpt-5.5", Headers: http.Header{}, Body: []byte(`{"model":"gpt-5.5"}`)}, policy, false, time.Now())
	if result.Decision != "rewrite" || result.RuleID != "route-1" || result.FinalModel != "openai/gpt-5.4" {
		t.Fatalf("result = %#v", result)
	}
}

func TestApplyRulesFallsBackToFallbackModel(t *testing.T) {
	policy := policyConfig{Enabled: true, Rules: []ruleConfig{{ID: "fallback-1", Enabled: true, OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{FallbackModels: []string{"openai/gpt-5.4-mini"}}}}}
	result := applyRules(pluginapi.RequestInterceptRequest{Model: "gpt-5.5", Headers: http.Header{}, Body: []byte(`{"model":"gpt-5.5"}`)}, policy, false, time.Now())
	if result.Reason != "fallback_model" || result.FinalModel != "openai/gpt-5.4-mini" {
		t.Fatalf("result = %#v", result)
	}
}

func TestAuditSummaryAggregatesByPolicy(t *testing.T) {
	state := newPluginState()
	state.auditLog = []auditEntry{{Decision: "rewrite", PolicyName: "Key A"}, {Decision: "reject", PolicyName: "Key A"}, {Decision: "pass", PolicyName: "Key B"}}
	summary := state.auditSummary(nil)
	if summary.TotalByPolicy["Key A"] != 2 || summary.TotalByPolicy["Key B"] != 1 {
		t.Fatalf("summary by policy = %#v", summary.TotalByPolicy)
	}
}

func TestRouteCloneTemplate(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.templates = []ruleTemplate{{ID: "tpl-1", Name: "Template 1", Category: "custom", Description: "desc", Rule: ruleConfig{ID: "rule-1", Enabled: true, Priority: 10, OnMatch: "stop"}}}
	resp, err := routeCloneTemplate(pluginapi.ManagementRequest{Query: map[string][]string{"template_id": {"tpl-1"}}})
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("routeCloneTemplate() = %#v, %v", resp, err)
	}
	resp, err = routeCloneTemplate(pluginapi.ManagementRequest{Query: map[string][]string{"template_id": {"tpl-1"}}})
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("routeCloneTemplate() second clone = %#v, %v", resp, err)
	}
	if len(gatewayState.templates) != 3 {
		t.Fatalf("templates len = %d, want 3", len(gatewayState.templates))
	}
	if gatewayState.templates[1].ID == gatewayState.templates[2].ID {
		t.Fatalf("clone template IDs should be unique: %#v", gatewayState.templates)
	}
}

func TestRouteImportAndExportTemplates(t *testing.T) {
	gatewayState = newPluginState()
	resp, err := routeImportTemplates(pluginapi.ManagementRequest{Body: []byte(`{"items":[{"id":"tpl-import","name":"Imported","category":"custom","description":"desc","rule":{"id":"rule-1","enabled":true,"priority":10,"on_match":"stop","match":{"models":["gpt-5.5"]},"actions":{"route_to_model":"openai/gpt-5.4"}}}]}`)})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("routeImportTemplates() = %#v, %v", resp, err)
	}
	resp, err = routeExportTemplates(pluginapi.ManagementRequest{})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("routeExportTemplates() = %#v, %v", resp, err)
	}
	if len(gatewayState.templates) == 0 {
		t.Fatal("expected templates after import")
	}
}

func TestRouteImportTemplatesAssignsUniqueIDs(t *testing.T) {
	gatewayState = newPluginState()
	resp, err := routeImportTemplates(pluginapi.ManagementRequest{Body: []byte(`{"items":[{"name":"A"},{"name":"B"}]}`)})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("routeImportTemplates() = %#v, %v", resp, err)
	}
	seen := map[string]struct{}{}
	for _, item := range gatewayState.templates {
		if strings.TrimSpace(item.ID) == "" {
			t.Fatalf("template missing ID: %#v", item)
		}
		if _, exists := seen[item.ID]; exists {
			t.Fatalf("duplicate template ID after import: %#v", gatewayState.templates)
		}
		seen[item.ID] = struct{}{}
	}
}

func TestRouteAddAndDeleteTemplate(t *testing.T) {
	gatewayState = newPluginState()
	resp, err := routeAddTemplate(pluginapi.ManagementRequest{Body: []byte(`{"id":"tpl-1","name":"Template 1","description":"desc","rule":{"id":"rule-1","enabled":true,"priority":10,"on_match":"stop","match":{"models":["gpt-5.5"]},"actions":{"route_to_model":"openai/gpt-5.4"}}}`)})
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("routeAddTemplate() = %#v, %v", resp, err)
	}
	if len(gatewayState.templates) == 0 {
		t.Fatal("expected template to be stored")
	}
	resp, err = routeDeleteTemplate(pluginapi.ManagementRequest{Query: map[string][]string{"template_id": {"tpl-1"}}})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("routeDeleteTemplate() = %#v, %v", resp, err)
	}
}

func TestGatewayUIIncludesContentSecurityPolicy(t *testing.T) {
	html := gatewayUIHTML()
	if !contains(html, "Content-Security-Policy") || !contains(html, "script-src 'nonce-gateway-ui'") || !contains(html, `<script nonce="gateway-ui">`) {
		t.Fatalf("gateway UI missing CSP nonce protection")
	}
}

func TestRouteUIIncludesContentSecurityPolicyHeader(t *testing.T) {
	gatewayState = newPluginState()
	resp, err := routeUI(pluginapi.ManagementRequest{})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("routeUI() = %#v, %v", resp, err)
	}
	if got := resp.Headers.Get("Content-Security-Policy"); got != gatewayContentSecurityPolicy {
		t.Fatalf("CSP header = %q, want %q", got, gatewayContentSecurityPolicy)
	}
}

func TestGatewayUIInstallsDynamicHTMLSanitizer(t *testing.T) {
	html := gatewayUIHTML()
	if !contains(html, "function sanitizeHTML") || !contains(html, "Object.defineProperty(Element.prototype, 'innerHTML'") {
		t.Fatalf("gateway UI missing dynamic HTML sanitizer")
	}
	if !contains(html, "name.startsWith('on')") || !contains(html, "value.startsWith('javascript:')") {
		t.Fatalf("gateway UI sanitizer is missing event or javascript URL filtering")
	}
	if !contains(html, "name === 'style'") || !contains(html, "style,svg,math,form") {
		t.Fatalf("gateway UI sanitizer is missing style or embedded-content filtering")
	}
}

func TestGatewayUIDefinesRouteBuilderHelpers(t *testing.T) {
	html := gatewayUIHTML()
	required := []string{
		"function parseCSV",
		"function weightedRoutes",
		"function renderWeightedRoutes",
		"function renderConditionGroups",
		"function apiURL",
		"gateway_token=' + encodeURIComponent(gatewayTokenParam)",
		"addFailoverHopBtn').addEventListener",
		"const failover = el('ruleFailoverChainInput').value.trim()",
		"const routePoolName = el('routePoolNameInput').value.trim()",
	}
	for _, item := range required {
		if !contains(html, item) {
			t.Fatalf("gateway UI missing route builder helper %q", item)
		}
	}
}

func TestGatewayUIAvoidsDynamicInnerHTMLRendering(t *testing.T) {
	html := gatewayUIHTML()
	forbidden := []string{"root.innerHTML", "node.innerHTML", "policyNode.innerHTML", "stageNode.innerHTML", "child.innerHTML"}
	for _, item := range forbidden {
		if strings.Contains(html, item) {
			t.Fatalf("gateway UI should not use dynamic innerHTML rendering: found %q", item)
		}
	}
}

func TestAuditSummaryAggregatesByDecision(t *testing.T) {
	state := newPluginState()
	state.auditLog = []auditEntry{{Decision: "rewrite", RuleID: "route-1", Reason: "route_to_model"}, {Decision: "reject", RuleID: "deny-1", Reason: "gateway_provider_denied"}, {Decision: "reject", RuleID: "deny-1", Reason: "gateway_provider_denied"}}
	summary := state.auditSummary(nil)
	if summary.TotalByDecision["reject"] != 2 || summary.TotalByDecision["rewrite"] != 1 {
		t.Fatalf("summary = %#v", summary)
	}
	if summary.TotalByRule["deny-1"] != 2 {
		t.Fatalf("summary rule counts = %#v", summary.TotalByRule)
	}
}

func TestRouteAuditDetailFindsEntry(t *testing.T) {
	gatewayState = newPluginState()
	entry := auditEntry{Time: time.Now(), Decision: "reject", RuleID: "deny-1", Reason: "gateway_provider_denied"}
	gatewayState.auditLog = []auditEntry{entry}
	resp, err := routeAuditDetail(pluginapi.ManagementRequest{Query: map[string][]string{"time": {entry.Time.Format(time.RFC3339Nano)}, "rule": {"deny-1"}, "reason": {"gateway_provider_denied"}}})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("routeAuditDetail() = %#v, %v", resp, err)
	}
}

func TestListAuditSupportsFilters(t *testing.T) {
	state := newPluginState()
	state.auditLog = []auditEntry{{Decision: "rewrite", RuleID: "route-1", Reason: "route_to_model", APIKey: "abc***xyz"}, {Decision: "reject", RuleID: "deny-1", Reason: "gateway_provider_denied", APIKey: "def***uvw"}}
	items := state.listAudit(10, map[string]string{"decision": "reject", "reason": "provider"})
	if len(items) != 1 || items[0].RuleID != "deny-1" {
		t.Fatalf("filtered audit = %#v", items)
	}
}

func TestPluginStateWritesAuditEntries(t *testing.T) {
	state := newPluginState()
	state.config = pluginConfig{Default: policyConfig{Enabled: true, Rules: []ruleConfig{{ID: "audit-1", Enabled: true, OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{RouteToModel: "openai/gpt-5.4"}}}}}
	resp := state.apply(pluginapi.RequestInterceptRequest{Model: "gpt-5.5", Headers: http.Header{}, Body: []byte(`{"model":"gpt-5.5"}`), Metadata: map[string]any{"request_path": "/v1/responses", "access.api_key": "secret-key"}}, false)
	if resp.Headers.Get("X-Gateway-Decision") == "" {
		t.Fatalf("missing gateway decision headers: %#v", resp.Headers)
	}
	items := state.listAudit(10, nil)
	if len(items) == 0 || items[0].Decision == "" {
		t.Fatalf("audit items = %#v", items)
	}
}

func TestMatchRuleSupportsNestedAnyOfPathAndModel(t *testing.T) {
	ok, reason := matchRule(pluginapi.RequestInterceptRequest{SourceFormat: "openai-response", Headers: http.Header{}, Metadata: map[string]any{"request_path": "/v1/responses"}}, "gpt-5.5", matchConfig{AnyOf: []matchConfig{{Models: []string{"other-model"}}, {Paths: []string{"/v1/responses"}}}}, time.Now())
	if !ok || reason != "matched" {
		t.Fatalf("matchRule() = %v, %q, want true/matched", ok, reason)
	}
}

func TestMatchRuleSupportsAnyOfAndAllOf(t *testing.T) {
	ok, reason := matchRule(pluginapi.RequestInterceptRequest{SourceFormat: "openai-response", Headers: http.Header{}, Body: []byte(`{"service_tier":"priority"}`), Metadata: map[string]any{"request_path": "/v1/responses", "request.query": "mode=strict&tenant=a"}}, "gpt-5.5", matchConfig{AnyOf: []matchConfig{{Models: []string{"other-model"}}, {Paths: []string{"/v1/responses"}}}, AllOf: []matchConfig{{Query: map[string]string{"mode": "strict"}}, {BodyContains: map[string]string{"service_tier": "priority"}}}}, time.Now())
	if !ok || reason != "matched" {
		t.Fatalf("matchRule() = %v, %q, want true/matched", ok, reason)
	}
}

func TestMatchRuleSupportsQueryAndBodyConditions(t *testing.T) {
	ok, reason := matchRule(pluginapi.RequestInterceptRequest{SourceFormat: "openai-response", Headers: http.Header{}, Body: []byte(`{"service_tier":"priority"}`), Metadata: map[string]any{"request_path": "/v1/responses", "request.query": "mode=strict&tenant=a"}}, "gpt-5.5", matchConfig{Query: map[string]string{"mode": "strict"}, BodyContains: map[string]string{"service_tier": "priority"}}, time.Now())
	if !ok || reason != "matched" {
		t.Fatalf("matchRule() = %v, %q, want true/matched", ok, reason)
	}
}

func TestMatchRuleSupportsPathHeaderAndMetadataConditions(t *testing.T) {
	ok, reason := matchRule(pluginapi.RequestInterceptRequest{SourceFormat: "openai-response", Headers: http.Header{"X-Test": {"yes"}}, Metadata: map[string]any{"request_path": "/v1/responses", "client.tag": "tenant-a"}}, "gpt-5.5", matchConfig{Paths: []string{"/v1/responses"}, Headers: map[string]string{"X-Test": "yes"}, MetadataContains: map[string]string{"client.tag": "tenant"}}, time.Now())
	if !ok || reason != "matched" {
		t.Fatalf("matchRule() = %v, %q, want true/matched", ok, reason)
	}
}

func TestMatchRuleSupportsOvernightTimeWindow(t *testing.T) {
	now := time.Date(2026, 7, 5, 1, 30, 0, 0, time.Local)
	ok, reason := matchRule(pluginapi.RequestInterceptRequest{Headers: http.Header{}}, "gpt-5.5", matchConfig{Start: "22:00", End: "02:00"}, now)
	if !ok || reason != "matched" {
		t.Fatalf("overnight matchRule() = %v, %q, want true/matched", ok, reason)
	}
}

func TestWithinSchedulesSupportsOvernightWindow(t *testing.T) {
	now := time.Date(2026, 7, 5, 1, 30, 0, 0, time.Local)
	if !withinSchedules([]scheduleConfig{{Start: "22:00", End: "02:00"}}, now) {
		t.Fatal("expected overnight schedule to match at 01:30")
	}
	daytime := time.Date(2026, 7, 5, 12, 0, 0, 0, time.Local)
	if withinSchedules([]scheduleConfig{{Start: "22:00", End: "02:00"}}, daytime) {
		t.Fatal("overnight schedule should not match at noon")
	}
}

func TestPluginStateUsageRecordReleasesInflight(t *testing.T) {
	state := newPluginState()
	state.config = pluginConfig{Default: policyConfig{Enabled: true, Limits: limitConfig{MaxInflight: 1}}}
	req := pluginapi.RequestInterceptRequest{Headers: http.Header{}, Body: []byte(`{"model":"gpt-5.5"}`), Metadata: map[string]any{"access.api_key": "abc"}}
	first := state.apply(req, false)
	if first.Reject {
		t.Fatalf("before-auth rejected unexpectedly: %#v", first)
	}
	if got := state.usage[stableKeyID("abc")].Inflight; got != 1 {
		t.Fatalf("inflight = %d, want 1", got)
	}
	second := state.apply(req, true)
	if second.Reject {
		t.Fatalf("after-auth rejected unexpectedly: %#v", second)
	}
	if got := state.usage[stableKeyID("abc")].Inflight; got != 1 {
		t.Fatalf("inflight after after-auth = %d, want 1 until usage record", got)
	}
	state.recordUsage(pluginapi.UsageRecord{APIKey: "abc", Detail: pluginapi.UsageDetail{InputTokens: 2, OutputTokens: 3}})
	if got := state.usage[stableKeyID("abc")].Inflight; got != 0 {
		t.Fatalf("inflight after usage record = %d, want 0", got)
	}
	if got := state.usage[stableKeyID("abc")].TokensToday; got != 5 {
		t.Fatalf("tokens today = %d, want 5", got)
	}
}

func TestRouteAddPatchAndDeleteRuleOnExistingPolicy(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = pluginConfig{KeyPolicies: []keyPolicyConfig{{KeyID: "policy-1", DisplayName: "Key 1", MatchAPIKey: "secret-1", Enabled: true}}}
	addResp, err := routeAddRule(pluginapi.ManagementRequest{Query: map[string][]string{"key_id": {"policy-1"}}, Body: []byte(`{"id":"rule-1","enabled":true,"priority":10,"on_match":"stop","match":{"models":["gpt-5.5"]},"actions":{"route_to_model":"openai/gpt-5.4"}}`)})
	if err != nil || addResp.StatusCode != http.StatusCreated {
		t.Fatalf("routeAddRule() = %#v, %v", addResp, err)
	}
	patchResp, err := routePatchRule(pluginapi.ManagementRequest{Query: map[string][]string{"key_id": {"policy-1"}, "rule_id": {"rule-1"}}, Body: []byte(`{"enabled":true,"priority":20,"on_match":"stop","match":{"paths":["/v1/responses"]},"actions":{"fallback_models":["openai/gpt-5.4-mini"]}}`)})
	if err != nil || patchResp.StatusCode != http.StatusOK {
		t.Fatalf("routePatchRule() = %#v, %v", patchResp, err)
	}
	delResp, err := routeDeleteRule(pluginapi.ManagementRequest{Query: map[string][]string{"key_id": {"policy-1"}, "rule_id": {"rule-1"}}})
	if err != nil || delResp.StatusCode != http.StatusOK {
		t.Fatalf("routeDeleteRule() = %#v, %v", delResp, err)
	}
}

func TestRouteAddRuleRejectsDuplicateRuleID(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = pluginConfig{KeyPolicies: []keyPolicyConfig{{KeyID: "policy-1", DisplayName: "Key 1", Enabled: true, Rules: []ruleConfig{{ID: "rule-1", Enabled: true}}}}}
	resp, err := routeAddRule(pluginapi.ManagementRequest{Query: map[string][]string{"key_id": {"policy-1"}}, Body: []byte(`{"id":"rule-1","enabled":true,"priority":10,"on_match":"stop"}`)})
	if err != nil || resp.StatusCode != http.StatusConflict {
		t.Fatalf("duplicate routeAddRule() = %#v, %v; want conflict", resp, err)
	}
}

func TestRouteAddAndDeletePolicy(t *testing.T) {
	gatewayState = newPluginState()
	addResp, err := routeAddPolicy(pluginapi.ManagementRequest{Body: []byte(`{"display_name":"Key 1","match_api_key":"secret-1","enabled":true}`)})
	if err != nil || addResp.StatusCode != http.StatusCreated {
		t.Fatalf("routeAddPolicy() = %#v, %v", addResp, err)
	}
	if len(gatewayState.config.KeyPolicies) != 1 {
		t.Fatalf("policies len = %d, want 1", len(gatewayState.config.KeyPolicies))
	}
	keyID := candidateKeyID(gatewayState.config.KeyPolicies[0])
	delResp, err := routeDeletePolicy(pluginapi.ManagementRequest{Query: map[string][]string{"key_id": {keyID}}})
	if err != nil || delResp.StatusCode != http.StatusOK {
		t.Fatalf("routeDeletePolicy() = %#v, %v", delResp, err)
	}
	if len(gatewayState.config.KeyPolicies) != 0 {
		t.Fatalf("policies len after delete = %d, want 0", len(gatewayState.config.KeyPolicies))
	}
}

func TestRouteAddPatchAndDeleteRuleViaAddedPolicy(t *testing.T) {
	gatewayState = newPluginState()
	addPolicyResp, err := routeAddPolicy(pluginapi.ManagementRequest{Body: []byte(`{"display_name":"Key 1","match_api_key":"secret-1","enabled":true}`)})
	if err != nil || addPolicyResp.StatusCode != http.StatusCreated {
		t.Fatalf("routeAddPolicy() = %#v, %v", addPolicyResp, err)
	}
	keyID := candidateKeyID(gatewayState.config.KeyPolicies[0])
	addRuleResp, err := routeAddRule(pluginapi.ManagementRequest{
		Query: map[string][]string{"key_id": {keyID}},
		Body:  []byte(`{"id":"rule-1","enabled":true,"priority":10,"on_match":"stop","match":{"models":["gpt-5.5"]},"actions":{"route_to_model":"openai/gpt-5.4"}}`),
	})
	if err != nil || addRuleResp.StatusCode != http.StatusCreated {
		t.Fatalf("routeAddRule() = %#v, %v", addRuleResp, err)
	}
	if got := len(gatewayState.config.KeyPolicies[0].Rules); got != 1 {
		t.Fatalf("rules len = %d, want 1", got)
	}
	patchRuleResp, err := routePatchRule(pluginapi.ManagementRequest{
		Query: map[string][]string{"key_id": {keyID}, "rule_id": {"rule-1"}},
		Body:  []byte(`{"enabled":true,"priority":20,"on_match":"stop","match":{"models":["gpt-5.5"]},"actions":{"fallback_models":["openai/gpt-5.4-mini"]}}`),
	})
	if err != nil || patchRuleResp.StatusCode != http.StatusOK {
		t.Fatalf("routePatchRule() = %#v, %v", patchRuleResp, err)
	}
	if got := gatewayState.config.KeyPolicies[0].Rules[0].Priority; got != 20 {
		t.Fatalf("patched rule priority = %d, want 20", got)
	}
	deleteRuleResp, err := routeDeleteRule(pluginapi.ManagementRequest{
		Query: map[string][]string{"key_id": {keyID}, "rule_id": {"rule-1"}},
	})
	if err != nil || deleteRuleResp.StatusCode != http.StatusOK {
		t.Fatalf("routeDeleteRule() = %#v, %v", deleteRuleResp, err)
	}
	if got := len(gatewayState.config.KeyPolicies[0].Rules); got != 0 {
		t.Fatalf("rules len after delete = %d, want 0", got)
	}
}

func TestListKeysMasksAPIKey(t *testing.T) {
	state := newPluginState()
	state.config = pluginConfig{KeyPolicies: []keyPolicyConfig{{DisplayName: "Key 1", MatchAPIKey: "abcdefghi", Enabled: true}}}
	keys := state.listKeys()
	if len(keys) != 1 {
		t.Fatalf("keys len = %d, want 1", len(keys))
	}
	if keys[0]["masked_key"] == "abcdefghi" {
		t.Fatalf("masked key leaked raw key: %#v", keys[0])
	}
}

func TestRoutePoliciesDoesNotReturnRawAPIKeys(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = pluginConfig{KeyPolicies: []keyPolicyConfig{{KeyID: "policy-1", DisplayName: "Key 1", MatchAPIKey: "secret-1", Enabled: true}}}
	resp, err := routePolicies(pluginapi.ManagementRequest{})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("routePolicies() = %#v, %v", resp, err)
	}
	body := string(resp.Body)
	if contains(body, "secret-1") {
		t.Fatalf("routePolicies leaked raw key: %s", body)
	}
	if !contains(body, "masked_key") {
		t.Fatalf("routePolicies missing masked key: %s", body)
	}
}

func TestRoutePutPoliciesPreservesExistingSecretForSanitizedPayload(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = pluginConfig{KeyPolicies: []keyPolicyConfig{{KeyID: "policy-1", DisplayName: "Key 1", MatchAPIKey: "secret-1", Enabled: true}}}
	resp, err := routePutPolicies(pluginapi.ManagementRequest{Body: []byte(`{"version":1,"key_policies":[{"key_id":"policy-1","display_name":"Updated","enabled":true,"rules":[]}]}`)})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("routePutPolicies() = %#v, %v", resp, err)
	}
	if got := gatewayState.config.KeyPolicies[0].MatchAPIKey; got != "secret-1" {
		t.Fatalf("preserved secret = %q, want secret-1", got)
	}
}

func TestRoutePutPoliciesRejectsDuplicatePolicyAndRuleIDs(t *testing.T) {
	gatewayState = newPluginState()
	resp, err := routePutPolicies(pluginapi.ManagementRequest{Body: []byte(`{"version":1,"key_policies":[{"key_id":"policy-1","display_name":"A","enabled":true},{"key_id":"policy-1","display_name":"B","enabled":true}]}`)})
	if err != nil || resp.StatusCode != http.StatusConflict {
		t.Fatalf("duplicate policy IDs = %#v, %v; want conflict", resp, err)
	}
	resp, err = routePutPolicies(pluginapi.ManagementRequest{Body: []byte(`{"version":1,"key_policies":[{"key_id":"policy-2","display_name":"A","enabled":true,"rules":[{"id":"rule-1","enabled":true},{"id":"rule-1","enabled":true}]}]}`)})
	if err != nil || resp.StatusCode != http.StatusConflict {
		t.Fatalf("duplicate rule IDs = %#v, %v; want conflict", resp, err)
	}
}

func TestPersistentStateSurvivesReload(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "gateway-state.json")
	gatewayState = newPluginState()
	gatewayState.config = pluginConfig{Persistence: persistenceConfig{StatePath: statePath}}
	resp, err := routeAddPolicy(pluginapi.ManagementRequest{Body: []byte(`{"key_id":"policy-persist","display_name":"Persisted","match_api_key":"secret-persist","enabled":true}`)})
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("routeAddPolicy() = %#v, %v", resp, err)
	}
	resp, err = routeAddTemplate(pluginapi.ManagementRequest{Body: []byte(`{"id":"tpl-persist","name":"Persisted Template","description":"desc","rule":{"id":"rule-1","enabled":true,"priority":10,"on_match":"stop","match":{"models":["gpt-5.5"]},"actions":{"route_to_model":"openai/gpt-5.4"}}}`)})
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("routeAddTemplate() = %#v, %v", resp, err)
	}

	gatewayState = newPluginState()
	gatewayState.config = pluginConfig{Persistence: persistenceConfig{StatePath: statePath}}
	if err := gatewayState.loadPersistentState(); err != nil {
		t.Fatalf("loadPersistentState(): %v", err)
	}
	if len(gatewayState.config.KeyPolicies) != 1 || gatewayState.config.KeyPolicies[0].MatchAPIKey != "secret-persist" {
		t.Fatalf("loaded policies = %#v", gatewayState.config.KeyPolicies)
	}
	foundTemplate := false
	for _, item := range gatewayState.templates {
		if item.ID == "tpl-persist" {
			foundTemplate = true
		}
	}
	if !foundTemplate {
		t.Fatalf("loaded templates = %#v", gatewayState.templates)
	}
}

func TestPersistentRuntimeUsageSurvivesReload(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "gateway-runtime.json")
	gatewayState = newPluginState()
	gatewayState.config = pluginConfig{Persistence: persistenceConfig{StatePath: statePath, PersistRuntime: true}, Default: policyConfig{Enabled: true}}
	req := pluginapi.RequestInterceptRequest{Headers: http.Header{}, Body: []byte(`{"model":"gpt-5.5"}`), Metadata: map[string]any{"access.api_key": "abc"}}
	resp := gatewayState.apply(req, false)
	if resp.Reject {
		t.Fatalf("apply rejected unexpectedly: %#v", resp)
	}
	gatewayState.recordUsage(pluginapi.UsageRecord{APIKey: "abc", Detail: pluginapi.UsageDetail{TotalTokens: 7}})

	gatewayState = newPluginState()
	gatewayState.config = pluginConfig{Persistence: persistenceConfig{StatePath: statePath, PersistRuntime: true}}
	if err := gatewayState.loadPersistentState(); err != nil {
		t.Fatalf("loadPersistentState(): %v", err)
	}
	usage := gatewayState.usage[stableKeyID("abc")]
	if usage == nil || usage.TokensToday != 7 || usage.RequestsToday != 1 {
		t.Fatalf("loaded runtime usage = %#v", usage)
	}
}

func TestManagementAuthorizationRequiresConfiguredTokens(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = pluginConfig{Security: securityConfig{RequireManagementToken: true, AdminTokens: []string{"admin-token"}, ReadTokens: []string{"read-token"}}}
	readHandler := authorizedHandler(roleRead, routeKeys)
	writeHandler := authorizedHandler(roleAdmin, routeAddPolicy)
	resp, err := readHandler(pluginapi.ManagementRequest{})
	if err != nil || resp.StatusCode != http.StatusForbidden {
		t.Fatalf("read without token = %#v, %v", resp, err)
	}
	resp, err = readHandler(pluginapi.ManagementRequest{Headers: http.Header{"Authorization": {"Bearer read-token"}}})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("read with read token = %#v, %v", resp, err)
	}
	resp, err = writeHandler(pluginapi.ManagementRequest{Headers: http.Header{"Authorization": {"Bearer read-token"}}, Body: []byte(`{"display_name":"Denied","match_api_key":"secret"}`)})
	if err != nil || resp.StatusCode != http.StatusForbidden {
		t.Fatalf("write with read token = %#v, %v", resp, err)
	}
	resp, err = writeHandler(pluginapi.ManagementRequest{Headers: http.Header{"Authorization": {"Bearer admin-token"}}, Body: []byte(`{"display_name":"Allowed","match_api_key":"secret"}`)})
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("write with admin token = %#v, %v", resp, err)
	}
}

func TestRouteUIRequiresConfiguredUIToken(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = pluginConfig{Security: securityConfig{UIAccessTokens: []string{"ui-token"}}}
	resp, err := routeUI(pluginapi.ManagementRequest{})
	if err != nil || resp.StatusCode != http.StatusForbidden {
		t.Fatalf("routeUI without token = %#v, %v", resp, err)
	}
	resp, err = routeUI(pluginapi.ManagementRequest{Query: map[string][]string{"gateway_token": {"ui-token"}}})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("routeUI with token = %#v, %v", resp, err)
	}

	gatewayState.config = pluginConfig{Security: securityConfig{RequireManagementToken: true, ReadTokens: []string{"read-token"}}}
	resp, err = routeUI(pluginapi.ManagementRequest{})
	if err != nil || resp.StatusCode != http.StatusForbidden {
		t.Fatalf("routeUI with management token required but no token = %#v, %v", resp, err)
	}
	resp, err = routeUI(pluginapi.ManagementRequest{Query: map[string][]string{"gateway_token": {"read-token"}}})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("routeUI with read token = %#v, %v", resp, err)
	}
}

func TestRouteHealthReportsLocalCountersAndNoSecrets(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = normalizeConfig(pluginConfig{
		Persistence: persistenceConfig{StatePath: filepath.Join(t.TempDir(), "gateway-state.json"), PersistRuntime: true},
		Security: securityConfig{
			RequireManagementToken: true,
			AdminTokens:            []string{"admin-secret"},
			ReadTokens:             []string{"read-secret"},
			UIAccessTokens:         []string{"ui-secret"},
		},
		Cluster: clusterConfig{Redis: redisConfig{Password: "redis-secret"}},
		KeyPolicies: []keyPolicyConfig{{
			KeyID:       "policy-1",
			DisplayName: "Policy 1",
			MatchAPIKey: "api-key-secret",
			Enabled:     true,
			Rules:       []ruleConfig{{ID: "rule-1", Enabled: true, Priority: 10, OnMatch: "stop"}},
		}},
	})
	resp, err := routeHealth(pluginapi.ManagementRequest{})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("routeHealth() = %#v, %v", resp, err)
	}
	var health gatewayHealth
	if err := json.Unmarshal(resp.Body, &health); err != nil {
		t.Fatalf("unmarshal health: %v", err)
	}
	if health.Counters.Backend != "local" || health.Counters.RedisRequired {
		t.Fatalf("counter health = %#v, want local and redis not required", health.Counters)
	}
	if health.Security.AdminTokenCount != 1 || health.Security.ReadTokenCount != 1 || health.Security.UITokenCount != 1 {
		t.Fatalf("security health = %#v, want token counts only", health.Security)
	}
	if health.Counts.KeyPolicies != 1 || health.Counts.Rules != 1 {
		t.Fatalf("counts = %#v, want one policy and one rule", health.Counts)
	}
	body := string(resp.Body)
	for _, secret := range []string{"api-key-secret", "admin-secret", "read-secret", "ui-secret", "redis-secret"} {
		if contains(body, secret) {
			t.Fatalf("routeHealth leaked secret %q in %s", secret, body)
		}
	}
}

func TestHealthSnapshotReportsRedisRequiredOnlyWhenConfigured(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = normalizeConfig(pluginConfig{Cluster: clusterConfig{Backend: "redis"}})
	health := gatewayState.healthSnapshot()
	if health.Counters.Backend != "redis" || !health.Counters.RedisRequired || health.Status != "degraded" {
		t.Fatalf("redis health = %#v, want redis required and degraded when store is missing", health)
	}
}

func TestHealthSeparatesUIAuthFromManagementAPIAuth(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = pluginConfig{Security: securityConfig{UIAccessTokens: []string{"ui-token"}}}
	health := gatewayState.healthSnapshot()
	if health.Security.ManagementAuthEnabled {
		t.Fatalf("management auth should be false when only UI tokens are configured: %#v", health.Security)
	}
	if !health.Security.UIAuthEnabled {
		t.Fatalf("UI auth should be true when UI tokens are configured: %#v", health.Security)
	}
}

func TestRedisCounterStoreIsConfiguredOnlyForRedisBackend(t *testing.T) {
	if store := newRedisCounterStore(clusterConfig{}); store != nil {
		t.Fatalf("default counter store = %#v, want nil", store)
	}
	store := newRedisCounterStore(clusterConfig{Backend: "redis", Redis: redisConfig{Addr: "127.0.0.1:6379", KeyPrefix: "test-gateway"}})
	if store == nil {
		t.Fatal("redis counter store is nil")
	}
	if got := store.keysKey(); got != "test-gateway:usage-keys" {
		t.Fatalf("redis keys key = %q", got)
	}
}

func TestNormalizeRedisFailureMode(t *testing.T) {
	cases := map[string]string{
		"":               "reject",
		"fail-open":      "allow",
		"fail_closed":    "reject",
		"local":          "local_fallback",
		"local-fallback": "local_fallback",
		"unexpected":     "reject",
	}
	for raw, want := range cases {
		got := normalizeConfig(pluginConfig{Cluster: clusterConfig{Backend: " Redis ", Redis: redisConfig{FailureMode: raw}}})
		if got.Cluster.Backend != "redis" || got.Cluster.Redis.FailureMode != want {
			t.Fatalf("normalize failure mode %q = backend %q mode %q, want redis/%s", raw, got.Cluster.Backend, got.Cluster.Redis.FailureMode, want)
		}
	}
}

func TestCounterBackendUnavailableFailureModes(t *testing.T) {
	reject := redisUnavailableResponse(errors.New("redis unavailable"))

	gatewayState = newPluginState()
	gatewayState.config = normalizeConfig(pluginConfig{Cluster: clusterConfig{Backend: "redis", Redis: redisConfig{FailureMode: "allow"}}})
	if got := gatewayState.handleCounterBackendUnavailable(reject, "key-1", "Key 1", "***", limitConfig{MaxInflight: 1}, time.Now(), true); got != nil {
		t.Fatalf("allow failure mode rejected: %#v", got)
	}

	gatewayState = newPluginState()
	gatewayState.config = normalizeConfig(pluginConfig{Cluster: clusterConfig{Backend: "redis", Redis: redisConfig{FailureMode: "local_fallback"}}})
	if got := gatewayState.handleCounterBackendUnavailable(reject, "key-1", "Key 1", "***", limitConfig{MaxInflight: 1}, time.Now(), true); got != nil {
		t.Fatalf("local fallback first request rejected: %#v", got)
	}
	if got := gatewayState.handleCounterBackendUnavailable(reject, "key-1", "Key 1", "***", limitConfig{MaxInflight: 1}, time.Now(), true); got == nil || got.RejectCode != "gateway_concurrency_exceeded" {
		t.Fatalf("local fallback did not enforce local inflight: %#v", got)
	}

	gatewayState = newPluginState()
	gatewayState.config = normalizeConfig(pluginConfig{Cluster: clusterConfig{Backend: "redis", Redis: redisConfig{FailureMode: "reject"}}})
	if got := gatewayState.handleCounterBackendUnavailable(reject, "key-1", "Key 1", "***", limitConfig{}, time.Now(), true); got == nil || got.RejectCode != "gateway_counter_backend_unavailable" {
		t.Fatalf("reject failure mode = %#v", got)
	}
}

func TestRedisUnavailableResponseDoesNotLeakErrorDetails(t *testing.T) {
	resp := redisUnavailableResponse(errors.New("dial tcp 10.0.0.12:6379: auth failed for redis-secret"))
	if got := resp.Headers.Get("X-Gateway-Counter-Error"); got != "unavailable" {
		t.Fatalf("counter error header = %q, want generic unavailable", got)
	}
	if contains(resp.RejectMessage, "10.0.0.12") || contains(resp.RejectMessage, "redis-secret") {
		t.Fatalf("redis unavailable response leaked details: %#v", resp)
	}
}

func TestRedisDryRunDoesNotRegisterUsageKey(t *testing.T) {
	saddIndex := strings.Index(redisEnforceLua, "redis.call('SADD'")
	enforceIndex := strings.Index(redisEnforceLua, "if ARGV[12] == '1' then")
	if saddIndex < 0 || enforceIndex < 0 || saddIndex < enforceIndex {
		t.Fatalf("redis enforce script should add usage keys only inside enforce block")
	}
}

func TestRouteImportPoliciesPreservesExistingSecretForSanitizedBundle(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = pluginConfig{KeyPolicies: []keyPolicyConfig{{KeyID: "policy-1", DisplayName: "Key 1", MatchAPIKey: "secret-1", Enabled: true}}}
	exportResp, err := routeExportPolicies(pluginapi.ManagementRequest{})
	if err != nil || exportResp.StatusCode != http.StatusOK {
		t.Fatalf("routeExportPolicies() = %#v, %v", exportResp, err)
	}
	if contains(string(exportResp.Body), "secret-1") {
		t.Fatalf("routeExportPolicies leaked raw key: %s", exportResp.Body)
	}
	importResp, err := routeImportPolicies(pluginapi.ManagementRequest{Query: map[string][]string{"mode": {"merge"}}, Body: exportResp.Body})
	if err != nil || importResp.StatusCode != http.StatusOK {
		t.Fatalf("routeImportPolicies() = %#v, %v", importResp, err)
	}
	if got := gatewayState.config.KeyPolicies[0].MatchAPIKey; got != "secret-1" {
		t.Fatalf("import preserved secret = %q, want secret-1", got)
	}
}

func TestLookupPolicySupportsKeyIDMetadata(t *testing.T) {
	state := newPluginState()
	state.config = pluginConfig{KeyPolicies: []keyPolicyConfig{{KeyID: "policy-1", DisplayName: "Key 1", MatchAPIKey: "secret-1", Enabled: true, Rules: []ruleConfig{{ID: "route-1", Enabled: true, OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{RouteToModel: "openai/gpt-5.4"}}}}}}
	result := state.evaluate(pluginapi.RequestInterceptRequest{Model: "gpt-5.5", Headers: http.Header{}, Body: []byte(`{"model":"gpt-5.5"}`), Metadata: map[string]any{"access.key_id": "policy-1"}}, false)
	if result.FinalModel != "openai/gpt-5.4" {
		t.Fatalf("final model = %q, want openai/gpt-5.4", result.FinalModel)
	}
}

func TestStringMetadataMissingKeyIsEmpty(t *testing.T) {
	if got := stringMetadata(map[string]any{}, "missing"); got != "" {
		t.Fatalf("stringMetadata missing = %q, want empty", got)
	}
}

func TestApplyRulesSelectsDeterministicWeightedRoute(t *testing.T) {
	policy := policyConfig{Enabled: true, Rules: []ruleConfig{{ID: "weighted-1", Enabled: true, OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{ShardBy: "api_key", WeightedRoutes: []weightedRoute{{Model: "openai/gpt-5.4", Weight: 1}, {Model: "codex/gpt-5.4", Weight: 3}}}}}}
	req := pluginapi.RequestInterceptRequest{Model: "gpt-5.5", Headers: http.Header{}, Body: []byte(`{"model":"gpt-5.5"}`), Metadata: map[string]any{"access.api_key": "tenant-a"}}
	first := applyRules(req, policy, false, time.Now())
	second := applyRules(req, policy, false, time.Now())
	if first.Reason != "weighted_route" || first.FinalModel == "" {
		t.Fatalf("first result = %#v", first)
	}
	if first.FinalModel != second.FinalModel {
		t.Fatalf("weighted route should be deterministic, got %q then %q", first.FinalModel, second.FinalModel)
	}
}

func TestApplyRulesSetsMirrorHeaderAndTags(t *testing.T) {
	policy := policyConfig{Enabled: true, Rules: []ruleConfig{{ID: "mirror-1", Enabled: true, OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{MirrorModels: []string{"openai/gpt-5.4-mini", "openai/gpt-4.1-mini"}, TagMetadata: map[string]string{"tenant.mode": "shadow"}}}}}
	result := applyRules(pluginapi.RequestInterceptRequest{Model: "gpt-5.5", Headers: http.Header{}, Body: []byte(`{"model":"gpt-5.5"}`)}, policy, false, time.Now())
	if got := result.Response.Headers.Get("X-Gateway-Mirror-Models"); got == "" || !contains(got, "openai/gpt-5.4-mini") {
		t.Fatalf("mirror header = %q", got)
	}
	if got := result.Response.Headers.Get("X-Gateway-Tag-tenant-mode"); got != "shadow" {
		t.Fatalf("tag header = %q, want shadow", got)
	}
}

func TestApplyRulesAllowOnlyProvidersChecksFinalModel(t *testing.T) {
	policy := policyConfig{Enabled: true, Rules: []ruleConfig{{ID: "force-allow", Enabled: true, OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{ForceProviderPrefix: "codex", AllowOnlyProviders: []string{"codex"}}}}}
	result := applyRules(pluginapi.RequestInterceptRequest{Model: "gpt-5.5", Headers: http.Header{}, Body: []byte(`{"model":"gpt-5.5"}`)}, policy, false, time.Now())
	if result.Response.Reject || result.FinalModel != "codex/gpt-5.5" {
		t.Fatalf("result = %#v, want forced codex model without reject", result)
	}
}

func TestAuditSummaryAggregatesByFinalModel(t *testing.T) {
	state := newPluginState()
	state.auditLog = []auditEntry{{Decision: "rewrite", FinalModel: "openai/gpt-5.4"}, {Decision: "rewrite", FinalModel: "openai/gpt-5.4"}, {Decision: "pass", FinalModel: "codex/gpt-5.4"}}
	summary := state.auditSummary(nil)
	if summary.TotalByModel["openai/gpt-5.4"] != 2 || summary.TotalByModel["codex/gpt-5.4"] != 1 {
		t.Fatalf("summary by model = %#v", summary.TotalByModel)
	}
}

func TestPluginStateAuditCapturesMirrorModels(t *testing.T) {
	state := newPluginState()
	state.config = pluginConfig{Default: policyConfig{Enabled: true, Rules: []ruleConfig{{ID: "mirror-audit", Enabled: true, OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{MirrorModels: []string{"openai/gpt-5.4-mini"}}}}}}
	state.apply(pluginapi.RequestInterceptRequest{Model: "gpt-5.5", Headers: http.Header{}, Body: []byte(`{"model":"gpt-5.5"}`), Metadata: map[string]any{"access.api_key": "secret-key"}}, false)
	items := state.listAudit(10, nil)
	if len(items) == 0 || len(items[0].Mirrors) != 1 || items[0].Mirrors[0] != "openai/gpt-5.4-mini" {
		t.Fatalf("audit mirrors = %#v", items)
	}
}

func TestAuditEntryMatchesSupportsPolicyModelProviderAndTime(t *testing.T) {
	entry := auditEntry{Time: time.Date(2026, 7, 4, 9, 30, 0, 0, time.UTC), PolicyName: "Tenant A", RequestedModel: "gpt-5.5", FinalModel: "openai/gpt-5.4", Provider: "openai"}
	filters := map[string]string{"policy": "tenant", "model": "gpt-5.4", "provider": "openai", "from": "2026-07-04T09:00:00Z", "to": "2026-07-04T10:00:00Z"}
	if !auditEntryMatches(entry, filters) {
		t.Fatalf("auditEntryMatches should match filters: %#v", filters)
	}
}

func TestAuditSummaryAggregatesProviderAndTimeline(t *testing.T) {
	state := newPluginState()
	state.auditLog = []auditEntry{
		{Time: time.Date(2026, 7, 4, 9, 30, 0, 0, time.UTC), Decision: "rewrite", FinalModel: "openai/gpt-5.4", Provider: "openai"},
		{Time: time.Date(2026, 7, 4, 9, 30, 30, 0, time.UTC), Decision: "pass", FinalModel: "codex/gpt-5.4", Provider: "codex"},
	}
	summary := state.auditSummary(nil)
	if summary.TotalByProvider["openai"] != 1 || summary.TotalByProvider["codex"] != 1 {
		t.Fatalf("summary by provider = %#v", summary.TotalByProvider)
	}
	if len(summary.Timeline) == 0 {
		t.Fatalf("expected timeline buckets, got %#v", summary.Timeline)
	}
}

func TestRouteTemplatesSupportsScenarioMaturityAndTagFilters(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.templates = []ruleTemplate{
		{ID: "tpl-route", Name: "Route", Category: "routing", Scenario: "traffic-split", Maturity: "beta", Tags: []string{"weighted", "routing"}, Description: "desc", Rule: ruleConfig{ID: "r1", Enabled: true}},
		{ID: "tpl-fallback", Name: "Fallback", Category: "fallback", Scenario: "cost-control", Maturity: "stable", Tags: []string{"fallback"}, Description: "desc", Rule: ruleConfig{ID: "r2", Enabled: true}},
	}
	resp, err := routeTemplates(pluginapi.ManagementRequest{Query: map[string][]string{"scenario": {"traffic-split"}, "maturity": {"beta"}, "tag": {"weighted"}}})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("routeTemplates() = %#v, %v", resp, err)
	}
	if !contains(string(resp.Body), "tpl-route") || contains(string(resp.Body), "tpl-fallback") {
		t.Fatalf("filtered templates body = %s", string(resp.Body))
	}
}

func TestRouteExportPoliciesReturnsBundle(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = pluginConfig{Default: policyConfig{Enabled: true}, KeyPolicies: []keyPolicyConfig{{KeyID: "policy-1", DisplayName: "Key 1", MatchAPIKey: "secret-1", Enabled: true}}}
	resp, err := routeExportPolicies(pluginapi.ManagementRequest{Query: map[string][]string{"name": {"bundle-a"}}})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("routeExportPolicies() = %#v, %v", resp, err)
	}
	if !contains(string(resp.Body), "bundle-a") || !contains(string(resp.Body), "policy-1") {
		t.Fatalf("bundle body = %s", string(resp.Body))
	}
}

func TestRouteImportPoliciesMergeAndReplace(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = pluginConfig{KeyPolicies: []keyPolicyConfig{{KeyID: "policy-1", DisplayName: "Key 1", MatchAPIKey: "secret-1", Enabled: true}}}
	mergeResp, err := routeImportPolicies(pluginapi.ManagementRequest{Query: map[string][]string{"mode": {"merge"}}, Body: []byte(`{"version":1,"name":"bundle","key_policies":[{"key_id":"policy-2","display_name":"Key 2","match_api_key":"secret-2","enabled":true}]}`)})
	if err != nil || mergeResp.StatusCode != http.StatusOK {
		t.Fatalf("routeImportPolicies(merge) = %#v, %v", mergeResp, err)
	}
	if len(gatewayState.config.KeyPolicies) != 2 {
		t.Fatalf("merge policies len = %d, want 2", len(gatewayState.config.KeyPolicies))
	}
	replaceResp, err := routeImportPolicies(pluginapi.ManagementRequest{Query: map[string][]string{"mode": {"replace"}}, Body: []byte(`{"version":1,"name":"bundle","key_policies":[{"key_id":"policy-3","display_name":"Key 3","match_api_key":"secret-3","enabled":true}]}`)})
	if err != nil || replaceResp.StatusCode != http.StatusOK {
		t.Fatalf("routeImportPolicies(replace) = %#v, %v", replaceResp, err)
	}
	if len(gatewayState.config.KeyPolicies) != 1 || gatewayState.config.KeyPolicies[0].KeyID != "policy-3" {
		t.Fatalf("replace policies = %#v", gatewayState.config.KeyPolicies)
	}
}

func TestRouteClonePolicyClearsBoundAPIKey(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = pluginConfig{KeyPolicies: []keyPolicyConfig{{KeyID: "policy-1", DisplayName: "Key 1", MatchAPIKey: "secret-1", Enabled: true, Rules: []ruleConfig{{ID: "r1", Enabled: true}}}}}
	resp, err := routeClonePolicy(pluginapi.ManagementRequest{Query: map[string][]string{"key_id": {"policy-1"}}})
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("routeClonePolicy() = %#v, %v", resp, err)
	}
	if len(gatewayState.config.KeyPolicies) != 2 {
		t.Fatalf("policies len = %d, want 2", len(gatewayState.config.KeyPolicies))
	}
	clone := gatewayState.config.KeyPolicies[1]
	if clone.MatchAPIKey != "" || clone.DisplayName == "Key 1" {
		t.Fatalf("cloned policy = %#v", clone)
	}
}

func TestApplyRulesWithStagesHonorsPreCheckBeforeRoute(t *testing.T) {
	policy := policyConfig{Enabled: true, Rules: []ruleConfig{
		{ID: "deny-pre", Enabled: true, Priority: 5, Stage: "pre-check", OnMatch: "stop", Match: matchConfig{Providers: []string{"claude"}}, Actions: actionConfig{Deny: &denyConfig{StatusCode: 403, Message: "blocked", Code: "gateway_provider_denied"}}},
		{ID: "route-late", Enabled: true, Priority: 10, Stage: "route", OnMatch: "stop", Match: matchConfig{Models: []string{"claude/gpt-5.5"}}, Actions: actionConfig{RouteToModel: "openai/gpt-5.4"}},
	}}
	result := applyRulesWithStages(pluginapi.RequestInterceptRequest{Model: "claude/gpt-5.5", Headers: http.Header{}, Body: []byte(`{"model":"claude/gpt-5.5"}`)}, policy, false, time.Now())
	if result.Decision != "reject" || result.RuleID != "deny-pre" {
		t.Fatalf("staged result = %#v", result)
	}
}

func TestApplyRulesWithStagesAllowsMirrorAfterRouteRewrite(t *testing.T) {
	policy := policyConfig{Enabled: true, Rules: []ruleConfig{
		{ID: "route-first", Enabled: true, Priority: 10, Stage: "route", OnMatch: "continue", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{RouteToModel: "openai/gpt-5.4"}},
		{ID: "mirror-second", Enabled: true, Priority: 20, Stage: "mirror", OnMatch: "stop", Match: matchConfig{Models: []string{"openai/gpt-5.4"}}, Actions: actionConfig{MirrorModels: []string{"openai/gpt-5.4-mini"}}},
	}}
	result := applyRulesWithStages(pluginapi.RequestInterceptRequest{Model: "gpt-5.5", Headers: http.Header{}, Body: []byte(`{"model":"gpt-5.5"}`)}, policy, false, time.Now())
	if result.FinalModel != "openai/gpt-5.4" {
		t.Fatalf("final model = %q, want openai/gpt-5.4", result.FinalModel)
	}
	if got := result.Response.Headers.Get("X-Gateway-Mirror-Models"); got != "openai/gpt-5.4-mini" {
		t.Fatalf("mirror header = %q", got)
	}
	if len(result.StageTrace) < 2 || result.StageTrace[0].Stage != "pre-check" {
		t.Fatalf("stage trace = %#v", result.StageTrace)
	}
}

func TestNormalizeRuleStageDefaultsToPreCheck(t *testing.T) {
	if got := normalizeRuleStage(""); got != "pre-check" {
		t.Fatalf("normalizeRuleStage('') = %q", got)
	}
	if got := normalizeRuleStage("post_audit"); got != "post-audit" {
		t.Fatalf("normalizeRuleStage(post_audit) = %q", got)
	}
}

func TestStageModeForDefaultsByStage(t *testing.T) {
	policy := policyConfig{}
	if got := stageModeFor(policy, "route"); got != "first-match" {
		t.Fatalf("route default mode = %q", got)
	}
	if got := stageModeFor(policy, "mirror"); got != "continue-all" {
		t.Fatalf("mirror default mode = %q", got)
	}
}

func TestApplyRulesWithStagesContinueAllAllowsMultipleRewrites(t *testing.T) {
	policy := policyConfig{Enabled: true, StagePolicy: map[string]stagePolicy{"rewrite": {Mode: "continue-all"}}, Rules: []ruleConfig{
		{ID: "rewrite-1", Enabled: true, Priority: 10, Stage: "rewrite", OnMatch: "continue", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{RewriteModel: "openai/gpt-5.4"}},
		{ID: "rewrite-2", Enabled: true, Priority: 20, Stage: "rewrite", OnMatch: "continue", Match: matchConfig{Models: []string{"openai/gpt-5.4"}}, Actions: actionConfig{ForceProviderPrefix: "codex"}},
	}}
	result := applyRulesWithStages(pluginapi.RequestInterceptRequest{Model: "gpt-5.5", Headers: http.Header{}, Body: []byte(`{"model":"gpt-5.5"}`)}, policy, false, time.Now())
	if result.FinalModel != "codex/gpt-5.4" {
		t.Fatalf("final model = %q, want codex/gpt-5.4", result.FinalModel)
	}
	if result.Decision != "rewrite" || result.RuleID != "rewrite-2" {
		t.Fatalf("decision = %s/%s, want rewrite/rewrite-2", result.Decision, result.RuleID)
	}
}

func TestApplyRulesWithStagesFirstMatchStopsWithinStage(t *testing.T) {
	policy := policyConfig{Enabled: true, Rules: []ruleConfig{
		{ID: "rewrite-1", Enabled: true, Priority: 10, Stage: "rewrite", OnMatch: "continue", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{RewriteModel: "openai/gpt-5.4"}},
		{ID: "rewrite-2", Enabled: true, Priority: 20, Stage: "rewrite", OnMatch: "continue", Match: matchConfig{Models: []string{"openai/gpt-5.4"}}, Actions: actionConfig{ForceProviderPrefix: "codex"}},
	}}
	result := applyRulesWithStages(pluginapi.RequestInterceptRequest{Model: "gpt-5.5", Headers: http.Header{}, Body: []byte(`{"model":"gpt-5.5"}`)}, policy, false, time.Now())
	if result.FinalModel != "openai/gpt-5.4" {
		t.Fatalf("final model = %q, want openai/gpt-5.4", result.FinalModel)
	}
	if result.Decision != "rewrite" || result.RuleID != "rewrite-1" {
		t.Fatalf("decision = %s/%s, want rewrite/rewrite-1", result.Decision, result.RuleID)
	}
	rewriteTraceCount := 0
	for _, item := range result.StageTrace {
		if item.Stage == "rewrite" {
			rewriteTraceCount++
			if item.Decision != "rewrite" {
				t.Fatalf("rewrite trace decision = %q, want rewrite: %#v", item.Decision, result.StageTrace)
			}
		}
	}
	if rewriteTraceCount != 1 {
		t.Fatalf("rewrite trace count = %d, want 1: %#v", rewriteTraceCount, result.StageTrace)
	}
}

func TestApplyRulesRoutePoolUsesMembers(t *testing.T) {
	policy := policyConfig{Enabled: true, Rules: []ruleConfig{{ID: "pool-1", Enabled: true, Priority: 10, Stage: "route", OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{ShardBy: "api_key", RoutePool: &routePoolConfig{Name: "primary", Mode: "weighted", Members: []weightedRoute{{Model: "openai/gpt-5.4", Weight: 2}, {Model: "codex/gpt-5.4", Weight: 1}}}}}}}
	result := applyRulesWithStages(pluginapi.RequestInterceptRequest{Model: "gpt-5.5", Headers: http.Header{}, Body: []byte(`{"model":"gpt-5.5"}`), Metadata: map[string]any{"access.api_key": "tenant-a"}}, policy, false, time.Now())
	if result.Reason != "route_pool" || result.FinalModel == "" {
		t.Fatalf("route pool result = %#v", result)
	}
}

func TestApplyRulesFailoverChainOverridesFallbackReason(t *testing.T) {
	policy := policyConfig{Enabled: true, Rules: []ruleConfig{{ID: "failover-1", Enabled: true, Priority: 10, Stage: "route", OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{FailoverChain: []string{"openai/gpt-5.4", "openai/gpt-4.1-mini"}}}}}
	result := applyRulesWithStages(pluginapi.RequestInterceptRequest{Model: "gpt-5.5", Headers: http.Header{}, Body: []byte(`{"model":"gpt-5.5"}`)}, policy, false, time.Now())
	if result.Reason != "failover_chain" || result.FinalModel != "openai/gpt-5.4" {
		t.Fatalf("failover result = %#v", result)
	}
}

func TestMergedFailoverChainDeduplicatesFallbacks(t *testing.T) {
	chain := mergedFailoverChain(actionConfig{FailoverChain: []string{"openai/gpt-5.4"}, FallbackModels: []string{"openai/gpt-5.4", "openai/gpt-4.1-mini"}})
	if len(chain) != 2 || chain[0] != "openai/gpt-5.4" || chain[1] != "openai/gpt-4.1-mini" {
		t.Fatalf("merged chain = %#v", chain)
	}
}

func TestNormalizedWeightedRoutesBuildsModelFromProviderAndSuffix(t *testing.T) {
	routes := normalizedWeightedRoutes([]weightedRoute{{Provider: "openai", Suffix: "gpt-5.4", Weight: 2}})
	if len(routes) != 1 || routes[0].Model != "openai/gpt-5.4" {
		t.Fatalf("normalized routes = %#v", routes)
	}
}

func TestSelectRouteTargetHonorsProviderAffinity(t *testing.T) {
	actions := actionConfig{ShardBy: "api_key", RoutePool: &routePoolConfig{Name: "pool", ProviderAffinity: "codex", Members: []weightedRoute{{Model: "openai/gpt-5.4", Weight: 1}, {Provider: "codex", Suffix: "gpt-5.4", Weight: 1}}}}
	model := selectRouteTarget(pluginapi.RequestInterceptRequest{Model: "gpt-5.5", Metadata: map[string]any{"access.api_key": "tenant-a"}}, actions)
	if providerFromModel(model) != "codex" {
		t.Fatalf("affinity model = %q, want codex provider", model)
	}
}

func TestNormalizedWeightedRoutesSkipsDisabledMembers(t *testing.T) {
	enabled := true
	disabled := false
	routes := normalizedWeightedRoutes([]weightedRoute{{Model: "openai/gpt-5.4", Weight: 1, Enabled: &disabled}, {Model: "codex/gpt-5.4", Weight: 1, Enabled: &enabled}})
	if len(routes) != 1 || routes[0].Model != "codex/gpt-5.4" {
		t.Fatalf("normalized routes = %#v", routes)
	}
}

func TestNormalizedWeightedRoutesSortsByPriority(t *testing.T) {
	routes := normalizedWeightedRoutes([]weightedRoute{{Model: "openai/gpt-5.4", Weight: 1, Priority: 200}, {Model: "codex/gpt-5.4", Weight: 1, Priority: 10}})
	if len(routes) != 2 || routes[0].Model != "codex/gpt-5.4" {
		t.Fatalf("priority-sorted routes = %#v", routes)
	}
}

func TestNormalizedWeightedRoutesSkipsBlockedStatuses(t *testing.T) {
	routes := normalizedWeightedRoutes([]weightedRoute{{Model: "openai/gpt-5.4", Weight: 1, Status: "drain"}, {Model: "codex/gpt-5.4", Weight: 1, Status: "active"}})
	if len(routes) != 1 || routes[0].Model != "codex/gpt-5.4" {
		t.Fatalf("status-filtered routes = %#v", routes)
	}
}

func TestMergedFailoverChainPrefersFailoverHops(t *testing.T) {
	enabled := true
	chain := mergedFailoverChain(actionConfig{FailoverHops: []failoverHop{{Provider: "openai", Suffix: "gpt-5.4", Reason: "quota", OnDecision: "reject", Enabled: &enabled}}, FailoverChain: []string{"codex/gpt-5.4"}, FallbackModels: []string{"openai/gpt-4.1-mini"}})
	if len(chain) != 3 || chain[0] != "openai/gpt-5.4" {
		t.Fatalf("merged hop chain = %#v", chain)
	}
}

func TestStatusBlocksRouteRecognizesOperationalStates(t *testing.T) {
	if !statusBlocksRoute("drain") || !statusBlocksRoute("offline") || !statusBlocksRoute("degraded") {
		t.Fatal("expected operational blocked statuses to be rejected")
	}
	if statusBlocksRoute("active") {
		t.Fatal("active status should not be blocked")
	}
}

func TestEffectiveRouteWeightAppliesHealthAndTrafficCap(t *testing.T) {
	got := effectiveRouteWeight(weightedRoute{Weight: 100, Health: 50, TrafficCap: 20})
	if got != 10 {
		t.Fatalf("effectiveRouteWeight = %d, want 10", got)
	}
}

func TestSelectWeightedModelSkipsZeroEffectiveRoutes(t *testing.T) {
	model := selectWeightedModel(pluginapi.RequestInterceptRequest{Model: "gpt-5.5", Metadata: map[string]any{"access.api_key": "tenant-a"}}, "api_key", []weightedRoute{{Model: "openai/gpt-5.4", Weight: 100, Health: 1, TrafficCap: 1}, {Model: "codex/gpt-5.4", Weight: 100, Health: 100, TrafficCap: 100}})
	if model != "codex/gpt-5.4" {
		t.Fatalf("selected model = %q, want codex/gpt-5.4", model)
	}
}

func TestRouteUsageReturnsMultiWindowSnapshots(t *testing.T) {
	gatewayState = newPluginState()
	now := time.Now()
	gatewayState.memberHitCounts["openai/gpt-5.4"] = 3
	gatewayState.ruleHitCounts["route-1"] = 5
	gatewayState.stageHitCounts["route"] = 8
	gatewayState.memberHitTimes["openai/gpt-5.4"] = []time.Time{now.Add(-2 * time.Minute), now.Add(-30 * time.Minute)}
	gatewayState.ruleHitTimes["route-1"] = []time.Time{now.Add(-4 * time.Minute), now.Add(-2 * time.Hour)}
	gatewayState.stageHitTimes["route"] = []time.Time{now.Add(-10 * time.Minute), now.Add(-23 * time.Hour)}
	resp, err := routeUsage(pluginapi.ManagementRequest{})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("routeUsage() = %#v, %v", resp, err)
	}
	body := string(resp.Body)
	if !contains(body, "member_hits_last_5m") || !contains(body, "rule_hits_last_24h") || !contains(body, "stage_hits_last_hour") {
		t.Fatalf("routeUsage body = %s", body)
	}
}

func contains(s, part string) bool {
	return len(s) >= len(part) && (s == part || len(s) > len(part) && (func() bool { return stringIndex(s, part) >= 0 })())
}

func stringIndex(s, sep string) int {
	for i := 0; i+len(sep) <= len(s); i++ {
		if s[i:i+len(sep)] == sep {
			return i
		}
	}
	return -1
}

func applyMemberOperationForTest(t *testing.T, body string) (pluginapi.ManagementResponse, error) {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatalf("unmarshal member operation payload: %v", err)
	}
	payload["preview_only"] = true
	previewBody, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal preview payload: %v", err)
	}
	previewResp, err := routeMemberPreview(pluginapi.ManagementRequest{Body: previewBody})
	if err != nil {
		t.Fatalf("routeMemberPreview() err = %v", err)
	}
	if previewResp.StatusCode != http.StatusOK {
		t.Fatalf("routeMemberPreview() status = %d, body = %s", previewResp.StatusCode, previewResp.Body)
	}
	var preview memberOperationPreview
	if err := json.Unmarshal(previewResp.Body, &preview); err != nil {
		t.Fatalf("unmarshal member operation preview: %v", err)
	}
	if preview.PreviewToken == "" {
		t.Fatal("routeMemberPreview() returned empty preview token")
	}
	payload["preview_token"] = preview.PreviewToken
	delete(payload, "preview_only")
	applyBody, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal apply payload: %v", err)
	}
	return routeMemberOperation(pluginapi.ManagementRequest{Body: applyBody})
}

func TestRouteMemberOperationUpdatesWeightsAndPoolStates(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = normalizeConfig(pluginConfig{KeyPolicies: []keyPolicyConfig{{KeyID: "policy-1", DisplayName: "Policy 1", Enabled: true, Rules: []ruleConfig{{ID: "rule-1", Enabled: true, Priority: 10, Stage: "route", OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{RoutePool: &routePoolConfig{Name: "primary", Mode: "weighted", Members: []weightedRoute{{Model: "openai/gpt-5.4", Weight: 90, Priority: 100, Health: 100, TrafficCap: 100}, {Model: "codex/gpt-5.4", Weight: 10, Priority: 100, Health: 100, TrafficCap: 100}}}}}}}}})

	resp, err := applyMemberOperationForTest(t, `{"key_id":"policy-1","rule_id":"rule-1","member":"openai/gpt-5.4","operation":"weight-up","delta":5}`)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("weight-up = %#v, %v", resp, err)
	}
	member := gatewayState.config.KeyPolicies[0].Rules[0].Actions.RoutePool.Members[0]
	if member.Weight != 95 {
		t.Fatalf("weight after weight-up = %d, want 95", member.Weight)
	}

	resp, err = applyMemberOperationForTest(t, `{"key_id":"policy-1","rule_id":"rule-1","operation":"pool-drain"}`)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("pool-drain = %#v, %v", resp, err)
	}
	for _, item := range gatewayState.config.KeyPolicies[0].Rules[0].Actions.RoutePool.Members {
		if item.Status != "drain" {
			t.Fatalf("pool-drain status = %q, want drain", item.Status)
		}
	}

	resp, err = applyMemberOperationForTest(t, `{"key_id":"policy-1","rule_id":"rule-1","operation":"canary-split","member":"openai/gpt-5.4","secondary":"codex/gpt-5.4","primary_weight":80,"canary_weight":20}`)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("canary-split = %#v, %v", resp, err)
	}
	members := gatewayState.config.KeyPolicies[0].Rules[0].Actions.RoutePool.Members
	if members[0].Weight != 80 || members[1].Weight != 20 {
		t.Fatalf("canary weights = %d/%d, want 80/20", members[0].Weight, members[1].Weight)
	}
}

func TestRouteMemberOperationAppendsOperatorAudit(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = normalizeConfig(pluginConfig{KeyPolicies: []keyPolicyConfig{{KeyID: "policy-audit", DisplayName: "Audit Policy", Enabled: true, Rules: []ruleConfig{{ID: "rule-audit", Enabled: true, Priority: 10, Stage: "route", OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{RoutePool: &routePoolConfig{Name: "primary", Mode: "weighted", Members: []weightedRoute{{Model: "openai/gpt-5.4", Weight: 100, Priority: 100, Health: 100, TrafficCap: 100}, {Model: "codex/gpt-5.4", Weight: 0, Priority: 100, Health: 100, TrafficCap: 100}}}}}}}}})

	resp, err := applyMemberOperationForTest(t, `{"key_id":"policy-audit","rule_id":"rule-audit","member":"openai/gpt-5.4","secondary":"codex/gpt-5.4","operation":"canary-split","primary_weight":95,"canary_weight":5,"reason":"operator-check"}`)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("routeMemberOperation() = %#v, %v", resp, err)
	}
	items := gatewayState.listAudit(10, map[string]string{"event_type": "operator", "operator": "canary-split"})
	if len(items) != 1 {
		t.Fatalf("operator audit items = %d, want 1", len(items))
	}
	if items[0].TargetMember != "openai/gpt-5.4" || items[0].Secondary != "codex/gpt-5.4" {
		t.Fatalf("operator audit target = %#v", items[0])
	}
	if items[0].BeforeState == "" || items[0].AfterState == "" || items[0].BeforeState == items[0].AfterState {
		t.Fatalf("operator audit snapshots = before %q after %q", items[0].BeforeState, items[0].AfterState)
	}
}

func TestDiffWeightedRoutesCapturesChangedMembers(t *testing.T) {
	before := []weightedRoute{{Model: "openai/gpt-5.4", Weight: 90, Priority: 100, Health: 100, TrafficCap: 100, Status: "active"}, {Model: "codex/gpt-5.4", Weight: 10, Priority: 100, Health: 90, TrafficCap: 100, Status: "active"}}
	after := []weightedRoute{{Model: "openai/gpt-5.4", Weight: 80, Priority: 100, Health: 100, TrafficCap: 100, Status: "active"}, {Model: "codex/gpt-5.4", Weight: 20, Priority: 100, Health: 90, TrafficCap: 100, Status: "active"}}
	diff := diffWeightedRoutes(before, after)
	if len(diff) != 2 {
		t.Fatalf("diffWeightedRoutes len = %d, want 2", len(diff))
	}
	if diff[0].Before.Weight == diff[0].After.Weight && diff[1].Before.Weight == diff[1].After.Weight {
		t.Fatalf("diffWeightedRoutes did not capture changed weights: %#v", diff)
	}
}

func TestRouteMemberOperationValidatesPreviewToken(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = normalizeConfig(pluginConfig{KeyPolicies: []keyPolicyConfig{{KeyID: "policy-preview", DisplayName: "Preview Policy", Enabled: true, Rules: []ruleConfig{{ID: "rule-preview", Enabled: true, Priority: 10, Stage: "route", OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{RoutePool: &routePoolConfig{Name: "primary", Mode: "weighted", Members: []weightedRoute{{Model: "openai/gpt-5.4", Weight: 100, Priority: 100, Health: 100, TrafficCap: 100}, {Model: "codex/gpt-5.4", Weight: 0, Priority: 100, Health: 100, TrafficCap: 100}}}}}}}}})

	previewResp, err := routeMemberPreview(pluginapi.ManagementRequest{Body: []byte(`{"key_id":"policy-preview","rule_id":"rule-preview","operation":"pool-drain","reason":"pool-drain","preview_only":true}`)})
	if err != nil || previewResp.StatusCode != http.StatusOK {
		t.Fatalf("routeMemberPreview() = %#v, %v", previewResp, err)
	}
	var preview memberOperationPreview
	if err := json.Unmarshal(previewResp.Body, &preview); err != nil {
		t.Fatalf("unmarshal preview: %v", err)
	}
	if preview.PreviewToken == "" {
		t.Fatal("expected preview token")
	}

	missingResp, err := routeMemberOperation(pluginapi.ManagementRequest{Body: []byte(`{"key_id":"policy-preview","rule_id":"rule-preview","operation":"pool-drain","reason":"pool-drain"}`)})
	if err != nil {
		t.Fatalf("routeMemberOperation missing token err = %v", err)
	}
	if missingResp.StatusCode != http.StatusForbidden {
		t.Fatalf("missing token status = %d, want %d", missingResp.StatusCode, http.StatusForbidden)
	}

	badResp, err := routeMemberOperation(pluginapi.ManagementRequest{Body: []byte(`{"key_id":"policy-preview","rule_id":"rule-preview","operation":"pool-drain","reason":"pool-drain","preview_token":"bad-token"}`)})
	if err != nil {
		t.Fatalf("routeMemberOperation bad token err = %v", err)
	}
	if badResp.StatusCode != http.StatusForbidden {
		t.Fatalf("bad token status = %d, want %d", badResp.StatusCode, http.StatusForbidden)
	}

	goodBody := []byte(`{"key_id":"policy-preview","rule_id":"rule-preview","operation":"pool-drain","reason":"pool-drain","preview_token":"` + preview.PreviewToken + `"}`)
	goodResp, err := routeMemberOperation(pluginapi.ManagementRequest{Body: goodBody})
	if err != nil {
		t.Fatalf("routeMemberOperation good token err = %v", err)
	}
	if goodResp.StatusCode != http.StatusOK {
		t.Fatalf("good token status = %d, want %d", goodResp.StatusCode, http.StatusOK)
	}
	reuseResp, err := routeMemberOperation(pluginapi.ManagementRequest{Body: goodBody})
	if err != nil {
		t.Fatalf("routeMemberOperation reuse token err = %v", err)
	}
	if reuseResp.StatusCode != http.StatusForbidden {
		t.Fatalf("reused token status = %d, want %d", reuseResp.StatusCode, http.StatusForbidden)
	}
}

func TestRouteMemberOperationRejectsPreviewTokenForDifferentTarget(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = normalizeConfig(pluginConfig{KeyPolicies: []keyPolicyConfig{{KeyID: "policy-preview-target", DisplayName: "Preview Policy", Enabled: true, Rules: []ruleConfig{{ID: "rule-preview-target", Enabled: true, Priority: 10, Stage: "route", OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{RoutePool: &routePoolConfig{Name: "primary", Mode: "weighted", Members: []weightedRoute{{Model: "openai/gpt-5.4", Weight: 90, Priority: 100, Health: 100, TrafficCap: 100}, {Model: "codex/gpt-5.4", Weight: 10, Priority: 100, Health: 100, TrafficCap: 100}}}}}}}}})

	previewResp, err := routeMemberPreview(pluginapi.ManagementRequest{Body: []byte(`{"key_id":"policy-preview-target","rule_id":"rule-preview-target","operation":"canary-split","member":"openai/gpt-5.4","secondary":"codex/gpt-5.4","primary_weight":80,"canary_weight":20,"preview_only":true}`)})
	if err != nil || previewResp.StatusCode != http.StatusOK {
		t.Fatalf("routeMemberPreview() = %#v, %v", previewResp, err)
	}
	var preview memberOperationPreview
	if err := json.Unmarshal(previewResp.Body, &preview); err != nil {
		t.Fatalf("unmarshal preview: %v", err)
	}
	badBody := []byte(`{"key_id":"policy-preview-target","rule_id":"rule-preview-target","operation":"canary-split","member":"codex/gpt-5.4","secondary":"openai/gpt-5.4","primary_weight":80,"canary_weight":20,"preview_token":"` + preview.PreviewToken + `"}`)
	badResp, err := routeMemberOperation(pluginapi.ManagementRequest{Body: badBody})
	if err != nil {
		t.Fatalf("routeMemberOperation bad target err = %v", err)
	}
	if badResp.StatusCode != http.StatusForbidden {
		t.Fatalf("bad target status = %d, want %d", badResp.StatusCode, http.StatusForbidden)
	}
}

func TestRouteMemberOperationSupportsWeightedRoutesWithoutRoutePool(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = normalizeConfig(pluginConfig{KeyPolicies: []keyPolicyConfig{{KeyID: "policy-weighted", DisplayName: "Weighted Policy", Enabled: true, Rules: []ruleConfig{{ID: "rule-weighted", Enabled: true, Priority: 10, Stage: "route", OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{WeightedRoutes: []weightedRoute{{Model: "openai/gpt-5.4", Weight: 50}, {Model: "codex/gpt-5.4", Weight: 50}}}}}}}})

	resp, err := applyMemberOperationForTest(t, `{"key_id":"policy-weighted","rule_id":"rule-weighted","member":"openai/gpt-5.4","operation":"weight-up","delta":5}`)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("weighted route member operation = %#v, %v", resp, err)
	}
	members := gatewayState.config.KeyPolicies[0].Rules[0].Actions.WeightedRoutes
	if members[0].Weight != 55 {
		t.Fatalf("weighted route weight = %d, want 55", members[0].Weight)
	}
}

func TestRouteMemberOperationResultMismatchDoesNotMutatePolicy(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = normalizeConfig(pluginConfig{KeyPolicies: []keyPolicyConfig{{KeyID: "policy-mismatch", DisplayName: "Mismatch Policy", Enabled: true, Rules: []ruleConfig{{ID: "rule-mismatch", Enabled: true, Priority: 10, Stage: "route", OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{WeightedRoutes: []weightedRoute{{Model: "openai/gpt-5.4", Weight: 50}, {Model: "codex/gpt-5.4", Weight: 50}}}}}}}})
	before := summarizeWeightedRoutes(gatewayState.config.KeyPolicies[0].Rules[0].Actions.WeightedRoutes)
	issuedAt := time.Now()
	token := signPreviewToken("policy-mismatch", "rule-mismatch", "weight-up", "openai/gpt-5.4", "", before, "wrong-after-state", issuedAt)
	gatewayState.previewTokens[token] = previewTokenRecord{KeyID: "policy-mismatch", RuleID: "rule-mismatch", Operation: "weight-up", Target: "openai/gpt-5.4", BeforeState: before, AfterState: "wrong-after-state", Token: token, IssuedAt: issuedAt}

	resp, err := routeMemberOperation(pluginapi.ManagementRequest{Body: []byte(`{"key_id":"policy-mismatch","rule_id":"rule-mismatch","member":"openai/gpt-5.4","operation":"weight-up","delta":5,"preview_token":"` + token + `"}`)})
	if err != nil {
		t.Fatalf("routeMemberOperation() err = %v", err)
	}
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("result mismatch status = %d, want %d", resp.StatusCode, http.StatusConflict)
	}
	after := summarizeWeightedRoutes(gatewayState.config.KeyPolicies[0].Rules[0].Actions.WeightedRoutes)
	if after != before {
		t.Fatalf("policy mutated after rejected operation: before %q after %q", before, after)
	}
}

func TestSignPreviewTokenIsDeterministic(t *testing.T) {
	issuedAt := time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC)
	a := signPreviewToken("k1", "r1", "pool-drain", "openai/gpt-5.4", "", "before", "after", issuedAt)
	b := signPreviewToken("k1", "r1", "pool-drain", "openai/gpt-5.4", "", "before", "after", issuedAt)
	if a == "" || a != b {
		t.Fatalf("signPreviewToken() = %q / %q", a, b)
	}
}

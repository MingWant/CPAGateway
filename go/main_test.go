package main

import (
	"encoding/json"
	"net/http"
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
	if len(gatewayState.templates) != 2 {
		t.Fatalf("templates len = %d, want 2", len(gatewayState.templates))
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

func TestPluginStateAfterAuthReleasesInflight(t *testing.T) {
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
	if got := state.usage[stableKeyID("abc")].Inflight; got != 0 {
		t.Fatalf("inflight after release = %d, want 0", got)
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

func TestRouteMemberOperationUpdatesWeightsAndPoolStates(t *testing.T) {
	gatewayState = newPluginState()
	gatewayState.config = normalizeConfig(pluginConfig{KeyPolicies: []keyPolicyConfig{{KeyID: "policy-1", DisplayName: "Policy 1", Enabled: true, Rules: []ruleConfig{{ID: "rule-1", Enabled: true, Priority: 10, Stage: "route", OnMatch: "stop", Match: matchConfig{Models: []string{"gpt-5.5"}}, Actions: actionConfig{RoutePool: &routePoolConfig{Name: "primary", Mode: "weighted", Members: []weightedRoute{{Model: "openai/gpt-5.4", Weight: 90, Priority: 100, Health: 100, TrafficCap: 100}, {Model: "codex/gpt-5.4", Weight: 10, Priority: 100, Health: 100, TrafficCap: 100}}}}}}}}})

	resp, err := routeMemberOperation(pluginapi.ManagementRequest{Body: []byte(`{"key_id":"policy-1","rule_id":"rule-1","member":"openai/gpt-5.4","operation":"weight-up","delta":5}`)})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("weight-up = %#v, %v", resp, err)
	}
	member := gatewayState.config.KeyPolicies[0].Rules[0].Actions.RoutePool.Members[0]
	if member.Weight != 95 {
		t.Fatalf("weight after weight-up = %d, want 95", member.Weight)
	}

	resp, err = routeMemberOperation(pluginapi.ManagementRequest{Body: []byte(`{"key_id":"policy-1","rule_id":"rule-1","operation":"pool-drain"}`)})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("pool-drain = %#v, %v", resp, err)
	}
	for _, item := range gatewayState.config.KeyPolicies[0].Rules[0].Actions.RoutePool.Members {
		if item.Status != "drain" {
			t.Fatalf("pool-drain status = %q, want drain", item.Status)
		}
	}

	resp, err = routeMemberOperation(pluginapi.ManagementRequest{Body: []byte(`{"key_id":"policy-1","rule_id":"rule-1","operation":"canary-split","member":"openai/gpt-5.4","secondary":"codex/gpt-5.4","primary_weight":80,"canary_weight":20}`)})
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

	resp, err := routeMemberOperation(pluginapi.ManagementRequest{Body: []byte(`{"key_id":"policy-audit","rule_id":"rule-audit","member":"openai/gpt-5.4","secondary":"codex/gpt-5.4","operation":"canary-split","primary_weight":95,"canary_weight":5,"reason":"operator-check"}`)})
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

func TestSignPreviewTokenIsDeterministic(t *testing.T) {
	issuedAt := time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC)
	a := signPreviewToken("k1", "r1", "pool-drain", "openai/gpt-5.4", "", "before", "after", issuedAt)
	b := signPreviewToken("k1", "r1", "pool-drain", "openai/gpt-5.4", "", "before", "after", issuedAt)
	if a == "" || a != b {
		t.Fatalf("signPreviewToken() = %q / %q", a, b)
	}
}

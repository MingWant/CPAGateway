package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/router-for-me/CLIProxyAPI/v7/sdk/pluginapi"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type actionTraceDetail struct {
	RoutePool       string
	RouteTarget     string
	FallbackTarget  string
	MirrorModels    []string
	FailoverChain   []string
	FailoverReasons []string
}

func applyRules(req pluginapi.RequestInterceptRequest, policy policyConfig, afterAuth bool, now time.Time) dryRunResult {
	result := applyRulesWithStages(req, policy, afterAuth, now)
	return dryRunResult{Decision: result.Decision, RuleID: result.RuleID, Reason: result.Reason, MatchedRules: result.MatchedRules, FinalModel: result.FinalModel, Response: result.Response, StageTrace: result.StageTrace}
}

func applyRulesWithStages(req pluginapi.RequestInterceptRequest, policy policyConfig, afterAuth bool, now time.Time) stageRunResult {
	currentHeaders := cloneHeader(req.Headers)
	currentBody := cloneBytes(req.Body)
	currentModel := req.Model
	allMatched := make([]string, 0)
	trace := make([]stageTrace, 0)
	finalDecision := "pass"
	finalRuleID := ""
	finalReason := ""
	for _, stage := range orderedStages(policy.Rules) {
		stageMatched := make([]string, 0)
		stageMode := stageModeFor(policy, stage)
		stageDecision := "pass"
		stageReason := ""
		stageDetail := actionTraceDetail{}
		for _, rule := range rulesForStage(policy.Rules, stage) {
			if !rule.Enabled {
				continue
			}
			ok, reason := matchRule(req, currentModel, rule.Match, now)
			if !ok {
				_ = reason
				continue
			}
			stageMatched = append(stageMatched, rule.ID)
			allMatched = append(allMatched, rule.ID)
			resp, nextBody, nextModel, actionReason, detail := applyRuleActions(req, currentHeaders, currentBody, currentModel, rule.Actions, afterAuth)
			currentHeaders = cloneHeader(resp.Headers)
			if len(nextBody) > 0 {
				currentBody = cloneBytes(nextBody)
			}
			currentModel = nextModel
			decision := "rewrite"
			if resp.Reject {
				decision = "reject"
			}
			stageDecision = decision
			stageReason = actionReason
			stageDetail = detail
			finalDecision = decision
			finalRuleID = rule.ID
			finalReason = actionReason
			shouldStop := resp.Reject || strings.EqualFold(strings.TrimSpace(rule.OnMatch), "stop") || strings.EqualFold(stageMode, "first-match")
			if shouldStop {
				trace = append(trace, stageTrace{Stage: stage, Mode: stageMode, MatchedRules: append([]string(nil), stageMatched...), MatchedCount: len(stageMatched), FinalModel: currentModel, Decision: decision, Reason: actionReason, RoutePool: detail.RoutePool, RouteTarget: detail.RouteTarget, FallbackTarget: detail.FallbackTarget, MirrorModels: append([]string(nil), detail.MirrorModels...), FailoverChain: append([]string(nil), detail.FailoverChain...), FailoverReasons: append([]string(nil), detail.FailoverReasons...)})
				if resp.Reject || strings.EqualFold(strings.TrimSpace(rule.OnMatch), "stop") {
					return stageRunResult{Decision: decision, RuleID: rule.ID, Reason: actionReason, MatchedRules: allMatched, FinalModel: currentModel, Response: resp, StageTrace: trace}
				}
				break
			}
		}
		if len(stageMatched) == 0 || stageDecision == "pass" {
			trace = append(trace, stageTrace{Stage: stage, Mode: stageMode, MatchedRules: append([]string(nil), stageMatched...), MatchedCount: len(stageMatched), FinalModel: currentModel, Decision: "pass"})
		} else if len(trace) == 0 || trace[len(trace)-1].Stage != stage {
			trace = append(trace, stageTrace{Stage: stage, Mode: stageMode, MatchedRules: append([]string(nil), stageMatched...), MatchedCount: len(stageMatched), FinalModel: currentModel, Decision: stageDecision, Reason: stageReason, RoutePool: stageDetail.RoutePool, RouteTarget: stageDetail.RouteTarget, FallbackTarget: stageDetail.FallbackTarget, MirrorModels: append([]string(nil), stageDetail.MirrorModels...), FailoverChain: append([]string(nil), stageDetail.FailoverChain...), FailoverReasons: append([]string(nil), stageDetail.FailoverReasons...)})
		}
	}
	return stageRunResult{Decision: finalDecision, RuleID: finalRuleID, Reason: finalReason, MatchedRules: allMatched, FinalModel: currentModel, Response: requestInterceptResponse{Headers: currentHeaders, Body: currentBody}, StageTrace: trace}
}

func matchRule(req pluginapi.RequestInterceptRequest, currentModel string, match matchConfig, now time.Time) (bool, string) {
	if len(match.AllOf) > 0 {
		for _, nested := range match.AllOf {
			ok, _ := matchRule(req, currentModel, nested, now)
			if !ok {
				return false, "all_of_mismatch"
			}
		}
	}
	if len(match.AnyOf) > 0 {
		matched := false
		for _, nested := range match.AnyOf {
			ok, _ := matchRule(req, currentModel, nested, now)
			if ok {
				matched = true
				break
			}
		}
		if !matched {
			return false, "any_of_mismatch"
		}
	}
	if match.Stream != nil && req.Stream != *match.Stream {
		return false, "stream_mismatch"
	}
	path := stringMetadata(req.Metadata, "request_path")
	if path == "" {
		path = stringMetadata(req.Metadata, "request.path")
	}
	pathKind := requestPathKind(path, req.SourceFormat)
	if len(match.PathKinds) > 0 && !containsFold(match.PathKinds, pathKind) {
		return false, "path_kind_mismatch"
	}
	if len(match.Paths) > 0 && !containsFold(match.Paths, path) {
		return false, "path_mismatch"
	}
	if len(match.Models) > 0 && !containsFold(match.Models, currentModel) {
		return false, "model_mismatch"
	}
	if len(match.ModelPrefixes) > 0 && !hasAnyPrefixFold(currentModel, match.ModelPrefixes) {
		return false, "model_prefix_mismatch"
	}
	provider := providerFromModel(currentModel)
	if len(match.Providers) > 0 && !containsFold(match.Providers, provider) {
		return false, "provider_mismatch"
	}
	for key, want := range match.Headers {
		if strings.TrimSpace(req.Headers.Get(key)) != strings.TrimSpace(want) {
			return false, "header_mismatch"
		}
	}
	for key, want := range match.Query {
		queryValue := strings.TrimSpace(stringMetadata(req.Metadata, "query."+key))
		if queryValue == "" {
			queryValue = queryValueFromMetadata(req.Metadata, key)
		}
		if !strings.Contains(queryValue, strings.TrimSpace(want)) {
			return false, "query_mismatch"
		}
	}
	for key, want := range match.BodyContains {
		value := strings.TrimSpace(gjson.GetBytes(req.Body, key).String())
		if !strings.Contains(value, strings.TrimSpace(want)) {
			return false, "body_mismatch"
		}
	}
	for key, want := range match.MetadataContains {
		if !strings.Contains(strings.TrimSpace(stringMetadata(req.Metadata, key)), strings.TrimSpace(want)) {
			return false, "metadata_mismatch"
		}
	}
	if len(match.Days) > 0 {
		matched := false
		for _, day := range match.Days {
			if strings.EqualFold(strings.TrimSpace(day), now.Weekday().String()) {
				matched = true
				break
			}
		}
		if !matched {
			return false, "day_mismatch"
		}
	}
	if match.Start != "" && match.End != "" {
		if !withinClockWindow(now.Format("15:04"), match.Start, match.End) {
			return false, "time_window_mismatch"
		}
	}
	return true, "matched"
}

func applyRuleActions(req pluginapi.RequestInterceptRequest, headers http.Header, body []byte, currentModel string, actions actionConfig, afterAuth bool) (requestInterceptResponse, []byte, string, string, actionTraceDetail) {
	resp := requestInterceptResponse{Headers: cloneHeader(headers), Body: cloneBytes(body)}
	model := currentModel
	detail := actionTraceDetail{}
	if actions.RoutePool != nil {
		detail.RoutePool = strings.TrimSpace(actions.RoutePool.Name)
	}
	if actions.Deny != nil {
		status := actions.Deny.StatusCode
		if status == 0 {
			status = http.StatusForbidden
		}
		return requestInterceptResponse{Reject: true, RejectStatusCode: status, RejectMessage: actions.Deny.Message, RejectCode: actions.Deny.Code}, body, currentModel, "deny", detail
	}
	reason := "rewrite"
	if strings.TrimSpace(actions.RewriteModel) != "" {
		model = strings.TrimSpace(actions.RewriteModel)
		detail.RouteTarget = model
		reason = "rewrite_model"
	}
	if routed := selectRouteTarget(req, actions); routed != "" {
		model = routed
		detail.RouteTarget = routed
		if actions.RoutePool != nil && len(actions.RoutePool.Members) > 0 {
			reason = "route_pool"
		} else {
			reason = "weighted_route"
		}
	}
	if strings.TrimSpace(actions.RouteToModel) != "" {
		model = strings.TrimSpace(actions.RouteToModel)
		detail.RouteTarget = model
		reason = "route_to_model"
	}
	if len(actions.FallbackModels) > 0 || len(actions.FailoverChain) > 0 || len(actions.FailoverHops) > 0 {
		merged := mergedFailoverChain(actions)
		if len(merged) > 0 {
			detail.FailoverChain = append([]string(nil), merged...)
		}
		if fallback := selectFallbackModel(currentModel, model, merged); fallback != "" {
			model = fallback
			detail.FallbackTarget = fallback
			if len(actions.FailoverChain) > 0 || len(actions.FailoverHops) > 0 {
				reason = "failover_chain"
			} else {
				reason = "fallback_model"
			}
		}
		if len(actions.FailoverHops) > 0 {
			reasons := make([]string, 0, len(actions.FailoverHops))
			for _, hop := range actions.FailoverHops {
				if strings.TrimSpace(hop.Reason) == "" {
					continue
				}
				reasons = append(reasons, strings.TrimSpace(hop.Reason))
			}
			detail.FailoverReasons = reasons
		}
	}
	if strings.TrimSpace(actions.ForceProviderPrefix) != "" {
		model = forceProviderPrefix(model, strings.TrimSpace(actions.ForceProviderPrefix))
		detail.RouteTarget = model
		if reason == "rewrite" {
			reason = "force_provider_prefix"
		}
	}
	if actions.RewriteEndpoint != nil {
		candidate := strings.TrimSpace(actions.RewriteEndpoint.TargetModel)
		if candidate != "" {
			model = candidate
			detail.RouteTarget = model
		}
		if reason == "rewrite" {
			reason = "rewrite_endpoint_semantics"
		}
	}
	if len(actions.AllowOnlyProviders) > 0 && !containsFold(actions.AllowOnlyProviders, providerFromModel(model)) {
		return requestInterceptResponse{Reject: true, RejectStatusCode: http.StatusForbidden, RejectMessage: "gateway provider policy rejected request", RejectCode: "gateway_provider_denied"}, body, currentModel, "allow_only_providers", detail
	}
	if len(actions.AllowOnlyModels) > 0 && !containsFold(actions.AllowOnlyModels, model) {
		return requestInterceptResponse{Reject: true, RejectStatusCode: http.StatusForbidden, RejectMessage: "gateway model policy rejected request", RejectCode: "gateway_model_denied"}, body, currentModel, "allow_only_models", detail
	}
	if model != currentModel {
		if updated, ok := rewriteModelInBody(resp.Body, model); ok {
			resp.Body = updated
		}
	}
	if actions.RewriteEndpoint != nil && !afterAuth && strings.EqualFold(req.SourceFormat, "openai-response") {
		if updated, ok := rewriteResponsesCompatibility(resp.Body, model); ok {
			resp.Body = updated
		}
	}
	for key, value := range actions.SetHeaders {
		if resp.Headers == nil {
			resp.Headers = make(http.Header)
		}
		resp.Headers.Set(key, value)
	}
	if len(actions.MirrorModels) > 0 {
		models := compactUnique(actions.MirrorModels)
		detail.MirrorModels = append([]string(nil), models...)
		if resp.Headers == nil {
			resp.Headers = make(http.Header)
		}
		resp.Headers.Set("X-Gateway-Mirror-Models", strings.Join(models, ","))
		if reason == "rewrite" {
			reason = "mirror_models"
		}
	}
	if len(actions.TagMetadata) > 0 {
		if resp.Headers == nil {
			resp.Headers = make(http.Header)
		}
		for key, value := range actions.TagMetadata {
			trimmedKey := strings.TrimSpace(key)
			trimmedValue := strings.TrimSpace(value)
			if trimmedKey == "" || trimmedValue == "" {
				continue
			}
			resp.Headers.Set("X-Gateway-Tag-"+sanitizeHeaderToken(trimmedKey), trimmedValue)
		}
	}
	if len(actions.ClearHeaders) > 0 {
		resp.ClearHeaders = append([]string(nil), actions.ClearHeaders...)
	}
	return resp, resp.Body, model, reason, detail
}

func rewriteModelInBody(body []byte, model string) ([]byte, bool) {
	if len(body) == 0 || !gjson.ValidBytes(body) {
		return body, false
	}
	updated, err := sjson.SetBytes(body, "model", model)
	if err != nil {
		return body, false
	}
	return updated, true
}

func rewriteResponsesCompatibility(body []byte, model string) ([]byte, bool) {
	if len(body) == 0 || !gjson.ValidBytes(body) {
		return body, false
	}
	updated, ok := rewriteModelInBody(body, model)
	if !ok {
		return body, false
	}
	var payload map[string]any
	if err := json.Unmarshal(updated, &payload); err != nil {
		return updated, true
	}
	if _, exists := payload["input"]; exists {
		return updated, true
	}
	if messages, exists := payload["messages"]; exists {
		payload["input"] = messages
		delete(payload, "messages")
	}
	finalBody, err := json.Marshal(payload)
	if err != nil {
		return updated, true
	}
	return finalBody, true
}

func requestPathKind(path, sourceFormat string) string {
	path = strings.ToLower(strings.TrimSpace(path))
	switch {
	case strings.Contains(path, "/responses"):
		return "/responses"
	case strings.Contains(path, "/messages"):
		return "/messages"
	case strings.Contains(path, "/chat/completions"):
		return "/chat/completions"
	}
	if strings.EqualFold(sourceFormat, "openai-response") {
		return "/responses"
	}
	return sourceFormat
}

func providerFromModel(model string) string {
	trimmed := strings.TrimSpace(model)
	if idx := strings.Index(trimmed, "/"); idx > 0 {
		return strings.ToLower(trimmed[:idx])
	}
	return ""
}

func modelSuffix(model string) string {
	trimmed := strings.TrimSpace(model)
	if idx := strings.Index(trimmed, "/"); idx >= 0 && idx+1 < len(trimmed) {
		return trimmed[idx+1:]
	}
	return trimmed
}

func forceProviderPrefix(model, provider string) string {
	provider = strings.Trim(strings.TrimSpace(provider), "/")
	if provider == "" {
		return model
	}
	base := model
	if idx := strings.Index(base, "/"); idx >= 0 {
		base = base[idx+1:]
	}
	return provider + "/" + base
}

func containsFold(items []string, target string) bool {
	for _, item := range items {
		if strings.EqualFold(strings.TrimSpace(item), strings.TrimSpace(target)) {
			return true
		}
	}
	return false
}

func hasAnyPrefixFold(value string, prefixes []string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	for _, prefix := range prefixes {
		if strings.HasPrefix(value, strings.ToLower(strings.TrimSpace(prefix))) {
			return true
		}
	}
	return false
}

func selectFallbackModel(currentModel, resolvedModel string, fallbacks []string) string {
	if len(fallbacks) == 0 {
		return ""
	}
	if strings.TrimSpace(resolvedModel) != "" && !strings.EqualFold(strings.TrimSpace(resolvedModel), strings.TrimSpace(currentModel)) {
		return ""
	}
	for _, candidate := range fallbacks {
		trimmed := strings.TrimSpace(candidate)
		if trimmed == "" {
			continue
		}
		if strings.EqualFold(trimmed, strings.TrimSpace(currentModel)) {
			continue
		}
		return trimmed
	}
	return ""
}

func queryValueFromMetadata(meta map[string]any, key string) string {
	raw := strings.TrimSpace(stringMetadata(meta, "request.query"))
	if raw == "" {
		return ""
	}
	for _, pair := range strings.Split(raw, "&") {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 && strings.EqualFold(strings.TrimSpace(parts[0]), strings.TrimSpace(key)) {
			return strings.TrimSpace(parts[1])
		}
	}
	return ""
}

func selectWeightedRoute(req pluginapi.RequestInterceptRequest, actions actionConfig) string {
	return selectRouteTarget(req, actions)
}

func selectWeightedModel(req pluginapi.RequestInterceptRequest, shardBy string, routes []weightedRoute) string {
	routes = normalizedWeightedRoutes(routes)
	if len(routes) == 0 {
		return ""
	}
	effective := make([]weightedRoute, 0, len(routes))
	for _, route := range routes {
		adjusted := effectiveRouteWeight(route)
		if adjusted <= 0 {
			continue
		}
		route.Weight = adjusted
		effective = append(effective, route)
	}
	if len(effective) == 0 {
		return ""
	}
	seed := routingSeed(req, shardBy)
	if seed == "" {
		seed = req.Model + "|" + stringMetadata(req.Metadata, "access.api_key") + "|" + stringMetadata(req.Metadata, "request_path")
	}
	hash := stableHash32(seed)
	total := 0
	for _, route := range effective {
		total += route.Weight
	}
	if total <= 0 {
		return effective[0].Model
	}
	pick := int(hash % uint32(total))
	cursor := 0
	for _, route := range effective {
		cursor += route.Weight
		if pick < cursor {
			return route.Model
		}
	}
	return effective[len(effective)-1].Model
}

func normalizedWeightedRoutes(routes []weightedRoute) []weightedRoute {
	out := make([]weightedRoute, 0, len(routes))
	for _, route := range routes {
		if route.Enabled != nil && !*route.Enabled {
			continue
		}
		if statusBlocksRoute(route.Status) {
			continue
		}
		model := strings.TrimSpace(route.Model)
		provider := strings.TrimSpace(route.Provider)
		suffix := strings.TrimSpace(route.Suffix)
		if model == "" {
			if provider != "" && suffix != "" {
				model = provider + "/" + strings.TrimPrefix(suffix, "/")
			} else if provider != "" && suffix == "" {
				model = providerFromModel(route.Model)
			}
		}
		if model == "" {
			continue
		}
		weight := route.Weight
		if weight <= 0 {
			weight = 1
		}
		out = append(out, weightedRoute{Model: model, Provider: providerFromModel(model), Suffix: modelSuffix(model), Weight: weight, Priority: route.Priority, Enabled: route.Enabled, Status: strings.ToLower(strings.TrimSpace(route.Status)), Reason: route.Reason, Health: route.Health, TrafficCap: route.TrafficCap})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Priority == out[j].Priority {
			return out[i].Model < out[j].Model
		}
		return out[i].Priority < out[j].Priority
	})
	return out
}

func routingSeed(req pluginapi.RequestInterceptRequest, shardBy string) string {
	switch strings.ToLower(strings.TrimSpace(shardBy)) {
	case "api_key":
		return stringMetadata(req.Metadata, "access.api_key")
	case "model":
		if strings.TrimSpace(req.RequestedModel) != "" {
			return req.RequestedModel
		}
		return req.Model
	case "path":
		return stringMetadata(req.Metadata, "request_path")
	case "user_agent":
		return req.Headers.Get("User-Agent")
	case "ip", "client_ip":
		return stringMetadata(req.Metadata, "client.ip")
	case "query":
		return stringMetadata(req.Metadata, "request.query")
	case "header":
		return req.Headers.Get("X-Request-Id")
	default:
		return stringMetadata(req.Metadata, "access.api_key")
	}
}

func sanitizeHeaderToken(value string) string {
	replacer := strings.NewReplacer(" ", "-", ".", "-", "_", "-", "/", "-", ":", "-")
	value = replacer.Replace(strings.TrimSpace(value))
	value = strings.Trim(value, "-")
	if value == "" {
		return "Tag"
	}
	return value
}

func compactUnique(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	out := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, trimmed)
	}
	return out
}

func stableHash32(value string) uint32 {
	var hash uint32 = 2166136261
	for i := 0; i < len(value); i++ {
		hash ^= uint32(value[i])
		hash *= 16777619
	}
	return hash
}

func orderedStages(rules []ruleConfig) []string {
	seen := map[string]struct{}{}
	base := []string{"pre-check", "rewrite", "route", "mirror", "post-audit"}
	out := make([]string, 0, len(base))
	for _, stage := range base {
		out = append(out, stage)
		seen[stage] = struct{}{}
	}
	for _, rule := range rules {
		stage := normalizeRuleStage(rule.Stage)
		if _, exists := seen[stage]; exists {
			continue
		}
		seen[stage] = struct{}{}
		out = append(out, stage)
	}
	return out
}

func rulesForStage(rules []ruleConfig, stage string) []ruleConfig {
	out := make([]ruleConfig, 0)
	for _, rule := range rules {
		if normalizeRuleStage(rule.Stage) != stage {
			continue
		}
		out = append(out, rule)
	}
	return out
}

func normalizeRuleStage(stage string) string {
	switch strings.ToLower(strings.TrimSpace(stage)) {
	case "", "precheck", "pre-check":
		return "pre-check"
	case "rewrite":
		return "rewrite"
	case "route":
		return "route"
	case "mirror":
		return "mirror"
	case "postaudit", "post-audit", "post_audit":
		return "post-audit"
	default:
		return strings.ToLower(strings.TrimSpace(stage))
	}
}

func stageModeFor(policy policyConfig, stage string) string {
	if policy.StagePolicy != nil {
		if cfg, ok := policy.StagePolicy[stage]; ok {
			if normalized := normalizeStageMode(cfg.Mode); normalized != "" {
				return normalized
			}
		}
	}
	switch stage {
	case "mirror", "post-audit":
		return "continue-all"
	default:
		return "first-match"
	}
}

func normalizeStageMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "", "first", "first-match", "first_match":
		return "first-match"
	case "continue", "continue-all", "continue_all":
		return "continue-all"
	default:
		return ""
	}
}

func selectRouteTarget(req pluginapi.RequestInterceptRequest, actions actionConfig) string {
	if actions.RoutePool != nil && len(actions.RoutePool.Members) > 0 {
		members := normalizedWeightedRoutes(actions.RoutePool.Members)
		if affinity := strings.TrimSpace(actions.RoutePool.ProviderAffinity); affinity != "" {
			affinityMembers := make([]weightedRoute, 0, len(members))
			for _, member := range members {
				if strings.EqualFold(providerFromModel(member.Model), affinity) {
					affinityMembers = append(affinityMembers, member)
				}
			}
			if len(affinityMembers) > 0 {
				members = affinityMembers
			}
		}
		return selectWeightedModel(req, actions.ShardBy, members)
	}
	if len(actions.WeightedRoutes) > 0 {
		return selectWeightedModel(req, actions.ShardBy, actions.WeightedRoutes)
	}
	return ""
}

func mergedFailoverChain(actions actionConfig) []string {
	chain := make([]string, 0, len(actions.FailoverHops)+len(actions.FailoverChain)+len(actions.FallbackModels))
	for _, hop := range actions.FailoverHops {
		if hop.Enabled != nil && !*hop.Enabled {
			continue
		}
		model := strings.TrimSpace(hop.Model)
		if model == "" && strings.TrimSpace(hop.Provider) != "" && strings.TrimSpace(hop.Suffix) != "" {
			model = strings.TrimSpace(hop.Provider) + "/" + strings.TrimPrefix(strings.TrimSpace(hop.Suffix), "/")
		}
		if model == "" || containsFold(chain, model) {
			continue
		}
		chain = append(chain, model)
	}
	for _, item := range actions.FailoverChain {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" || containsFold(chain, trimmed) {
			continue
		}
		chain = append(chain, trimmed)
	}
	for _, item := range actions.FallbackModels {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" || containsFold(chain, trimmed) {
			continue
		}
		chain = append(chain, trimmed)
	}
	return chain
}

func statusBlocksRoute(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "disabled", "drain", "drained", "degraded", "offline":
		return true
	default:
		return false
	}
}

func effectiveRouteWeight(route weightedRoute) int {
	weight := route.Weight
	if weight <= 0 {
		weight = 1
	}
	if route.TrafficCap > 0 && route.TrafficCap < 100 {
		weight = weight * route.TrafficCap / 100
	}
	if route.Health > 0 && route.Health < 100 {
		weight = weight * route.Health / 100
	}
	if route.Health == 0 && route.TrafficCap == 0 {
		return weight
	}
	if weight <= 0 {
		return 0
	}
	return weight
}

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/router-for-me/CLIProxyAPI/v7/sdk/pluginapi"
)

func normalizeConfig(cfg pluginConfig) pluginConfig {
	cfg.Default = normalizePolicy(cfg.Default)
	cfg.KeyPolicies = normalizeKeyPolicies(cfg.KeyPolicies)
	cfg.Cluster = normalizeClusterConfig(cfg.Cluster)
	return cfg
}

func normalizeClusterConfig(cfg clusterConfig) clusterConfig {
	cfg.Backend = strings.ToLower(strings.TrimSpace(cfg.Backend))
	cfg.Redis.Addr = strings.TrimSpace(cfg.Redis.Addr)
	cfg.Redis.Username = strings.TrimSpace(cfg.Redis.Username)
	cfg.Redis.KeyPrefix = strings.TrimSpace(cfg.Redis.KeyPrefix)
	switch strings.ToLower(strings.TrimSpace(cfg.Redis.FailureMode)) {
	case "", "reject", "fail_closed", "fail-closed":
		cfg.Redis.FailureMode = "reject"
	case "allow", "fail_open", "fail-open":
		cfg.Redis.FailureMode = "allow"
	case "local", "local_fallback", "local-fallback":
		cfg.Redis.FailureMode = "local_fallback"
	default:
		cfg.Redis.FailureMode = "reject"
	}
	return cfg
}

func normalizePolicy(policy policyConfig) policyConfig {
	rules := append([]ruleConfig(nil), policy.Rules...)
	sort.SliceStable(rules, func(i, j int) bool { return rules[i].Priority < rules[j].Priority })
	policy.Rules = rules
	return policy
}

func policyHasContent(policy policyConfig) bool {
	return policy.Enabled || len(policy.Rules) > 0 || len(policy.StagePolicy) > 0 || !limitConfigEmpty(policy.Limits)
}

func limitConfigEmpty(limits limitConfig) bool {
	return limits.RequestsPerDay == 0 &&
		limits.TokensPerDay == 0 &&
		limits.RequestsPerMin == 0 &&
		limits.MaxInflight == 0 &&
		strings.TrimSpace(limits.NotBefore) == "" &&
		strings.TrimSpace(limits.NotAfter) == "" &&
		len(limits.Schedules) == 0
}

func normalizeKeyPolicies(items []keyPolicyConfig) []keyPolicyConfig {
	out := make([]keyPolicyConfig, 0, len(items))
	for _, item := range items {
		item.KeyID = candidateKeyID(item)
		item.MaskedKey = ""
		item.Rules = normalizePolicy(policyConfig{Rules: item.Rules}).Rules
		out = append(out, item)
	}
	return out
}

func sanitizedStoredPolicy(policy storedPolicy) storedPolicy {
	policy.DefaultPolicy = clonePolicyConfig(policy.DefaultPolicy)
	policy.KeyPolicies = sanitizedKeyPolicies(policy.KeyPolicies)
	return policy
}

func sanitizedPolicyBundle(bundle policyBundle) policyBundle {
	bundle.DefaultPolicy = clonePolicyConfig(bundle.DefaultPolicy)
	bundle.KeyPolicies = sanitizedKeyPolicies(bundle.KeyPolicies)
	return bundle
}

func sanitizedKeyPolicies(items []keyPolicyConfig) []keyPolicyConfig {
	out := make([]keyPolicyConfig, 0, len(items))
	for _, item := range items {
		clone := cloneKeyPolicyConfig(item)
		clone.MaskedKey = maskKey(clone.MatchAPIKey)
		clone.MatchAPIKey = ""
		out = append(out, clone)
	}
	return out
}

func (s *pluginState) preservePolicySecretsLocked(items []keyPolicyConfig) []keyPolicyConfig {
	if len(items) == 0 {
		return nil
	}
	secrets := make(map[string]string, len(s.config.KeyPolicies))
	for _, item := range s.config.KeyPolicies {
		if keyID := candidateKeyID(item); keyID != "" && strings.TrimSpace(item.MatchAPIKey) != "" {
			secrets[keyID] = strings.TrimSpace(item.MatchAPIKey)
		}
	}
	out := make([]keyPolicyConfig, len(items))
	for i, item := range items {
		if strings.TrimSpace(item.MatchAPIKey) == "" {
			if secret := secrets[candidateKeyID(item)]; secret != "" {
				item.MatchAPIKey = secret
			}
		}
		out[i] = item
	}
	return out
}

func clonePolicyConfig(policy policyConfig) policyConfig {
	policy.Rules = cloneRuleConfigs(policy.Rules)
	if policy.StagePolicy != nil {
		policy.StagePolicy = cloneMap(policy.StagePolicy)
	}
	return policy
}

func cloneKeyPolicyConfig(item keyPolicyConfig) keyPolicyConfig {
	item.Rules = cloneRuleConfigs(item.Rules)
	item.StagePolicy = cloneMap(item.StagePolicy)
	return item
}

func cloneRuleConfigs(items []ruleConfig) []ruleConfig {
	if len(items) == 0 {
		return nil
	}
	out := make([]ruleConfig, len(items))
	for i, item := range items {
		out[i] = cloneRuleConfig(item)
	}
	return out
}

func cloneRuleConfig(rule ruleConfig) ruleConfig {
	rule.Match = cloneMatchConfig(rule.Match)
	rule.Actions = cloneActionConfig(rule.Actions)
	return rule
}

func cloneMatchConfig(match matchConfig) matchConfig {
	match.PathKinds = append([]string(nil), match.PathKinds...)
	match.Paths = append([]string(nil), match.Paths...)
	match.Models = append([]string(nil), match.Models...)
	match.ModelPrefixes = append([]string(nil), match.ModelPrefixes...)
	match.Providers = append([]string(nil), match.Providers...)
	match.Headers = cloneStringMap(match.Headers)
	match.Query = cloneStringMap(match.Query)
	match.BodyContains = cloneStringMap(match.BodyContains)
	match.MetadataContains = cloneStringMap(match.MetadataContains)
	match.Days = append([]string(nil), match.Days...)
	match.AnyOf = cloneMatchConfigs(match.AnyOf)
	match.AllOf = cloneMatchConfigs(match.AllOf)
	if match.Stream != nil {
		value := *match.Stream
		match.Stream = &value
	}
	return match
}

func cloneMatchConfigs(items []matchConfig) []matchConfig {
	if len(items) == 0 {
		return nil
	}
	out := make([]matchConfig, len(items))
	for i, item := range items {
		out[i] = cloneMatchConfig(item)
	}
	return out
}

func cloneActionConfig(actions actionConfig) actionConfig {
	actions.WeightedRoutes = cloneWeightedRoutes(actions.WeightedRoutes)
	if actions.RoutePool != nil {
		pool := *actions.RoutePool
		pool.Members = cloneWeightedRoutes(pool.Members)
		actions.RoutePool = &pool
	}
	actions.FailoverChain = append([]string(nil), actions.FailoverChain...)
	actions.FailoverHops = append([]failoverHop(nil), actions.FailoverHops...)
	actions.MirrorModels = append([]string(nil), actions.MirrorModels...)
	actions.AllowOnlyProviders = append([]string(nil), actions.AllowOnlyProviders...)
	actions.AllowOnlyModels = append([]string(nil), actions.AllowOnlyModels...)
	actions.FallbackModels = append([]string(nil), actions.FallbackModels...)
	if actions.Deny != nil {
		deny := *actions.Deny
		actions.Deny = &deny
	}
	if actions.RewriteEndpoint != nil {
		rewrite := *actions.RewriteEndpoint
		actions.RewriteEndpoint = &rewrite
	}
	actions.SetHeaders = cloneStringMap(actions.SetHeaders)
	actions.ClearHeaders = append([]string(nil), actions.ClearHeaders...)
	actions.TagMetadata = cloneStringMap(actions.TagMetadata)
	return actions
}

func cloneStringMap(src map[string]string) map[string]string {
	if len(src) == 0 {
		return nil
	}
	out := make(map[string]string, len(src))
	for key, value := range src {
		out[key] = value
	}
	return out
}

func cloneMap[V any](src map[string]V) map[string]V {
	if len(src) == 0 {
		return nil
	}
	out := make(map[string]V, len(src))
	for key, value := range src {
		out[key] = value
	}
	return out
}

func candidateKeyID(item keyPolicyConfig) string {
	if strings.TrimSpace(item.KeyID) != "" {
		return strings.TrimSpace(item.KeyID)
	}
	return stableKeyID(item.MatchAPIKey)
}

func stableKeyID(apiKey string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(apiKey)))
	return hex.EncodeToString(sum[:8])
}
func previewTokenSecret() string {
	return "gateway-preview-v1"
}

func signPreviewToken(keyID, ruleID, operation, target, secondary, beforeState, afterState string, issuedAt time.Time) string {
	payload := strings.Join([]string{previewTokenSecret(), strings.TrimSpace(keyID), strings.TrimSpace(ruleID), strings.TrimSpace(operation), strings.TrimSpace(target), strings.TrimSpace(secondary), strings.TrimSpace(beforeState), strings.TrimSpace(afterState), issuedAt.UTC().Format(time.RFC3339Nano)}, "|")
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:16])
}

func maskKey(apiKey string) string {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return ""
	}
	if len(apiKey) <= 6 {
		return strings.Repeat("*", len(apiKey))
	}
	return apiKey[:3] + strings.Repeat("*", len(apiKey)-6) + apiKey[len(apiKey)-3:]
}

func stringListHeader(headers http.Header, key string) []string {
	if headers == nil {
		return nil
	}
	values := headers.Values(key)
	if len(values) == 0 {
		raw := strings.TrimSpace(headers.Get(key))
		if raw == "" {
			return nil
		}
		values = []string{raw}
	}
	out := make([]string, 0, len(values))
	for _, value := range values {
		for _, item := range strings.Split(value, ",") {
			trimmed := strings.TrimSpace(item)
			if trimmed == "" {
				continue
			}
			out = append(out, trimmed)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func withGatewayMetadata(resp pluginapi.RequestInterceptResponse, attrs map[string]string) pluginapi.RequestInterceptResponse {
	if len(attrs) == 0 {
		return resp
	}
	if resp.Headers == nil {
		resp.Headers = make(http.Header)
	}
	for key, value := range attrs {
		if strings.TrimSpace(value) == "" {
			continue
		}
		resp.Headers.Set("X-"+strings.NewReplacer(".", "-", "_", "-").Replace(strings.Title(strings.ReplaceAll(key, "gateway.", "gateway-"))), value)
	}
	return resp
}

func cloneHeader(src http.Header) http.Header {
	if src == nil {
		return nil
	}
	return src.Clone()
}

func cloneBytes(src []byte) []byte {
	if src == nil {
		return nil
	}
	return append([]byte(nil), src...)
}

func stringMetadata(meta map[string]any, key string) string {
	if meta == nil {
		return ""
	}
	value, ok := meta[key]
	if !ok || value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case fmt.Stringer:
		return strings.TrimSpace(typed.String())
	default:
		return strings.TrimSpace(fmt.Sprint(typed))
	}
}

func pruneRecent(entries []time.Time, cutoff time.Time) []time.Time {
	out := entries[:0]
	for _, entry := range entries {
		if entry.After(cutoff) {
			out = append(out, entry)
		}
	}
	return out
}

func withinAbsoluteWindow(limits limitConfig, now time.Time) bool {
	if strings.TrimSpace(limits.NotBefore) != "" {
		if ts, err := time.Parse(time.RFC3339, strings.TrimSpace(limits.NotBefore)); err == nil && now.Before(ts) {
			return false
		}
	}
	if strings.TrimSpace(limits.NotAfter) != "" {
		if ts, err := time.Parse(time.RFC3339, strings.TrimSpace(limits.NotAfter)); err == nil && now.After(ts) {
			return false
		}
	}
	return true
}

func withinSchedules(schedules []scheduleConfig, now time.Time) bool {
	if len(schedules) == 0 {
		return true
	}
	for _, schedule := range schedules {
		if matchesSchedule(schedule, now) {
			return true
		}
	}
	return false
}

func matchesSchedule(schedule scheduleConfig, now time.Time) bool {
	if len(schedule.Days) > 0 {
		matchDay := false
		for _, day := range schedule.Days {
			if strings.EqualFold(strings.TrimSpace(day), now.Weekday().String()) {
				matchDay = true
				break
			}
		}
		if !matchDay {
			return false
		}
	}
	if schedule.Start == "" || schedule.End == "" {
		return true
	}
	return withinClockWindow(now.Format("15:04"), schedule.Start, schedule.End)
}

func withinClockWindow(current, start, end string) bool {
	current = strings.TrimSpace(current)
	start = strings.TrimSpace(start)
	end = strings.TrimSpace(end)
	if start == "" || end == "" {
		return true
	}
	if start <= end {
		return current >= start && current <= end
	}
	return current >= start || current <= end
}

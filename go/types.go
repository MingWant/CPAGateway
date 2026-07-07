package main

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/router-for-me/CLIProxyAPI/v7/sdk/pluginapi"
)

type envelope struct {
	OK     bool            `json:"ok"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *envelopeError  `json:"error,omitempty"`
}

type envelopeError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type lifecycleRequest struct {
	ConfigYAML []byte `json:"config_yaml"`
}

type registration struct {
	SchemaVersion uint32                 `json:"schema_version"`
	Metadata      pluginapi.Metadata     `json:"metadata"`
	Capabilities  registrationCapability `json:"capabilities"`
}

type registrationCapability struct {
	RequestInterceptor bool `json:"request_interceptor"`
	UsagePlugin        bool `json:"usage_plugin"`
	ManagementAPI      bool `json:"management_api"`
}

type pluginConfig struct {
	Enabled     bool              `yaml:"enabled" json:"enabled"`
	Priority    int               `yaml:"priority" json:"priority"`
	Persistence persistenceConfig `yaml:"persistence" json:"persistence,omitempty"`
	Security    securityConfig    `yaml:"security" json:"security,omitempty"`
	Cluster     clusterConfig     `yaml:"cluster" json:"cluster,omitempty"`
	Default     policyConfig      `yaml:"default_policy" json:"default_policy"`
	KeyPolicies []keyPolicyConfig `yaml:"key_policies" json:"key_policies"`
}

type persistenceConfig struct {
	StatePath      string `yaml:"state_path" json:"state_path,omitempty"`
	PersistRuntime bool   `yaml:"persist_runtime" json:"persist_runtime,omitempty"`
}

type securityConfig struct {
	RequireManagementToken bool     `yaml:"require_management_token" json:"require_management_token,omitempty"`
	AdminTokens            []string `yaml:"admin_tokens" json:"admin_tokens,omitempty"`
	ReadTokens             []string `yaml:"read_tokens" json:"read_tokens,omitempty"`
	UIAccessTokens         []string `yaml:"ui_access_tokens" json:"ui_access_tokens,omitempty"`
}

type clusterConfig struct {
	Backend string      `yaml:"backend" json:"backend,omitempty"`
	Redis   redisConfig `yaml:"redis" json:"redis,omitempty"`
}

type redisConfig struct {
	Addr        string `yaml:"addr" json:"addr,omitempty"`
	Username    string `yaml:"username" json:"username,omitempty"`
	Password    string `yaml:"password" json:"password,omitempty"`
	DB          int    `yaml:"db" json:"db,omitempty"`
	KeyPrefix   string `yaml:"key_prefix" json:"key_prefix,omitempty"`
	FailureMode string `yaml:"failure_mode" json:"failure_mode,omitempty"`
}

type policyConfig struct {
	Enabled     bool                   `yaml:"enabled" json:"enabled"`
	Limits      limitConfig            `yaml:"limits" json:"limits"`
	StagePolicy map[string]stagePolicy `yaml:"stage_policy" json:"stage_policy,omitempty"`
	Rules       []ruleConfig           `yaml:"rules" json:"rules"`
}

type stagePolicy struct {
	Mode string `yaml:"mode" json:"mode,omitempty"`
}

type keyPolicyConfig struct {
	KeyID       string                 `yaml:"key_id" json:"key_id"`
	MatchAPIKey string                 `yaml:"match_api_key" json:"match_api_key,omitempty"`
	MaskedKey   string                 `yaml:"-" json:"masked_key,omitempty"`
	DisplayName string                 `yaml:"display_name" json:"display_name"`
	Enabled     bool                   `yaml:"enabled" json:"enabled"`
	Limits      limitConfig            `yaml:"limits" json:"limits"`
	StagePolicy map[string]stagePolicy `yaml:"stage_policy" json:"stage_policy,omitempty"`
	Rules       []ruleConfig           `yaml:"rules" json:"rules"`
}

type limitConfig struct {
	RequestsPerDay int              `yaml:"requests_per_day" json:"requests_per_day"`
	TokensPerDay   int              `yaml:"tokens_per_day" json:"tokens_per_day"`
	RequestsPerMin int              `yaml:"requests_per_minute" json:"requests_per_minute"`
	MaxInflight    int              `yaml:"max_inflight" json:"max_inflight"`
	NotBefore      string           `yaml:"not_before" json:"not_before,omitempty"`
	NotAfter       string           `yaml:"not_after" json:"not_after,omitempty"`
	Schedules      []scheduleConfig `yaml:"allowed_time_ranges" json:"allowed_time_ranges,omitempty"`
}

type scheduleConfig struct {
	Days  []string `yaml:"days" json:"days,omitempty"`
	Start string   `yaml:"start" json:"start,omitempty"`
	End   string   `yaml:"end" json:"end,omitempty"`
}

type ruleConfig struct {
	ID       string       `yaml:"id" json:"id"`
	Enabled  bool         `yaml:"enabled" json:"enabled"`
	Priority int          `yaml:"priority" json:"priority"`
	Stage    string       `yaml:"stage" json:"stage,omitempty"`
	OnMatch  string       `yaml:"on_match" json:"on_match"`
	Match    matchConfig  `yaml:"match" json:"match"`
	Actions  actionConfig `yaml:"actions" json:"actions"`
}

type matchConfig struct {
	PathKinds        []string          `yaml:"path_kinds" json:"path_kinds,omitempty"`
	Paths            []string          `yaml:"paths" json:"paths,omitempty"`
	Models           []string          `yaml:"models" json:"models,omitempty"`
	ModelPrefixes    []string          `yaml:"model_prefixes" json:"model_prefixes,omitempty"`
	Providers        []string          `yaml:"providers" json:"providers,omitempty"`
	AnyOf            []matchConfig     `yaml:"any_of" json:"any_of,omitempty"`
	AllOf            []matchConfig     `yaml:"all_of" json:"all_of,omitempty"`
	Stream           *bool             `yaml:"stream" json:"stream,omitempty"`
	Headers          map[string]string `yaml:"headers" json:"headers,omitempty"`
	Query            map[string]string `yaml:"query" json:"query,omitempty"`
	BodyContains     map[string]string `yaml:"body_contains" json:"body_contains,omitempty"`
	MetadataContains map[string]string `yaml:"metadata_contains" json:"metadata_contains,omitempty"`
	Days             []string          `yaml:"days" json:"days,omitempty"`
	Start            string            `yaml:"start" json:"start,omitempty"`
	End              string            `yaml:"end" json:"end,omitempty"`
}

type actionConfig struct {
	RewriteModel        string            `yaml:"rewrite_model" json:"rewrite_model,omitempty"`
	RouteToModel        string            `yaml:"route_to_model" json:"route_to_model,omitempty"`
	WeightedRoutes      []weightedRoute   `yaml:"weighted_routes" json:"weighted_routes,omitempty"`
	RoutePool           *routePoolConfig  `yaml:"route_pool" json:"route_pool,omitempty"`
	FailoverChain       []string          `yaml:"failover_chain" json:"failover_chain,omitempty"`
	FailoverHops        []failoverHop     `yaml:"failover_hops" json:"failover_hops,omitempty"`
	ShardBy             string            `yaml:"shard_by" json:"shard_by,omitempty"`
	MirrorModels        []string          `yaml:"mirror_models" json:"mirror_models,omitempty"`
	ForceProviderPrefix string            `yaml:"force_provider_prefix" json:"force_provider_prefix,omitempty"`
	AllowOnlyProviders  []string          `yaml:"allow_only_providers" json:"allow_only_providers,omitempty"`
	AllowOnlyModels     []string          `yaml:"allow_only_models" json:"allow_only_models,omitempty"`
	FallbackModels      []string          `yaml:"fallback_models" json:"fallback_models,omitempty"`
	Deny                *denyConfig       `yaml:"deny" json:"deny,omitempty"`
	RewriteEndpoint     *endpointRewrite  `yaml:"rewrite_endpoint_semantics" json:"rewrite_endpoint_semantics,omitempty"`
	SetHeaders          map[string]string `yaml:"set_headers" json:"set_headers,omitempty"`
	ClearHeaders        []string          `yaml:"clear_headers" json:"clear_headers,omitempty"`
	TagMetadata         map[string]string `yaml:"tag_metadata" json:"tag_metadata,omitempty"`
}

type routePoolConfig struct {
	Name             string          `yaml:"name" json:"name,omitempty"`
	Mode             string          `yaml:"mode" json:"mode,omitempty"`
	ProviderAffinity string          `yaml:"provider_affinity" json:"provider_affinity,omitempty"`
	Members          []weightedRoute `yaml:"members" json:"members,omitempty"`
}

type weightedRoute struct {
	Model      string `yaml:"model" json:"model,omitempty"`
	Provider   string `yaml:"provider" json:"provider,omitempty"`
	Suffix     string `yaml:"suffix" json:"suffix,omitempty"`
	Weight     int    `yaml:"weight" json:"weight"`
	Priority   int    `yaml:"priority" json:"priority,omitempty"`
	Enabled    *bool  `yaml:"enabled" json:"enabled,omitempty"`
	Status     string `yaml:"status" json:"status,omitempty"`
	Reason     string `yaml:"reason" json:"reason,omitempty"`
	Health     int    `yaml:"health" json:"health,omitempty"`
	TrafficCap int    `yaml:"traffic_cap" json:"traffic_cap,omitempty"`
}

type failoverHop struct {
	Model      string `yaml:"model" json:"model,omitempty"`
	Provider   string `yaml:"provider" json:"provider,omitempty"`
	Suffix     string `yaml:"suffix" json:"suffix,omitempty"`
	Reason     string `yaml:"reason" json:"reason,omitempty"`
	OnDecision string `yaml:"on_decision" json:"on_decision,omitempty"`
	Enabled    *bool  `yaml:"enabled" json:"enabled,omitempty"`
}

type denyConfig struct {
	StatusCode int    `yaml:"status_code" json:"status_code"`
	Message    string `yaml:"message" json:"message"`
	Code       string `yaml:"code" json:"code,omitempty"`
}

type endpointRewrite struct {
	Mode          string `yaml:"mode" json:"mode,omitempty"`
	TargetModel   string `yaml:"target_model" json:"target_model,omitempty"`
	Reason        string `yaml:"reason" json:"reason,omitempty"`
	ResponseError string `yaml:"response_error" json:"response_error,omitempty"`
}

type storedPolicy struct {
	Version       int               `json:"version"`
	DefaultPolicy policyConfig      `json:"default_policy"`
	KeyPolicies   []keyPolicyConfig `json:"key_policies"`
}

type policyBundle struct {
	Version       int               `json:"version"`
	Name          string            `json:"name"`
	Description   string            `json:"description,omitempty"`
	Tags          []string          `json:"tags,omitempty"`
	DefaultPolicy policyConfig      `json:"default_policy"`
	KeyPolicies   []keyPolicyConfig `json:"key_policies"`
	ExportedAt    time.Time         `json:"exported_at,omitempty"`
}

type persistedState struct {
	Version       int                      `json:"version"`
	DefaultPolicy policyConfig             `json:"default_policy"`
	KeyPolicies   []keyPolicyConfig        `json:"key_policies"`
	Templates     []ruleTemplate           `json:"templates,omitempty"`
	Usage         map[string]*usageCounter `json:"usage,omitempty"`
	RequestWindow map[string][]time.Time   `json:"request_window,omitempty"`
	AuditLog      []auditEntry             `json:"audit_log,omitempty"`
	MemberHits    map[string]int           `json:"member_hits,omitempty"`
	RuleHits      map[string]int           `json:"rule_hits,omitempty"`
	StageHits     map[string]int           `json:"stage_hits,omitempty"`
	MemberTimes   map[string][]time.Time   `json:"member_hit_times,omitempty"`
	RuleTimes     map[string][]time.Time   `json:"rule_hit_times,omitempty"`
	StageTimes    map[string][]time.Time   `json:"stage_hit_times,omitempty"`
	SavedAt       time.Time                `json:"saved_at,omitempty"`
}

type usageEntry struct {
	KeyID          string    `json:"key_id"`
	DisplayName    string    `json:"display_name"`
	MaskedKey      string    `json:"masked_key"`
	RequestsToday  int       `json:"requests_today"`
	TokensToday    int64     `json:"tokens_today"`
	RequestsMinute int       `json:"requests_minute"`
	Inflight       int       `json:"inflight"`
	LastSeenAt     time.Time `json:"last_seen_at,omitempty"`
}

type dryRunResult struct {
	Decision     string                             `json:"decision"`
	RuleID       string                             `json:"rule_id,omitempty"`
	Reason       string                             `json:"reason,omitempty"`
	MatchedRules []string                           `json:"matched_rules,omitempty"`
	FinalModel   string                             `json:"final_model,omitempty"`
	Response     pluginapi.RequestInterceptResponse `json:"response"`
	StageTrace   []stageTrace                       `json:"stage_trace,omitempty"`
}

type stageRunResult struct {
	Decision     string                             `json:"decision"`
	RuleID       string                             `json:"rule_id,omitempty"`
	Reason       string                             `json:"reason,omitempty"`
	MatchedRules []string                           `json:"matched_rules,omitempty"`
	FinalModel   string                             `json:"final_model,omitempty"`
	Response     pluginapi.RequestInterceptResponse `json:"response"`
	StageTrace   []stageTrace                       `json:"stage_trace,omitempty"`
}

type stageTrace struct {
	Stage           string   `json:"stage"`
	Mode            string   `json:"mode,omitempty"`
	MatchedRules    []string `json:"matched_rules,omitempty"`
	FinalModel      string   `json:"final_model,omitempty"`
	Decision        string   `json:"decision,omitempty"`
	Reason          string   `json:"reason,omitempty"`
	MatchedCount    int      `json:"matched_count,omitempty"`
	RoutePool       string   `json:"route_pool,omitempty"`
	RouteTarget     string   `json:"route_target,omitempty"`
	FallbackTarget  string   `json:"fallback_target,omitempty"`
	MirrorModels    []string `json:"mirror_models,omitempty"`
	FailoverChain   []string `json:"failover_chain,omitempty"`
	FailoverReasons []string `json:"failover_reasons,omitempty"`
}

type previewTokenRecord struct {
	KeyID       string
	RuleID      string
	Operation   string
	Target      string
	Secondary   string
	BeforeState string
	AfterState  string
	Token       string
	IssuedAt    time.Time
}

type pluginState struct {
	mu              sync.RWMutex
	config          pluginConfig
	redisCounters   *redisCounterStore
	usage           map[string]*usageCounter
	requestWindow   map[string][]time.Time
	auditLog        []auditEntry
	templates       []ruleTemplate
	memberHitCounts map[string]int
	ruleHitCounts   map[string]int
	stageHitCounts  map[string]int
	memberHitTimes  map[string][]time.Time
	ruleHitTimes    map[string][]time.Time
	stageHitTimes   map[string][]time.Time
	previewTokens   map[string]previewTokenRecord
}

type usageCounter struct {
	DisplayName    string
	MaskedKey      string
	RequestsToday  int
	TokensToday    int64
	RequestsMinute int
	Inflight       int
	LastSeenAt     time.Time
	DayBucket      string
}

type auditEntry struct {
	Time           time.Time             `json:"time"`
	PolicyID       string                `json:"policy_id,omitempty"`
	PolicyName     string                `json:"policy_name,omitempty"`
	Decision       string                `json:"decision"`
	RuleID         string                `json:"rule_id,omitempty"`
	Reason         string                `json:"reason,omitempty"`
	RequestedModel string                `json:"requested_model,omitempty"`
	FinalModel     string                `json:"final_model,omitempty"`
	Mirrors        []string              `json:"mirrors,omitempty"`
	Path           string                `json:"path,omitempty"`
	APIKey         string                `json:"api_key,omitempty"`
	Provider       string                `json:"provider,omitempty"`
	EventType      string                `json:"event_type,omitempty"`
	OperatorAction string                `json:"operator_action,omitempty"`
	TargetMember   string                `json:"target_member,omitempty"`
	Secondary      string                `json:"secondary,omitempty"`
	BeforeState    string                `json:"before_state,omitempty"`
	AfterState     string                `json:"after_state,omitempty"`
	Diff           []memberChangeSummary `json:"diff,omitempty"`
}

type memberChangeSummary struct {
	Member string        `json:"member"`
	Before weightedRoute `json:"before"`
	After  weightedRoute `json:"after"`
}

type auditSummary struct {
	TotalByDecision map[string]int `json:"total_by_decision"`
	TotalByReason   map[string]int `json:"total_by_reason"`
	TotalByRule     map[string]int `json:"total_by_rule"`
	TotalByPolicy   map[string]int `json:"total_by_policy"`
	TotalByModel    map[string]int `json:"total_by_model"`
	TotalByProvider map[string]int `json:"total_by_provider"`
	Timeline        []auditBucket  `json:"timeline,omitempty"`
}

type auditBucket struct {
	Window string `json:"window"`
	Count  int    `json:"count"`
}

type memberOperationRequest struct {
	KeyID         string `json:"key_id"`
	RuleID        string `json:"rule_id"`
	Member        string `json:"member,omitempty"`
	MemberType    string `json:"member_type,omitempty"`
	Operation     string `json:"operation"`
	Delta         int    `json:"delta,omitempty"`
	Reason        string `json:"reason,omitempty"`
	Secondary     string `json:"secondary,omitempty"`
	PrimaryWeight int    `json:"primary_weight,omitempty"`
	CanaryWeight  int    `json:"canary_weight,omitempty"`
	PreviewOnly   bool   `json:"preview_only,omitempty"`
	PreviewToken  string `json:"preview_token,omitempty"`
}

type memberOperationResult struct {
	OK           bool                  `json:"ok"`
	KeyID        string                `json:"key_id,omitempty"`
	RuleID       string                `json:"rule_id,omitempty"`
	Operation    string                `json:"operation,omitempty"`
	Reason       string                `json:"reason,omitempty"`
	TargetMember string                `json:"target_member,omitempty"`
	Secondary    string                `json:"secondary,omitempty"`
	BeforeState  string                `json:"before_state,omitempty"`
	AfterState   string                `json:"after_state,omitempty"`
	Diff         []memberChangeSummary `json:"diff,omitempty"`
	Members      []weightedRoute       `json:"members,omitempty"`
}

type memberOperationPreview struct {
	memberOperationResult
	PreviewToken string `json:"preview_token,omitempty"`
}

type ruleTemplate struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Category    string     `json:"category,omitempty"`
	Description string     `json:"description"`
	Scenario    string     `json:"scenario,omitempty"`
	Maturity    string     `json:"maturity,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	Rule        ruleConfig `json:"rule"`
}

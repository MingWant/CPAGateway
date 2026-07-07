package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/router-for-me/CLIProxyAPI/v7/sdk/pluginapi"
)

const redisHealthTimeout = 500 * time.Millisecond

type gatewayHealth struct {
	Status      string                   `json:"status"`
	Version     string                   `json:"version"`
	GeneratedAt time.Time                `json:"generated_at"`
	Config      gatewayConfigHealth      `json:"config"`
	Counts      gatewayCountsHealth      `json:"counts"`
	Persistence gatewayPersistenceHealth `json:"persistence"`
	Security    gatewaySecurityHealth    `json:"security"`
	Counters    gatewayCountersHealth    `json:"counters"`
	Warnings    []string                 `json:"warnings,omitempty"`
}

type gatewayConfigHealth struct {
	Enabled  bool `json:"enabled"`
	Priority int  `json:"priority"`
}

type gatewayCountsHealth struct {
	KeyPolicies  int `json:"key_policies"`
	Rules        int `json:"rules"`
	Templates    int `json:"templates"`
	AuditEntries int `json:"audit_entries"`
	UsageEntries int `json:"usage_entries"`
}

type gatewayPersistenceHealth struct {
	Configured            bool   `json:"configured"`
	StatePathSet          bool   `json:"state_path_set"`
	PersistRuntime        bool   `json:"persist_runtime"`
	RuntimeUsagePersisted bool   `json:"runtime_usage_persisted"`
	Storage               string `json:"storage"`
}

type gatewaySecurityHealth struct {
	ManagementAuthEnabled bool `json:"management_auth_enabled"`
	UIAuthEnabled         bool `json:"ui_auth_enabled"`
	RequireToken          bool `json:"require_token"`
	AdminTokenCount       int  `json:"admin_token_count"`
	ReadTokenCount        int  `json:"read_token_count"`
	UITokenCount          int  `json:"ui_token_count"`
}

type gatewayCountersHealth struct {
	Backend         string `json:"backend"`
	RedisConfigured bool   `json:"redis_configured"`
	RedisRequired   bool   `json:"redis_required"`
	RedisStatus     string `json:"redis_status"`
	FailureMode     string `json:"failure_mode,omitempty"`
	KeyPrefixSet    bool   `json:"key_prefix_set,omitempty"`
}

func routeHealth(_ pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	return jsonResponse(http.StatusOK, gatewayState.healthSnapshot())
}

func (s *pluginState) healthSnapshot() gatewayHealth {
	s.mu.RLock()
	cfg := s.config
	redisCounters := s.redisCounters
	counts := gatewayCountsHealth{
		KeyPolicies:  len(cfg.KeyPolicies),
		Rules:        countConfiguredRules(cfg),
		Templates:    len(s.templates),
		AuditEntries: len(s.auditLog),
		UsageEntries: len(s.usage),
	}
	s.mu.RUnlock()

	clusterBackend := strings.ToLower(strings.TrimSpace(cfg.Cluster.Backend))
	backend := "local"
	redisStatus := "not_configured"
	redisRequired := clusterBackend == "redis"
	warnings := make([]string, 0, 4)
	if clusterBackend != "" && !redisRequired {
		warnings = append(warnings, "unsupported cluster backend; local counters are active")
	}
	if redisRequired {
		backend = "redis"
		redisStatus = "not_initialized"
		if redisCounters != nil {
			redisStatus = "ok"
			if err := redisCounters.ping(redisHealthTimeout); err != nil {
				redisStatus = "unavailable"
				warnings = append(warnings, "redis counter backend is unavailable")
			}
		} else {
			warnings = append(warnings, "redis backend is configured but the counter store is not initialized")
		}
	}

	failureMode := strings.TrimSpace(cfg.Cluster.Redis.FailureMode)
	if redisRequired && failureMode == "" {
		failureMode = "reject"
	}
	if redisRequired && failureMode == "reject" && redisStatus != "ok" {
		warnings = append(warnings, "redis failure_mode is reject; requests can fail closed while redis is unavailable")
	}
	if strings.TrimSpace(cfg.Persistence.StatePath) == "" {
		warnings = append(warnings, "state_path is not configured; policy changes are memory-only")
	}
	if !managementAPIAuthConfigured(cfg.Security) {
		warnings = append(warnings, "plugin-local management tokens are not configured")
	}

	status := "ok"
	if redisRequired && redisStatus != "ok" {
		status = "degraded"
	}
	if clusterBackend != "" && !redisRequired {
		status = "degraded"
	}

	return gatewayHealth{
		Status:      status,
		Version:     gatewayPluginVersion,
		GeneratedAt: time.Now().UTC(),
		Config:      gatewayConfigHealth{Enabled: cfg.Enabled, Priority: cfg.Priority},
		Counts:      counts,
		Persistence: gatewayPersistenceHealth{
			Configured:            strings.TrimSpace(cfg.Persistence.StatePath) != "",
			StatePathSet:          strings.TrimSpace(cfg.Persistence.StatePath) != "",
			PersistRuntime:        cfg.Persistence.PersistRuntime,
			RuntimeUsagePersisted: cfg.Persistence.PersistRuntime && !redisRequired,
			Storage:               persistenceStorageName(cfg.Persistence),
		},
		Security: gatewaySecurityHealth{
			ManagementAuthEnabled: managementAPIAuthConfigured(cfg.Security),
			UIAuthEnabled:         uiAuthConfigured(cfg.Security),
			RequireToken:          cfg.Security.RequireManagementToken,
			AdminTokenCount:       countNonEmpty(cfg.Security.AdminTokens),
			ReadTokenCount:        countNonEmpty(cfg.Security.ReadTokens),
			UITokenCount:          countNonEmpty(cfg.Security.UIAccessTokens),
		},
		Counters: gatewayCountersHealth{
			Backend:         backend,
			RedisConfigured: redisRequired,
			RedisRequired:   redisRequired,
			RedisStatus:     redisStatus,
			FailureMode:     failureMode,
			KeyPrefixSet:    strings.TrimSpace(cfg.Cluster.Redis.KeyPrefix) != "",
		},
		Warnings: warnings,
	}
}

func countConfiguredRules(cfg pluginConfig) int {
	total := len(cfg.Default.Rules)
	for _, policy := range cfg.KeyPolicies {
		total += len(policy.Rules)
	}
	return total
}

func managementAPIAuthConfigured(cfg securityConfig) bool {
	return cfg.RequireManagementToken || countNonEmpty(cfg.AdminTokens) > 0 || countNonEmpty(cfg.ReadTokens) > 0
}

func uiAuthConfigured(cfg securityConfig) bool {
	return countNonEmpty(cfg.UIAccessTokens) > 0 || managementAPIAuthConfigured(cfg)
}

func countNonEmpty(items []string) int {
	count := 0
	for _, item := range items {
		if strings.TrimSpace(item) != "" {
			count++
		}
	}
	return count
}

func persistenceStorageName(cfg persistenceConfig) string {
	if strings.TrimSpace(cfg.StatePath) == "" {
		return "memory"
	}
	return "json_file"
}

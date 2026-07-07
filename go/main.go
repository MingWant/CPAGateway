package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/router-for-me/CLIProxyAPI/v7/sdk/pluginabi"
	"github.com/router-for-me/CLIProxyAPI/v7/sdk/pluginapi"
	"gopkg.in/yaml.v3"
)

var gatewayState = newPluginState()

const gatewayPluginVersion = "0.1.0"

func main() {}

func handleMethod(method string, request []byte) ([]byte, error) {
	switch method {
	case pluginabi.MethodPluginRegister, pluginabi.MethodPluginReconfigure:
		if err := configure(request); err != nil {
			return nil, err
		}
		return okEnvelope(pluginRegistration())
	case pluginabi.MethodRequestInterceptBefore:
		return interceptBefore(request)
	case pluginabi.MethodRequestInterceptAfter:
		return interceptAfter(request)
	case pluginabi.MethodUsageHandle:
		return handleUsage(request)
	case pluginabi.MethodManagementRegister:
		return okEnvelope(managementRegistration())
	case pluginabi.MethodManagementHandle:
		return handleManagement(request)
	default:
		return errorEnvelope("unknown_method", "unknown method: "+method), nil
	}
}

func configure(raw []byte) error {
	var req lifecycleRequest
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &req); err != nil {
			return err
		}
	}
	cfg := pluginConfig{}
	if len(req.ConfigYAML) > 0 {
		if err := yaml.Unmarshal(req.ConfigYAML, &cfg); err != nil {
			return err
		}
	}
	cfg = normalizeConfig(cfg)
	redisCounters := newRedisCounterStore(cfg.Cluster)
	gatewayState.mu.Lock()
	oldRedisCounters := gatewayState.redisCounters
	gatewayState.config = cfg
	gatewayState.redisCounters = redisCounters
	gatewayState.ensureRuntimeMapsLocked()
	gatewayState.mu.Unlock()
	if oldRedisCounters != nil && oldRedisCounters != redisCounters {
		oldRedisCounters.close()
	}
	return gatewayState.loadPersistentState()
}

func pluginRegistration() registration {
	return registration{
		SchemaVersion: pluginabi.SchemaVersion,
		Metadata: pluginapi.Metadata{
			Name:             "gateway",
			Version:          gatewayPluginVersion,
			Author:           "router-for-me",
			GitHubRepository: "https://github.com/router-for-me/CLIProxyAPI",
			Logo:             "https://raw.githubusercontent.com/router-for-me/CLIProxyAPI/main/docs/logo.png",
			ConfigFields: []pluginapi.ConfigField{{
				Name:        "default_policy",
				Type:        pluginapi.ConfigFieldTypeObject,
				Description: "Default policy applied when no API key specific policy matches.",
			}, {
				Name:        "key_policies",
				Type:        pluginapi.ConfigFieldTypeArray,
				Description: "Per API key gateway policies bound to top-level CPA api-keys.",
			}, {
				Name:        "persistence",
				Type:        pluginapi.ConfigFieldTypeObject,
				Description: "Optional plugin state persistence settings, including state_path and persist_runtime.",
			}, {
				Name:        "cluster",
				Type:        pluginapi.ConfigFieldTypeObject,
				Description: "Optional shared counter backend settings. Set backend to redis for multi-instance quota accounting.",
			}, {
				Name:        "security",
				Type:        pluginapi.ConfigFieldTypeObject,
				Description: "Optional plugin-level read/admin tokens layered on top of host management authentication.",
			}},
		},
		Capabilities: registrationCapability{RequestInterceptor: true, UsagePlugin: true, ManagementAPI: true},
	}
}

func managementRegistration() pluginapi.ManagementRegistrationResponse {
	return pluginapi.ManagementRegistrationResponse{
		Routes: []pluginapi.ManagementRoute{
			{Method: http.MethodGet, Path: "/plugins/gateway/keys", Handler: managementHandlerFunc(authorizedHandler(roleRead, routeKeys))},
			{Method: http.MethodGet, Path: "/plugins/gateway/health", Handler: managementHandlerFunc(authorizedHandler(roleRead, routeHealth))},
			{Method: http.MethodGet, Path: "/plugins/gateway/policies", Handler: managementHandlerFunc(authorizedHandler(roleRead, routePolicies))},
			{Method: http.MethodPut, Path: "/plugins/gateway/policies", Handler: managementHandlerFunc(authorizedHandler(roleAdmin, routePutPolicies))},
			{Method: http.MethodGet, Path: "/plugins/gateway/policies/export", Handler: managementHandlerFunc(authorizedHandler(roleRead, routeExportPolicies))},
			{Method: http.MethodPost, Path: "/plugins/gateway/policies/import", Handler: managementHandlerFunc(authorizedHandler(roleAdmin, routeImportPolicies))},
			{Method: http.MethodPost, Path: "/plugins/gateway/policies/clone", Handler: managementHandlerFunc(authorizedHandler(roleAdmin, routeClonePolicy))},
			{Method: http.MethodPost, Path: "/plugins/gateway/policies/add", Handler: managementHandlerFunc(authorizedHandler(roleAdmin, routeAddPolicy))},
			{Method: http.MethodPatch, Path: "/plugins/gateway/policies", Handler: managementHandlerFunc(authorizedHandler(roleAdmin, routePatchPolicy))},
			{Method: http.MethodDelete, Path: "/plugins/gateway/policies", Handler: managementHandlerFunc(authorizedHandler(roleAdmin, routeDeletePolicy))},
			{Method: http.MethodPost, Path: "/plugins/gateway/rules/add", Handler: managementHandlerFunc(authorizedHandler(roleAdmin, routeAddRule))},
			{Method: http.MethodPatch, Path: "/plugins/gateway/rules", Handler: managementHandlerFunc(authorizedHandler(roleAdmin, routePatchRule))},
			{Method: http.MethodDelete, Path: "/plugins/gateway/rules", Handler: managementHandlerFunc(authorizedHandler(roleAdmin, routeDeleteRule))},
			{Method: http.MethodPost, Path: "/plugins/gateway/route-members/op", Handler: managementHandlerFunc(authorizedHandler(roleAdmin, routeMemberOperation))},
			{Method: http.MethodPost, Path: "/plugins/gateway/route-members/preview", Handler: managementHandlerFunc(authorizedHandler(roleRead, routeMemberPreview))},
			{Method: http.MethodGet, Path: "/plugins/gateway/usage", Handler: managementHandlerFunc(authorizedHandler(roleRead, routeUsage))},
			{Method: http.MethodGet, Path: "/plugins/gateway/audit", Handler: managementHandlerFunc(authorizedHandler(roleRead, routeAudit))},
			{Method: http.MethodGet, Path: "/plugins/gateway/audit/detail", Handler: managementHandlerFunc(authorizedHandler(roleRead, routeAuditDetail))},
			{Method: http.MethodGet, Path: "/plugins/gateway/audit/summary", Handler: managementHandlerFunc(authorizedHandler(roleRead, routeAuditSummary))},
			{Method: http.MethodGet, Path: "/plugins/gateway/templates", Handler: managementHandlerFunc(authorizedHandler(roleRead, routeTemplates))},
			{Method: http.MethodPost, Path: "/plugins/gateway/templates", Handler: managementHandlerFunc(authorizedHandler(roleAdmin, routeAddTemplate))},
			{Method: http.MethodPatch, Path: "/plugins/gateway/templates", Handler: managementHandlerFunc(authorizedHandler(roleAdmin, routePatchTemplate))},
			{Method: http.MethodDelete, Path: "/plugins/gateway/templates", Handler: managementHandlerFunc(authorizedHandler(roleAdmin, routeDeleteTemplate))},
			{Method: http.MethodPost, Path: "/plugins/gateway/templates/clone", Handler: managementHandlerFunc(authorizedHandler(roleAdmin, routeCloneTemplate))},
			{Method: http.MethodGet, Path: "/plugins/gateway/templates/export", Handler: managementHandlerFunc(authorizedHandler(roleRead, routeExportTemplates))},
			{Method: http.MethodPost, Path: "/plugins/gateway/templates/import", Handler: managementHandlerFunc(authorizedHandler(roleAdmin, routeImportTemplates))},
			{Method: http.MethodPost, Path: "/plugins/gateway/usage/reset", Handler: managementHandlerFunc(authorizedHandler(roleAdmin, routeResetUsage))},
			{Method: http.MethodPost, Path: "/plugins/gateway/dry-run", Handler: managementHandlerFunc(authorizedHandler(roleRead, routeDryRun))},
		},
		Resources: []pluginapi.ResourceRoute{{
			Path:        "/ui",
			Menu:        "Gateway",
			Description: "Manage gateway policies, limits, dry-run, and usage.",
			Handler:     managementHandlerFunc(routeUI),
		}},
	}
}

func interceptBefore(raw []byte) ([]byte, error) {
	var req pluginapi.RequestInterceptRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, err
	}
	resp := gatewayState.apply(req, false)
	return okEnvelope(resp)
}

func interceptAfter(raw []byte) ([]byte, error) {
	var req pluginapi.RequestInterceptRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, err
	}
	resp := gatewayState.apply(req, true)
	if resp.Reject {
		gatewayState.releaseInflightForRequest(req)
	}
	return okEnvelope(resp)
}

func handleUsage(raw []byte) ([]byte, error) {
	var record pluginapi.UsageRecord
	if err := json.Unmarshal(raw, &record); err != nil {
		return nil, err
	}
	gatewayState.recordUsage(record)
	return okEnvelope(struct{}{})
}

func handleManagement(raw []byte) ([]byte, error) {
	var req pluginapi.ManagementRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, err
	}
	for _, route := range managementRegistration().Routes {
		path := route.Path
		if !strings.HasPrefix(path, "/v0/management") {
			path = "/v0/management" + path
		}
		if strings.EqualFold(route.Method, req.Method) && path == req.Path {
			resp, err := route.Handler.HandleManagement(nil, req)
			if err != nil {
				return nil, err
			}
			return okEnvelope(resp)
		}
	}
	return okEnvelope(pluginapi.ManagementResponse{StatusCode: http.StatusNotFound, Body: []byte(`{"error":"not found"}`), Headers: http.Header{"Content-Type": []string{"application/json"}}})
}

func okEnvelope(v any) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return json.Marshal(envelope{OK: true, Result: raw})
}

func errorEnvelope(code, message string) []byte {
	raw, _ := json.Marshal(envelope{OK: false, Error: &envelopeError{Code: code, Message: message}})
	return raw
}

type managementHandlerFunc func(pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error)

func (f managementHandlerFunc) HandleManagement(_ context.Context, req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
	return f(req)
}

func jsonResponse(status int, payload any) (pluginapi.ManagementResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return pluginapi.ManagementResponse{}, err
	}
	return pluginapi.ManagementResponse{StatusCode: status, Headers: http.Header{"Content-Type": []string{"application/json"}}, Body: body}, nil
}

package main

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/router-for-me/CLIProxyAPI/v7/sdk/pluginapi"
)

type managementRole int

const (
	roleRead managementRole = iota
	roleAdmin
)

func authorizedHandler(role managementRole, next managementHandlerFunc) managementHandlerFunc {
	return func(req pluginapi.ManagementRequest) (pluginapi.ManagementResponse, error) {
		if !gatewayState.managementAuthorized(req, role) {
			return jsonResponse(http.StatusForbidden, map[string]any{"error": "gateway management permission denied"})
		}
		return next(req)
	}
}

func (s *pluginState) managementAuthorized(req pluginapi.ManagementRequest, role managementRole) bool {
	s.mu.RLock()
	cfg := s.config.Security
	s.mu.RUnlock()
	if !cfg.RequireManagementToken && len(cfg.AdminTokens) == 0 && len(cfg.ReadTokens) == 0 {
		return true
	}
	token := managementRequestToken(req)
	if token == "" {
		return false
	}
	if tokenMatchesAny(token, cfg.AdminTokens) {
		return true
	}
	return role == roleRead && tokenMatchesAny(token, cfg.ReadTokens)
}

func (s *pluginState) uiAuthorized(req pluginapi.ManagementRequest) bool {
	s.mu.RLock()
	cfg := s.config.Security
	s.mu.RUnlock()
	if !cfg.RequireManagementToken && len(cfg.AdminTokens) == 0 && len(cfg.ReadTokens) == 0 && len(cfg.UIAccessTokens) == 0 {
		return true
	}
	token := managementRequestToken(req)
	if token == "" {
		return false
	}
	return tokenMatchesAny(token, cfg.UIAccessTokens) || tokenMatchesAny(token, cfg.AdminTokens) || tokenMatchesAny(token, cfg.ReadTokens)
}

func managementRequestToken(req pluginapi.ManagementRequest) string {
	if raw := strings.TrimSpace(req.Headers.Get("Authorization")); raw != "" {
		if token, ok := strings.CutPrefix(raw, "Bearer "); ok {
			return strings.TrimSpace(token)
		}
	}
	for _, header := range []string{"X-Gateway-Admin-Token", "X-Gateway-Token", "X-Plugin-Token"} {
		if token := strings.TrimSpace(req.Headers.Get(header)); token != "" {
			return token
		}
	}
	return strings.TrimSpace(req.Query.Get("gateway_token"))
}

func tokenMatchesAny(token string, candidates []string) bool {
	token = strings.TrimSpace(token)
	if token == "" {
		return false
	}
	for _, candidate := range candidates {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" {
			continue
		}
		if subtle.ConstantTimeCompare([]byte(token), []byte(candidate)) == 1 {
			return true
		}
	}
	return false
}

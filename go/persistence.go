package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (s *pluginState) loadPersistentState() error {
	s.mu.RLock()
	path := strings.TrimSpace(s.config.Persistence.StatePath)
	runtimeEnabled := s.config.Persistence.PersistRuntime
	usingRedis := s.redisCounters != nil
	s.mu.RUnlock()
	if path == "" {
		return nil
	}
	raw, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	var snapshot persistedState
	if err := json.Unmarshal(raw, &snapshot); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if policyHasContent(snapshot.DefaultPolicy) || len(snapshot.KeyPolicies) > 0 {
		s.config.Default = normalizePolicy(snapshot.DefaultPolicy)
		s.config.KeyPolicies = normalizeKeyPolicies(snapshot.KeyPolicies)
	}
	if len(snapshot.Templates) > 0 {
		s.templates = append([]ruleTemplate(nil), snapshot.Templates...)
	}
	if runtimeEnabled {
		if len(snapshot.AuditLog) > 0 {
			s.auditLog = append([]auditEntry(nil), snapshot.AuditLog...)
		}
		if len(snapshot.MemberHits) > 0 {
			s.memberHitCounts = cloneIntMap(snapshot.MemberHits)
		}
		if len(snapshot.RuleHits) > 0 {
			s.ruleHitCounts = cloneIntMap(snapshot.RuleHits)
		}
		if len(snapshot.StageHits) > 0 {
			s.stageHitCounts = cloneIntMap(snapshot.StageHits)
		}
		if len(snapshot.MemberTimes) > 0 {
			s.memberHitTimes = cloneTimeSliceMap(snapshot.MemberTimes)
		}
		if len(snapshot.RuleTimes) > 0 {
			s.ruleHitTimes = cloneTimeSliceMap(snapshot.RuleTimes)
		}
		if len(snapshot.StageTimes) > 0 {
			s.stageHitTimes = cloneTimeSliceMap(snapshot.StageTimes)
		}
		if !usingRedis {
			if len(snapshot.Usage) > 0 {
				s.usage = cloneUsageMap(snapshot.Usage)
			}
			if len(snapshot.RequestWindow) > 0 {
				s.requestWindow = cloneTimeSliceMap(snapshot.RequestWindow)
			}
		}
	}
	s.ensureRuntimeMapsLocked()
	return nil
}

func (s *pluginState) savePersistentStateLocked() error {
	path := strings.TrimSpace(s.config.Persistence.StatePath)
	if path == "" {
		return nil
	}
	snapshot := persistedState{
		Version:       1,
		DefaultPolicy: clonePolicyConfig(s.config.Default),
		KeyPolicies:   cloneKeyPolicyConfigs(s.config.KeyPolicies),
		Templates:     append([]ruleTemplate(nil), s.templates...),
		SavedAt:       time.Now(),
	}
	if s.config.Persistence.PersistRuntime {
		snapshot.AuditLog = append([]auditEntry(nil), s.auditLog...)
		snapshot.MemberHits = cloneIntMap(s.memberHitCounts)
		snapshot.RuleHits = cloneIntMap(s.ruleHitCounts)
		snapshot.StageHits = cloneIntMap(s.stageHitCounts)
		snapshot.MemberTimes = cloneTimeSliceMap(s.memberHitTimes)
		snapshot.RuleTimes = cloneTimeSliceMap(s.ruleHitTimes)
		snapshot.StageTimes = cloneTimeSliceMap(s.stageHitTimes)
		if s.redisCounters == nil {
			snapshot.Usage = cloneUsageMap(s.usage)
			snapshot.RequestWindow = cloneTimeSliceMap(s.requestWindow)
		}
	}
	return writePersistentState(path, snapshot)
}

func writePersistentState(path string, snapshot persistedState) error {
	raw, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".gateway-state-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()
	if _, err := tmp.Write(raw); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpName, 0o600); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

func (s *pluginState) persistRuntimeLocked() error {
	if !s.config.Persistence.PersistRuntime {
		return nil
	}
	return s.savePersistentStateLocked()
}

func (s *pluginState) ensureRuntimeMapsLocked() {
	if s.usage == nil {
		s.usage = make(map[string]*usageCounter)
	}
	if s.requestWindow == nil {
		s.requestWindow = make(map[string][]time.Time)
	}
	if s.auditLog == nil {
		s.auditLog = make([]auditEntry, 0, 128)
	}
	if s.templates == nil {
		s.templates = builtInRuleTemplates()
	}
	if s.memberHitCounts == nil {
		s.memberHitCounts = make(map[string]int)
	}
	if s.ruleHitCounts == nil {
		s.ruleHitCounts = make(map[string]int)
	}
	if s.stageHitCounts == nil {
		s.stageHitCounts = make(map[string]int)
	}
	if s.memberHitTimes == nil {
		s.memberHitTimes = make(map[string][]time.Time)
	}
	if s.ruleHitTimes == nil {
		s.ruleHitTimes = make(map[string][]time.Time)
	}
	if s.stageHitTimes == nil {
		s.stageHitTimes = make(map[string][]time.Time)
	}
	if s.previewTokens == nil {
		s.previewTokens = make(map[string]previewTokenRecord)
	}
}

func cloneKeyPolicyConfigs(items []keyPolicyConfig) []keyPolicyConfig {
	if len(items) == 0 {
		return nil
	}
	out := make([]keyPolicyConfig, len(items))
	for i, item := range items {
		out[i] = cloneKeyPolicyConfig(item)
	}
	return out
}

func cloneUsageMap(src map[string]*usageCounter) map[string]*usageCounter {
	if len(src) == 0 {
		return nil
	}
	out := make(map[string]*usageCounter, len(src))
	for key, value := range src {
		if value == nil {
			continue
		}
		clone := *value
		out[key] = &clone
	}
	return out
}

func cloneIntMap(src map[string]int) map[string]int {
	if len(src) == 0 {
		return nil
	}
	out := make(map[string]int, len(src))
	for key, value := range src {
		out[key] = value
	}
	return out
}

func cloneTimeSliceMap(src map[string][]time.Time) map[string][]time.Time {
	if len(src) == 0 {
		return nil
	}
	out := make(map[string][]time.Time, len(src))
	for key, value := range src {
		out[key] = append([]time.Time(nil), value...)
	}
	return out
}

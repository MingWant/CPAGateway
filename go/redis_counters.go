package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

const redisCounterTimeout = 2 * time.Second

const redisEnforceLua = `
local day = redis.call('HGET', KEYS[1], 'day_bucket')
local requests = tonumber(redis.call('HGET', KEYS[1], 'requests_today') or '0')
local tokens = tonumber(redis.call('HGET', KEYS[1], 'tokens_today') or '0')
local inflight = tonumber(redis.call('HGET', KEYS[1], 'inflight') or '0')
if day ~= ARGV[4] then
  requests = 0
  tokens = 0
end
redis.call('ZREMRANGEBYSCORE', KEYS[2], '-inf', ARGV[7])
local minute = tonumber(redis.call('ZCARD', KEYS[2]) or '0')
local reqPerMin = tonumber(ARGV[8])
local reqPerDay = tonumber(ARGV[9])
local tokensPerDay = tonumber(ARGV[10])
local maxInflight = tonumber(ARGV[11])
if reqPerMin > 0 and minute >= reqPerMin then
  return {'reject','gateway_rate_limit_exceeded','429','gateway rate limit exceeded', requests, tokens, inflight, minute}
end
if reqPerDay > 0 and requests >= reqPerDay then
  return {'reject','gateway_quota_exceeded','403','gateway daily quota exceeded', requests, tokens, inflight, minute}
end
if tokensPerDay > 0 and tokens >= tokensPerDay then
  return {'reject','gateway_token_quota_exceeded','403','gateway token quota exceeded', requests, tokens, inflight, minute}
end
if maxInflight > 0 and inflight >= maxInflight then
  return {'reject','gateway_concurrency_exceeded','429','gateway concurrency limit exceeded', requests, tokens, inflight, minute}
end
if ARGV[12] == '1' then
  redis.call('SADD', KEYS[3], ARGV[1])
  requests = requests + 1
  inflight = inflight + 1
  minute = minute + 1
  redis.call('HSET', KEYS[1], 'display_name', ARGV[2], 'masked_key', ARGV[3], 'day_bucket', ARGV[4], 'requests_today', requests, 'tokens_today', tokens, 'inflight', inflight, 'last_seen_unix', ARGV[5])
  redis.call('ZADD', KEYS[2], ARGV[6], ARGV[13])
end
return {'ok','', '0', '', requests, tokens, inflight, minute}
`

var redisEnforceScript = redis.NewScript(redisEnforceLua)

var redisReleaseScript = redis.NewScript(`
redis.call('SADD', KEYS[3], ARGV[1])
local day = redis.call('HGET', KEYS[1], 'day_bucket')
local requests = tonumber(redis.call('HGET', KEYS[1], 'requests_today') or '0')
local tokens = tonumber(redis.call('HGET', KEYS[1], 'tokens_today') or '0')
local inflight = tonumber(redis.call('HGET', KEYS[1], 'inflight') or '0')
if day ~= ARGV[4] then
  requests = 0
  tokens = 0
end
if inflight > 0 then
  inflight = inflight - 1
end
tokens = tokens + tonumber(ARGV[8])
redis.call('ZREMRANGEBYSCORE', KEYS[2], '-inf', ARGV[7])
local minute = tonumber(redis.call('ZCARD', KEYS[2]) or '0')
redis.call('HSET', KEYS[1], 'display_name', ARGV[2], 'masked_key', ARGV[3], 'day_bucket', ARGV[4], 'requests_today', requests, 'tokens_today', tokens, 'inflight', inflight, 'last_seen_unix', ARGV[5])
return {'ok','', '0', '', requests, tokens, inflight, minute}
`)

type redisCounterStore struct {
	client *redis.Client
	prefix string
	seq    atomic.Int64
}

func newRedisCounterStore(cfg clusterConfig) *redisCounterStore {
	if !strings.EqualFold(strings.TrimSpace(cfg.Backend), "redis") {
		return nil
	}
	addr := strings.TrimSpace(cfg.Redis.Addr)
	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	prefix := strings.TrimSpace(cfg.Redis.KeyPrefix)
	if prefix == "" {
		prefix = "gateway"
	}
	return &redisCounterStore{
		client: redis.NewClient(&redis.Options{
			Addr:                  addr,
			Username:              strings.TrimSpace(cfg.Redis.Username),
			Password:              cfg.Redis.Password,
			DB:                    cfg.Redis.DB,
			DialTimeout:           redisCounterTimeout,
			ReadTimeout:           redisCounterTimeout,
			WriteTimeout:          redisCounterTimeout,
			MaxRetries:            1,
			ContextTimeoutEnabled: true,
		}),
		prefix: prefix,
	}
}

func (s *redisCounterStore) close() {
	if s == nil || s.client == nil {
		return
	}
	_ = s.client.Close()
}

func (s *redisCounterStore) ping(timeout time.Duration) error {
	if s == nil || s.client == nil {
		return fmt.Errorf("redis client not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.client.Ping(ctx).Err()
}

func (s *redisCounterStore) enforce(keyID, displayName, maskedKey string, limits limitConfig, now time.Time, enforce bool) *requestInterceptResponse {
	result, err := s.runEnforceScript(keyID, displayName, maskedKey, limits, now, enforce)
	if err != nil {
		return redisUnavailableResponse(err)
	}
	return redisCounterReject(result)
}

func (s *redisCounterStore) release(keyID, displayName, maskedKey string, now time.Time, tokens int64) error {
	_, err := s.runReleaseScript(keyID, displayName, maskedKey, now, tokens)
	return err
}

func (s *redisCounterStore) reset(keyID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), redisCounterTimeout)
	defer cancel()
	_, err := s.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Del(ctx, s.usageKey(keyID), s.windowKey(keyID))
		pipe.SRem(ctx, s.keysKey(), keyID)
		return nil
	})
	return err
}

func (s *redisCounterStore) listUsage() []usageEntry {
	ctx, cancel := context.WithTimeout(context.Background(), redisCounterTimeout)
	defer cancel()
	keyIDs, err := s.client.SMembers(ctx, s.keysKey()).Result()
	if err != nil {
		return nil
	}
	now := time.Now()
	out := make([]usageEntry, 0, len(keyIDs))
	for _, keyID := range keyIDs {
		keyID = strings.TrimSpace(keyID)
		if keyID == "" {
			continue
		}
		usageKey := s.usageKey(keyID)
		windowKey := s.windowKey(keyID)
		_, _ = s.client.ZRemRangeByScore(ctx, windowKey, "-inf", strconv.FormatInt(now.Add(-time.Minute).UnixMilli(), 10)).Result()
		fields, errFields := s.client.HGetAll(ctx, usageKey).Result()
		if errFields != nil || len(fields) == 0 {
			continue
		}
		minute, _ := s.client.ZCard(ctx, windowKey).Result()
		out = append(out, usageEntry{
			KeyID:          keyID,
			DisplayName:    fields["display_name"],
			MaskedKey:      fields["masked_key"],
			RequestsToday:  atoiDefault(fields["requests_today"], 0),
			TokensToday:    int64Default(fields["tokens_today"], 0),
			RequestsMinute: int(minute),
			Inflight:       atoiDefault(fields["inflight"], 0),
			LastSeenAt:     unixTimeDefault(fields["last_seen_unix"]),
		})
	}
	return out
}

func (s *redisCounterStore) runEnforceScript(keyID, displayName, maskedKey string, limits limitConfig, now time.Time, enforce bool) ([]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), redisCounterTimeout)
	defer cancel()
	member := strconv.FormatInt(now.UnixNano(), 10) + ":" + strconv.FormatInt(s.seq.Add(1), 10)
	enforceValue := "0"
	if enforce {
		enforceValue = "1"
	}
	raw, err := redisEnforceScript.Run(ctx, s.client, []string{s.usageKey(keyID), s.windowKey(keyID), s.keysKey()},
		keyID,
		strings.TrimSpace(displayName),
		strings.TrimSpace(maskedKey),
		now.Format("2006-01-02"),
		strconv.FormatInt(now.Unix(), 10),
		strconv.FormatInt(now.UnixMilli(), 10),
		strconv.FormatInt(now.Add(-time.Minute).UnixMilli(), 10),
		limits.RequestsPerMin,
		limits.RequestsPerDay,
		limits.TokensPerDay,
		limits.MaxInflight,
		enforceValue,
		member,
	).Result()
	if err != nil {
		return nil, err
	}
	values, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("unexpected redis counter response %T", raw)
	}
	return values, nil
}

func (s *redisCounterStore) runReleaseScript(keyID, displayName, maskedKey string, now time.Time, tokens int64) ([]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), redisCounterTimeout)
	defer cancel()
	raw, err := redisReleaseScript.Run(ctx, s.client, []string{s.usageKey(keyID), s.windowKey(keyID), s.keysKey()},
		keyID,
		strings.TrimSpace(displayName),
		strings.TrimSpace(maskedKey),
		now.Format("2006-01-02"),
		strconv.FormatInt(now.Unix(), 10),
		strconv.FormatInt(now.UnixMilli(), 10),
		strconv.FormatInt(now.Add(-time.Minute).UnixMilli(), 10),
		tokens,
	).Result()
	if err != nil {
		return nil, err
	}
	values, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("unexpected redis counter response %T", raw)
	}
	return values, nil
}

func redisCounterReject(values []any) *requestInterceptResponse {
	if len(values) < 4 || redisValueString(values[0]) != "reject" {
		return nil
	}
	status := atoiDefault(redisValueString(values[2]), http.StatusForbidden)
	return &requestInterceptResponse{
		Reject:           true,
		RejectStatusCode: status,
		RejectMessage:    redisValueString(values[3]),
		RejectCode:       redisValueString(values[1]),
	}
}

func redisUnavailableResponse(_ error) *requestInterceptResponse {
	return &requestInterceptResponse{
		Reject:           true,
		RejectStatusCode: http.StatusServiceUnavailable,
		RejectMessage:    "gateway counter backend unavailable",
		RejectCode:       "gateway_counter_backend_unavailable",
		Headers:          http.Header{"X-Gateway-Counter-Error": []string{"unavailable"}},
	}
}

func (s *redisCounterStore) usageKey(keyID string) string {
	return s.prefix + ":usage:" + strings.TrimSpace(keyID)
}

func (s *redisCounterStore) windowKey(keyID string) string {
	return s.prefix + ":window:" + strings.TrimSpace(keyID)
}

func (s *redisCounterStore) keysKey() string {
	return s.prefix + ":usage-keys"
}

func redisValueString(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case []byte:
		return string(typed)
	default:
		return fmt.Sprint(typed)
	}
}

func atoiDefault(value string, fallback int) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return fallback
	}
	return parsed
}

func int64Default(value string, fallback int64) int64 {
	parsed, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return fallback
	}
	return parsed
}

func unixTimeDefault(value string) time.Time {
	parsed := int64Default(value, 0)
	if parsed <= 0 {
		return time.Time{}
	}
	return time.Unix(parsed, 0)
}

package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"rizkysr90-pos/internal/store"
	"time"

	"github.com/redis/go-redis/v9"
)

type Session struct {
	client      *redis.Client
	PrefixState string
	defaultTTL  time.Duration
}

func NewSessionRedisManager(rds *redis.Client) *Session {
	return &Session{
		client:      rds,
		PrefixState: "pos:session",
		defaultTTL:  5 * time.Minute,
	}
}
func (s *Session) buildKey(sessionID string) string {
	return fmt.Sprintf("%s:%s", s.PrefixState, sessionID)
}

// Insert stores session data in Redis
func (s *Session) Insert(ctx context.Context, sessionData *store.SessionRedisData) error {
	if sessionData == nil {
		return fmt.Errorf("session data cannot be nil")
	}

	key := s.buildKey(sessionData.SessionID)

	// Convert roles array to JSON string
	rolesJSON, err := json.Marshal(sessionData.UserRoles)
	if err != nil {
		return fmt.Errorf("failed to marshal user roles: %w", err)
	}

	// Store as hash fields
	fields := map[string]interface{}{
		"user:id":       sessionData.UserID,
		"user:type":     sessionData.UserType,
		"user:email":    sessionData.UserEmail,
		"user:fullName": sessionData.UserFullName,
		"tenant:id":     sessionData.UserTenantID,
		"branch:id":     sessionData.UserBranchID,
		"user:authType": sessionData.UserAuthType,
		"user:roles":    string(rolesJSON),
		"created_at":    sessionData.CreatedAt.UTC().Format(time.RFC3339Nano),
	}

	// Use HSet with expiration
	pipe := s.client.Pipeline()
	pipe.HSet(ctx, key, fields)
	pipe.Expire(ctx, key, s.defaultTTL)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to save session to redis: %w", err)
	}

	return nil
}

func (s *Session) Get(ctx context.Context, sessionID string) (*store.SessionRedisData, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("session ID cannot be empty")
	}

	key := s.buildKey(sessionID)

	// Get all hash fields
	fields, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get session from redis: %w", err)
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("session not found")
	}

	// Parse roles from JSON string
	var roles []string
	if err := json.Unmarshal([]byte(fields["user:roles"]), &roles); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user roles: %w", err)
	}

	// Parse created_at timestamp
	createdAt, err := time.Parse(time.RFC3339Nano, fields["created_at"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	return &store.SessionRedisData{
		SessionID:    sessionID,
		UserID:       fields["user:id"],
		UserType:     fields["user:type"],
		UserEmail:    fields["user:email"],
		UserFullName: fields["user:fullName"],
		UserTenantID: fields["tenant:id"],
		UserBranchID: fields["branch:id"],
		UserAuthType: fields["user:authType"],
		UserRoles:    roles,
		CreatedAt:    createdAt,
	}, nil
}

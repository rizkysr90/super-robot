package store

import (
	"context"
	"time"
)

type SessionRedisData struct {
	SessionID    string
	UserID       string
	UserType     string
	UserEmail    string
	UserFullName string
	UserTenantID string
	UserBranchID string
	UserAuthType string
	UserRoles    []string
	CreatedAt    time.Time
}

type SessionRedis interface {
	Insert(ctx context.Context, sessionData *SessionRedisData) error
	Get(ctx context.Context, sessionID string) (*SessionRedisData, error)
}

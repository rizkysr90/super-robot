package store

import (
	"context"
	"time"
)

type SessionRedisData struct {
	CreatedAt    time.Time
	SessionID    string
	UserID       string
	UserType     string
	UserEmail    string
	UserFullName string
	UserTenantID string
	UserBranchID string
	UserAuthType string
	UserRoles    []string
}

type SessionRedis interface {
	Insert(ctx context.Context, sessionData *SessionRedisData) error
	Get(ctx context.Context, sessionID string) (*SessionRedisData, error)
}

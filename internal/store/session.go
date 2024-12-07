package store

import (
	"context"
	"database/sql"
	"time"
)

type SessionData struct {
	SessionID    string
	AccessToken  string
	RefreshToken string
	UserEmail    string
	UserID       string
	UserFullName string
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    sql.NullTime
}

type Session interface {
	Insert(ctx context.Context, sessionData *SessionData) error
}

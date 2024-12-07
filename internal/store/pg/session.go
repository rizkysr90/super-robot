package pg

import (
	"context"
	"database/sql"
	"rizkysr90-pos/internal/store"

	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type Session struct {
	db *sql.DB
}

func NewSession(db *sql.DB) *Session {
	return &Session{
		db: db,
	}
}

func (s *Session) Insert(ctx context.Context, sessionData *store.SessionData) error {
	query := `
	INSERT INTO sessions (
		session_id, access_token, refresh_token,
		user_email, user_id, user_fullname, created_at, expires_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
`
	insertFunc := func(tx sqldb.QueryExecutor) error {
		_, err := tx.ExecContext(ctx, query,
			sessionData.SessionID,
			sessionData.AccessToken,
			sessionData.RefreshToken,
			sessionData.UserEmail,
			sessionData.UserID,
			sessionData.UserFullName,
			sessionData.CreatedAt,
			sessionData.ExpiresAt,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, insertFunc)
}

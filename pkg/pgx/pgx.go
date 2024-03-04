package pgx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	Username    string
	Password    string
	Database    string
	Host        string
	Port        int
	ConnMaxOpen int
	ConnMaxIdle int
}

func (c Config) DSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		c.Username, c.Password, c.Host, c.Port, c.Database,
	)
}
func NewDB(cfg Config, ctx context.Context) (*sql.DB, error) {
	conCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("pgx: parse DSN failed: %w", err)
	}
	// Set the maximum connection lifetime to 5 minutes
	maxLifetime := 5 * time.Minute
	conCfg.MaxConnLifetime = maxLifetime
	conCfg.MaxConns = 20
	var pool *pgxpool.Pool
	pool, err = pgxpool.NewWithConfig(ctx, conCfg)

	if err != nil {
		return nil, fmt.Errorf("pgx: open db is failed: %w", err)
	}
	db := stdlib.OpenDBFromPool(pool)
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("pgx: DB ping failed: %w", err)
	}
	return db, nil
}

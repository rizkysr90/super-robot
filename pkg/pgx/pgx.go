package pgx

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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
func NewDB(cfg Config, ctx context.Context) (*pgxpool.Pool, error) {
	conCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("pgx: parse DSN failed: %w", err)
	}
	// Set the maximum connection lifetime to 5 minutes
	maxLifetime := 5 * time.Minute
	conCfg.MaxConnLifetime = maxLifetime
	conCfg.MaxConns = 20
	var db *pgxpool.Pool
	db, err = pgxpool.NewWithConfig(ctx, conCfg)
	if err != nil {
		return nil, fmt.Errorf("pgx: open db is failed: %w", err)
	}
	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pgx: DB ping failed: %w", err)
	}
	return db, nil
}

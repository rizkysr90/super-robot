package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
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
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		c.Username, c.Password, c.Host, c.Port, c.Database,
	)
}
func NewDB(cfg Config) (*sql.DB, error) {
	conCfg, err := mysql.ParseDSN(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("mysql: parse DSN failed: %w", err)
	}
	var db *sql.DB
	db, err = sql.Open("mysql", conCfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("mysql: open db is failed: %w", err)
	}
	// Set the maximum connection lifetime to 5 minutes
	maxLifetime := 5 * time.Minute
	db.SetConnMaxLifetime(maxLifetime)
	db.SetMaxOpenConns(cfg.ConnMaxOpen)
	db.SetMaxIdleConns(cfg.ConnMaxIdle)
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("mysql: DB ping failed: %w", err)
	}
	return db, nil
}

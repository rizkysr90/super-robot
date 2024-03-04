package sqldb

import (
	"context"
	"database/sql"
	"fmt"
)

// QueryExecutor is common methods that is intended to be scoped to a DB transaction implemented by *sql.Tx.
//
// It is also a set of methods that is implemented by *sql.DB, so we can use this interface for both.
type QueryExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// TxFunc is a func that is scoped to single DB transaction from WithinTx.
type TxFunc func(QueryExecutor) error

type txContextKey struct{}

// WithTxContext create a copy of given context with given DB transaction interface.
// Use case of this method is if we want execute single DB transaction for multiple SQL DB datastore/repository call.
// Data store or repository should explicitly check whether context contains transaction Query.
//
// See WithinTxContextOrError or WithinTxContextOrDB as the helper function for this.
func WithTxContext(ctx context.Context, tx QueryExecutor) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

// TxFromContext extract *sql.Tx that implements QueryExecutor and make sure it's not original *sql.DB.
func TxFromContext(ctx context.Context) QueryExecutor {
	val := ctx.Value(txContextKey{})
	switch val.(type) {
	case *sql.DB:
		// *sql.DB also implement QueryExecutor, but we make sure to return it was a DB transaction.
		// See WithinTxContextOrDB() to get executor with fallback.
		return nil
	case *sql.Tx, QueryExecutor:
		return val.(QueryExecutor)
	default:
		return nil
	}
}

// WithinTxContextOrDB is a helper function to get the SQL database executor with fallback using given sqlDB.
//
// If there is a database transaction found in given context by the upper application layer, use it,
// it means that the method call is part of multiple data store or repository call in a single transaction,
// otherwise return given *sql.DB as the executor as it is also implements QueryExecutor.
func WithinTxContextOrDB(ctx context.Context, sqlDB *sql.DB) QueryExecutor {
	if tx := TxFromContext(ctx); tx != nil {
		return tx
	}
	return sqlDB
}

// WithinTxContextOrError is a helper function to execute given txFunc only within db transaction from given context,
// if a Tx could NOT be fetched from given context, it will return error.
//
// This operation ideally called on data change: INSERT, UPDATE, DELETE.
func WithinTxContextOrError(ctx context.Context, txFunc TxFunc) error {
	if tx := TxFromContext(ctx); tx != nil {
		return txFunc(tx)
	}
	return fmt.Errorf("sqldb: WithinTxContextOrError() could not fetch a *sql.Tx from given context")
}

type txOption struct {
	txOpt *sql.TxOptions
}

type TxOption func(*txOption)

func TxIsolationLevelDefault() TxOption {
	return func(opt *txOption) { opt.txOpt.Isolation = sql.LevelDefault }
}
func TxIsolationLevelReadUncommitted() TxOption {
	return func(opt *txOption) { opt.txOpt.Isolation = sql.LevelReadUncommitted }
}
func TxIsolationLevelReadCommitted() TxOption {
	return func(opt *txOption) { opt.txOpt.Isolation = sql.LevelReadCommitted }
}
func TxIsolationLevelWriteCommitted() TxOption {
	return func(opt *txOption) { opt.txOpt.Isolation = sql.LevelWriteCommitted }
}
func TxIsolationLevelRepeatableRead() TxOption {
	return func(opt *txOption) { opt.txOpt.Isolation = sql.LevelRepeatableRead }
}
func TxIsolationLevelSnapshot() TxOption {
	return func(opt *txOption) { opt.txOpt.Isolation = sql.LevelSnapshot }
}
func TxIsolationLevelSerializable() TxOption {
	return func(opt *txOption) { opt.txOpt.Isolation = sql.LevelSerializable }
}
func TxIsolationLevelLinearizable() TxOption {
	return func(opt *txOption) { opt.txOpt.Isolation = sql.LevelLinearizable }
}

// TxReadOnly pass read only flag for current transaction.
func TxReadOnly() TxOption { return func(opt *txOption) { opt.txOpt.ReadOnly = true } }

// WithinTx create a scoped func that wrap DB transaction begin, commit, and rollback mechanism.
// Any error returned by txFunc will rollback a transaction.
//
// To execute a DB transaction in multiple data store, see WithTxContext().
func WithinTx(ctx context.Context, db *sql.DB, txFunc TxFunc, txOpts ...TxOption) (err error) {
	var txOpt *sql.TxOptions
	if len(txOpts) > 0 {
		opt := &txOption{txOpt: &sql.TxOptions{}}
		for _, o := range txOpts {
			o(opt)
		}
		txOpt = opt.txOpt
	}
	sqlTx, err := db.BeginTx(ctx, txOpt)
	if err != nil {
		return fmt.Errorf("sqldb: WithinTx begin SQL transaction failed: %w", err)
	}
	defer func() {
		p := recover()
		if p != nil {
			// Rollback immediately, ignore error and continue original panicking.
			_ = sqlTx.Rollback()
			panic(p)
		}
		if err != nil {
			rollErr := sqlTx.Rollback()
			if rollErr != nil {
				err = fmt.Errorf("sqldb: WithinTx rollback failed before commit: %w", rollErr)
				return
			}
			err = fmt.Errorf("sqldb: WithinTx failed before commit: %w", err)
			return
		}
		err = sqlTx.Commit()
		if err != nil {
			rollErr := sqlTx.Rollback()
			if rollErr != nil {
				err = fmt.Errorf("sqldb: WithinTx commit rollback error after commit failed: %w", rollErr)
				return
			}
			err = fmt.Errorf("sqldb: WithinTx commit failed: %w", err)
		}
	}()
	err = txFunc(sqlTx)
	return err
}

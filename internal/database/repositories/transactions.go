package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/antelman107/filestorage/pkg/domain"
)

const txContextKey = "sqx_tx"

type dbExecutor interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

// WithTransaction gets transaction from db.InitTransaction(),
// runs function f with transaction stored in context.
// Depending on function return there is either COMMIT or ROLLBACK after running the function.
// WithTransaction is separated from PostgresDB because no need to mock this logic.
func WithTransaction(ctx context.Context, repo domain.TransactionalRepository, f func(ctx context.Context) error) error {
	tx, err := repo.InitTransaction()
	if err != nil {
		return err
	}

	// nil tx can only happen if repo is a mock
	if tx != nil {
		ctx = context.WithValue(ctx, txContextKey, tx)
	}

	if err := f(ctx); err != nil {
		if tx == nil {
			// nil tx can only happen if repo is a mock
			return err
		}

		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("tx rollback failed: %w", err)
		}
		return err
	}

	if tx == nil {
		// nil tx can only happen if repo is a mock
		return nil
	}

	return tx.Commit()
}

// getDBExecutorFromCtx extracts transaction from context. if no transaction is stored, defaultStmt is returned.
// to be used database methods to work both with/without transactions
func getDBExecutorFromCtx(ctx context.Context, defaultExecutor dbExecutor) dbExecutor {
	val := ctx.Value(txContextKey)
	if val == nil {
		return defaultExecutor
	}

	return val.(dbExecutor)
}

package dbexecutor

import (
	"context"
	"database/sql"
)

type (
	IDBExecutor interface {
		RunWithTX(ctx context.Context, fn func(tx *sql.Tx) error) error
		RunWithDB(ctx context.Context, fn func(db *sql.DB) error) error
		GetDB() *sql.DB
	}
	DBExecutor struct {
		db *sql.DB
	}
)

func NewDBExecutor(db *sql.DB) IDBExecutor {
	return &DBExecutor{db: db}
}

// Get a Transaction callback
func (e *DBExecutor) RunWithTX(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// Get a db callback
func (e *DBExecutor) RunWithDB(ctx context.Context, fn func(db *sql.DB) error) error {
	err := fn(e.db)
	if err != nil {
		return err
	}
	return nil
}

func (e *DBExecutor) GetDB() *sql.DB {
	return e.db
}

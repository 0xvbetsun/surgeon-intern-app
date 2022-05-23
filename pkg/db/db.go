package db

import (
	"database/sql"
)

type Postgresql struct {
	conn *sql.DB
}

func NewSqlConnection(pgConnectionParameters string) (*Postgresql, error) {
	pgConnection, err := connect(pgConnectionParameters)
	if err != nil {
		return nil, err
	}
	return &Postgresql{conn: pgConnection}, nil
}

// Connect ...
func connect(pgConnectionParameters string) (*sql.DB, error) {
	db, err := sql.Open("postgres", pgConnectionParameters)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// Disconnect ...
func (psqlConnection *Postgresql) Disconnect() error {
	err := psqlConnection.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (PsqlConnection *Postgresql) Connection() *sql.DB {
	return PsqlConnection.conn
}

package db

import (
	"database/sql"
)

type Postgresql struct {
	Conn *sql.DB
}

type ConnectionParameters string

func NewSqlConnection(pgConnectionParameters ConnectionParameters) (*Postgresql, error) {
	pgConnection, err := connect(pgConnectionParameters)
	if err != nil {
		return nil, err
	}
	return &Postgresql{Conn: pgConnection}, nil
}

// Connect ...
func connect(pgConnectionParameters ConnectionParameters) (*sql.DB, error) {
	db, err := sql.Open("postgres", string(pgConnectionParameters))
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
	err := psqlConnection.Conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (PsqlConnection *Postgresql) Connection() *sql.DB {
	return PsqlConnection.Conn
}

package storage

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type StoreDBMS struct {
	DB  *sql.DB
	CTX context.Context
}

func (dbms *StoreDBMS) GetConn() (*sql.Conn, error) {
	conn, err := dbms.DB.Conn(dbms.CTX)

	return conn, err
}

func NewDBMS(connString string) (*StoreDBMS, error) {
	ctx := context.Background()
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	return &StoreDBMS{
		DB:  db,
		CTX: ctx,
	}, nil
}

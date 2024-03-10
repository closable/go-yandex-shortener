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

func NewDBMS(connString string) (*StoreDBMS, error) {
	ctx := context.Background()

	db, err := sql.Open("pgx", connString)

	return &StoreDBMS{
		DB:  db,
		CTX: ctx,
	}, err
}

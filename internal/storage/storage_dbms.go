package storage

import (
	"context"
	"database/sql"
	"fmt"

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
	// db, err := sql.Open("pgx", connString)
	ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		`localhost`, `postgres`, `postgres`, `praktikum`)
	fmt.Println(connString)

	db, err := sql.Open("pgx", ps)
	if err != nil {
		return nil, err
	}

	return &StoreDBMS{
		DB:  db,
		CTX: ctx,
	}, nil
}

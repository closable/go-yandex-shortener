package storage

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type StoreDBMS struct {
	DB  *sql.DB
	CTX context.Context
}

func NewDBMS(connString string, logger zap.Logger) *StoreDBMS {
	ctx := context.Background()
	sugar := *logger.Sugar()
	db, err := sql.Open("pgx", connString)

	if err != nil {
		sugar.Panicln("Unable to connection to database", err)
		os.Exit(1)
	}
	//defer db.Close()

	return &StoreDBMS{
		DB:  db,
		CTX: ctx,
	}
}

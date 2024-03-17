package storage

import (
	"context"
	"database/sql"
	"errors"
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

// create schena if schema doesn't exists
func (dbms *StoreDBMS) CreateSchema() error {
	sql := "CREATE SCHEMA IF NOT EXISTS ya AUTHORIZATION postgres"
	_, err := dbms.DB.ExecContext(dbms.CTX, sql)
	if err != nil {
		return err
	}

	return nil
}

// create table if table doesn't exists
func (dbms *StoreDBMS) CreateTable() error {
	// sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (key varchar(10) not null, url text)", name)

	sql := `
	CREATE TABLE IF NOT EXISTS ya.shortener (
		key varchar(10) not null, 
		url text)
	`
	_, err := dbms.DB.ExecContext(dbms.CTX, sql)
	if err != nil {
		return fmt.Errorf("error during creating schema - %w", err)
	}

	return nil
}

func (dbms *StoreDBMS) CreateIndex() error {
	var cnt int
	sql := `select count(*) from pg_indexes where tablename = 'shortener' and indexname = 'url'`
	// check index
	err := dbms.DB.QueryRowContext(dbms.CTX, sql).Scan(&cnt)
	if err != nil {
		return fmt.Errorf("error during check index - %w", err)
	}

	create := `CREATE UNIQUE INDEX url ON ya.shortener USING btree (url ASC NULLS LAST) TABLESPACE pg_default`
	// create index if not exists
	if cnt == 0 {
		_, err = dbms.DB.ExecContext(dbms.CTX, create)
		if err != nil {
			return fmt.Errorf("error during creating index %w", err)
		}
	}
	return nil
}

// add new shortener and return key
func (dbms *StoreDBMS) GetShortener(url string) (string, error) {
	sqlBefore := "SELECT key, url FROM ya.shortener WHERE url like '" + url + "%' order by length(url) asc limit 1"

	sql := `MERGE INTO ya.shortener ys using
				(SELECT $1 url) res ON (ys.url = res.url) 
				WHEN NOT MATCHED 
				THEN INSERT (key, url) VALUES (substr(md5(random()::text), 1, 10), res.url)`

	tx, err := dbms.DB.BeginTx(dbms.CTX, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var short, existingURL string
	_ = tx.QueryRowContext(dbms.CTX, sqlBefore).Scan(&short, &existingURL)

	// if URL exists then return key and mark error 409
	if len(existingURL) > 0 {
		return short, errors.New("409")
	}

	_, err = tx.ExecContext(dbms.CTX, sql, url)
	if err != nil {
		return "", err
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	shortener, find := dbms.FindKeyByValue(url)
	if !find {
		//err := errors.New("key by url not foundß")
		return "", err
	}
	return shortener, nil
}

// get shortener by url
func (dbms *StoreDBMS) FindExistingKey(key string) (string, bool) {
	sql := "SELECT url FROM ya.shortener WHERE key = $1"
	var url string

	err := dbms.DB.QueryRowContext(dbms.CTX, sql, key).Scan(&url)
	if err != nil {
		return "", false
	}

	return url, true
}

func (dbms *StoreDBMS) FindKeyByValue(url string) (string, bool) {
	sql := "SELECT key FROM ya.shortener WHERE url is not null and url = $1"
	var key string

	err := dbms.DB.QueryRowContext(dbms.CTX, sql, url).Scan(&key)
	if err != nil {
		return "", false
	}

	return key, true
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

func (dbms *StoreDBMS) Ping() bool {
	conn, err := dbms.GetConn()
	if err != nil {
		return false
	}
	defer conn.Close()
	err = conn.PingContext(dbms.CTX)
	return err == nil
}

func (dbms *StoreDBMS) PrepareStore() {

	err := dbms.CreateSchema()
	if err != nil {
		fmt.Println(err)
	}

	err = dbms.CreateTable()
	if err != nil {
		fmt.Println(err)
	}

	err = dbms.CreateIndex()
	if err != nil {
		fmt.Println(err)
	}

}

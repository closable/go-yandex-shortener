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

// create schena if schema doesn't exists
func (dbms *StoreDBMS) CreateSchema(name string) error {
	sql := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s AUTHORIZATION postgres", name)
	_, err := dbms.DB.ExecContext(dbms.CTX, sql)
	if err != nil {
		return err
	}

	return nil
}

// create table if table doesn't exists
func (dbms *StoreDBMS) CreateTable(name string) error {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (key varchar(10) not null, url text)", name)
	_, err := dbms.DB.ExecContext(dbms.CTX, sql)
	if err != nil {
		return err
	}

	return nil
}

// add new shortener and return key
func (dbms *StoreDBMS) GetShortener(url string) string {
	sql := "MERGE INTO ya.shortener ys using (SELECT $1 url) res ON (ys.url = res.url) WHEN NOT MATCHED THEN INSERT (key, url) VALUES (substr(md5(random()::text), 1, 10), res.url)"

	_, err := dbms.DB.ExecContext(dbms.CTX, sql, url)
	if err != nil {
		return ""
	}

	shortener, find := dbms.FindKeyByValue(url)
	if !find {
		//err := errors.New("key by url not fount")
		return ""
	}

	return shortener
}

// get shortener by url
func (dbms *StoreDBMS) FindExistingKey(key string) (string, bool) {
	sql := "SELECT url FROM ya.shortener WHERE key = $1"
	var url string = ""

	err := dbms.DB.QueryRowContext(dbms.CTX, sql, key).Scan(&url)
	if err != nil {
		return "", false
	}

	return url, true
}

func (dbms *StoreDBMS) FindKeyByValue(url string) (string, bool) {
	sql := "SELECT key FROM ya.shortener WHERE url = $1"
	var key string = ""

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

	err := dbms.CreateSchema("ya")
	if err != nil {
		fmt.Println("!!! ошибка создания схемы", err)
	}

	err = dbms.CreateTable("ya.shortener")
	if err != nil {
		fmt.Println("!!! ошибка создания таблицы", err)
	}

	// key, err := dbms.AddItem("", "www.yandex.ru/124")
	// if err != nil {
	// 	fmt.Println("!!! ошибка при добавлении", err)
	// }
}

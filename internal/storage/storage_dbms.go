package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Описание структур данных
type (
	// StoreDBMS модель данных
	StoreDBMS struct {
		DB *sql.DB
		//CTX context.Context
	}
	//authURLs модель хранения
	authURLs struct {
		sync.RWMutex
		urls map[string]string
	}
	// Semaphore структуа модели
	Semaphore struct {
		semaCh chan struct{}
	}
)

// NewSemaphore создает семафор с буферизованным каналом емкостью maxReq
func NewSemaphore(maxReq int) *Semaphore {
	return &Semaphore{
		semaCh: make(chan struct{}, maxReq),
	}
}

// Acquire когда горутина запускается, отправляем пустую структуру в канал semaCh
func (s *Semaphore) Acquire() {
	s.semaCh <- struct{}{}
}

// Release когда горутина завершается, из канала semaCh убирается пустая структура
func (s *Semaphore) Release() {
	<-s.semaCh
}

// GetConn получение коннекта к СУБД
func (dbms *StoreDBMS) GetConn() (*sql.Conn, error) {
	ctx := context.Background()
	conn, err := dbms.DB.Conn(ctx)

	return conn, err
}

// CreateSchema create schena if schema doesn't exists
func (dbms *StoreDBMS) CreateSchema() error {
	sql := "CREATE SCHEMA IF NOT EXISTS ya AUTHORIZATION postgres"
	ctx := context.Background()
	_, err := dbms.DB.ExecContext(ctx, sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateTable create table if table doesn't exists
func (dbms *StoreDBMS) CreateTable() error {
	// sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (key varchar(10) not null, url text)", name)
	ctx := context.Background()
	sql := `
	CREATE TABLE IF NOT EXISTS ya.shortener (
		key varchar(10) not null, 
		url text,
		user_id int,
		is_deleted bool
	)
	`
	_, err := dbms.DB.ExecContext(ctx, sql)
	if err != nil {
		return fmt.Errorf("error during creating schema - %w", err)
	}
	return nil
}

// CreateIndex создает индех при необходимости
func (dbms *StoreDBMS) CreateIndex() error {
	var cnt int
	sql := `select count(*) from pg_indexes where tablename = 'shortener' and indexname = 'url'`

	ctx := context.Background()
	// check index
	err := dbms.DB.QueryRowContext(ctx, sql).Scan(&cnt)
	if err != nil {
		return fmt.Errorf("error during check index - %w", err)
	}

	create := `CREATE UNIQUE INDEX url ON ya.shortener USING btree (url ASC NULLS LAST) TABLESPACE pg_default`
	// create index if not exists
	if cnt == 0 {
		_, err = dbms.DB.ExecContext(ctx, create)
		if err != nil {
			return fmt.Errorf("error during creating index %w", err)
		}
	}
	return nil
}

// GetShortener add new shortener and return key
func (dbms *StoreDBMS) GetShortener(userID int, url string) (string, error) {
	sqlBefore := "SELECT key, url FROM ya.shortener WHERE url like '" + url + "%' order by length(url) asc limit 1"

	sql := `MERGE INTO ya.shortener ys using
				(SELECT $1 url) res ON (ys.url = res.url) 
				WHEN NOT MATCHED 
				THEN INSERT (key, url, user_id) 
				VALUES (
					substr(md5(random()::text), 1, 10), 
					res.url, 
					case when $2 <> 0 then $2 else floor(random() * (20-1+1) + 1)::int end
					)`

	ctx := context.Background()
	tx, err := dbms.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var short, existingURL string
	_ = tx.QueryRowContext(ctx, sqlBefore).Scan(&short, &existingURL)
	// if URL exists then return key and mark error 409
	if len(existingURL) > 0 {
		return short, errors.New("409")
	}

	_, err = tx.ExecContext(ctx, sql, url, userID)
	if err != nil {
		// fmt.Println("!!!!", err) TODO !!! ошибка не обработанаы
		return "", err
	}
	if err = tx.Commit(); err != nil {
		// fmt.Println("!!!!", err) // TODO !!! ошибка не обработанаы
		return "", err
	}
	shortener, find := dbms.FindKeyByValue(url)
	if !find {
		//err := errors.New("key by url not found")
		return "", err
	}
	return shortener, nil
}

// FindExistingKey  get shortener by url
func (dbms *StoreDBMS) FindExistingKey(key string) (string, bool) {
	sql := "SELECT url, coalesce(is_deleted, false) is_deleted FROM ya.shortener WHERE key = $1"
	var url string
	var isDeleted bool
	ctx := context.Background()
	err := dbms.DB.QueryRowContext(ctx, sql, key).Scan(&url, &isDeleted)
	if err != nil {
		return "", false
	}

	if isDeleted {
		return "", true
	}
	return url, true
}

// FindKeyByValue поиск ключа по значению
func (dbms *StoreDBMS) FindKeyByValue(url string) (string, bool) {
	sql := "SELECT key FROM ya.shortener WHERE url is not null and url = $1"
	var key string
	ctx := context.Background()
	err := dbms.DB.QueryRowContext(ctx, sql, url).Scan(&key)
	if err != nil {
		return "", false
	}

	return key, true
}

// NewDBMS новый экземпляр СУБД
func NewDBMS(connString string) (*StoreDBMS, error) {
	//ctx := context.Background()
	db, err := sql.Open("pgx", connString)

	if err != nil {
		return nil, err
	}

	return &StoreDBMS{
		DB: db,
		// CTX: ctx,
	}, nil
}

// Ping функция контроля работоспособности СУБД
func (dbms *StoreDBMS) Ping() bool {
	conn, err := dbms.GetConn()
	if err != nil {
		return false
	}
	defer conn.Close()
	ctx := context.Background()
	err = conn.PingContext(ctx)
	return err == nil
}

// PrepareStore функция подготовки хранилища при 1м запуске
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

// GetURLs получение urls пользователя
func (dbms *StoreDBMS) GetURLs(userID int) (map[string]string, error) {
	var result = &authURLs{urls: make(map[string]string)} // make(map[string]string)
	ctx := context.Background()
	sql := "SELECT key, url FROM ya.shortener where user_id = $1"

	stmt, err := dbms.DB.PrepareContext(ctx, sql)
	if err != nil {
		return result.urls, err
	}

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil || rows.Err() != nil {
		return result.urls, err
	}
	var key, url string

	for rows.Next() {
		if err := rows.Scan(&key, &url); err != nil {
			return result.urls, err
		}
		result.RWMutex.Lock()
		result.urls[key] = url
		result.RWMutex.Unlock()
	}
	return result.urls, nil
}

// SoftDeleteURLs мягкое уалене занечений из СУБД
func (dbms *StoreDBMS) SoftDeleteURLs(userID int, keys ...string) error {
	var wg sync.WaitGroup
	ctx := context.Background()
	var errList = make([]string, 0)
	semaphore := NewSemaphore(4)

	for _, keyString := range keys {
		wg.Add(1)
		go func(idKey string) {
			semaphore.Acquire()
			defer wg.Done()
			defer semaphore.Release()

			tx, _ := dbms.DB.BeginTx(ctx, nil)

			defer tx.Rollback()

			sql := "UPDATE ya.shortener SET is_deleted=true where key=$1 and user_id=$2 and not coalesce(is_deleted, false)"
			stmt, err := tx.PrepareContext(ctx, sql)
			if err != nil {
				errList = append(errList, idKey)
				fmt.Println("\nerror with sql prepare for key=", idKey, userID)
			}

			_, err = stmt.ExecContext(ctx, idKey, userID)
			if err != nil {
				errList = append(errList, idKey)
				fmt.Println("\nerror during execute for key=", idKey)
			}

			if err = tx.Commit(); err != nil {
				fmt.Println("tx error ", err)
			}

		}(keyString)
	}
	wg.Wait()

	// del := "delete from ya.shortener where is_deleted"
	// _, err := dbms.DB.ExecContext(ctx, del)
	// if err != nil {
	// 	return err
	// }

	if len(errList) > 0 {
		return errors.New("during soft delete records where found some errors")
	}

	return nil
}

// GetStats получене статистики  сохранненных данных
func (dbms *StoreDBMS) GetStats() (int, int) {
	ctx := context.Background()
	sql := "SELECT count(*) urls, count(distinct user_id) useres FROM ya.shortener"

	var urls, users int
	err := dbms.DB.QueryRowContext(ctx, sql).Scan(&urls, &users)
	if err != nil {
		return 0, 0
	}

	return urls, users
}

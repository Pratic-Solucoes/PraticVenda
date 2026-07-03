package db

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type BancoDados struct {
	ConnectionString string
	Driver           string
	MaxOpenConns     int
	MaxIdleConns     int
	MaxIdleTime      time.Duration
}

func (bd *BancoDados) Conectar() (*sql.DB, error) {

	db, err := sql.Open(bd.Driver, bd.ConnectionString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(bd.MaxOpenConns)
	db.SetMaxIdleConns(bd.MaxIdleConns)
	db.SetConnMaxIdleTime(bd.MaxIdleTime)

	return db, nil
}

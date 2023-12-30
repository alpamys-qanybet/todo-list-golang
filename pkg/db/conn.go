package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

var connectionPool *pgxpool.Pool

func Connect(url string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}

	connectionPool = dbpool
	connectionPool.Ping(context.Background())
	return dbpool, nil
}

func ConnectionPool() (*pgxpool.Pool, error) {
	if connectionPool == nil {
		return nil, errors.New("postgres db connection pool not connected")
	}
	return connectionPool, nil
}

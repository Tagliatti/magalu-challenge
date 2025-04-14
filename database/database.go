package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func Connect(ctx context.Context, minConns int32, maxConns int32) (*pgxpool.Pool, error) {
	dataSource := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	poolConfig, err := pgxpool.ParseConfig(dataSource)

	if err != nil {
		return nil, err
	}

	poolConfig.MinConns = minConns
	poolConfig.MaxConns = maxConns

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)

	if err != nil {
		return nil, err
	}

	return pool, nil
}

func ConnectTest(ctx context.Context, connectionStr string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(connectionStr)

	if err != nil {
		return nil, err
	}

	return pgxpool.NewWithConfig(ctx, poolConfig)
}

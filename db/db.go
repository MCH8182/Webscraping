package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	Pool *pgxpool.Pool
)

func StartConnection() error {
	connStr := "postgres://postgres.jfugaikxhuxsryzqpres:mtlfztox1987@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres"
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return err
	}
	Pool = pool
	return nil
}

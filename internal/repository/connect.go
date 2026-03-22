package repository

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	conn := os.Getenv("DB_CONNECTION")
	return pgxpool.New(ctx, conn)
}

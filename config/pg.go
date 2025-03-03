package config

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPg(ctx context.Context, url string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}
	slog.Info("db connection established.")
	_, err = conn.Exec(ctx, "SET TIMEZONE TO 'Asia/Jakarta';")
	if err != nil {
		return nil, err
	}
	slog.Info("db timezone set to Asia/Jakarta.")
	return conn, nil
}

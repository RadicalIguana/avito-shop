package database

import (
	"context"
	"fmt"
	"os"
	"time"

	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool // Глобальная переменная для пула соединений

func Connect() error {
	// TODO: Добавить .env
	dsn := fmt.Sprintf(
		// "postgres://%s:%s@%s:%s/%s?sslmode=disable",
		"postgres://postgres:postgres@db:5432/avito_shop_db?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse DSN: %w", err)
	}

	config.MaxConns = 50
	config.MinConns = 5
	config.MaxConnLifetime = 2 * time.Hour
    config.MaxConnIdleTime = 10 * time.Minute
    config.HealthCheckPeriod = 30 * time.Second
    config.ConnConfig.ConnectTimeout = 10 * time.Second

	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := DB.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Connected to database successfully")
	return nil
}
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

	config.MaxConns = 500
	config.MinConns = 50
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

// var DB *pgx.Conn // Глобальная переменная для соединения с базой данных

// func Connect() error {
//     // Формируем строку подключения (DSN)
//     dsn := fmt.Sprintf(
//         // "postgres://%s:%s@%s:%s/%s?sslmode=disable",
// 		"postgres://postgres:postgres@localhost:5432/avito_shop_db?sslmode=disable",
//         os.Getenv("DB_USER"),
//         os.Getenv("DB_PASSWORD"),
//         os.Getenv("DB_HOST"),
//         os.Getenv("DB_PORT"),
//         os.Getenv("DB_NAME"),
//     )

//     // Создаем контекст с таймаутом для подключения
//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()

//     // Устанавливаем соединение с базой данных
//     conn, err := pgx.Connect(ctx, dsn)
//     if err != nil {
//         return fmt.Errorf("failed to connect to database: %w", err)
//     }

//     // Проверяем соединение с базой данных
//     if err := conn.Ping(ctx); err != nil {
//         return fmt.Errorf("failed to ping database: %w", err)
//     }

//     // Сохраняем соединение в глобальную переменную
//     DB = conn

//     fmt.Println("Connected to database successfully")
//     return nil
// }
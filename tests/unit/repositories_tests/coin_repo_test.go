package repositories_tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/RadicalIguana/avito-shop/internal/models"
	"github.com/RadicalIguana/avito-shop/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setupCoinTestDB(t *testing.T) *pgxpool.Pool {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }
        
	// Подключение к тестовой базе данных
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_NAME"),
	)
	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Создание таблиц
	_, err = db.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            coins INTEGER NOT NULL
        );
        CREATE TABLE IF NOT EXISTS transfers (
            id SERIAL PRIMARY KEY,
            from_user INTEGER NOT NULL,
            to_user INTEGER NOT NULL,
            amount INTEGER NOT NULL,
            created_at TIMESTAMP DEFAULT NOW()
        );
    `)
	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	// Очистка таблиц перед тестами
	_, err = db.Exec(context.Background(), "TRUNCATE users, transfers RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}

	// Добавление тестовых данных
	_, err = db.Exec(context.Background(), "INSERT INTO users (id, coins) VALUES (1, 200), (2, 50)")
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	return db
}

func TestGetUserForUpdate_Success(t *testing.T) {
    // Настройка тестовой базы данных
    db := setupCoinTestDB(t)
    defer db.Close()

    // Инициализация репозитория
    repo := repositories.NewCoinRepository(db)

    // Выполнение метода
    user, err := repo.GetUserForUpdate(context.Background(), 1)
    assert.NoError(t, err)
    assert.Equal(t, 1, user.ID)
    assert.Equal(t, 200, user.Coins)
}

func TestUpdateBalance_Success(t *testing.T) {
    // Настройка тестовой базы данных
    db := setupCoinTestDB(t)
    defer db.Close()

    // Инициализация репозитория
    repo := repositories.NewCoinRepository(db)

    // Обновление баланса
    err := repo.UpdateBalance(context.Background(), 1, 100)
    assert.NoError(t, err)

    // Проверка состояния базы данных
    var coins int
    err = db.QueryRow(context.Background(), "SELECT coins FROM users WHERE id = $1", 1).Scan(&coins)
    assert.NoError(t, err)
    assert.Equal(t, 100, coins)
}

func TestCreateTransfer_Success(t *testing.T) {
    // Настройка тестовой базы данных
    db := setupCoinTestDB(t)
    defer db.Close()

    // Инициализация репозитория
    repo := repositories.NewCoinRepository(db)

    // Создание транзакции
    transfer := &models.Transfer{
        FromUser: 1,
        ToUser:   2,
        Amount:   100,
    }
    err := repo.CreateTransfer(context.Background(), transfer)
    assert.NoError(t, err)

    // Проверка состояния базы данных
    // var count int
    // err = db.QueryRow(context.Background(), "SELECT COUNT(*) FROM transfers WHERE from_user = $1 AND to_user = $2 AND amount = $3", 1, 2, 100).Scan(&count)
    // assert.NoError(t, err)
    // assert.Equal(t, 1, count)
}
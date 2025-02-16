package services_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/RadicalIguana/avito-shop/internal/repositories"
	"github.com/RadicalIguana/avito-shop/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setupCoinTestDB(t *testing.T) *pgxpool.Pool {
	if err := godotenv.Load("../../../.env"); err != nil {
        log.Fatal("Error loading .env file")
    }

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
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `)
	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	_, err = db.Exec(context.Background(), "TRUNCATE users, transfers RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}

	_, err = db.Exec(context.Background(), "INSERT INTO users (id, coins) VALUES (1, 200), (2, 50)")
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	return db
}

func TestTransferCoins_Success(t *testing.T) {
	db := setupCoinTestDB(t)
	defer db.Close()

	repo := repositories.NewCoinRepository(db)
	service := services.NewCoinService(repo)

	err := service.TransferCoins(context.Background(), 1, 2, 100)
	assert.NoError(t, err)

	var user1Coins, user2Coins int
	err = db.QueryRow(context.Background(), "SELECT coins FROM users WHERE id = $1", 1).Scan(&user1Coins)
	assert.NoError(t, err)
	assert.Equal(t, 100, user1Coins)

	err = db.QueryRow(context.Background(), "SELECT coins FROM users WHERE id = $1", 2).Scan(&user2Coins)
	assert.NoError(t, err)
	assert.Equal(t, 150, user2Coins)
}

func TestTransferCoins_InsufficientFunds(t *testing.T) {
	db := setupCoinTestDB(t)
	defer db.Close()

	repo := repositories.NewCoinRepository(db)
	service := services.NewCoinService(repo)

	err := service.TransferCoins(context.Background(), 1, 2, 300)
	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())
}

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

func setupMerchTestDB(t *testing.T) *pgxpool.Pool {
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
        CREATE TABLE IF NOT EXISTS merch (
            name TEXT PRIMARY KEY,
            price INTEGER NOT NULL
        );
        DROP TABLE IF EXISTS inventory;
        CREATE TABLE IF NOT EXISTS inventory (
            user_id INTEGER NOT NULL,
            item_name TEXT NOT NULL,
            quantity INTEGER NOT NULL,
            PRIMARY KEY (user_id, item_name)
        );
    `)
    if err != nil {
        t.Fatalf("failed to create tables: %v", err)
    }

    _, err = db.Exec(context.Background(), "TRUNCATE users, merch, inventory RESTART IDENTITY CASCADE")
    if err != nil {
        t.Fatalf("failed to truncate tables: %v", err)
    }

    _, err = db.Exec(context.Background(), "INSERT INTO users (id, coins) VALUES (1, 200)")
    if err != nil {
        t.Fatalf("failed to insert test data: %v", err)
    }
    _, err = db.Exec(context.Background(), "INSERT INTO merch (name, price) VALUES ('sword', 100)")
    if err != nil {
        t.Fatalf("failed to insert test data: %v", err)
    }

    return db
}

func TestPurchaseItem_Success(t *testing.T) {
    db := setupMerchTestDB(t)
    defer db.Close()

    repo := repositories.NewMerchRepository(db)
    service := services.NewMerchService(repo)

    err := service.PurchaseItem(context.Background(), 1, "sword")
    assert.NoError(t, err)

    var userCoins int
    err = db.QueryRow(context.Background(), "SELECT coins FROM users WHERE id = $1", 1).Scan(&userCoins)
    assert.NoError(t, err)
    assert.Equal(t, 100, userCoins)

    var inventoryQuantity int
    err = db.QueryRow(context.Background(), "SELECT quantity FROM inventory WHERE user_id = $1 AND item_name = $2", 1, "sword").Scan(&inventoryQuantity)
    assert.NoError(t, err)
    assert.Equal(t, 1, inventoryQuantity)
}

func TestPurchaseItem_NotEnoughCoins(t *testing.T) {
    db := setupMerchTestDB(t)
    defer db.Close()

    repo := repositories.NewMerchRepository(db)
    service := services.NewMerchService(repo)

    err := service.PurchaseItem(context.Background(), 1, "sword")
    assert.NoError(t, err)

    err = service.PurchaseItem(context.Background(), 1, "sword")
    assert.Nil(t, err)
    assert.Equal(t, "not enough coins", err.Error())
}

func TestPurchaseItem_ItemNotFound(t *testing.T) {
    db := setupMerchTestDB(t)
    defer db.Close()

    repo := repositories.NewMerchRepository(db)
    service := services.NewMerchService(repo)

    err := service.PurchaseItem(context.Background(), 1, "shield")
    assert.Error(t, err)
    assert.Equal(t, "item not found", err.Error())
}

func TestPurchaseItem_UserNotFound(t *testing.T) {
    db := setupMerchTestDB(t)
    defer db.Close()

    repo := repositories.NewMerchRepository(db)
    service := services.NewMerchService(repo)

    err := service.PurchaseItem(context.Background(), 999, "sword")
    assert.Error(t, err)
    assert.Equal(t, "user not found", err.Error())
}

func TestPurchaseItem_TransactionError(t *testing.T) {
    db := setupMerchTestDB(t)
    defer db.Close()

    repo := repositories.NewMerchRepository(db)
    service := services.NewMerchService(repo)

    db.Close()

    err := service.PurchaseItem(context.Background(), 1, "sword")
    assert.Error(t, err) 
    assert.Contains(t, err.Error(), "closed pool")
}
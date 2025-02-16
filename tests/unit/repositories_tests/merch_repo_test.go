package repositories_tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/RadicalIguana/avito-shop/internal/repositories"
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
            user_id SERIAL REFERENCES users(id),
            item_name VARCHAR(255) REFERENCES merch(name),
            quantity INT NOT NULL DEFAULT 0,
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


func TestGetUserById_Success(t *testing.T) {
    db := setupMerchTestDB(t)
    defer db.Close()

    repo := repositories.NewMerchRepository(db)

    user, err := repo.GetUserById(context.Background(), 1)
    assert.NoError(t, err)
    assert.Equal(t, 1, user.Id)
    assert.Equal(t, 200, user.Coins)
}

func TestGetItemByName_Success(t *testing.T) {
    db := setupMerchTestDB(t)
    defer db.Close()

    repo := repositories.NewMerchRepository(db)

    item, err := repo.GetItemByName(context.Background(), "sword")
    assert.NoError(t, err)
    assert.Equal(t, "sword", item.Name)
    assert.Equal(t, 100, item.Price)
}

func TestUpdateUserBalance_Success(t *testing.T) {
    db := setupMerchTestDB(t)
    defer db.Close()

    repo := repositories.NewMerchRepository(db)

    err := repo.UpdateUserBalance(context.Background(), 1, 100)
    assert.NoError(t, err)

    var coins int
    err = db.QueryRow(context.Background(), "SELECT coins FROM users WHERE id = $1", 1).Scan(&coins)
    assert.NoError(t, err)
    assert.Equal(t, 100, coins)
}

func TestAddOrUpdateItemToInventory_Success(t *testing.T) {
    db := setupMerchTestDB(t)
    defer db.Close()

    repo := repositories.NewMerchRepository(db)

    err := repo.AddOrUpdateItemToInventory(context.Background(), 1, "sword")
    assert.NoError(t, err)

    var quantity int
    err = db.QueryRow(context.Background(), "SELECT quantity FROM inventory WHERE user_id = $1 AND item_name = $2", 1, "sword").Scan(&quantity)
    assert.NoError(t, err)
    assert.Equal(t, 1, quantity)
}
package services_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/RadicalIguana/avito-shop/internal/models"
	"github.com/RadicalIguana/avito-shop/internal/repositories"
	"github.com/RadicalIguana/avito-shop/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setupDB() (*pgxpool.Pool, error) {
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
    pool, err := pgxpool.New(context.Background(), connString)
    if err != nil {
        return nil, err
    }
    return pool, nil
}

func seedDB(pool *pgxpool.Pool) error {
    _, err := pool.Exec(context.Background(), `
        DROP TABLE IF EXISTS transfers;
        DROP TABLE IF EXISTS inventory;
        DROP TABLE IF EXISTS users;

        CREATE TABLE users (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255),
            password VARCHAR(255),
            coins INT
        );

        CREATE TABLE inventory (
            id SERIAL PRIMARY KEY,
            user_id INT,
            item_name VARCHAR(255),
            quantity INT,
            FOREIGN KEY (user_id) REFERENCES users(id)
        );

        CREATE TABLE transfers (
            id SERIAL PRIMARY KEY,
            from_user INT,
            to_user INT,
            amount INT,
            FOREIGN KEY (from_user) REFERENCES users(id),
            FOREIGN KEY (to_user) REFERENCES users(id)
        );

        INSERT INTO users (id, name, password, coins) VALUES (1, 'TestUser1', 'password', 100);
		INSERT INTO users (id, name, password, coins) VALUES (2, 'TestUser2', 'password', 100);
		INSERT INTO users (id, name, password, coins) VALUES (3, 'TestUser3', 'password', 100);
		INSERT INTO users (id, name, password, coins) VALUES (4, 'TestUser4', 'password', 100);

        INSERT INTO inventory (user_id, item_name, quantity) VALUES (1, 'Item1', 1);
        INSERT INTO inventory (user_id, item_name, quantity) VALUES (1, 'Item2', 2);
        INSERT INTO transfers (from_user, to_user, amount) VALUES (1, 2, 10);
        INSERT INTO transfers (from_user, to_user, amount) VALUES (3, 1, 5);
    `)
    return err
}

func TestGetUserInfo(t *testing.T) {
    pool, err := setupDB()
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    defer pool.Close()

    err = seedDB(pool)
    if err != nil {
        t.Fatalf("Failed to seed database: %v", err)
    }

    repo := repositories.NewUserInfoRepository(pool)
    service := services.NewUserInfoService(repo)

    ctx := context.Background()
    userId := 1

    expectedInventory := []models.InventoryItem{
        {Name: "Item1", Quantiry: 1},
        {Name: "Item2", Quantiry: 2},
    }
    expectedCoinHistory := models.CoinHistory{
        Received: []models.Transfer{
            {FromUser: 3, Amount: 5},
        },
        Sent: []models.Transfer{
            {ToUser: 2, Amount: 10},
        },
    }

    userInfo, err := service.GetUserInfo(ctx, userId)

    assert.NoError(t, err)
    assert.Equal(t, 100, userInfo.Coins)
    assert.ElementsMatch(t, expectedInventory, userInfo.Inventory)
    assert.ElementsMatch(t, expectedCoinHistory.Received, userInfo.CoinHistory.Received)
    assert.ElementsMatch(t, expectedCoinHistory.Sent, userInfo.CoinHistory.Sent)
}
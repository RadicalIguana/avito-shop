package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var responseDataPurchase struct {
	Token string `json:"token"`
}

func setupPurchaseDB(t *testing.T) {
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
	assert.NoError(t, err)

	_, err = db.Exec(context.Background(), `
		DROP TABLE IF EXISTS inventory;
		DROP TABLE IF EXISTS transfers;
		DROP TABLE IF EXISTS merch;
		DROP TABLE IF EXISTS users;
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			password VARCHAR(25) NOT NULL,
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
    assert.NoError(t, err)

	_, err = db.Exec(context.Background(), "TRUNCATE TABLE users, merch, inventory CASCADE")
	assert.NoError(t, err)

	_, err = db.Exec(context.Background(), "INSERT INTO users (id, username, password, coins) VALUES ($1, $2, $3, $4)", 1, "testuser", "testpass", 100)
	assert.NoError(t, err)

	_, err = db.Exec(context.Background(), "INSERT INTO merch (name, price) VALUES ($1, $2)", "pen", 50)
	assert.NoError(t, err)
}

func TestPurchaseItem(t *testing.T) {
	setupPurchaseDB(t)

	client := &http.Client{}

	authReq := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	authBody, _ := json.Marshal(authReq)
	authReqHTTP, _ := http.NewRequest("POST", "http://localhost:8080/api/auth", bytes.NewBuffer(authBody))
	authReqHTTP.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(authReqHTTP)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &responseDataPurchase); err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	token := responseDataPurchase.Token

	purchaseReq, _ := http.NewRequest("GET", "http://localhost:8080/api/buy/book", bytes.NewBuffer(nil))
	purchaseReq.Header.Set("Authorization", "Bearer " + token)
	purchaseResp, err := client.Do(purchaseReq)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, purchaseResp.StatusCode)
}

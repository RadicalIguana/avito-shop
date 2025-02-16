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

	"github.com/RadicalIguana/avito-shop/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var responseDataSend struct {
	Token string `json:"token"`
}

func setupSendDB(t *testing.T) {
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
            coins INTEGER NOT NULL DEFAULT 1000
        );
		CREATE TABLE transfers (
			id SERIAL PRIMARY KEY,
			from_user SERIAL REFERENCES users(id),
			to_user SERIAL REFERENCES users(id),
			amount INT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
    `)
    assert.NoError(t, err)

	_, err = db.Exec(context.Background(), "TRUNCATE TABLE users, merch, inventory CASCADE")
	assert.NoError(t, err)

	_, err = db.Exec(context.Background(), "INSERT INTO users (id, username, password, coins) VALUES ($1, $2, $3)", 1, "first", "password")
	assert.NoError(t, err)

	_, err = db.Exec(context.Background(), "INSERT INTO users (id, username, password, coins) VALUES ($1, $2, $3)", 2, "second", "password")
	assert.NoError(t, err)
}

func TestSendCoin(t *testing.T) {
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

	sendReq := models.TransferRequest{
		ToUser: 1,
		Amount: 100,
	}
	sendBody, _ := json.Marshal(sendReq)
	purchaseReq, _ := http.NewRequest("GET", "http://localhost:8080/api/buy/book", bytes.NewBuffer(sendBody))
	purchaseReq.Header.Set("Authorization", "Bearer " + token)
	purchaseResp, err := client.Do(purchaseReq)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, purchaseResp.StatusCode)
}

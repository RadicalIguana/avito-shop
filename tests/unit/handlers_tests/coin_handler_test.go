package handlers_tests

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
	"github.com/gin-gonic/gin"
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

func TestSendCoins(t *testing.T) {
    db := setupCoinTestDB(t)
    defer db.Close()

    gin.SetMode(gin.TestMode)

	requestBody, _ := json.Marshal(map[string]string{
		"username": "MainTestUser",
		"password": "MainTestUserPassword",
	})
	url := "http://localhost:8080/api/auth"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var responseData struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(body, &responseData); err != nil {
		panic(err)
	}

    sendUrl := "http://localhost:8080/api/sendCoin"
    reqBody := models.TransferRequest{
        ToUser: 1,
        Amount: 100,
    }
    reqBodyBytes, _ := json.Marshal(reqBody)
    sendReq, err := http.NewRequest(http.MethodPost, sendUrl, bytes.NewBuffer(reqBodyBytes))
    sendReq.Header.Set("Content-Type", "application/json")
    sendReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", responseData.Token))

    sendResp, err := client.Do(sendReq)
	if err != nil {
		panic(err)
	}
	defer sendResp.Body.Close()
	infoBody, err := io.ReadAll(sendResp.Body)
	if err != nil {
		panic(err)
	}
    var response map[string]interface{}
    if err := json.Unmarshal(infoBody, &response); err != nil {
        panic(err)
    }

    responseJSON, err := json.Marshal(response)
    if err != nil {
        panic(err)
    }

    fmt.Println(sendResp.Body)
    
    assert.Equal(t, http.StatusBadRequest, sendResp.StatusCode)
    assert.JSONEq(t, `{"error":"insufficient funds"}`, string(responseJSON))

    var user1Coins, user2Coins int
    err = db.QueryRow(context.Background(), "SELECT coins FROM users WHERE id = $1", 1).Scan(&user1Coins)
    assert.NoError(t, err)
    assert.Equal(t, 200, user1Coins)

    err = db.QueryRow(context.Background(), "SELECT coins FROM users WHERE id = $1", 2).Scan(&user2Coins)
    assert.NoError(t, err)
    assert.Equal(t, 50, user2Coins)
}
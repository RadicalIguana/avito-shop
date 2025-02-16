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

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setupMerchTestDB(t *testing.T) *pgxpool.Pool {
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
        DROP TABLE IF EXISTS inventory;
        DROP TABLE IF EXISTS transfers;
        DROP TABLE IF EXISTS users;
        DROP TABLE IF EXISTS merch;

        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(255),
            password VARCHAR(255),
            coins INTEGER NOT NULL DEFAULT 1000
        );
        CREATE TABLE IF NOT EXISTS merch (
            name TEXT PRIMARY KEY,
            price INTEGER NOT NULL
        );
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

    // Очистка таблиц перед тестами
    _, err = db.Exec(context.Background(), "TRUNCATE users, merch, inventory RESTART IDENTITY CASCADE")
    if err != nil {
        t.Fatalf("failed to truncate tables: %v", err)
    }

    // Добавление тестовых данных
    _, err = db.Exec(context.Background(), "INSERT INTO users (id, username, password, coins) VALUES (1, 'first', 'first', 1000)")
    if err != nil {
        t.Fatalf("failed to insert test data: %v", err)
    }
    _, err = db.Exec(context.Background(), "INSERT INTO merch (name, price) VALUES ('pen', 100)")
    _, err = db.Exec(context.Background(), "INSERT INTO merch (name, price) VALUES ('book', 100000)")

    return db
}

func TestPurchaseItem_Success(t *testing.T) {
    db := setupMerchTestDB(t)
    defer db.Close()

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

    buyUrl := "http://localhost:8080/api/buy/pen"
    buyReq, err := http.NewRequest(http.MethodGet, buyUrl, nil)
    buyReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", responseData.Token))

    buyResp, err := client.Do(buyReq)
	if err != nil {
		panic(err)
	}
	defer buyResp.Body.Close()
	infoBody, err := io.ReadAll(buyResp.Body)
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

    // Проверяем результат
    assert.Equal(t, http.StatusOK, buyResp.StatusCode)
    assert.JSONEq(t, `{"message":"Item purchased successfully"}`, string(responseJSON))
}

func TestPurchaseItem_NotEnoughCoins(t *testing.T) {
    // Настройка тестовой базы данных
    db := setupMerchTestDB(t)
    defer db.Close()

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
    
    buyUrl := "http://localhost:8080/api/buy/book"
    buyReq, err := http.NewRequest(http.MethodGet, buyUrl, nil)
    buyReq.Header.Set("Content-Type", "application/json")
    buyReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", responseData.Token))

    buyResp, err := client.Do(buyReq)
	if err != nil {
		panic(err)
	}
	defer buyResp.Body.Close()
	infoBody, err := io.ReadAll(buyResp.Body)
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

    assert.Equal(t, http.StatusBadRequest, buyResp.StatusCode)
    assert.Contains(t, string(responseJSON), "not enough coins")
}
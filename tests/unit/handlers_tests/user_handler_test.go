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

func setupDB() (*pgxpool.Pool, error) {
	if err := godotenv.Load(); err != nil {
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
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(255),
            password VARCHAR(255),
            coins INT NOT NULL DEFAULT 1000 
        );

        CREATE TABLE IF NOT EXISTS inventory (
            id SERIAL PRIMARY KEY,
            user_id INT,
            item_name VARCHAR(255),
            quantity INT,
            FOREIGN KEY (user_id) REFERENCES users(id)
        );

        CREATE TABLE IF NOT EXISTS transfers (
            id SERIAL PRIMARY KEY,
            from_user INT,
            to_user INT,
            amount INT,
            FOREIGN KEY (from_user) REFERENCES users(id),
            FOREIGN KEY (to_user) REFERENCES users(id)
        );

        INSERT INTO users (username, password, coins) VALUES ('TestUser', 'password', 100);
        INSERT INTO users (username, password, coins) VALUES ('User2', 'password2', 200);
        INSERT INTO users (username, password, coins) VALUES ('User3', 'password3', 300);

        INSERT INTO inventory (user_id, item_name, quantity) VALUES (1, 'pen', 1);
        INSERT INTO inventory (user_id, item_name, quantity) VALUES (1, 'book', 2);

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

	infoUrl := "http://localhost:8080/api/info"
	infoReq, _ := http.NewRequest("GET", infoUrl, bytes.NewBuffer(nil))
	infoReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", responseData.Token))

	infoResp, err := client.Do(infoReq)
	if err != nil {
		panic(err)
	}
	defer infoResp.Body.Close()

	infoBody, err := io.ReadAll(infoResp.Body)
	if err != nil {
		panic(err)
	}

	var infoResponse models.UserResponse
	if err := json.Unmarshal(infoBody, &infoResponse); err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusOK, infoResp.StatusCode)

	var response models.UserResponse

	fmt.Println(infoResponse)

	expectedResponse := &models.UserResponse{
		Coins: 0,
		Inventory: []models.InventoryItem{},
		CoinHistory: models.CoinHistory{
			Received: []models.Transfer{},
			Sent: []models.Transfer{},
		},
	}

	assert.Equal(t, expectedResponse.Coins, infoResponse.Coins)
	assert.ElementsMatch(t, expectedResponse.Inventory, response.Inventory)
	assert.ElementsMatch(t, expectedResponse.CoinHistory.Received, response.CoinHistory.Received)
	assert.ElementsMatch(t, expectedResponse.CoinHistory.Sent, response.CoinHistory.Sent)
}

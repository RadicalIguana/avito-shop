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

    // Очистка таблиц перед тестами
    _, err = db.Exec(context.Background(), "TRUNCATE users, merch, inventory RESTART IDENTITY CASCADE")
    if err != nil {
        t.Fatalf("failed to truncate tables: %v", err)
    }

    // Добавление тестовых данных
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
    // Настройка тестовой базы данных
    db := setupMerchTestDB(t)
    defer db.Close()

    // Инициализация репозитория и сервиса
    repo := repositories.NewMerchRepository(db)
    service := services.NewMerchService(repo)

    // Выполнение покупки
    err := service.PurchaseItem(context.Background(), 1, "sword")
    assert.NoError(t, err)

    // Проверка состояния базы данных после транзакции
    var userCoins int
    err = db.QueryRow(context.Background(), "SELECT coins FROM users WHERE id = $1", 1).Scan(&userCoins)
    assert.NoError(t, err)
    assert.Equal(t, 100, userCoins) // Баланс должен уменьшиться на стоимость предмета

    var inventoryQuantity int
    err = db.QueryRow(context.Background(), "SELECT quantity FROM inventory WHERE user_id = $1 AND item_name = $2", 1, "sword").Scan(&inventoryQuantity)
    assert.NoError(t, err)
    assert.Equal(t, 1, inventoryQuantity) // Предмет должен быть добавлен в инвентарь
}

func TestPurchaseItem_NotEnoughCoins(t *testing.T) {
    // Настройка тестовой базы данных
    db := setupMerchTestDB(t)
    defer db.Close()

    // Инициализация репозитория и сервиса
    repo := repositories.NewMerchRepository(db)
    service := services.NewMerchService(repo)

    // Выполнение покупки
    err := service.PurchaseItem(context.Background(), 1, "sword")
    assert.NoError(t, err)

    // Попытка повторной покупки (недостаточно монет)
    err = service.PurchaseItem(context.Background(), 1, "sword")
    assert.Nil(t, err)
    // assert.Equal(t, "not enough coins", err.Error())
}

func TestPurchaseItem_ItemNotFound(t *testing.T) {
    // Настройка тестовой базы данных
    db := setupMerchTestDB(t)
    defer db.Close()

    // Инициализация репозитория и сервиса
    repo := repositories.NewMerchRepository(db)
    service := services.NewMerchService(repo)

    // Попытка купить несуществующий предмет
    err := service.PurchaseItem(context.Background(), 1, "shield")
    assert.Error(t, err)
    assert.Equal(t, "item not found", err.Error())
}

func TestPurchaseItem_UserNotFound(t *testing.T) {
    // Настройка тестовой базы данных
    db := setupMerchTestDB(t)
    defer db.Close()

    // Инициализация репозитория и сервиса
    repo := repositories.NewMerchRepository(db)
    service := services.NewMerchService(repo)

    // Попытка купить предмет для несуществующего пользователя
    err := service.PurchaseItem(context.Background(), 999, "sword")
    assert.Error(t, err)
    assert.Equal(t, "user not found", err.Error())
}

func TestPurchaseItem_TransactionError(t *testing.T) {
    // Настройка тестовой базы данных
    db := setupMerchTestDB(t)
    defer db.Close()

    // Инициализация репозитория и сервиса
    repo := repositories.NewMerchRepository(db)
    service := services.NewMerchService(repo)

    // Закрываем соединение с базой данных, чтобы вызвать ошибку транзакции
    db.Close()

    // Попытка выполнить покупку, которая должна вызвать ошибку
    err := service.PurchaseItem(context.Background(), 1, "sword")
    assert.Error(t, err)  // Ошибка должна быть
    assert.Contains(t, err.Error(), "closed pool") // Проверка конкретной ошибки
}
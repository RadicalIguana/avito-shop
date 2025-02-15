package repositories

import (
	"context"

	"github.com/RadicalIguana/avito-shop/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserInfoRepository struct {
	db *pgxpool.Pool
}

func NewUserInfoRepository(db *pgxpool.Pool) *UserInfoRepository {
	return &UserInfoRepository{db: db}
}

// TODO: Может убрать эту функцию?
func (r *UserInfoRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
    return r.db.Begin(ctx)
}

func (r *UserInfoRepository) GetUserCoins(ctx context.Context, userId int) (int, error) {
	var coins int
	query := "SELECT coins FROM users WHERE id = $1"

	err := r.db.QueryRow(ctx, query, userId).Scan(&coins)
	if err != nil {
		return 0, err
	}
	return coins, nil
}

func (r *UserInfoRepository) GetUserInventory(ctx context.Context, userId int) ([]models.InventoryItem, error) {
	query := `
		SELECT item_name, quantity
		FROM inventory
		WHERE user_id = $1
	`
	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventory []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.Name, &item.Quantiry); err != nil {
			return nil, err
		}
		inventory = append(inventory, item)
	}

	if inventory == nil {
		inventory = []models.InventoryItem{}
	}

	return inventory, nil
}

func (r *UserInfoRepository) GetCoinTransfers(ctx context.Context, userId int) (models.CoinHistory, error) {
	var history models.CoinHistory

	sentRows, err := r.db.Query(
		ctx, 
		"SELECT to_user, amount FROM transfers WHERE from_user = $1", 
		userId,
	)
	if err != nil {
		return history, err
	}
	defer sentRows.Close()

	for sentRows.Next() {
		var transfer models.Transfer
		if err := sentRows.Scan(&transfer.ToUser, &transfer.Amount); err!= nil {
			return history, err
		}
		history.Sent = append(history.Sent, transfer)
	}

	receivedRows, err := r.db.Query(
		ctx,
		"SELECT from_user, amount FROM transfers WHERE to_user = $1",
		userId,
	)
	if err != nil {
        return history, err
    }
	defer receivedRows.Close()

	for receivedRows.Next() {
		var transfer models.Transfer
		if err := receivedRows.Scan(&transfer.FromUser, &transfer.Amount); err!= nil {
			return history, err
		}
		history.Received = append(history.Received, transfer)
	}

	if history.Sent == nil {
		history.Sent = []models.Transfer{}
	}

	if history.Received == nil {
		history.Received = []models.Transfer{}
	}

	return history, nil
}


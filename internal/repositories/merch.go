package repositories

import (
	"context"
	"fmt"

	"github.com/RadicalIguana/avito-shop/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MerchRepository struct {
	db *pgxpool.Pool
}

func NewMerchRepository(db *pgxpool.Pool) *MerchRepository {
    return &MerchRepository{db: db}
}

func (r *MerchRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
    return r.db.Begin(ctx)
}

func (r *MerchRepository) GetUserById(ctx context.Context, userID int) (*models.User, error) {
	query := `
		SELECT id, coins
		FROM users
		WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, userID)
	var user models.User
	if err := row.Scan(&user.Id, &user.Coins); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func (r *MerchRepository) GetItemByName(ctx context.Context, name string) (*models.Merch, error) {
	query := `
		SELECT name, price
		FROM merch
		WHERE name = $1
	`
	row := r.db.QueryRow(ctx, query, name)
	var merch models.Merch
	if err := row.Scan(&merch.Name, &merch.Price); err != nil {
        return nil, fmt.Errorf("item not found: %w", err)
    }
	return &merch, nil
}

func (r *MerchRepository) UpdateUserBalance(ctx context.Context, userID int, newCoins int) error {
	query := `UPDATE users SET coins = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, newCoins, userID)
	return err
}

func (r *MerchRepository) AddOrUpdateItemToInventory(ctx context.Context, userID int, itemName string) error {
	query := `
		INSERT INTO inventory (user_id, item_name, quantity) 
		VALUES ($1, $2, 1)
		ON CONFLICT (user_id, item_name) DO UPDATE
		SET quantity = inventory.quantity + 1
	`
	_, err := r.db.Exec(ctx, query, userID, itemName)
	return err
}
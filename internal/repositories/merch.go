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

func (r *MerchRepository) GetUserById(ctx context.Context, userId string) (*models.User, error) {
	query := `
		SELECT id, coins
		FROM users
		WHERE id = $1
		FOR UPDATE
	`
	row := r.db.QueryRow(ctx, query, userId)
	var user models.User
	if err := row.Scan(&user.Id, &user.Coins); err != nil {
		// TODO: в чем разница errors и fmt.Errorf? Что использовать?
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

// Получение стоимости предмета
// TODO: Почему мы возвращаем *models.Merch, а не models.Merch?
func (r *MerchRepository) GetItemByName(ctx context.Context, name string) (*models.Merch, error) {
	query := `
		SELECT name, price
		FROM merch
		WHERE name = $1
		FOR UPDATE
	`
	row := r.db.QueryRow(ctx, query, name)
	var merch models.Merch
	if err := row.Scan(&merch.Name, &merch.Price); err != nil {
        return nil, fmt.Errorf("item not found: %w", err)
    }
	return &merch, nil
}

func (r *MerchRepository) UpdateUserBalance(ctx context.Context, userId string, newCoins int) error {
	query := `UPDATE users SET coins = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, newCoins, userId)
	return err
}

func (r *MerchRepository) AddOrUpdateItemToInventory(ctx context.Context, userId string, itemName string) error {
	query := `
		INSERT INTO inventory (user_id, item_name, quantity) 
		VALUES ($1, $2, 1)
		ON CONFLICT (user_id, item_name) DO UPDATE
		SET quantity = inventory.quantity + 1
	`
	_, err := r.db.Exec(ctx, query, userId, itemName)
	return err
}
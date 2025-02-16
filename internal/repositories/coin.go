package repositories

import (
    "context"
    "fmt"
    "github.com/RadicalIguana/avito-shop/internal/models"
    "github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
)

type CoinRepository struct {
    db *pgxpool.Pool
}

func NewCoinRepository(db *pgxpool.Pool) *CoinRepository {
    return &CoinRepository{db: db}
}

func (r *CoinRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
    return r.db.Begin(ctx)
}

func (r *CoinRepository) GetUserForUpdate(ctx context.Context, userID int) (*models.UserCoin, error) {
    query := `
        SELECT id, coins
        FROM users
        WHERE id = $1
        FOR UPDATE
    `
    row := r.db.QueryRow(ctx, query, userID)
    var user models.UserCoin
    if err := row.Scan(&user.ID, &user.Coins); err != nil {
        return nil, fmt.Errorf("user not found: %w", err)
    }
    return &user, nil
}

func (r *CoinRepository) UpdateBalance(ctx context.Context, userID, newCoins int) error {
    query := `UPDATE users SET coins = $1 WHERE id = $2`
    _, err := r.db.Exec(ctx, query, newCoins, userID)
    return err
}

func (r *CoinRepository) CreateTransfer(ctx context.Context, transfer *models.Transfer) error {
    query := `
        INSERT INTO transfers (from_user, to_user, amount)
        VALUES ($1, $2, $3)
    `
    _, err := r.db.Exec(ctx, query,
        transfer.FromUser,
        transfer.ToUser,
        transfer.Amount,
    )
    return err
}
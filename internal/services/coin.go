package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/RadicalIguana/avito-shop/internal/models"
	"github.com/RadicalIguana/avito-shop/internal/repositories"
	"github.com/google/uuid"
)

type CoinService struct {
    repo *repositories.CoinRepository
}

func NewCoinService(repo *repositories.CoinRepository) *CoinService {
    return &CoinService{repo: repo}
}

func (s *CoinService) TransferCoins(ctx context.Context, fromUserID, toUserID string, amount int) error {
    // Валидация
    if amount <= 0 {
        return errors.New("amount must be positive")
    }

    // Начало транзакции
    tx, err := s.repo.BeginTx(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

	fromUser, err := s.repo.GetUserForUpdate(ctx, fromUserID)
	if err != nil {
        return err
    }

	toUser, err := s.repo.GetUserForUpdate(ctx, toUserID)
	if err != nil {
        return err
    }

	if fromUser.Coins < amount {
		return errors.New("insufficient funds")
	}

	if err := s.repo.UpdateBalance(ctx, fromUserID, fromUser.Coins-amount); err != nil {
		return err
	}

	if err := s.repo.UpdateBalance(ctx, toUserID, toUser.Coins+amount); err != nil {
        return err
    }

	transfer := &models.Transfer{
		ID: uuid.New().String(),
		FromUser: fromUserID,
		ToUser: toUserID,
        Amount: amount,
	}
	if err := s.repo.CreateTransfer(ctx, transfer); err != nil {
		return err
	}
	
	if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("transaction commit error: %w", err)
    }

    return nil
}
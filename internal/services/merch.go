package services

import (
	"context"
	"errors"

	"github.com/RadicalIguana/avito-shop/internal/repositories"
)

type MerchService struct {
	repo *repositories.MerchRepository
}

func NewMerchService(repo *repositories.MerchRepository) *MerchService {
    return &MerchService{repo: repo}
}

func (s *MerchService) PurchaseItem(ctx context.Context, userID int, itemName string) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
        return err
    }
	defer tx.Rollback(ctx)

	user, err := s.repo.GetUserById(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	item, err := s.repo.GetItemByName(ctx, itemName)
	if err != nil {
		return errors.New("item not found")
	}

	if user.Coins < item.Price {
		return errors.New("not enough coins")
	}

	if err := s.repo.UpdateUserBalance(ctx, userID, user.Coins-item.Price); err!= nil {
		return err
	}

	if err := s.repo.AddOrUpdateItemToInventory(ctx, userID, item.Name); err != nil {
        return err
    }

    return tx.Commit(ctx)
}

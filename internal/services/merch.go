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

	// Получение пользователя
	user, err := s.repo.GetUserById(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Получение стоимости покупаемого предмета
	item, err := s.repo.GetItemByName(ctx, itemName)
	if err != nil {
		return errors.New("item not found")
	}

	if user.Coins < item.Price {
		return errors.New("not enough coins")
	}

    // TODO: Может изменить уменьшение монет?
    // Обновление монет
	if err := s.repo.UpdateUserBalance(ctx, userID, user.Coins-item.Price); err!= nil {
		return err
	}

	// Добавить предмет в инвентарь
	if err := s.repo.AddOrUpdateItemToInventory(ctx, userID, item.Name); err != nil {
        return err
    }

    return tx.Commit(ctx)
}

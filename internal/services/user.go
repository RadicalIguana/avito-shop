package services

import (
	"context"
	// "errors"

	"github.com/RadicalIguana/avito-shop/internal/models"
	"github.com/RadicalIguana/avito-shop/internal/repositories"
)

type UserInfoService struct {
	repo *repositories.UserInfoRepository
}

func NewUserInfoService(repo *repositories.UserInfoRepository) *UserInfoService {
    return &UserInfoService{repo: repo}
}

func (s *UserInfoService) GetUserInfo(ctx context.Context, userId int) (*models.UserResponse, error) {
	// get coins
	coins, err := s.repo.GetUserCoins(ctx, userId)
	if err != nil {
		return nil, err
	}

	// get inventory
	inventory, err := s.repo.GetUserInventory(ctx, userId)
	if err != nil {
        return nil, err
    }

	// get history
	coinHistory, err := s.repo.GetCoinTransfers(ctx, userId)
	if err != nil {
		return nil, err
	}

	// format and return the result
	response := &models.UserResponse {
		Coins: coins,
		Inventory: inventory,
		CoinHistory: coinHistory,
	}

	return response, nil
}
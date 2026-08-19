package service

import (
	"context"

	"github.com/waves2k/go-wallet/monolith/internal/wallet/model"
	walRepo "github.com/waves2k/go-wallet/monolith/internal/wallet/repository"
)

type WalletService interface {
	GetByUserID(ctx context.Context, userID string) (*model.Wallet, error)
}

type walletService struct {
	walRepo walRepo.WalletRepository
}

func NewWalletService(walRepo walRepo.WalletRepository) WalletService {
	return &walletService{
		walRepo: walRepo,
	}
}

func (s *walletService) GetByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	const op = "walletService.GetByUserID"

	wallet, err := s.walRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

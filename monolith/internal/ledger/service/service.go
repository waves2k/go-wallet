package service

import (
	"context"

	"github.com/waves2k/go-wallet/monolith/internal/ledger/model"
	lRepo "github.com/waves2k/go-wallet/monolith/internal/ledger/repository"
	wRepo "github.com/waves2k/go-wallet/monolith/internal/wallet/repository"
)

type LedgerService interface {
	RecordWalletBalance(ctx context.Context, userID string) (bool, float64, float64, error)
	GetMutationHistory(ctx context.Context, userID string) ([]model.LedgerEntry, error)
}

type ledgerService struct {
	ledRepo lRepo.LedgerRepository
	walRepo wRepo.WalletRepository
}

func NewLedgerService(lRepo lRepo.LedgerRepository, wRepo wRepo.WalletRepository) LedgerService {
	return &ledgerService{
		ledRepo: lRepo,
		walRepo: wRepo,
	}
}

func (s *ledgerService) RecordWalletBalance(ctx context.Context, userID string) (bool, float64, float64, error) {
	wallet, err := s.walRepo.GetByUserID(ctx, userID)
	if err != nil {
		return false, 0, 0, err
	}

	ledgerBalance, err := s.ledRepo.GetBallanceByWalletID(ctx, wallet.ID)
	if err != nil {
		return false, 0, 0, err
	}

	isConsist := wallet.Balance == ledgerBalance

	return isConsist, wallet.Balance, ledgerBalance, nil
}

func (s *ledgerService) GetMutationHistory(ctx context.Context, userID string) ([]model.LedgerEntry, error) {
	wallet, err := s.walRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	entries, err := s.ledRepo.GetEntriesByWalletID(ctx, wallet.ID)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

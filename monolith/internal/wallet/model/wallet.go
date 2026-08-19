package model

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Balance   float64   `db:"balance"`
	Currency  string    `db:"currency"`
	Status    string    `db:"status"`
	Version   int       `db:"version"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func New(userID string) (*Wallet, error) {
	now := time.Now()

	return &Wallet{
		ID:        uuid.New().String(),
		UserID:    userID,
		Balance:   0.00,
		Currency:  "RUB",
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil

}

type WalletResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (w *Wallet) ToResponse() *WalletResponse {
	return &WalletResponse{
		ID:        w.ID,
		UserID:    w.UserID,
		Balance:   w.Balance,
		Currency:  w.Currency,
		Status:    w.Status,
		Version:   w.Version,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}
}

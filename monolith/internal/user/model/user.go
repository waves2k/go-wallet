package model

import (
	"time"

	"github.com/google/uuid"
	walletModel "github.com/waves2k/go-wallet/monolith/internal/wallet/model"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string    `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	DeletedAt    time.Time `db:"deleted_at"`
}

type UserWithWallet struct {
	User   *User
	Wallet *walletModel.Wallet
}

func New(username, email, password string) (*User, error) {

	// Validation

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &User{
		ID:           uuid.New().String(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
		DeletedAt:    time.Time{},
	}, nil
}

func NewUserWithWallet(u *User, w *walletModel.Wallet) *UserWithWallet {
	return &UserWithWallet{
		User:   u,
		Wallet: w,
	}
}

type UserWithWalletResponse struct {
	User   UserResponse               `json:"user"`
	Wallet walletModel.WalletResponse `json:"wallet"`
}

type UserResponse struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// type CreateUserRequest struct {
// 	Username string `json:"username" validate:"required,min=2,max=100"`
// 	Email    string `json:"email" validate:"required,email"`
// 	Password string `json:"password" validate:"required,min=8"`
// }

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func (uw *UserWithWallet) ToResponse() *UserWithWalletResponse {
	return &UserWithWalletResponse{
		User:   uw.User.ToResponse(),
		Wallet: *uw.Wallet.ToResponse(),
	}
}

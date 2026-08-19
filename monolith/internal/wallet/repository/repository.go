package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/waves2k/go-wallet/monolith/internal/wallet/model"
)

const (
	walletsTableName = "wallets"
)

var (
	ErrWalletNotFound  = errors.New("wallet not found")
	ErrDuplicate       = errors.New("duplicate error")
	ErrInternalFailure = errors.New("internal server error")
)

type WalletRepository interface {
	CreateTx(ctx context.Context, tx pgx.Tx, w *model.Wallet) error
	GetByUserID(ctx context.Context, userID string) (*model.Wallet, error)
}

type postgresqlWalletRepository struct {
	db *pgxpool.Pool
}

func NewPostgresqlUserRepository(db *pgxpool.Pool) WalletRepository {
	return &postgresqlWalletRepository{
		db: db,
	}
}

func (r *postgresqlWalletRepository) CreateTx(ctx context.Context, tx pgx.Tx, w *model.Wallet) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (id, user_id, balance, currency, status, version)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, walletsTableName)

	_, err := tx.Exec(
		ctx, query,
		w.ID, w.UserID, w.Balance, w.Currency, w.Status, w.Version,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return fmt.Errorf("%v: %w", ErrDuplicate, err)
		}
		return fmt.Errorf("%v: %w", ErrInternalFailure, err)
	}

	return nil
}

func (r *postgresqlWalletRepository) GetByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	query := fmt.Sprintf(`
		SELECT id, user_id, balance, currency, status, version, created_at, updated_at
		FROM %s
		WHERE user_id = $1
	`, walletsTableName)

	row := r.db.QueryRow(ctx, query, userID)

	wallet, err := scanIntoWallet(row)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, ErrWalletNotFound
		}
		return nil, ErrInternalFailure
	}

	return wallet, err
}

func scanIntoWallet(row pgx.Row) (*model.Wallet, error) {
	var (
		wallet model.Wallet
		err    error
	)
	err = row.Scan(
		&wallet.ID, &wallet.UserID, &wallet.Balance, &wallet.Currency,
		&wallet.Status, &wallet.Version, &wallet.CreatedAt, &wallet.UpdatedAt,
	)
	return &wallet, err
}

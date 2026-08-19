package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	userModel "github.com/waves2k/go-wallet/monolith/internal/user/model"
	walletModel "github.com/waves2k/go-wallet/monolith/internal/wallet/model"
	walletRepo "github.com/waves2k/go-wallet/monolith/internal/wallet/repository"
)

const (
	usersTableName = "accounts"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrDuplicate       = errors.New("duplicate error")
	ErrInternalFailure = errors.New("internal server error")
)

type UserRepository interface {
	Create(ctx context.Context, u *userModel.User) error
	FindByID(ctx context.Context, id string) (*userModel.User, error)
	FindByEmail(ctx context.Context, email string) (*userModel.User, error)
	Update(ctx context.Context, u *userModel.User) error
	Delete(ctx context.Context, id string) error

	CreateTx(ctx context.Context, tx pgx.Tx, u *userModel.User) error
	CreateWithWallet(ctx context.Context, u *userModel.User, w *walletModel.Wallet) (*userModel.UserWithWallet, error)
}

type postgresqlUserRepository struct {
	db *pgxpool.Pool

	walletRepo walletRepo.WalletRepository
}

func NewPostgresqlUserRepository(
	db *pgxpool.Pool,
	walRepo walletRepo.WalletRepository,
) UserRepository {
	return &postgresqlUserRepository{
		db:         db,
		walletRepo: walRepo,
	}
}

func (r *postgresqlUserRepository) Create(ctx context.Context, u *userModel.User) error {
	const op = "postgresqlUserRepository.Create"

	query := fmt.Sprintf(`
		INSERT INTO %s (id, username, email, password_hash)
		VALUES($1, $2, $3, $4)
	`, usersTableName)

	_, err := r.db.Exec(
		ctx,
		query,
		u.ID, u.Username, u.Email, u.PasswordHash,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return fmt.Errorf("%v: %w", ErrDuplicate, err)
		}
		return fmt.Errorf("%v: %w", ErrInternalFailure, err)
	}

	return nil
}

func (r *postgresqlUserRepository) FindByID(ctx context.Context, id string) (*userModel.User, error) {
	const op = "postgresqlUserRepository.FindByID"

	query := fmt.Sprintf(`
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM  %s
		WHERE id = $1
	`, usersTableName)

	row := r.db.QueryRow(ctx, query, id)

	user, err := scanIntoUser(row)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, ErrUserNotFound
		}

		return nil, ErrInternalFailure
	}

	return user, nil
}

func (r *postgresqlUserRepository) FindByEmail(ctx context.Context, email string) (*userModel.User, error) {
	const op = "postgresqlUserRepository.FindByEmail"

	query := fmt.Sprintf(`
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM  %s
		WHERE email = $1
	`, usersTableName)

	row := r.db.QueryRow(ctx, query, email)

	user, err := scanIntoUser(row)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternalFailure
	}

	return user, nil
}

func (r *postgresqlUserRepository) Update(ctx context.Context, u *userModel.User) error {
	const op = "postgresqlUserRepository.Update"

	query := fmt.Sprintf(`
		UPDATE  %s
		SET username = $1
		WHERE id = $2
	`, usersTableName)

	_, err := r.db.Exec(ctx, query, u.Username, u.ID)
	return err
}

func (r *postgresqlUserRepository) Delete(ctx context.Context, id string) error {
	const op = "postgresqlUserRepository.Delete"

	query := fmt.Sprintf(`
		DELETE FROM  %s
		WHERE id = $1
	`, usersTableName)

	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *postgresqlUserRepository) CreateTx(ctx context.Context, tx pgx.Tx, u *userModel.User) error {
	const op = "postgresqlUserRepository.CreateTx"

	query := fmt.Sprintf(`
		INSERT INTO %s (id, username, email, password_hash)
		VALUES($1, $2, $3, $4)
	`, usersTableName)

	_, err := tx.Exec(
		ctx,
		query,
		u.ID, u.Username, u.Email, u.PasswordHash,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return fmt.Errorf("%v: %w", ErrDuplicate, err)
		}
		return fmt.Errorf("%v: %w", ErrInternalFailure, err)
	}

	return nil
}

func (r *postgresqlUserRepository) CreateWithWallet(ctx context.Context, u *userModel.User, w *walletModel.Wallet) (*userModel.UserWithWallet, error) {
	const op = "postgresqlUserRepository.CreateWithWallet"

	var err error

	txCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	tx, err := r.db.BeginTx(txCtx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err = r.CreateTx(txCtx, tx, u); err != nil {
		return nil, err
	}

	if err = r.walletRepo.CreateTx(txCtx, tx, w); err != nil {
		return nil, err
	}

	if err = tx.Commit(txCtx); err != nil {
		return nil, err
	}

	return userModel.NewUserWithWallet(u, w), nil
}

func scanIntoUser(row pgx.Row) (*userModel.User, error) {
	var (
		user userModel.User
		err  error
	)
	err = row.Scan(
		&user.ID, &user.Username, &user.Email,
		&user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)
	return &user, err
}

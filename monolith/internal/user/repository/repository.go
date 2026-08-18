package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/waves2k/go-wallet/monolith/internal/user/model"
)

const dbSchema = "public"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrDuplicate       = errors.New("duplicate error")
	ErrInternalFailure = errors.New("internal server error")
)

type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, u *model.User) error
	Delete(ctx context.Context, id string) error
}

type postgresqlUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresqlUserRepository(db *pgxpool.Pool) UserRepository {
	return &postgresqlUserRepository{
		db: db,
	}
}

func (r *postgresqlUserRepository) Create(ctx context.Context, u *model.User) error {
	query := fmt.Sprintf(`
		INSERT INTO accounts (id, username, email, password_hash)
		VALUES($1, $2, $3, $4)
	`)

	_, err := r.db.Exec(
		ctx,
		query,
		u.ID, u.FullName, u.Email, u.PasswordHash,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return fmt.Errorf("%v: %w", ErrDuplicate, err)
		}
		return fmt.Errorf("%v: %w", ErrInternalFailure, err)
	}

	return nil
}

func (r *postgresqlUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	query := fmt.Sprintf(`
		SELECT id, username, email, password_hash, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1
	`)

	var user model.User
	if err := r.db.QueryRow(ctx, query, id).Scan(&user); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, ErrUserNotFound
		}

		return nil, ErrInternalFailure
	}

	return &user, nil
}

func (r *postgresqlUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := fmt.Sprintf(`
		SELECT id, username, email, password_hash, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1
	`)

	var user model.User
	if err := r.db.QueryRow(ctx, query, email).Scan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *postgresqlUserRepository) Update(ctx context.Context, u *model.User) error {
	query := fmt.Sprintf(`
		UPDATE users
		SET username = $1
		WHERE id = $2
	`)
	_, err := r.db.Exec(ctx, query, u.FullName, u.ID)
	return err
}

func (r *postgresqlUserRepository) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`
		DELETE FROM users
		WHERE id = $1
	`)
	_, err := r.db.Exec(ctx, query, id)
	return err
}

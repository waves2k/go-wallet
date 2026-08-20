package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/waves2k/go-wallet/monolith/internal/ledger/model"
)

const (
	ledgerEntriesTableName = "ledger_entries"
)

type LedgerRepository interface {
	CreateTx(ctx context.Context, tx pgx.Tx, l *model.LedgerEntry) error
	GetBallanceByWalletID(ctx context.Context, walletID string) (float64, error)
	GetEntriesByWalletID(ctx context.Context, walletID string) ([]model.LedgerEntry, error)
}

type postgresqlLedgerRepository struct {
	db *pgxpool.Pool
}

func NewPostgresqlLedgerRepository(db *pgxpool.Pool) LedgerRepository {
	return &postgresqlLedgerRepository{
		db: db,
	}
}

func (r *postgresqlLedgerRepository) CreateTx(ctx context.Context, tx pgx.Tx, l *model.LedgerEntry) error {
	const op = "postgresqlLedgerRepository.CreateTx"

	query := fmt.Sprintf(`
		INSERT INTO (id, wallet_id, transaction_id, entry_type, amount)
		VALUES ($1, $2, $3, $4, $5)
	`)
	_, err := tx.Exec(ctx, query,
		&l.ID, &l.WalletID, &l.TransactionID,
		&l.EntryType, &l.Amount)

	return err
}

func (r *postgresqlLedgerRepository) GetBallanceByWalletID(ctx context.Context, walletID string) (float64, error) {
	const op = "postgresqlLedgerRepository.GetBallanceByWalletID"

	query := fmt.Sprintf(`
		SELECT
			COALESCE(SUM(CASE WHEN entry_type = 'credit' THEN amount ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN entry_type = 'debit' THEN amount ELSE 0 END), 0)
		FROM %s
		WHERE id = $1
	`, ledgerEntriesTableName)

	var balance float64
	err := r.db.QueryRow(ctx, query, walletID).Scan(&balance)
	if err != nil {
		return 0, nil
	}
	return balance, nil
}

func (r *postgresqlLedgerRepository) GetEntriesByWalletID(ctx context.Context, walletID string) ([]model.LedgerEntry, error) {
	const op = "postgresqlLedgerRepository.GetEntriesByWalletID"

	query := fmt.Sprintf(`
		SELECT id, wallet_id, transaction_id, entry_type, amount, created_at
		FROM %s
		WHERE wallet_id = $1
		ORDER BY created_at DESC
	`, ledgerEntriesTableName)

	rows, err := r.db.Query(ctx, query, walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []model.LedgerEntry

	for rows.Next() {
		entry, err := scanIntoLedger(rows)
		if err != nil {
			return nil, err
		}

		entries = append(entries, *entry)
	}

	return entries, nil
}

func scanIntoLedger(row pgx.Row) (*model.LedgerEntry, error) {
	var (
		entry model.LedgerEntry
		err   error
	)
	err = row.Scan(&entry.ID, &entry.WalletID, &entry.TransactionID,
		&entry.EntryType, &entry.Amount, &entry.CreatedAt)

	return &entry, err
}

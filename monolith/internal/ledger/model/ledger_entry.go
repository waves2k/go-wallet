package model

import "time"

type LedgerEntry struct {
	ID            string    `db:"id"`
	WalletID      string    `db:"wallet_id"`
	TransactionID string    `db:"transaction_id"`
	EntryType     string    `db:"entry_type"`
	Amount        float64   `db:"amount"`
	CreatedAt     time.Time `db:"created_at"`
}

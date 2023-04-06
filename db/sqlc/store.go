package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore returns a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("Tx error: %v , rb Error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input to create a transfer transaction
type TransferTxParams struct {
	FromAcctID int64 `json:"from_acct_id"`
	ToAcctID   int64 `json:"to_acct_id"`
	Amount     int64 `json:"amount"`
}

// TransferTxResult os the result of the transfer transaction
type TransferTxResult struct {
	Transfer  Transfer `json:"transfer"`
	FromAcct  Account  `json:"from_acct"`
	ToAcct    Account  `json:"to_acct"`
	FromEntry Entry    `json:"from_entry"`
	ToEntry   Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// TransferTx performs a money transfer from one account to another.
// It creates a new transfer record, add account entries and update account balance record within a single db transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {

		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAcctID: arg.FromAcctID,
			ToAcctID:   arg.ToAcctID,
			Amount:     arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AcctID: arg.FromAcctID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AcctID: arg.ToAcctID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		// get account -> update the balance
		result.FromAcct, err = q.AddAccountBalnce(ctx, AddAccountBalnceParams{
			ID:     arg.FromAcctID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToAcct, err = q.AddAccountBalnce(ctx, AddAccountBalnceParams{
			ID:     arg.ToAcctID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

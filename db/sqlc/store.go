package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries
// and transactions: Composition
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// Result of the transfer function
type TransferTxResults struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// Performs money transfer from one account to the other
// It creates the transfer, add account entries, and update account's balance within a database transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResults, error) {
	var result TransferTxResults

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// get account -> update accounts' balance -> use locks

		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountId)
		if err != nil {
			return err
		}

		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountId,
			Balance: account1.Balance - arg.Amount,
		})

		if err != nil {
			return err
		}

		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountId)
		if err != nil {
			return err
		}

		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ToAccountId,
			Balance: account2.Balance + arg.Amount,
		})

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

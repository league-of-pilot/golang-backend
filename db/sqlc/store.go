package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	db *sql.DB
	*Queries
}

// Queries cung cấp ko đủ interface
// Tuy nhiên db có vẻ là cha của Queries tạo sao tạo bọc lẫn nhau ?
// -> để gọn interface dễ hiểu ?
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// ExecTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}
var DEBUG = true

// TransferTx performs a money transfer from one account to the other.
// It creates the transfer, add account entries, and update accounts' balance within a database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)
		if DEBUG {
			fmt.Print(txName, "Create transfer\n")
		}

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		if DEBUG {
			fmt.Print(txName, "Create entry 1\n")
		}
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		if DEBUG {
			fmt.Print(txName, "Create entry 2\n")
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// BAD - DEADLOCK IMPLEMENT
		if DEBUG {
			fmt.Print(txName, "Get account 1 for update\n")
		}
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		if DEBUG {
			fmt.Print(txName, "Update account 1\n")
		}
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      account1.ID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		if DEBUG {
			fmt.Print(txName, "Get account 2 for update\n")
		}
		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		if DEBUG {
			fmt.Print(txName, "Update account 2\n")
		}

		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      account2.ID,
			Balance: account2.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}

		// TODO: update accounts' balance

		return nil
	})
	return result, err
}

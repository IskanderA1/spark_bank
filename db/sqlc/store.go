package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) exexTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbError := tx.Rollback(); rbError != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbError)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	ToAccountID   int64 `json:"to_account_id"`
	FromAccountID int64 `json:"from_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"trasfer"`
	ToAccount   Account  `json:"to_account"`
	FromAccount Account  `json:"from_account"`
	ToEntry     Entry    `json:"to_entry"`
	FromEntry   Entry    `json:"from_entry"`
}

var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	//txName := ctx.Value(txKey)
	err := store.exexTx(ctx, func(q *Queries) error {

		var err error
		result.Transfer, err = q.CreateTransfer(ctx,
			CreateTransferParams{
				FromAccountID: arg.FromAccountID,
				ToAccountID:   arg.ToAccountID,
				Amount:        arg.Amount,
			})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx,
			CreateEntryParams{
				AccountID: arg.FromAccountID,
				Amount:    -arg.Amount,
			})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx,
			CreateEntryParams{
				AccountID: arg.ToAccountID,
				Amount:    arg.Amount,
			})
		if err != nil {
			return err
		}

		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

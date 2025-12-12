package Anuskh

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface{
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	Querier
}
type RealStore struct {
	*Queries
	db *sql.DB
}

func NewTxConn(db *sql.DB) Store {
	return &RealStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *RealStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback Err: %v, TxErr: %v", rbErr, err)
		}
		return err
	}
	return tx.Commit()

}

type TransferTxParams struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
}

type TransferTxResult struct {
	Transfer    Transfer
	FromAccount Account
	ToAccount   Account
	FromEntry   Entry
	ToEntry     Entry
}

func (store *RealStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfers(ctx, CreateTransfersParams(arg))
		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}
		//
		//
		//
		//Update Account and Balance

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = q.AddBalance(ctx, AddBalanceParams{
				ID:      arg.FromAccountID,
				Balance: -arg.Amount,
			})
			if err != nil {
				return err
			}

			result.ToAccount, err = q.AddBalance(ctx, AddBalanceParams{
				ID:      arg.ToAccountID,
				Balance: +arg.Amount,
			})
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, err = q.AddBalance(ctx, AddBalanceParams{
				ID:      arg.ToAccountID,
				Balance: +arg.Amount,
			})
			if err != nil {
				return err
			}

			result.FromAccount, err = q.AddBalance(ctx, AddBalanceParams{
				ID:      arg.FromAccountID,
				Balance: -arg.Amount,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return result, err
}

package Anuskh

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/nilesh0729/OrdinaryBank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {
	user:= CreateRandomUser(t)
	arg := CreateAccountsParams{
		Owner:    user.Username,
		Currency: util.RandomCurrency(),
		Balance:  util.RandomBalance(),
	}

	Account, err := testQueries.CreateAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, Account)

	require.Equal(t, arg.Owner, Account.Owner)
	require.Equal(t, arg.Balance, Account.Balance)
	require.Equal(t, arg.Currency, Account.Currency)

	require.NotZero(t, Account.ID)
	require.NotZero(t, Account.CreatedAt)

	return Account
}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	account2, err := testQueries.GetAccounts(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Owner, account2.Owner)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	arg := UpdateAccountsParams{
		ID:      account1.ID,
		Balance: util.RandomBalance(),
	}

	account2, err := testQueries.UpdateAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account2.ID, account1.ID)
	require.Equal(t, account2.Balance, arg.Balance)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	err := testQueries.DeleteAccounts(context.Background(), account1.ID)

	require.NoError(t, err)

	account2, err := testQueries.GetAccounts(context.Background(), account1.ID)

	require.EqualError(t, err, sql.ErrNoRows.Error())

	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
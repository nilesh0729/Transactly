package Anuskh

import (
	"context"
	"testing"
	"time"

	"github.com/nilesh0729/Transactly/internal/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomTransfers(t *testing.T, account1, account2 Account) Transfer {
	arg := CreateTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomBalance(),
	}

	Transfer, err := testQueries.CreateTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, Transfer)

	require.Equal(t, Transfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, Transfer.ToAccountID, arg.ToAccountID)
	require.Equal(t, Transfer.Amount, arg.Amount)

	require.NotZero(t, Transfer.ID)
	require.NotZero(t, Transfer.CreatedAt)

	return Transfer
}

func TestCreateTransfers(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	CreateRandomTransfers(t, account1, account2)
}

func TestGetTransfers(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	Transfer1 := CreateRandomTransfers(t, account1, account2)

	Transfer2, err := testQueries.GetTransfers(context.Background(), Transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, Transfer2)

	require.Equal(t, Transfer1.ID, Transfer2.ID)
	require.Equal(t, Transfer1.FromAccountID, Transfer2.FromAccountID)
	require.Equal(t, Transfer1.ToAccountID, Transfer2.ToAccountID)
	require.Equal(t, Transfer1.Amount, Transfer2.Amount)

	require.WithinDuration(t, Transfer1.CreatedAt, Transfer2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	for i := 0; i < 5; i++ {
		CreateRandomTransfers(t, account1, account2)
		CreateRandomTransfers(t, account2, account1)
	}
	arg := ListTransfersParams{
		Limit:         5,
		Offset:        5,
		FromAccountID: account1.ID,
		ToAccountID:   account1.ID,
	}

	Transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, Transfers)

	for _, transfer := range Transfers {
		require.NoError(t, err)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)
	}
}

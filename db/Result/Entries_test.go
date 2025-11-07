package Anuskh

import (
	"context"
	"testing"
	"time"

	"github.com/nilesh0729/Transactly/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomEntries(t *testing.T, account Account) Entry {
	arg := CreateEntriesParams{
		AccountID: account.ID,
		Amount:    util.RandomBalance(),
	}

	Entry, err := testQueries.CreateEntries(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, Entry)

	require.Equal(t, Entry.AccountID, arg.AccountID)
	require.Equal(t, Entry.Amount, arg.Amount)

	require.NotZero(t, Entry.ID)
	require.NotZero(t, Entry.CreatedAt)

	return Entry
}

func TestCreateEntries(t *testing.T) {
	account1 := CreateRandomAccount(t)
	CreateRandomEntries(t, account1)
}

func TestGetEntries(t *testing.T) {
	account1 := CreateRandomAccount(t)
	entry1 := CreateRandomEntries(t, account1)

	entry2, err := testQueries.GetEntries(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.Amount, entry2.Amount)

	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	account := CreateRandomAccount(t)

	for i := 0; i < 10; i++ {
		CreateRandomEntries(t, account)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Offset:    5,
		Limit:     5,
	}

	Entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, Entries)

	require.Len(t, Entries, 5)

	for _, entry := range Entries {
		require.NoError(t, err)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}
}

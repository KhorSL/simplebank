package db

import (
	"context"
	"testing"

	"github.com/khorsl/simple_bank/db/util"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(account Account) (Entry, CreateEntryParams, error) {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)

	return entry, arg, err
}

func TestCreateEntry(t *testing.T) {
	account, _, _ := createRandomAccount()
	entry, expected, err := createRandomEntry(account)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	require.Equal(t, expected.Amount, entry.Amount)
	require.Equal(t, expected.AccountID, entry.AccountID)
}

func TestGetEntry(t *testing.T) {
	account, _, _ := createRandomAccount()
	entry1, _, _ := createRandomEntry(account)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
}

func TestListEntries(t *testing.T) {
	account, _, _ := createRandomAccount()

	for i := 0; i < 10; i++ {
		createRandomEntry(account)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, account.ID, entry.AccountID)
	}
}

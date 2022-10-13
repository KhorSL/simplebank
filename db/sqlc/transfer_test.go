package db

import (
	"context"
	"testing"

	"github.com/khorsl/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	account1, _, err := createRandomAccount(t)
	require.NoError(t, err)

	account2, _, err := createRandomAccount(t)
	require.NoError(t, err)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.NotZero(t, transfer.ID)
	require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, transfer.ToAccountID, arg.ToAccountID)
	require.Equal(t, transfer.Amount, arg.Amount)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestListTransfers(t *testing.T) {
	var lastTransfer Transfer
	for i := 0; i < 10; i++ {
		lastTransfer = createRandomTransfer(t)
	}

	arg := ListTransfersParams{
		FromAccountID: lastTransfer.FromAccountID,
		ToAccountID:   lastTransfer.ToAccountID,
		Limit:         5,
		Offset:        0,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	// require.Len(t, accounts, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, lastTransfer.FromAccountID, transfer.FromAccountID)
		require.Equal(t, lastTransfer.ToAccountID, transfer.ToAccountID)
	}
}

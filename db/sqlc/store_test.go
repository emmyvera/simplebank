package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println("<<>> before: ", account1.Balance, account2.Balance)

	// run a concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {

		go func() {
			cxt := context.Background()
			result, err := store.TransferTx(cxt, TransferTxParams{
				FromAcctID: account1.ID,
				ToAcctID:   account2.ID,
				Amount:     amount,
			})

			errs <- err
			results <- result
		}()
	}

	// Check the result channel
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check the transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)

		require.Equal(t, account1.ID, transfer.FromAcctID)
		require.Equal(t, account2.ID, transfer.ToAcctID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check the from entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AcctID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// check the to entries
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AcctID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check the accounts
		fromAcct := result.FromAcct
		require.NotEmpty(t, fromAcct)
		require.Equal(t, account1.ID, fromAcct.ID)

		toAcct := result.ToAcct
		require.NotEmpty(t, toAcct)
		require.Equal(t, account2.ID, toAcct.ID)

		// check the account balance
		fmt.Println("<<>> tx: ", fromAcct.Balance, toAcct.Balance)
		diff1 := account1.Balance - fromAcct.Balance
		diff2 := toAcct.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}

	// Update the balance of an account
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println("<<>> after: ", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)
}

func TestTransferDeadlockTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println("<<>> before: ", account1.Balance, account2.Balance)

	// run a concurrent transfer transactions
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAcctID := account1.ID
		toAcctID := account2.ID

		if i%2 == 0 {
			fromAcctID = account2.ID
			toAcctID = account1.ID
		}

		go func() {
			cxt := context.Background()
			_, err := store.TransferTx(cxt, TransferTxParams{
				FromAcctID: fromAcctID,
				ToAcctID:   toAcctID,
				Amount:     amount,
			})

			errs <- err

		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	// Update the balance of an account
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println("<<>> after: ", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
}

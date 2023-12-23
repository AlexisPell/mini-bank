package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	fmt.Println(">> Balance before:", acc1.Balance, acc2.Balance)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	// goroutines
	errs := make(chan error)
	results := make(chan TransferTxResult)

	// Test in goroutines
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// Check the results
	existed := make(map[int]bool) //

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, acc1.ID, transfer.FromAccountID.Int64)
		require.Equal(t, acc2.ID, transfer.ToAccountID.Int64)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.TODO(), transfer.ID)
		require.NoError(t, err)

		// Check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, acc1.ID, fromEntry.AccountID.Int64)
		require.Equal(t, -amount, fromEntry.Amount)

		_, err = store.GetEntry(context.TODO(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, acc2.ID, toEntry.AccountID.Int64)
		require.Equal(t, amount, toEntry.Amount)

		_, err = store.GetEntry(context.TODO(), toEntry.ID)
		require.NoError(t, err)

		// Check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, acc1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, acc2.ID, toAccount.ID)

		// Check accounts balance
		fmt.Println(">> Tx: ", fromAccount.Balance, toAccount.Balance)
		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated account balances
	updatedAcc1, err := testQueries.GetAccount(context.TODO(), acc1.ID)
	require.NoError(t, err)

	updatedAcc2, err := testQueries.GetAccount(context.TODO(), acc2.ID)
	require.NoError(t, err)

	fmt.Println(">> Balance after:", updatedAcc1.Balance, updatedAcc2.Balance)
	require.Equal(t, acc1.Balance-int64(n)*amount, updatedAcc1.Balance)
	require.Equal(t, acc2.Balance+int64(n)*amount, updatedAcc2.Balance)
}

// TestTransferTxDeadlock
// We dont check the results here, but only
// verify that there are no db deadlocks
func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	fmt.Println(">> Balance before:", acc1.Balance, acc2.Balance)

	// run n concurrent transfer transactions
	n := 10
	amount := int64(10)

	// goroutines
	errs := make(chan error)

	// Test in goroutines
	for i := 0; i < n; i++ {
		fromAccId := acc1.ID
		toAccId := acc2.ID

		if i%2 == 1 {
			fromAccId = acc2.ID
			toAccId = acc1.ID
		}

		go func() {
			_, err := store.TransferTx(context.TODO(), TransferTxParams{
				FromAccountID: fromAccId,
				ToAccountID:   toAccId,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	// Check the results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated account balances
	updatedAcc1, err := testQueries.GetAccount(context.TODO(), acc1.ID)
	require.NoError(t, err)

	updatedAcc2, err := testQueries.GetAccount(context.TODO(), acc2.ID)
	require.NoError(t, err)

	fmt.Println(">> Balance after:", updatedAcc1.Balance, updatedAcc2.Balance)
	require.Equal(t, acc1.Balance, updatedAcc1.Balance)
	require.Equal(t, acc2.Balance, updatedAcc2.Balance)
}

package db

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    gofakeit.Name(),
		Currency: gofakeit.Currency().Short,
		Balance:  int64(gofakeit.Number(0, 500)),
	}

	account, err := testQueries.CreateAccount(context.TODO(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2, err := testQueries.GetAccount(context.TODO(), acc1.ID)
	require.NoError(t, err)
	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.TODO(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, acc := range accounts {
		require.NotEmpty(t, acc)
	}
}

func TestUpdateAccount(t *testing.T) {
	acc1 := createRandomAccount(t)
	newBalance := gofakeit.Number(0, 500)

	arg := UpdateAccountParams{
		ID:      acc1.ID,
		Balance: int64(newBalance),
	}
	_, err := testQueries.UpdateAccount(context.TODO(), arg)
	require.NoError(t, err)
	updatedAcc, err := testQueries.GetAccount(context.TODO(), acc1.ID)
	require.NoError(t, err)
	require.Equal(t, updatedAcc.Balance, int64(newBalance))
}

func TestDeleteAccount(t *testing.T) {
	acc1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.TODO(), acc1.ID)
	require.NoError(t, err)
	deletedAcc, err := testQueries.GetAccount(context.TODO(), acc1.ID)
	require.Empty(t, deletedAcc)
	require.Error(t, err)
}

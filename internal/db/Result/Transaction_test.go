package Anuskh


import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransaction(t *testing.T) {
	TxConn := NewTxConn(TestDb)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	fmt.Println(">>before tx : ", account1.Balance, account2.Balance)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	Results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			Result, err := TxConn.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			Results <- Result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <- errs
		require.NoError(t, err)

		result := <-Results
		require.NotEmpty(t, result)

		//check Transfer

		Transfer := result.Transfer

		require.Equal(t, account1.ID, Transfer.FromAccountID)
		require.Equal(t, account2.ID, Transfer.ToAccountID)

		require.Equal(t, amount, Transfer.Amount)

		require.NotZero(t, Transfer.ID)
		require.NotZero(t, Transfer.CreatedAt)

		_, err = TxConn.GetTransfers(context.Background(), Transfer.ID)
		require.NoError(t, err)
		//
		//
		//Check Entries

		//FromEntries
		fromEntries := result.FromEntry
		require.NotEmpty(t, fromEntries)

		require.Equal(t, account1.ID, fromEntries.AccountID)
		require.Equal(t, -amount, fromEntries.Amount)

		require.NotZero(t, fromEntries.ID)
		require.NotZero(t, fromEntries.CreatedAt)

		_, err = TxConn.GetEntries(context.Background(), fromEntries.ID)
		require.NoError(t, err)

		//ToEntries
		toEntries := result.ToEntry
		require.NotEmpty(t, toEntries)

		require.Equal(t, account2.ID, toEntries.AccountID)
		require.Equal(t, amount, toEntries.Amount)

		require.NotZero(t, toEntries.ID)
		require.NotZero(t, toEntries.CreatedAt)

		_, err = TxConn.GetEntries(context.Background(), toEntries.ID)
		require.NoError(t, err)

		//
		//
		//Check Updated Balance

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		fmt.Println(">>tx : ", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1 % amount == 0)

		k := int(diff1/amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}

	//All the indivudual check are done 
	// Now let's check the final updated Accounts

	UpdatedAccount1, err := testQueries.GetAccounts(context.Background(),account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, UpdatedAccount1)

	UpdatedAccount2 ,err :=testQueries.GetAccounts(context.Background(),account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, UpdatedAccount2)

	fmt.Println(">>After Tx : ", UpdatedAccount1.Balance,UpdatedAccount2.Balance)

	require.Equal(t, account1.Balance - int64(n)*amount, UpdatedAccount1.Balance)
	require.Equal(t, account2.Balance + int64(n)*amount, UpdatedAccount2.Balance)
}

func TestTransactionDeadlock(t *testing.T) {
	TxConn := NewTxConn(TestDb)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	fmt.Println(">>before tx : ", account1.Balance, account2.Balance)

	n := 10
	amount := int64(10)

	errs := make(chan error)
	

	for i := 0; i < n; i++ {

		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2==0 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			
			ctx := context.Background()
			_, err := TxConn.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
			
		}()
	}

	
	for i := 0; i < n; i++ {
		err := <- errs
		require.NoError(t, err)
	}

	

	UpdatedAccount1, err := testQueries.GetAccounts(context.Background(),account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, UpdatedAccount1)

	UpdatedAccount2 ,err :=testQueries.GetAccounts(context.Background(),account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, UpdatedAccount2)

	fmt.Println(">>After Tx : ", UpdatedAccount1.Balance,UpdatedAccount2.Balance)

	require.Equal(t, account1.Balance, UpdatedAccount1.Balance)
	require.Equal(t, account2.Balance, UpdatedAccount2.Balance)
	
}

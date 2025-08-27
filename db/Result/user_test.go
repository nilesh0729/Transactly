package Anuskh

import (
	"context"
	"testing"

	"github.com/nilesh0729/OrdinaryBank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User{
	hashPassword, err:= util.HashedPassword(util.RandomString(8))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username: util.RandomOwner(),
		HashedPassword: hashPassword,
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, user.HashedPassword, arg.HashedPassword)
	require.Equal(t, user.FullName, arg.FullName)
	require.Equal(t, user.Email, arg.Email)

	require.True(t,user.PasswordChangedAt.IsZero())
	require.NotZero(t , user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T){
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T){
	user := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(),user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user2.Username, user.Username)
	require.Equal(t, user2.HashedPassword, user.HashedPassword)
	require.Equal(t, user2.FullName, user.FullName)
	require.Equal(t, user2.Email, user.Email)

	require.NotZero(t, user2.CreatedAt, user.CreatedAt)
}

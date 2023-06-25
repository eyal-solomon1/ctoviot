package db

import (
	"context"
	"testing"

	"github.com/eyal-solomon1/ctoviot/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	pass, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: pass,
		FullName:       util.RandomOwner(),
		Balance:        100,
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	if err != nil {
		t.Error(err)
	}

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Balance, user.Balance)
	require.Equal(t, arg.Username, user.Username)

	require.NotZero(t, user.PasswordChangedAt.IsZero()) // checks if the user.PasswordChangedAt timestampz is a "zero" timestamptz
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.PasswordChangedAt, user2.PasswordChangedAt)
	require.Equal(t, user1.Balance, user2.Balance)
	require.Equal(t, user1.Username, user2.Username)

}

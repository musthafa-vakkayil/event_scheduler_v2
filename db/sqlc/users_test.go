package db

import (
	"context"
	"testing"

	"github.com/musthafa-vakkayil/event_scheduler_v2/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	password := utils.RandomString(8)
	passwordHash, err := utils.HashPassword(password)
	require.NoError(t, err)

	args := CreateUserParams{
		Username:       utils.RandomString(6),
		FullName:       utils.RandomString(10),
		Email:          utils.RandomEmail(),
		HashedPassword: passwordHash,
	}

	user, err := testQueries.CreateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.Email, user.Email)
	require.NotEmpty(t, user.CreatedAt)
	require.NotEmpty(t, user.UpdatedAt)

	return user
}

func TestUser(t *testing.T) {
	user := CreateRandomUser(t)

	// test get user
	user1, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user.Username, user1.Username)
	require.Equal(t, user.FullName, user1.FullName)
	require.Equal(t, user.Email, user1.Email)

	// test list users
	arg := ListUsersParams{
		Limit:  1,
		Offset: 0,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)

	// test update user
	newName := utils.RandomString(10)
	updateArgs := UpdateUserParams{
		FullName: newName,
		Username: user.Username,
	}

	user2, err := testQueries.UpdateUser(context.Background(), updateArgs)

	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, newName, user2.FullName)

	// Test delete user
	err = testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)
}

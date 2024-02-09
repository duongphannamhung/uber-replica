package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPhoneOrCreateUser(t *testing.T) {
	phone := "+1234567890"
	user1, err := testQueries.GetPhoneOrCreateUser(context.Background(), phone)
	require.NoError(t, err)
	require.NotEmpty(t, user1)
	require.Equal(t, user1.Phone, phone)

	user2, err := testQueries.GetPhoneOrCreateUser(context.Background(), phone)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user2.Phone, user1.Phone)
	require.Equal(t, user2.ID, user1.ID)

	testQueries.DeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)
	testQueries.DeleteUser(context.Background(), user2.ID)
	require.NoError(t, err)
}

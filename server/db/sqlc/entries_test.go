package db

import (
	"context"
	"testing"

	"github.com/eyal-solomon1/ctoviot/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) (Entry, User) {

	user := createRandomUser(t)
	video := createRandomVideo(t, user)

	arg := CreateEntryParams{
		Username:  user.Username,
		VideoName: video.VideoName,
		Amount:    util.RandomInt(1, 100),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	if err != nil {
		t.Error(err)
	}

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.Username, entry.Username)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.Username)
	require.NotZero(t, entry.CreatedAt)

	return entry, user
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1, _ := createRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.Username, entry2.Username)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)

}

func TestListEntries(t *testing.T) {

	_, user := createRandomEntry(t)

	arg := ListEntriesParams{
		Username: user.Username,
		Limit:    5,
		Offset:   0,
	}

	entires, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entires, 1)

}

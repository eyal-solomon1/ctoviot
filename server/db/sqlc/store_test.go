package db

import (
	"context"
	"github.com/eyal-solomon1/ctoviot/util"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)
	user := createRandomUser(t)

	// run n concurrent transfer transactions

	n := 3
	var reducedAmount = 10
	var totalReducedAmount float64

	errors := make(chan error)
	results := make(chan VideoTxResult)

	videoArg := CreateVideoParams{
		Owner:       user.Username,
		VideoLength: int64(reducedAmount),
		VideoDecs:   util.RandomString(20),
	}

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.VideoTx(context.Background(), VideoTxParam{Username: user.Username, Video: CreateVideoParams{
				Owner:           videoArg.Owner,
				VideoLength:     videoArg.VideoLength,
				VideoName:       util.RandomString(20),
				VideoRemotePath: util.RandomString(20),
				VideoDecs:       videoArg.VideoDecs,
			}})
			errors <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check Video
		video := result.Video
		require.Equal(t, video.Owner, user.Username)
		require.Equal(t, video.VideoDecs, videoArg.VideoDecs)
		require.Equal(t, int(video.VideoLength), 10)
		require.NotZero(t, video.ID)
		require.NotZero(t, video.CreatedAt)

		_, err = store.GetVideo(context.Background(), GetVideoParams{
			VideoIdentifier: video.VideoIdentifier,
			Owner:           user.Username,
		})
		require.NoError(t, err)

		// check Entry
		entry := result.Entry
		require.NotEmpty(t, entry)
		require.Equal(t, math.Abs(float64(entry.Amount)), float64(reducedAmount))
		require.NotZero(t, entry.CreatedAt)
		require.NotZero(t, entry.ID)

		_, err = store.GetEntry(context.Background(), entry.ID)
		require.NoError(t, err)

		// check User
		user2 := result.User
		require.NotEmpty(t, user)
		require.Equal(t, user, user)
		_, err = store.GetUser(context.Background(), user2.Username)
		require.NoError(t, err)

		totalReducedAmount += math.Abs(float64(entry.Amount))

	}

	require.Equal(t, int(totalReducedAmount), int(math.Abs(float64(n*reducedAmount))))

}

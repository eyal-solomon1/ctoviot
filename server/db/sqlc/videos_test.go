package db

import (
	"context"
	"testing"

	"github.com/eyal-solomon1/ctoviot/util"
	"github.com/stretchr/testify/require"
)

func createRandomVideo(t *testing.T, user User) Video {

	arg2 := CreateVideoParams{
		Owner:           user.Username,
		VideoName:       util.RandomString(6),
		VideoLength:     util.RandomInt(1, 20),
		VideoRemotePath: util.RandomString(6),
		VideoDecs:       util.RandomString(15),
	}

	video, err := testQueries.CreateVideo(context.Background(), arg2)

	if err != nil {
		t.Error(err)
	}

	require.NoError(t, err)
	require.NotEmpty(t, video)
	require.Equal(t, arg2.Owner, video.Owner)
	require.Equal(t, arg2.VideoDecs, video.VideoDecs)
	require.Equal(t, arg2.VideoName, video.VideoName)
	require.Equal(t, arg2.VideoRemotePath, video.VideoRemotePath)

	return video
}

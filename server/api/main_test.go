package api

import (
	"os"
	"testing"

	db "github.com/eyal-solomon1/ctoviot/db/sqlc"
	"github.com/eyal-solomon1/ctoviot/internal/aws"
	"github.com/eyal-solomon1/ctoviot/internal/ffmpeg"
	"github.com/eyal-solomon1/ctoviot/internal/openai"
	"github.com/eyal-solomon1/ctoviot/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type TesetServer struct {
}

func newTestServer(t *testing.T, store db.Store) *Server {
	config, err := util.LoadConfig(util.WithConfigFileName("conf.test"), util.WithConfigFilePath("../"))
	require.NoError(t, err)

	server, err := NewServer(config, store, &aws.AWSservice{}, &openai.OpenAIService{}, ffmpeg.FFMPEGService{}, util.InitializeLogger())

	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

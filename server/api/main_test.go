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

type initServerOption func(*Server)

func WithStore(store db.Store) initServerOption {
	return func(c *Server) {
		c.store = store
	}
}

func WithAWSService(aws aws.AWS) initServerOption {
	return func(c *Server) {
		c.awsService = aws
	}
}

func WithOpenAIService(openai openai.OpenAI) initServerOption {
	return func(c *Server) {
		c.openaiService = openai
	}
}
func WithFFMPEGService(ffmpeg ffmpeg.FFMPEG) initServerOption {
	return func(c *Server) {
		c.ffmpegService = ffmpeg
	}
}

type TesetServer struct {
}

func newTestServer(t *testing.T, options ...initServerOption) *Server {
	config, err := util.LoadConfig(util.WithConfigFileName("conf.test"), util.WithConfigFilePath("../"))
	require.NoError(t, err)

	var serverConfig = &Server{
		logger: util.InitializeLogger(),
	}

	for _, opt := range options {
		opt(serverConfig)
	}

	server, err := NewServer(config, serverConfig.store, serverConfig.awsService, serverConfig.openaiService, serverConfig.ffmpegService, serverConfig.logger)

	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"

	db "github.com/eyal-solomon1/ctoviot/db/sqlc"
	"github.com/eyal-solomon1/ctoviot/internal/aws"
	"github.com/eyal-solomon1/ctoviot/internal/ffmpeg"
	"github.com/eyal-solomon1/ctoviot/internal/openai"
	"github.com/eyal-solomon1/ctoviot/token"
	"github.com/eyal-solomon1/ctoviot/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type LogGroupName struct {
	name string
}

var VideoAPILogGroup = LogGroupName{name: "VideoAPI"}
var UserAPILogGroup = LogGroupName{name: "UserAPI"}

// Server type with all dependencies
type Server struct {
	config        util.Config
	store         db.Store
	tokenMaker    token.Maker
	router        *gin.Engine
	awsService    aws.AWS
	openaiService openai.OpenAI
	ffmpegService ffmpeg.FFMPEG
	logger        zerolog.Logger
}

// Creates a new server Instnace
func NewServer(cfg util.Config,
	store db.Store,
	awsService aws.AWS,
	openaiService openai.OpenAI,
	ffmpegService ffmpeg.FFMPEG,
	logger zerolog.Logger,
) (*Server, error) {

	tokenMaker, err := token.NewJWTMaker(cfg.TokenSymmatricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		store:         store,
		config:        cfg,
		tokenMaker:    tokenMaker,
		awsService:    awsService,
		openaiService: openaiService,
		ffmpegService: ffmpegService,
		logger:        logger,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(stringLoggerMiddleware())

	router.Use(cors.New(
		cors.Config{
			AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders: []string{"Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"},
			MaxAge:       12 * time.Hour,
			AllowOriginFunc: func(origin string) bool {
				return origin == server.config.AllowOriginEndpoint
			},
			AllowCredentials: true,
		}))

	router.GET("/health", server.Health)
	router.POST("/users/register", server.CreateUser)
	router.POST("/users/login", server.LoginUser)
	router.POST("/users/refresh_token", server.renewAccessToken)
	router.POST("/users/logout", server.LogoutUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/videos", server.getVideos)
	authRoutes.POST("/new_video", server.createVideo)
	authRoutes.POST("/delete-video", server.deleteVideo)

	server.router = router
}

// Starts the server Instance
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) errLogger(service LogGroupName, err error) error {
	server.logger.Error().Stack().Str("service", service.name).Err(err).Msg(" ")
	return nil
}

func (server *Server) infoLogger(service LogGroupName, logs ...string) error {
	var builder strings.Builder
	for _, log := range logs {
		if _, err := builder.WriteString(fmt.Sprintf("%v ,", log)); err != nil {
			server.logger.Error().Msg(err.Error())
		}
	}
	str := strings.TrimSpace(builder.String())
	server.logger.Info().Str("service", service.name).Msg(str)
	return nil
}

func errorResponse(err error) gin.H {
	return gin.H{"ok": false, "error": err.Error()}
}

func okResponse(payload interface{}) gin.H {
	return gin.H{"ok": true, "payload": payload}
}

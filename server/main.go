package main

import (
	"context"
	"database/sql"
	"github.com/eyal-solomon1/ctoviot/api"
	db "github.com/eyal-solomon1/ctoviot/db/sqlc"
	"github.com/eyal-solomon1/ctoviot/internal/aws"
	"github.com/eyal-solomon1/ctoviot/internal/ffmpeg"
	"github.com/eyal-solomon1/ctoviot/internal/openai"
	util "github.com/eyal-solomon1/ctoviot/util"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
)

func main() {
	logger := util.InitializeLogger()

	config, err := util.LoadConfig()

	if err != nil {
		logger.Fatal().
			Err(err).
			Str("service", "config").
			Msgf("couldn't load config")
	}

	err = config.Validate()

	if err != nil {
		logger.Fatal().
			Err(err).
			Str("service", "config").Msg("config is missing values")
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		logger.Error().Msg(err.Error())
	}

	store := db.NewStore(conn)

	awsService := aws.Initialize(context.Background(), config.AwsRegion)
	openaiSerivce := openai.Initialize(config.OpenAIToken)
	ffmpegService := ffmpeg.Initialize()

	_, err = store.ValidateDB(context.Background())
	if err != nil {
		logger.Fatal().
			Err(err).
			Str("service", "db").
			Msgf("couldn't connect to db")
	}

	server, err := api.NewServer(config, store, awsService, openaiSerivce, ffmpegService, logger)

	if err != nil {
		logger.Error().Msg(err.Error())
	}

	logger.Printf("Starting server on %v ...", config.ServerAddress)
	server.Start(config.ServerAddress)

}

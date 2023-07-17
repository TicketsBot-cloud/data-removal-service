package main

import (
	"github.com/TicketsBot/data-removal-service/config"
	"github.com/TicketsBot/data-removal-service/internal/cache"
	"github.com/TicketsBot/data-removal-service/internal/database"
	"github.com/TicketsBot/data-removal-service/pkg/service"
	"go.uber.org/zap"
)

func main() {
	config := config.ParseConfig()

	var logger *zap.Logger
	var err error
	if config.ProductionMode {
		logger, err = zap.NewProduction(zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	} else {
		logger, err = zap.NewDevelopment(zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	}

	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	logger.Debug("Connecting to cache...")

	pgCache, err := cache.ConnectPostgres(config)
	if err != nil {
		logger.Fatal("Failed to connect to postgres (cache)", zap.Error(err))
		panic(err)
	}

	logger.Debug("Connected to cache!")

	logger.Debug("Connecting to database...")
	dbClient, err := database.ConnectPostgres(config)
	if err != nil {
		logger.Fatal("Failed to connect to postgres (db)", zap.Error(err))
		panic(err)
	}

	logger.Debug("Connected to database!")

	cacheExecutor := cache.NewPostgresExecutor(config, logger, pgCache)
	dbExecutor := database.NewPostgresExecutor(config, logger, dbClient)

	service.NewService(config, logger, cacheExecutor, dbExecutor).Run()
}

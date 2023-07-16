package main

import (
	dataremovalservice "github.com/TicketsBot/data-removal-service"
	"github.com/TicketsBot/data-removal-service/cache"
	"github.com/TicketsBot/data-removal-service/config"
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

	logger.Info("Connecting to cache...")

	pgCache, err := cache.ConnectPostgres(config)
	if err != nil {
		logger.Fatal("Failed to connect to postgres", zap.Error(err))
		panic(err)
	}

	logger.Info("Connected to cache!")

	cacheExecutor := cache.NewPostgresExecutor(config, logger, pgCache)

	dataremovalservice.NewService(config, logger, cacheExecutor).Run()
}

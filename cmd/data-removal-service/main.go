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
	if config.ProductionMode {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}

	pgCache, err := cache.ConnectPostgres(config)
	if err != nil {
		logger.Fatal("Failed to connect to postgres", zap.Error(err))
		panic(err)
	}

	cacheExecutor := cache.NewPostgresExecutor(config, logger, pgCache)

	dataremovalservice.NewService(config, logger, cacheExecutor).Run()
}

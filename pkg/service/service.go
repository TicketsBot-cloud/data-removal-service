package service

import (
	"github.com/TicketsBot/data-removal-service/config"
	"github.com/TicketsBot/data-removal-service/internal/cache"
	"github.com/TicketsBot/data-removal-service/internal/database"
	"go.uber.org/zap"
	"time"
)

type Service struct {
	config           config.Config
	logger           *zap.Logger
	CacheExecutor    cache.CacheExecutor
	DatabaseExecutor database.DatabaseExecutor
}

func NewService(
	config config.Config,
	logger *zap.Logger,
	cacheExecutor cache.CacheExecutor,
	dbExecutor database.DatabaseExecutor,
) *Service {
	return &Service{
		config:           config,
		logger:           logger,
		CacheExecutor:    cacheExecutor,
		DatabaseExecutor: dbExecutor,
	}
}

func (s *Service) Run() {
	if s.config.DaemonMode {
		s.logger.Info("Running in daemon mode", zap.Duration("frequency", s.config.DaemonExecutionFrequency))
		s.startDaemon()
	} else {
		s.logger.Info("Running in one-shot mode")
		if err := s.doPurge(); err != nil {
			s.logger.Error("Failed to purge", zap.Error(err))
		}
	}
}

func (s *Service) startDaemon() {
	for {
		time.Sleep(s.config.DaemonExecutionFrequency)

		s.logger.Info("Running purge...")

		if err := s.doPurge(); err != nil {
			s.logger.Error("Failed to purge", zap.Error(err))
		}

		s.logger.Info("Purge complete")
	}
}

func (s *Service) doPurge() error {
	purgeThreshold := time.Hour * 24 * time.Duration(s.config.PurgeThresholdDays)

	if err := s.CacheExecutor.PurgeUsers(purgeThreshold); err != nil {
		return err
	}

	if err := s.CacheExecutor.PurgeMembers(purgeThreshold); err != nil {
		return err
	}

	if err := s.DatabaseExecutor.Purge(purgeThreshold); err != nil {
		return err
	}

	return nil
}

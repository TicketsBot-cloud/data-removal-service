package database

import (
	"context"
	"github.com/TicketsBot/data-removal-service/config"
	"github.com/TicketsBot/database"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"time"
)

type PostgresExecutor struct {
	config config.Config
	logger *zap.Logger
	client *database.Database
}

func NewPostgresExecutor(config config.Config, logger *zap.Logger, client *database.Database) *PostgresExecutor {
	return &PostgresExecutor{
		config,
		logger,
		client,
	}
}

func ConnectPostgres(config config.Config) (*database.Database, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelFunc()

	pool, err := pgxpool.Connect(ctx, config.DatabaseUri)
	if err != nil {
		return nil, err
	}

	return database.NewDatabase(pool), nil
}

func (e *PostgresExecutor) Purge(threshold time.Duration) error {
	return e.purgeUserGuilds(threshold)
}

func (e *PostgresExecutor) purgeUserGuilds(threshold time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.config.QueryTimeout)
	defer cancel()

	var interval pgtype.Interval
	if err := interval.Set(threshold); err != nil {
		return err
	}

	e.logger.Info("Purging dashboard users")
	rowCount, err := e.client.DashboardUsers.PurgeOldUsers(ctx, threshold)
	if err != nil {
		return err
	}

	e.logger.Info("Purged dashboard users", zap.Int64("user_count", rowCount))
	return nil
}

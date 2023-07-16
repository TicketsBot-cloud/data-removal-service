package cache

import (
	"context"
	_ "embed"
	"github.com/TicketsBot/data-removal-service/config"
	"github.com/jackc/pgtype"
	"github.com/rxdn/gdl/cache"
	"go.uber.org/zap"
	"time"
)

type PostgresExecutor struct {
	config config.Config
	logger *zap.Logger
	client *cache.PgCache
}

func NewPostgresExecutor(config config.Config, logger *zap.Logger, client *cache.PgCache) *PostgresExecutor {
	return &PostgresExecutor{
		config,
		logger,
		client,
	}
}

var (
	//go:embed sql/purge_users.sql
	queryPurgeUsers string

	//go:embed sql/purge_members.sql
	queryPurgeMembers string
)

func (m *PostgresExecutor) PurgeUsers(threshold time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.config.QueryTimeout)
	defer cancel()

	var interval pgtype.Interval
	if err := interval.Set(threshold); err != nil {
		return err
	}

	metadata, err := m.client.Exec(ctx, queryPurgeUsers, interval)
	if err != nil {
		return err
	}

	m.logger.Info("Purged users", zap.Int64("user_count", metadata.RowsAffected()))

	return nil
}

func (m *PostgresExecutor) PurgeMembers(threshold time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.config.QueryTimeout)
	defer cancel()

	var interval pgtype.Interval
	if err := interval.Set(threshold); err != nil {
		return err
	}

	metadata, err := m.client.Exec(ctx, queryPurgeMembers, interval)
	if err != nil {
		return err
	}

	m.logger.Info("Purged members", zap.Int64("member_count", metadata.RowsAffected()))

	return nil
}

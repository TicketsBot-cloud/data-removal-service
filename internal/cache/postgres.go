package cache

import (
	"context"
	_ "embed"
	"github.com/TicketsBot/data-removal-service/config"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
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

func ConnectPostgres(config config.Config) (*cache.PgCache, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelFunc()

	pool, err := pgxpool.Connect(ctx, config.CacheUri)
	if err != nil {
		return nil, err
	}

	pgCache := cache.NewPgCache(pool, cache.CacheOptions{
		Guilds:   true,
		Users:    true,
		Members:  true,
		Channels: true,
		Roles:    true,
	})

	return &pgCache, nil
}

var (
	//go:embed sql/purge_users.sql
	queryPurgeUsers string

	//go:embed sql/purge_members.sql
	queryPurgeMembers string
)

func (e *PostgresExecutor) PurgeUsers(threshold time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.config.QueryTimeout)
	defer cancel()

	var interval pgtype.Interval
	if err := interval.Set(threshold); err != nil {
		return err
	}

	metadata, err := e.client.Exec(ctx, queryPurgeUsers, interval)
	if err != nil {
		return err
	}

	e.logger.Info("Purged users", zap.Int64("user_count", metadata.RowsAffected()))

	return nil
}

func (e *PostgresExecutor) PurgeMembers(threshold time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.config.QueryTimeout)
	defer cancel()

	var interval pgtype.Interval
	if err := interval.Set(threshold); err != nil {
		return err
	}

	metadata, err := e.client.Exec(ctx, queryPurgeMembers, interval)
	if err != nil {
		return err
	}

	e.logger.Info("Purged members", zap.Int64("member_count", metadata.RowsAffected()))

	return nil
}

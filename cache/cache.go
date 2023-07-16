package cache

import (
	"context"
	"github.com/TicketsBot/data-removal-service/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rxdn/gdl/cache"
	"time"
)

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

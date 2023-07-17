package cache

import "time"

type CacheExecutor interface {
	PurgeUsers(threshold time.Duration) error
	PurgeMembers(threshold time.Duration) error
}

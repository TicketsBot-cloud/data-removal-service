package database

import "time"

type DatabaseExecutor interface {
	Purge(threshold time.Duration) error
}

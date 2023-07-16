package config

import (
	"github.com/caarlos0/env/v6"
	"time"
)

type Config struct {
	ProductionMode bool `env:"PRODUCTION_MODE" envDefault:"false"`

	DaemonMode               bool          `env:"DAEMON_MODE" envDefault:"false"`
	DaemonExecutionFrequency time.Duration `env:"DAEMON_EXECUTION_FREQUENCY" envDefault:"1h"`

	CacheUri     string        `env:"CACHE_URI,notEmpty"`
	QueryTimeout time.Duration `env:"QUERY_TIMEOUT" envDefault:"10m"`

	PurgeThresholdDays int `env:"PURGE_THRESHOLD_DAYS,required,notEmpty"`
}

func ParseConfig() (conf Config) {
	if err := env.Parse(&conf); err != nil {
		panic(err)
	}

	return
}

package logger

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var (
	log  zerolog.Logger
	once sync.Once
)

func Get() *zerolog.Logger {
	once.Do(func() {
		zerolog.TimeFieldFormat = time.RFC3339Nano
		log = zerolog.New(os.Stdout).With().
			Timestamp().
			Str("service", "northstar-api").
			Str("version", "1.0.0").
			Str("environment", os.Getenv("APP_ENV")).
			Logger()
	})
	return &log
}
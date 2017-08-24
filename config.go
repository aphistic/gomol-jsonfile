package jsonstream

import (
	"github.com/efritz/glock"
	"io"
)

type loggerConfig func(cfg *Config)

func withClock(clock glock.Clock) loggerConfig {
	return func(cfg *Config) {
		cfg.clock = clock
	}
}

func WithHeaders(headers map[string]interface{}) loggerConfig {
	return func(cfg *Config) {
		if headers == nil {
			headers = make(map[string]interface{})
		}
		cfg.headers = headers
	}
}

type Config struct {
	clock glock.Clock

	w       io.Writer
	headers map[string]interface{}
}

func NewConfig(w io.Writer, configs ...loggerConfig) *Config {
	cfg := &Config{
		clock: glock.NewRealClock(),

		w:       w,
		headers: make(map[string]interface{}),
	}

	for _, f := range configs {
		f(cfg)
	}

	return cfg
}

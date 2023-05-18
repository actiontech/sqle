package retry

import (
	"time"
)

const (
	defaultAttempts = 3
	defaultDelay    = time.Second
)

type Config struct {
	attempts uint
	delay    time.Duration
}

type Option func(*Config)

// Option: Attempts
func Attempts(attempts uint) Option {
	return func(c *Config) {
		c.attempts = attempts
	}
}

// Option: Delay
func Delay(delay time.Duration) Option {
	return func(c *Config) {
		c.delay = delay
	}
}

func NewDefaultRetryConfig() *Config {
	return &Config{
		attempts: defaultAttempts,
		delay:    defaultDelay,
	}
}

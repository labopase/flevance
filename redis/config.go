package redis

import "time"

type Config struct {
	Addrs      []string `mapstructure:"addrs"`
	Password   string   `mapstructure:"password"`
	DB         int      `mapstructure:"db"`
	ClientName string   `mapstructure:"client_name"`

	MaxRetries      int           `mapstructure:"max_retries"`
	MinRetryBackoff time.Duration `mapstructure:"min_retry_backoff"`
	MaxRetryBackoff time.Duration `mapstructure:"max_retry_backoff"`
}

func DefaultConfig() *Config {
	return &Config{
		Addrs:      []string{"localhost:6379"},
		Password:   "",
		DB:         0,
		ClientName: "",
	}
}

func (c *Config) applyDefault() {

	// disable retries (set to -1)

	if c.MaxRetries == 0 {
		c.MaxRetries = -1
	}

	if c.MinRetryBackoff == 0 {
		c.MinRetryBackoff = -1
	}

	if c.MaxRetryBackoff == 0 {
		c.MaxRetryBackoff = -1
	}
}

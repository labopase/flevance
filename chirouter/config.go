package chirouter

import (
	"fmt"
	"time"
)

type Config struct {
	Host                string        `mapstructure:"host"`
	Port                int           `mapstructure:"port"`
	ReadTimeout         time.Duration `mapstructure:"read_timeout"`
	WriteTimeout        time.Duration `mapstructure:"write_timeout"`
	IdleTimeout         time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout     time.Duration `mapstructure:"shutdown_timeout"`
	MaxHeaderBytes      int64         `mapstructure:"max_header_bytes"`
	MaxRequestBodyBytes int64         `mapstructure:"max_request_body_bytes"`
	MaxUploadSize       int64         `mapstructure:"max_upload_size"`
}

func DefaultConfig() *Config {
	return &Config{
		Host: "localhost",
		Port: 8080,
	}
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *Config) applyDefaults() {
	if c.Host == "" {
		c.Host = "localhost"
	}
	if c.Port == 0 {
		c.Port = 8080
	}

	if c.ReadTimeout == 0 {
		c.ReadTimeout = time.Second * 15
	}

	if c.WriteTimeout == 0 {
		c.WriteTimeout = time.Second * 15
	}

	if c.IdleTimeout == 0 {
		c.IdleTimeout = time.Second * 60
	}

	if c.ShutdownTimeout == 0 {
		c.ShutdownTimeout = time.Second * 30
	}

	if c.MaxHeaderBytes == 0 {
		c.MaxHeaderBytes = 1 << 20
	}

	if c.MaxRequestBodyBytes == 0 {
		c.MaxRequestBodyBytes = 1 << 20
	}

	if c.MaxUploadSize == 0 {
		c.MaxUploadSize = 3 << 20
	}
}

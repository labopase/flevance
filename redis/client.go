package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type client struct {
	rdb    redis.UniversalClient
	config *Config
}

type Client interface {
	Db() redis.UniversalClient
	Ping(ctx context.Context) error
	Close() error
	Count(ctx context.Context, pattern string) (int64, error)
	Exists(ctx context.Context, keyPrefix string) (bool, error)
	Destroy(ctx context.Context, pattern string) error
}

func NewClient(cfg *Config) (Client, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	cfg.applyDefault()

	opts := &redis.UniversalOptions{
		Addrs:    cfg.Addrs,
		Password: cfg.Password,
	}

	if len(cfg.Addrs) == 1 {
		opts.DB = cfg.DB
	}

	rdb := redis.NewUniversalClient(opts)

	return &client{
		rdb:    rdb,
		config: cfg,
	}, nil
}

func (c *client) Db() redis.UniversalClient {
	return c.rdb
}

func (c *client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

func (c *client) Close() error {
	return c.rdb.Close()
}

func (c *client) Count(ctx context.Context, pattern string) (int64, error) {
	iter := c.rdb.Scan(ctx, 0, pattern, 100).Iterator()

	var count int64
	for iter.Next(ctx) {
		count++
	}

	return count, nil
}

func (c *client) Exists(ctx context.Context, keyPrefix string) (bool, error) {
	res, err := c.rdb.Exists(ctx, keyPrefix).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, fmt.Errorf("exists: %w", err)
	}

	return res > 0, err
}

func (c *client) Destroy(ctx context.Context, pattern string) error {
	cursor := uint64(0)
	for {
		var keys []string
		var err error

		keys, cursor, err = c.rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("destroy: %w", err)
		}

		if len(keys) == 0 {
			break
		}

		if err := c.rdb.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("destroy: %w", err)
		}
	}

	return nil
}

package config

import (
	"context"
	"crypto/tls"
	"github.com/redis/go-redis/v9"
)

func NewRedisConn(c *RedisConfig) (*redis.Client, error) {

	var rdb *redis.Client
	var tlsConfig *tls.Config

	if c.UseTLS {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:      c.Host,
		Username:  c.Username,
		Password:  c.Password,
		DB:        c.DBName,
		TLSConfig: tlsConfig,
	})

	// testing
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	
	return rdb, nil
}
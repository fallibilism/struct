package config

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
)

const (
	ViewPath     = "views"
	ViewExt      = ".html"
	Developement = true
)

var (
	Prometheus = PrometheusConfig{
		Enabled:     true,
		Namespace:   "v",
		MetricsPath: "/metrics",
	}
	Server = ServerConfig{
		Port: "8080",
		Host: "localhost",
	}
	App     = &AppConfig{}
	Livekit = LivekitConfig{
		Host:   "http://localhost:7880",
		ApiKey: "api_key",
		Secret: "secret",
	}
)

type AppConfig struct {
	DB    *sql.DB
	Redis *redis.Client
}

type PrometheusConfig struct {
	Enabled bool
	// Namespace for prometheus metrics
	Namespace string
	// MetricsPath for prometheus metrics
	MetricsPath string
}

type ServerConfig struct {
	// Port for server
	Port string
	// Host for server
	Host string
}

type LivekitConfig struct {
	Host   string
	ApiKey string
	Secret string
}

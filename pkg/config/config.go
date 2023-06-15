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

var Redis RedisConfig
var Postgres PostgresConfig
var Openai OpenAIConfig

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
	Host   string `yaml:"host"`
	ApiKey string `yaml:"api_key"`
	Secret string `yaml:"secret"`
}

type RedisConfig struct {
	Host              string   `yaml:"host"`
	Username          string   `yaml:"username"`
	Password          string   `yaml:"password"`
	DBName            int      `yaml:"db"`
	UseTLS            bool     `yaml:"use_tls"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int32  `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db"`
	Prefix   string `yaml:"prefix"`
	SslMode	 string `yaml:"sslmode" default:"disable"`
}

// OPEN AI CONFIG 
type OpenAIConfig struct {
	ApiKey string `yaml:"api_key"`
	Secret string `yaml:"secret"`
}

type Config struct {
	App *AppConfig
	Openai OpenAIConfig `yaml:"open_ai"`
	Logging string `yaml:"logging"`
	Postgres PostgresConfig `yaml:"postgres"`
	Redis RedisConfig `yaml:"redis"`
	Livekit LivekitConfig `yaml:"livekit"`
}

func SetConfig(filename string) {
	conf := &Conf {

	}

	conf := &Config{}

	if content != "" {
		if err := yaml.Unmarshal([]byte(content), conf); err != nil {
			return nil, fmt.Errorf("could not parse config: %v", err)
		}
	}

	Openai = conf.Openai
	PostgresConfig = conf.Postgres
	Openai = conf.Openai
}
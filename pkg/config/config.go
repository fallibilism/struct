package config

import (
	"v/pkg/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
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
	TestConfig = &AppConfig{} //hack to get Config in test
	Livekit    = LivekitConfig{
		Host:   "http://localhost:7880",
		ApiKey: "api_key",
		Secret: "secret",
	}
)

var Redis RedisConfig
var Postgres PostgresConfig
var Openai OpenAIConfig

type AppConfig struct {
	DB    *gorm.DB
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
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   int    `yaml:"db"`
	UseTLS   bool   `yaml:"use_tls"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int32  `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"database"`
	Prefix   string `yaml:"prefix"`
	SslMode  string `yaml:"sslmode" default:"disable"`
	TimeZone string `yaml:"timezone" default:"Asia/Jakarta"`
}

// OPEN AI CONFIG
type OpenAIConfig struct {
	ApiKey string `yaml:"api_key"`
	Secret string `yaml:"secret"`
}

type Config struct {
	Name         string         `yaml:"name"`
	Developement bool           `yaml:"developement"`
	Port         uint           `yaml:"port"`
	Openai       OpenAIConfig   `yaml:"open_ai"`
	Logging      string         `yaml:"logging"`
	Postgres     PostgresConfig `yaml:"postgres"`
	Redis        RedisConfig    `yaml:"redis"`
	Livekit      LivekitConfig  `yaml:"livekit"`
}

func SetConfig(filename string) {
	conf, err := utils.ReadFile(filename, &Config{})
	if err != nil {
		panic(err)
	}

	Openai = conf.Openai
	Postgres = conf.Postgres
	Redis = conf.Redis
}

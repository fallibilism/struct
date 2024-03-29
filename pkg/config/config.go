package config

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
	"strconv"
	stt "cloud.google.com/go/speech/apiv1"
	tts "cloud.google.com/go/texttospeech/apiv1"
	"github.com/sashabaranov/go-openai"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	ViewPath     = "views"
	ViewExt      = ".html"
	Developement = true
	keySize = 2048
)

var (
	BotIdentity = "BOT"
	RsaPrivateKey *rsa.PrivateKey
	Prometheus = PrometheusConfig{
		Enabled:     true,
		Namespace:   "v",
		MetricsPath: "/metrics",
	}
	Server = ServerConfig{
		Port: "8080",
		Host: "0.0.0.0",
	}
	TestConfig = &AppConfig{} //hack to get Config in test
	App        = &AppConfig{}
	Livekit    = LivekitConfig{
		Host:   "http://localhost:7880",
		ApiKey: "api_key",
		Secret: "secret",
	}
	Conf = &Config{
	}
)

var Redis *RedisConfig
var Postgres *PostgresConfig
var Openai *OpenAIConfig

type AppConfig struct {
	DB    *gorm.DB
	Redis *redis.Client
	SpeechClient *stt.Client
	TextClient *tts.Client
	OpenaiClient *openai.Client
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
	Port     int32  `yaml:"port"`
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
	TimeZone string `yaml:"timezone"`

	External bool   `yaml:"external"`
	URI      string `env:"POSTGRES_URI"`
}

// OPEN AI CONFIG
type OpenAIConfig struct {
	ApiKey string `yaml:"api_key"`
	Secret string `yaml:"secret"`
	Token string `yaml:"token"`
}

type Config struct {
	Name         string         `yaml:"name"`
	Developement bool           `yaml:"developement"`
	Port         uint           `yaml:"port"`
	JWTSecret    string         `yaml:"jwt_secret"` // LTI and JWT secret
	JWTIssuer    string         `yaml:"jwt_issuer"`
	ConsumerKey  string         `yaml:"consumer_key"`
	Openai       OpenAIConfig   `yaml:"open_ai"`
	Logging      string         `yaml:"logging"`
	Sqlite	string		`yaml:"sqlite"`
	Postgres     PostgresConfig `yaml:"postgres"`
	Redis        RedisConfig    `yaml:"redis"`
	Livekit      LivekitConfig  `yaml:"livekit"`
}

func generateRSAKey(keySize int) (*rsa.PrivateKey, error) {
	// Generate an RSA private key of the specified key size
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func SetConfig(filename string) (conf *Config) {
	conf, err := readFile(filename, &Config{})
	if err != nil {
		panic("config: " + err.Error())
	}
	rp, err := generateRSAKey(keySize)
	if err != nil {
		panic("config: " + err.Error())
	}
	RsaPrivateKey = rp

	Openai = &conf.Openai
	Postgres = &conf.Postgres
	Redis = &conf.Redis
	Conf = conf

	// environmental variables override config
	if v := os.Getenv("PORT"); v != "" {
		port, err := strconv.Atoi(v)
		if err == nil {
			conf.Port = uint(port)
		}
	}

	if v := os.Getenv("JWT_SECRET"); v != "" {
		conf.JWTSecret = v
	}

	if v := os.Getenv("POSTGRES_URI"); v != "" {
		conf.Postgres.URI = v
	}

	return conf
}

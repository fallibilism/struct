package config

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
)

type Config struct {
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

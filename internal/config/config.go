// Package config provides configuration structures and loading logic for goincidentcli.
package config

// ServiceConfig holds the name and health-check URL for a monitored service.
type ServiceConfig struct {
	Name string `mapstructure:"name"`
	URL  string `mapstructure:"url"`
}

// Config stores the application's configuration loaded from ~/.incident.yaml.
// It uses mapstructure tags for Viper unmarshaling.
type Config struct {
	APIToken        string          `mapstructure:"api_token"`
	BaseURL         string          `mapstructure:"base_url"`
	SlackToken      string          `mapstructure:"slack_token"`
	PrometheusURL   string          `mapstructure:"prometheus_url"`
	PrometheusQuery string          `mapstructure:"prometheus_query"`
	Services        []ServiceConfig `mapstructure:"services"`
}

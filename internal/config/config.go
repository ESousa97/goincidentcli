package config

// Config stores the application's configuration loaded from ~/.incident.yaml.
// It uses mapstructure tags for Viper unmarshaling.
type Config struct {
	APIToken        string `mapstructure:"api_token"`
	BaseURL         string `mapstructure:"base_url"`
	SlackToken      string `mapstructure:"slack_token"`
	PrometheusURL   string `mapstructure:"prometheus_url"`
	PrometheusQuery string `mapstructure:"prometheus_query"`
}

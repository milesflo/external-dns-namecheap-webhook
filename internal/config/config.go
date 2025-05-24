package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Username   string       `mapstructure:"username"`
	APIKey     string       `mapstructure:"api_key"`
	ClientIP   string       `mapstructure:"client_ip"`
	UseSandbox bool         `mapstructure:"use_sandbox"`
	Server     ServerConfig `mapstructure:"server"`
}

type ServerConfig struct {
	WebhookPort int `mapstructure:"webhook_port"`
	HealthPort  int `mapstructure:"health_port"`
}

func LoadConfig() (*Config, error) {
	// Configure Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Environment variable settings
	viper.AutomaticEnv()

	// Define CLI flags
	pflag.String("username", "", "Namecheap username")
	pflag.String("api-key", "", "Namecheap api key")
	pflag.String("client-ip", "", "Allowlisted IP you are calling from")
	pflag.Bool("use-sandbox", false, "Optional: Target the Namecheap sandbox instance")
	pflag.Int("webhook-port", 8888, "Port for webhook server")
	pflag.Int("health-port", 8080, "Port for health check server")
	pflag.Parse()

	// Bind CLI flags to Viper
	if err := viper.BindPFlag("username", pflag.Lookup("username")); err != nil {
		return nil, fmt.Errorf("error binding username flag: %v", err)
	}
	if err := viper.BindPFlag("api_key", pflag.Lookup("api-key")); err != nil {
		return nil, fmt.Errorf("error binding api-key flag: %v", err)
	}
	if err := viper.BindPFlag("client_ip", pflag.Lookup("client-ip")); err != nil {
		return nil, fmt.Errorf("error binding client-ip flag: %v", err)
	}
	if err := viper.BindPFlag("use_sandbox", pflag.Lookup("use-sandbox")); err != nil {
		return nil, fmt.Errorf("error binding use-sandbox flag: %v", err)
	}
	if err := viper.BindPFlag("server.webhook_port", pflag.Lookup("webhook-port")); err != nil {
		return nil, fmt.Errorf("error binding webhook-port flag: %v", err)
	}
	if err := viper.BindPFlag("server.health_port", pflag.Lookup("health-port")); err != nil {
		return nil, fmt.Errorf("error binding health-port flag: %v", err)
	}

	// Set default values
	viper.SetDefault("server.webhook_port", 8888)
	viper.SetDefault("server.health_port", 8080)

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %v", err)
	}

	// Validate required fields
	if config.Username == "" {
		return nil, fmt.Errorf("username configuration is required")
	}

	if config.APIKey == "" {
		return nil, fmt.Errorf("api_key configuration is required")
	}
	if config.ClientIP == "" {
		return nil, fmt.Errorf("client_ip configuration is required")
	}

	return &config, nil
}

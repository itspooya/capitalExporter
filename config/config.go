package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Configuration struct to hold all configuration values
type Configuration struct {
	Email    string
	Password string
	APIKey   string
	Port     string
	Demo     bool
	Interval int
	Debug    bool
}

// LoadConfiguration loads the configuration from environment variables and flags
func LoadConfiguration() (Configuration, error) {
	viper.SetEnvPrefix("EXPORTER")

	if err := viper.BindEnv("email"); err != nil {
		return Configuration{}, err
	}
	if err := viper.BindEnv("password"); err != nil {
		return Configuration{}, err
	}
	if err := viper.BindEnv("apikey"); err != nil {
		return Configuration{}, err
	}
	if err := viper.BindEnv("port"); err != nil {
		return Configuration{}, err
	}
	if err := viper.BindEnv("demo"); err != nil {
		return Configuration{}, err
	}
	if err := viper.BindEnv("interval"); err != nil {
		return Configuration{}, err
	}

	viper.SetDefault("port", "9682")
	viper.SetDefault("demo", false)
	viper.SetDefault("interval", 60)

	pflag.String("email", viper.GetString("email"), "Email address")
	pflag.Bool("demo", viper.GetBool("demo"), "Demo account")
	pflag.String("password", viper.GetString("password"), "Password")
	pflag.String("apikey", viper.GetString("apikey"), "API Key")
	pflag.String("port", viper.GetString("port"), "Port to listen on")
	pflag.Int("interval", viper.GetInt("interval"), "Interval to update metrics")
	pflag.Bool("debug", viper.GetBool("debug"), "Debug mode")

	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return Configuration{}, err
	}

	return Configuration{
		Email:    viper.GetString("email"),
		Password: viper.GetString("password"),
		APIKey:   viper.GetString("apikey"),
		Port:     viper.GetString("port"),
		Demo:     viper.GetBool("demo"),
		Interval: viper.GetInt("interval"),
		Debug:    viper.GetBool("debug"),
	}, nil
}

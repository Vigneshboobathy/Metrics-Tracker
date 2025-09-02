package config

import (
	"os"
	"github.com/spf13/viper"
)

type NodeConfig struct {
	Type    string `mapstructure:"type"`
	URL     string `mapstructure:"url"`
	Address string `mapstructure:"address"`
}

type AppConfig struct {
	Node     NodeConfig `mapstructure:"node"`
	Receiver string     `mapstructure:"receiver"`
}

func LoadConfig() (*AppConfig, error) {
	if cfgFile := os.Getenv("CONFIG_FILE"); cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("json")
		viper.AddConfigPath("./config")
		viper.AddConfigPath(".")
		viper.AddConfigPath("..")
		viper.AddConfigPath("../..")
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

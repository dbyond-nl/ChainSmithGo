package config

import (
	"github.com/spf13/viper"
	"log"
)

type CertificateConfig struct {
	CertPath   string `mapstructure:"cert_path"`
	KeyPath    string `mapstructure:"key_path"`
	CommonName string `mapstructure:"common_name"`
}

type Config struct {
	RootCAPath         string                        `mapstructure:"root_ca_path"`
	IntermediateCAPath string                        `mapstructure:"intermediate_ca_path"`
	Certificates       map[string]CertificateConfig `mapstructure:"certificates"`
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv() // Allow environment variable overrides

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	log.Printf("Loaded configuration from %s", viper.ConfigFileUsed())
	return &cfg, nil
}

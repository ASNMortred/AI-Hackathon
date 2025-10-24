package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	QiNiu     QiNiuConfig     `mapstructure:"qiniu"`
	Upload    UploadConfig    `mapstructure:"upload"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type QiNiuConfig struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
}

type UploadConfig struct {
	MaxSize       int64    `mapstructure:"max_size"`
	AllowedTypes  []string `mapstructure:"allowed_types"`
	UploadDir     string   `mapstructure:"upload_dir"`
}

func LoadConfig() (*Config, error) {
	var configPath string
	pflag.StringVar(&configPath, "config", "configs/config.yaml", "path to config file")
	pflag.Parse()

	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

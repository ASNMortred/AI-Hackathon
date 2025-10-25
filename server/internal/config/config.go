package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	QiNiu    QiNiuConfig    `mapstructure:"qiniu"`
	Minio    MinioConfig    `mapstructure:"minio"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Database DatabaseConfig `mapstructure:"database"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type QiNiuConfig struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
}

type MinioConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	UseSSL    bool   `mapstructure:"use_ssl"`
}

type UploadConfig struct {
	MaxSize      int64    `mapstructure:"max_size"`
	AllowedTypes []string `mapstructure:"allowed_types"`
	UploadDir    string   `mapstructure:"upload_dir"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

func LoadConfig() (*Config, error) {
	var configPath string
	pflag.StringVar(&configPath, "config", "configs/config.yaml", "path to config file")
	pflag.Parse()

	viper.SetConfigFile(configPath)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Expand ${ENV_VAR} placeholders and validate required fields
	expandEnvInConfig(&config)
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func expandEnvInConfig(cfg *Config) {
	cfg.Server.Port = os.ExpandEnv(cfg.Server.Port)
	cfg.Database.Host = os.ExpandEnv(cfg.Database.Host)
	cfg.Database.Port = os.ExpandEnv(cfg.Database.Port)
	cfg.Database.User = os.ExpandEnv(cfg.Database.User)
	cfg.Database.Password = os.ExpandEnv(cfg.Database.Password)
	cfg.Database.Name = os.ExpandEnv(cfg.Database.Name)

	cfg.Minio.Endpoint = os.ExpandEnv(cfg.Minio.Endpoint)
	cfg.Minio.AccessKey = os.ExpandEnv(cfg.Minio.AccessKey)
	cfg.Minio.SecretKey = os.ExpandEnv(cfg.Minio.SecretKey)
	cfg.Minio.Bucket = os.ExpandEnv(cfg.Minio.Bucket)
}

func validateConfig(cfg *Config) error {
	missing := []string{}
	if cfg.Database.Host == "" {
		missing = append(missing, "DATABASE_HOST")
	}
	if cfg.Database.Port == "" {
		missing = append(missing, "DATABASE_PORT")
	}
	if cfg.Database.User == "" {
		missing = append(missing, "DATABASE_USER")
	}
	if cfg.Database.Name == "" {
		missing = append(missing, "DATABASE_NAME")
	}
	// MinIO basic validations (optional but helpful)
	if cfg.Minio.Endpoint == "" {
		missing = append(missing, "MINIO_ENDPOINT")
	}
	if cfg.Minio.AccessKey == "" {
		missing = append(missing, "MINIO_ACCESS_KEY")
	}
	if cfg.Minio.SecretKey == "" {
		missing = append(missing, "MINIO_SECRET_KEY")
	}
	if cfg.Minio.Bucket == "" {
		missing = append(missing, "MINIO_BUCKET")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required env vars: %v", missing)
	}
	return nil
}

package config

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env                 string        `mapstructure:"ENV"`
	Timeout             time.Duration `mapstructure:"TIMEOUT"`
	LogLevel            string        `mapstructure:"LOGLEVEL"`
	PaginationLimit     int           `mapstructure:"PAGINATION_LIMIT"`
	Path                string        `mapstructure:"SERVER_PATH"`
	Port                int           `mapstructure:"SERVER_PORT"`
	PostgresPath        string        `mapstructure:"POSTGRES_PATH"`
	PostgresPort        int           `mapstructure:"POSTGRES_PORT"`
	PostgreUser         string        `mapstructure:"POSTGRES_USER"`
	PostgrePassword     string        `mapstructure:"POSTGRES_PASSWORD"`
	PostgreDatabaseName string        `mapstructure:"POSTGRES_DATABASE"`
}

func (c *Config) ServerToString() string {
	return fmt.Sprintf("%s:%d", c.Path, c.Port)
}

func Load() (*Config, error) {

	path, file := fetchConfigPath()

	if path == "" || file == "" {
		return &Config{}, errors.New("fetch config path error")
	}

	var cfg Config

	cfg, err := initViper(path, file, cfg)
	if err != nil {
		return &Config{}, fmt.Errorf("failed to init viper: %w", err)
	}

	slog.Info("config file load successfully", "log_level", cfg.LogLevel)

	return &cfg, nil
}

func fetchConfigPath() (path, file string) {

	flag.StringVar(&path, "config_path", "", "config file path")
	flag.StringVar(&file, "config_file", "", "config file name")
	flag.Parse()

	return path, file
}

func initViper(path string, file string, cfg Config) (Config, error) {
	viper.SetConfigName(file)
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

package config

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

var (
	DefaultGatewayConfig = GatewayConfig{
		MinFileSizeToSplit: 1024,
		NumChunks:          int64(6),
		Concurrency:        10,
		HTTP:               HTTP{ListenPort: ":8080"},
		DB: DB{
			Host:     "localhost",
			User:     "postgres",
			Database: "filestorage_gateway",
		},
	}

	DefaultStorageConfig = StorageConfig{
		StoragePath: "./local/storage_data1",
		HTTP:        HTTP{ListenPort: ":8081"},
	}
)

func Load(name string, in interface{}, configPaths ...string) error {
	viper.SetConfigName(name + ".config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./local/configs/default")

	for _, configPath := range configPaths {
		viper.AddConfigPath(configPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read in config: %w", err)
	}

	if err := viper.Unmarshal(&in); err != nil {
		return fmt.Errorf("failed to unmarshal config into struct: %w", err)
	}

	if err := validator.New().Struct(in); err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	return nil
}

type GatewayConfig struct {
	MinFileSizeToSplit int64
	NumChunks          int64
	Concurrency        int
	HTTP               HTTP
	DB                 DB
}

type StorageConfig struct {
	StoragePath string
	HTTP
}

type Storage struct {
	URL string `validate:"required"`
}

type HTTP struct {
	ListenPort string `validate:"required"`
}

type DB struct {
	Host     string `validate:"required"`
	Port     string
	User     string `validate:"required"`
	Password string
	Database string `validate:"required"`
}

func (db DB) GetConnectionString() string {
	s := fmt.Sprintf("host=%s", db.Host)
	if db.Port != "" {
		s = fmt.Sprintf("%s port=%s", s, db.Port)
	}
	if db.User != "" {
		s = fmt.Sprintf("%s user=%s", s, db.User)
	}
	if db.Database != "" {
		s = fmt.Sprintf("%s database=%s", s, db.Database)
	}
	if db.Password != "" {
		s = fmt.Sprintf("%s password=%s", s, db.Password)
	}

	return s
}

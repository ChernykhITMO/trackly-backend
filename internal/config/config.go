package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
)

type Config struct {
	AppName     string      `yaml:"app_name"`
	Port        string      `yaml:"port"`
	Database    DbConfig    `yaml:"database"`
	JwtSecret   string      `yaml:"jwt_secret"`
	MinioConfig MinioConfig `yaml:"minio"`
}

type DbConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbname"`
}

type MinioConfig struct {
	MinioEndpoint     string `yaml:"endpoint"`
	BucketName        string `yaml:"bucket_name"`
	MinioRootUser     string `yaml:"root_user"`
	MinioRootPassword string `yaml:"root_password"`
	MinioUseSSL       bool   `yaml:"use_ssl"`
}

func LoadConfig(filePath string) (*Config, error) {

	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading configs file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing configs file: %w", err)
	}
	config = Config{
		AppName: getEnv("APP_NAME", config.AppName),
		Port:    getEnv("APP_PORT", config.Port),
		Database: DbConfig{
			Host:     getEnv("DB_HOST", config.Database.Host),
			Port:     getEnv("DB_PORT", config.Database.Port),
			Username: getEnv("POSTGRES_USER", config.Database.Username),
			Password: getEnv("POSTGRES_PASSWORD", config.Database.Password),
			DbName:   getEnv("POSTGRES_DB", config.Database.DbName),
		},
		JwtSecret: getEnv("JWT_SECRET", config.JwtSecret),
		MinioConfig: MinioConfig{
			MinioEndpoint:     getEnv("MINIO_ENDPOINT", config.MinioConfig.MinioEndpoint),
			BucketName:        getEnv("MINIO_BUCKET_NAME", config.MinioConfig.BucketName),
			MinioRootUser:     getEnv("MINIO_ROOT_USER", config.MinioConfig.MinioRootUser),
			MinioRootPassword: getEnv("MINIO_ROOT_PASSWORD", config.MinioConfig.MinioRootPassword),
			MinioUseSSL:       getEnvAsBool("MINIO_USE_SSL", config.MinioConfig.MinioUseSSL),
		},
	}

	log.Printf("Loaded config: %+v", config.MinioConfig)

	return &config, nil
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if valueStr := getEnv(key, ""); valueStr != "" {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

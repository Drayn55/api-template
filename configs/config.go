// configs/config.go
package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config — all application configuration
type Config struct {
	DB         DBConfig
	JWT        JWTConfig
	MinIO      MinIOConfig
	ServerAddr string
}

type DBConfig struct {
	Host   string
	Port   string
	User   string
	Pass   string
	DBName string
}

type JWTConfig struct {
	Secret string
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
}

// DSN builds the MySQL/MariaDB connection string
func (d DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		d.User, d.Pass, d.Host, d.Port, d.DBName,
	)
}

// Load reads .env file and environment variables
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DB: DBConfig{
			Host:   getEnv("DB_HOST", "127.0.0.1"),
			Port:   getEnv("DB_PORT", "3306"),
			User:   getEnv("DB_USER", "root"),
			Pass:   getEnv("DB_PASSWORD", ""),
			DBName: getEnv("DB_NAME", "api_template_db"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
		},
		MinIO: MinIOConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", ""),
			AccessKey: getEnv("MINIO_ACCESS_KEY", ""),
			SecretKey: getEnv("MINIO_SECRET_KEY", ""),
			Bucket:    getEnv("MINIO_BUCKET", "api-template"),
		},
		ServerAddr: getEnv("SERVER_ADDR", ":8080"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.DB.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.DB.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.DB.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}


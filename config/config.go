package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AppEnv   string
	Port     string
	IP       string
	BcryptSalt int

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret string
	JWTExpire string

	EmailFrom string
	EmailUser string
	EmailPass string
	EmailHost string
	EmailPort int

	SuperAdminEmail    string
	SuperAdminPassword string
}

var cfg *Config

func Load() *Config {
	if cfg != nil {
		return cfg
	}

	salt, _ := strconv.Atoi(getEnv("BCRYPT_SALT_ROUNDS", "12"))
	emailPort, _ := strconv.Atoi(getEnv("EMAIL_PORT", "587"))

	cfg = &Config{
		AppEnv:     getEnv("APP_ENV", "development"),
		Port:       getEnv("PORT", "8080"),
		IP:         getEnv("IP_ADDRESS", "0.0.0.0"),
		BcryptSalt: salt,

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASS", "postgres"),
		DBName:     getEnv("DB_NAME", "go_backend"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret: getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpire: getEnv("JWT_EXPIRE_IN", "7d"),

		EmailFrom: getEnv("EMAIL_FROM", ""),
		EmailUser: getEnv("EMAIL_USER", ""),
		EmailPass: getEnv("EMAIL_PASS", ""),
		EmailHost: getEnv("EMAIL_HOST", "smtp.gmail.com"),
		EmailPort: emailPort,

		SuperAdminEmail:    getEnv("SUPER_ADMIN_EMAIL", "superadmin@gmail.com"),
		SuperAdminPassword: getEnv("SUPER_ADMIN_PASSWORD", "password@123"),
	}

	return cfg
}

func Get() *Config {
	if cfg == nil {
		return Load()
	}
	return cfg
}

func (c *Config) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

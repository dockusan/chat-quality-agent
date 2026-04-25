package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	// Server
	ServerPort string
	ServerHost string
	AppURL     string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Security
	JWTSecret     string
	EncryptionKey string // 32 bytes for AES-256-GCM

	// Rate limiting
	RateLimitPerIP   int // requests per minute
	RateLimitPerUser int

	// AI
	AIMaxTokens int // max tokens for AI responses

	// Environment
	Env string // "development" | "production"

	// Storage
	FileStoragePath string
}

func Load() (*Config, error) {
	dbURL := getEnvFirst([]string{"DATABASE_URL", "MYSQL_URL"}, "")
	dbURLCfg := parseDatabaseURL(dbURL)

	cfg := &Config{
		ServerPort:       getEnvFirst([]string{"SERVER_PORT", "PORT"}, "8080"),
		ServerHost:       getEnv("SERVER_HOST", "0.0.0.0"),
		AppURL:           GetAppURL(),
		DBHost:           getEnvFirst([]string{"DB_HOST", "MYSQLHOST"}, firstNonEmpty(dbURLCfg.host, "localhost")),
		DBPort:           getEnvFirst([]string{"DB_PORT", "MYSQLPORT"}, firstNonEmpty(dbURLCfg.port, "3306")),
		DBUser:           getEnvFirst([]string{"DB_USER", "MYSQLUSER"}, firstNonEmpty(dbURLCfg.user, "cqa")),
		DBPassword:       getEnvFirst([]string{"DB_PASSWORD", "MYSQLPASSWORD"}, dbURLCfg.password),
		DBName:           getEnvFirst([]string{"DB_NAME", "MYSQLDATABASE"}, firstNonEmpty(dbURLCfg.name, "cqa")),
		JWTSecret:        getEnv("JWT_SECRET", ""),
		EncryptionKey:    getEnv("ENCRYPTION_KEY", ""),
		RateLimitPerIP:   getEnvInt("RATE_LIMIT_PER_IP", 500),
		RateLimitPerUser: getEnvInt("RATE_LIMIT_PER_USER", 1000),
		AIMaxTokens:      getEnvInt("AI_MAX_TOKENS", 16384),
		Env:              getEnv("APP_ENV", "development"),
		FileStoragePath:  GetFileStoragePath(),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if len(cfg.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters for HS256 security")
	}
	if cfg.EncryptionKey == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY is required")
	}
	if len(cfg.EncryptionKey) != 32 {
		return nil, fmt.Errorf("ENCRYPTION_KEY must be exactly 32 bytes for AES-256-GCM")
	}
	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}

	return cfg, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func (c *Config) ListenAddr() string {
	return fmt.Sprintf("%s:%s", c.ServerHost, c.ServerPort)
}

func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvFirst(keys []string, fallback string) string {
	for _, key := range keys {
		if v := strings.TrimSpace(os.Getenv(key)); v != "" {
			return v
		}
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

type databaseURLConfig struct {
	host     string
	port     string
	user     string
	password string
	name     string
}

func parseDatabaseURL(raw string) databaseURLConfig {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return databaseURLConfig{}
	}

	u, err := url.Parse(raw)
	if err != nil {
		return databaseURLConfig{}
	}

	password, _ := u.User.Password()
	cfg := databaseURLConfig{
		host:     u.Hostname(),
		port:     u.Port(),
		user:     u.User.Username(),
		password: password,
		name:     strings.TrimPrefix(u.Path, "/"),
	}
	return cfg
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func GetAppURL() string {
	if v := strings.TrimSpace(os.Getenv("APP_URL")); v != "" {
		return strings.TrimRight(v, "/")
	}

	if domain := strings.TrimSpace(os.Getenv("RAILWAY_PUBLIC_DOMAIN")); domain != "" {
		domain = strings.TrimPrefix(domain, "https://")
		domain = strings.TrimPrefix(domain, "http://")
		return "https://" + strings.TrimRight(domain, "/")
	}

	return ""
}

func GetFileStoragePath() string {
	basePath := getEnvFirst([]string{"FILE_STORAGE_PATH", "RAILWAY_VOLUME_MOUNT_PATH"}, "/var/lib/cqa/files")
	cleaned := filepath.Clean(strings.TrimSpace(basePath))
	if cleaned == "." || cleaned == "" {
		return "/var/lib/cqa/files"
	}
	return cleaned
}

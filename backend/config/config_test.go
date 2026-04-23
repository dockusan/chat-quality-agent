package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Set required env vars
	os.Unsetenv("PORT")
	os.Unsetenv("MYSQLHOST")
	os.Unsetenv("MYSQLPORT")
	os.Unsetenv("MYSQLUSER")
	os.Unsetenv("MYSQLPASSWORD")
	os.Unsetenv("MYSQLDATABASE")
	os.Unsetenv("FILE_STORAGE_PATH")
	os.Unsetenv("RAILWAY_VOLUME_MOUNT_PATH")
	os.Unsetenv("APP_URL")
	os.Unsetenv("RAILWAY_PUBLIC_DOMAIN")
	os.Setenv("JWT_SECRET", "test-jwt-secret-at-least-32-chars-long")
	os.Setenv("ENCRYPTION_KEY", "12345678901234567890123456789012")
	os.Setenv("DB_PASSWORD", "testpassword")
	defer func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("ENCRYPTION_KEY")
		os.Unsetenv("DB_PASSWORD")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.ServerPort != "8080" {
		t.Errorf("Default ServerPort should be 8080, got %s", cfg.ServerPort)
	}
	if cfg.ServerHost != "0.0.0.0" {
		t.Errorf("Default ServerHost should be 0.0.0.0, got %s", cfg.ServerHost)
	}
	if cfg.DBName != "cqa" {
		t.Errorf("Default DBName should be cqa, got %s", cfg.DBName)
	}
	if cfg.FileStoragePath != "/var/lib/cqa/files" {
		t.Errorf("Default FileStoragePath should be /var/lib/cqa/files, got %s", cfg.FileStoragePath)
	}
}

func TestLoadConfigMissingRequired(t *testing.T) {
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("ENCRYPTION_KEY")
	os.Unsetenv("DB_PASSWORD")

	_, err := Load()
	if err == nil {
		t.Fatal("Load should fail with missing required vars")
	}
}

func TestDSN(t *testing.T) {
	cfg := &Config{
		DBUser:     "testuser",
		DBPassword: "testpass",
		DBHost:     "localhost",
		DBPort:     "3306",
		DBName:     "testdb",
	}
	dsn := cfg.DSN()
	expected := "testuser:testpass@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	if dsn != expected {
		t.Errorf("DSN = %q, want %q", dsn, expected)
	}
}

func TestLoadConfigSupportsRailwayEnv(t *testing.T) {
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")

	os.Setenv("PORT", "9090")
	os.Setenv("MYSQLHOST", "railway-mysql")
	os.Setenv("MYSQLPORT", "3307")
	os.Setenv("MYSQLUSER", "railway")
	os.Setenv("MYSQLPASSWORD", "railway-password")
	os.Setenv("MYSQLDATABASE", "railwaydb")
	os.Setenv("JWT_SECRET", "test-jwt-secret-at-least-32-chars-long")
	os.Setenv("ENCRYPTION_KEY", "12345678901234567890123456789012")
	os.Setenv("RAILWAY_PUBLIC_DOMAIN", "cqa-production.up.railway.app")
	os.Setenv("RAILWAY_VOLUME_MOUNT_PATH", "/data")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("MYSQLHOST")
		os.Unsetenv("MYSQLPORT")
		os.Unsetenv("MYSQLUSER")
		os.Unsetenv("MYSQLPASSWORD")
		os.Unsetenv("MYSQLDATABASE")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("ENCRYPTION_KEY")
		os.Unsetenv("RAILWAY_PUBLIC_DOMAIN")
		os.Unsetenv("RAILWAY_VOLUME_MOUNT_PATH")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.ServerPort != "9090" {
		t.Errorf("ServerPort = %q, want %q", cfg.ServerPort, "9090")
	}
	if cfg.DBHost != "railway-mysql" || cfg.DBPort != "3307" || cfg.DBUser != "railway" || cfg.DBPassword != "railway-password" || cfg.DBName != "railwaydb" {
		t.Errorf("Unexpected Railway DB config: %+v", cfg)
	}
	if cfg.AppURL != "https://cqa-production.up.railway.app" {
		t.Errorf("AppURL = %q, want %q", cfg.AppURL, "https://cqa-production.up.railway.app")
	}
	if cfg.FileStoragePath != "/data" {
		t.Errorf("FileStoragePath = %q, want %q", cfg.FileStoragePath, "/data")
	}
}

func TestIsProduction(t *testing.T) {
	cfg := &Config{Env: "production"}
	if !cfg.IsProduction() {
		t.Error("Should be production")
	}

	cfg.Env = "development"
	if cfg.IsProduction() {
		t.Error("Should not be production")
	}
}

func TestGetAppURLPrefersExplicitEnv(t *testing.T) {
	os.Setenv("APP_URL", "https://custom.example.com/")
	os.Setenv("RAILWAY_PUBLIC_DOMAIN", "ignored.up.railway.app")
	defer func() {
		os.Unsetenv("APP_URL")
		os.Unsetenv("RAILWAY_PUBLIC_DOMAIN")
	}()

	if got := GetAppURL(); got != "https://custom.example.com" {
		t.Errorf("GetAppURL() = %q, want %q", got, "https://custom.example.com")
	}
}

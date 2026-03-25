package main

import (
	"log"

	"github.com/vietbui/chat-quality-agent/api"
	"github.com/vietbui/chat-quality-agent/api/handlers"
	"github.com/vietbui/chat-quality-agent/api/middleware"
	"github.com/vietbui/chat-quality-agent/config"
	"github.com/vietbui/chat-quality-agent/db"
	"github.com/vietbui/chat-quality-agent/engine"
)

var version = "dev"

func main() {
	log.Printf("Chat Quality Agent %s", version)
	handlers.AppVersion = version

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize JWT
	middleware.SetJWTSecret(cfg.JWTSecret)

	// Connect database
	if err := db.Connect(cfg.DSN(), cfg.IsProduction()); err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Start scheduler
	scheduler, err := engine.NewScheduler(cfg)
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}
	scheduler.Start()
	defer scheduler.Stop()

	// Setup router
	router := api.SetupRouter(cfg)

	// Start server
	log.Printf("CQA server starting on %s (env: %s)", cfg.ListenAddr(), cfg.Env)
	if err := router.Run(cfg.ListenAddr()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

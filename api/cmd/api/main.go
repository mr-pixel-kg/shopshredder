package main

import (
	"log"

	"github.com/manuel/shopware-testenv-platform/api/internal/config"
	"github.com/manuel/shopware-testenv-platform/api/internal/database"
	httpserver "github.com/manuel/shopware-testenv-platform/api/internal/http"
)

func main() {
	cfg := config.MustLoad()

	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	server, err := httpserver.NewServer(cfg, db)
	if err != nil {
		log.Fatalf("create server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("start server: %v", err)
	}
}

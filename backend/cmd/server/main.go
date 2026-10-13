package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Max-65/libraryInterface/internal/config"
	"github.com/Max-65/libraryInterface/internal/db"
	"github.com/Max-65/libraryInterface/internal/httpapi"
	"github.com/Max-65/libraryInterface/internal/repository"
)

func main() {
	cfg := config.Load()

	pool, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	repo := repository.New(pool)
	handler := httpapi.NewRouter(repo, cfg.AllowedOrigin, logger)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("api listening on :" + cfg.Port)
	log.Fatal(srv.ListenAndServe())
}

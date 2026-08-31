package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/DragonZero000/url-shortener/internal/config"
	"github.com/DragonZero000/url-shortener/internal/handler"
	"github.com/DragonZero000/url-shortener/internal/repository"
	"github.com/DragonZero000/url-shortener/internal/service"
)

func main() {
	config, err := config.Load()
	if err != nil {
		slog.Error("fail to load config", "error", err)
		return
	}
	db, err := repository.Connection(config.DatabaseURL)
	if err != nil {
		slog.Error("fail to connect to database", "error", err)
		return
	}
	defer db.Close()
	slog.Info("database connected")
	err = repository.RunMigrations(config.DatabaseURL)
	if err != nil {
		slog.Error("fail to migrate", "error", err)
		return
	}
	slog.Info("migration completed")
	repos := repository.NewRepository(db)
	svc := service.NewService(repos)
	hand := handler.NewHandler(svc, config.WebBaseURL)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", hand.PostShorten)
	mux.HandleFunc("GET /{code}", hand.GetByShorten)
	err = http.ListenAndServe(config.ServerPort, mux)
	if err != nil {
		log.Fatal(err)
	}
}

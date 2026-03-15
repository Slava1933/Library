package main

import (
	"context"
	"library/internal/handlers"
	"library/internal/logger"
	"library/internal/repository"

	"github.com/gorilla/mux"
)

func main() {
	logger, logFileClose, err := logger.NewLogger("INFO")
	if err != nil {
		panic(err)
	}
	defer logFileClose()
	logger.Info("app started")

	ctx := context.Background()

	pool, err := repository.Connect(ctx)
	if err != nil {
		logger.Fatal("Не удалось подключиться к БД")
	}
	defer pool.Close()

	repo := repository.NewDocumentRepository(pool, logger)
	h := handlers.NewHandlers(repo, logger)
	router := mux.NewRouter()
	router.Path("/api/disciplines").HandlerFunc(h.GetDisciplinesHandler)
	router.Path("/api/disciplines/{id}/documents").HandlerFunc(h.GetDocsByDiscipline)
	router.Path("/api/documents/{id}/download").HandlerFunc(h.GetDocument)
}

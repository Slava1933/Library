package main

import (
	"context"
	"library/internal/handlers"
	"library/internal/logger"
	"library/internal/repository"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
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
	h := handlers.NewHandlers(repo, logger, pool)
	router := mux.NewRouter()
	router.Path("/api/disciplines").HandlerFunc(h.GetDisciplinesHandler)
	router.Path("/api/disciplines/{discipline_id}/documents").HandlerFunc(h.GetDocsByDiscipline)
	router.Path("/api/documents/{id}/download").HandlerFunc(h.GetDocument)
	router.Path("/api/admin/upload_document").Methods("POST")
	router.Path("/api/admin/upload_discipline").Methods("POST")
	router.Path("/api/admin/documents").Methods("GET")
	router.Path("/api/admin/documents/{id}").Methods("DELETE")
	router.Path("/api/admin/documents").Methods("PATCH")
	router.Path("/api/admin/disciplines/{id}").Methods("DELETE")
	router.Path("/api/admin/disciplines").Methods("PATCH")
	fs := http.FileServer(http.Dir("./front"))
	router.PathPrefix("/").Handler(fs)
	if er := http.ListenAndServe(":8080", router); er != nil {
		logger.Fatal("Cant start the server", zap.Error(er))
	}

}

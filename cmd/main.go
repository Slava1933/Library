package main

import (
	"context"
	"library/auth"
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

	adminrepo := repository.NewAdminRepo(pool, logger)
	a := handlers.NewAdminHandlers(adminrepo, logger, pool)

	router := mux.NewRouter()
	router.Path("/api/disciplines").HandlerFunc(h.GetDisciplinesHandler)
	router.Path("/api/disciplines/{discipline_id}/documents").HandlerFunc(h.GetDocsByDiscipline)
	router.Path("/api/documents/{id}/download").HandlerFunc(h.GetDocument)

	admin := router.PathPrefix("/api/admin").Subrouter()
	admin.Use(Auth)
	admin.Path("/upload_document").Methods("POST").HandlerFunc(a.UploadDocumentHandler)
	admin.Path("/upload_discipline").Methods("POST").HandlerFunc(a.UploadDisciplineHandler)
	admin.Path("/documents").Methods("GET").HandlerFunc(a.GetAllDocumentsHandler)
	admin.Path("/documents/{id}").Methods("DELETE").HandlerFunc(a.DeleteDocumentHandler)
	admin.Path("/documents").Methods("PATCH").HandlerFunc(a.UpdateDocumentHandler)
	admin.Path("/disciplines/{id}").Methods("DELETE").HandlerFunc(a.DeleteDisciplineHandler)
	admin.Path("/disciplines").Methods("PATCH").HandlerFunc(a.UpdateDisciplineHandler)

	router.Path("/api/admin/login").HandlerFunc(a.AuthHandler)

	fs := http.FileServer(http.Dir("./front"))
	router.PathPrefix("/").Handler(fs)
	if er := http.ListenAndServe(":8080", router); er != nil {
		logger.Fatal("Cant start the server", zap.Error(er))
	}

}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Admin-Token")
		if token != auth.CurrentToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

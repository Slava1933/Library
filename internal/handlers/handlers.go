package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"library/internal/interfaces"
	"library/internal/models"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Handlers struct {
	repo interfaces.Repository_interface
	log  *zap.Logger
	pool *pgxpool.Pool
}

func NewHandlers(repo interfaces.Repository_interface, log *zap.Logger, pool *pgxpool.Pool) *Handlers {
	return &Handlers{repo: repo, log: log, pool: pool}
}

func (h *Handlers) GetDisciplinesHandler(w http.ResponseWriter, r *http.Request) {
	disciplines, err := h.repo.FindAllDisciplines(r.Context())
	if err != nil {
		h.log.Error("Cant get the disciplines", zap.Error(err))
		http.Error(w, "Cant get the disciplines from server", http.StatusInternalServerError)
		return
	}

	if disciplines == nil {
		disciplines = []models.Discipline{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(disciplines); err != nil {
		h.log.Error("Failed to encode disciplines to json", zap.Error(err))
		http.Error(w, "Failed to encode disciplines to json", http.StatusInternalServerError)
		return
	}
	h.log.Info("Get disciplines was successfully ended", zap.Int("count:", len(disciplines)))
}

func (h *Handlers) GetDocsByDiscipline(w http.ResponseWriter, r *http.Request) {
	pathparts := strings.Split(r.URL.Path, "/")

	if len(pathparts) < 4 {
		h.log.Error("Bad request, small path")
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id := pathparts[3]
	DisciplineID, err := strconv.Atoi(id)
	if err != nil {
		h.log.Error("Cant convert DisciplineID to integer", zap.Error(err))
		http.Error(w, "Cant convert DisciplineID to integer", http.StatusInternalServerError)
		return
	}

	Documents, err := h.repo.FindDocumentsByDiscipline(r.Context(), DisciplineID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.log.Info("Documents not found", zap.Int("Discipline id", DisciplineID))
			http.Error(w, "Documents not found", http.StatusNotFound)
			return
		}
		h.log.Error("Cant get the documents of discipline",
			zap.Int("discipline id:", DisciplineID), zap.Error(err))
		http.Error(w, "Cant get the documents from server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(Documents); err != nil {
		h.log.Error("Failed to encode documents to json", zap.Error(err))
		http.Error(w, "Failed to encode documents to json", http.StatusInternalServerError)
		return
	}
	h.log.Info("Get documents was succcessfully ended", zap.Int("discipline:", DisciplineID),
		zap.Int("counts", len(Documents)))
}

func (h *Handlers) GetDocument(w http.ResponseWriter, r *http.Request) {
	query := `UPDATE documents 
	SET download_count = download_count + 1
	WHERE id = $1`
	pathparts := strings.Split(r.URL.Path, "/")

	if len(pathparts) < 4 {
		h.log.Error("Bad request, small path")
		http.Error(w, "Invaild URL", http.StatusBadRequest)
		return
	}

	id := pathparts[3]
	DocumentID, err := strconv.Atoi(id)
	if err != nil {
		h.log.Error("Cant convert DocumentID to integer", zap.Error(err))
		http.Error(w, "Cant convert DocumenID to integer", http.StatusInternalServerError)
		return
	}
	Document, err := h.repo.FindDocument(r.Context(), DocumentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.log.Info("Document not found", zap.Int("id", DocumentID))
			http.Error(w, "Document not found", http.StatusNotFound)
			return
		}
		h.log.Error("Cant get document to download", zap.Int("Document ID", DocumentID),
			zap.Error(err))
		http.Error(w, "Cant get the document to download from server", http.StatusInternalServerError)
		return
	}
	file, er := os.Open(Document.Filepath)
	if er != nil {
		h.log.Error("Failed to Open file", zap.Error(er))
		http.Error(w, "File not found", http.StatusNotFound)
	}
	defer file.Close()

	w.Header().Set("Content-type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment;filename="+Document.Title)

	_, err = io.Copy(w, file)
	if err != nil {
		h.log.Error("Failed to send file", zap.Error(err))
		return
	}
	_, err = h.pool.Exec(r.Context(), query, DocumentID)
	if err != nil {
		h.log.Error("Failed to up download_counts", zap.Error(err))
	}
	h.log.Info("Get document to download was successfully ended", zap.Int("Document ID", DocumentID))
}

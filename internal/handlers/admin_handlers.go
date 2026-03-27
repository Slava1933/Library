package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"library/auth"
	"library/internal/errs"
	"library/internal/interfaces"
	"library/internal/models"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type AdminHandlers struct {
	repo interfaces.Admin_interface
	log  *zap.Logger
	pool *pgxpool.Pool
}

func NewAdminHandlers(repo interfaces.Admin_interface, log *zap.Logger, pool *pgxpool.Pool) *AdminHandlers {
	return &AdminHandlers{repo: repo, log: log, pool: pool}
}

func generateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

func (a *AdminHandlers) AuthHandler(w http.ResponseWriter, r *http.Request) {
	var path models.Auth
	var login models.Login
	if err := json.NewDecoder(r.Body).Decode(&path); err != nil {
		a.log.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	adminPass := os.Getenv("ADMIN_PASS")
	if path.Pass == adminPass {
		auth.CurrentToken = generateToken()
		login.Token = auth.CurrentToken
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(login); err != nil {
			a.log.Error("Failed to encode", zap.Error(err))
			http.Error(w, "Failed to encode", http.StatusInternalServerError)
			return
		}

	} else {
		a.log.Info("Invalid password")
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}
}

func (a *AdminHandlers) GetAllDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	documents, err := a.repo.GetAllDocuments(r.Context())
	if err != nil {
		a.log.Error("Failed to get documents", zap.Error(err))
		http.Error(w, "Failed to get documents from server", http.StatusInternalServerError)
		return
	}

	if documents == nil {
		documents = []models.Document{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(documents); err != nil {
		a.log.Error("Failed to encode documents to json", zap.Error(err))
		http.Error(w, "Failed to encode documents to json", http.StatusInternalServerError)
		return
	}

	a.log.Info("Get documents was successfully ended", zap.Int("count:", len(documents)))
}

func (a *AdminHandlers) DeleteDocumentHandler(w http.ResponseWriter, r *http.Request) {
	pathparts := strings.Split(r.URL.Path, "/")

	if len(pathparts) < 5 {
		a.log.Error("Bad request, small URL path")
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id := pathparts[4]
	DocId, err := strconv.Atoi(id)
	if err != nil {
		a.log.Error("Failed to convert DocID to int", zap.Error(err))
		http.Error(w, "Failed to convert DocID to int", http.StatusInternalServerError)
		return
	}

	err = a.repo.DeleteDocument(r.Context(), DocId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			a.log.Info("Document not found")
			http.Error(w, "Document not found", http.StatusNotFound)
			return
		}
		a.log.Error("Failed to delete document", zap.Int("doc ID: ", DocId),
			zap.Error(err))
		http.Error(w, "Failed to delete document", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	a.log.Info("Delete document was successfully ended", zap.Int("DocID: ", DocId))
}

func (a *AdminHandlers) DeleteDisciplineHandler(w http.ResponseWriter, r *http.Request) {
	pathparts := strings.Split(r.URL.Path, "/")

	if len(pathparts) < 5 {
		a.log.Error("Bad request, small URL path")
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id := pathparts[4]
	DiscID, err := strconv.Atoi(id)
	if err != nil {
		a.log.Error("Failed to convert DiscID to int", zap.Error(err))
		http.Error(w, "Failed to convert DiscID to int", http.StatusInternalServerError)
		return
	}

	err = a.repo.DeleteDiscipline(r.Context(), DiscID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			a.log.Info("Discipline not found")
			http.Error(w, "Discipline not found", http.StatusNotFound)
			return
		}
		a.log.Error("Failed to delete discipline", zap.Int("DisID: ", DiscID), zap.Error(err))
		http.Error(w, "Failed to delete discsipline", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	a.log.Info("Delete discipline was successfully ended", zap.Int("DiscID: ", DiscID))
}

func (a *AdminHandlers) UpdateDisciplineHandler(w http.ResponseWriter, r *http.Request) {
	var disc models.Discipline
	err := json.NewDecoder(r.Body).Decode(&disc)
	if err != nil {
		a.log.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	discipline, er := a.repo.UpdateDiscipline(r.Context(), disc)
	if er != nil {
		if errors.Is(er, errs.ErrDiscNotFound) {
			a.log.Error("Discipline not found", zap.Int("ID:", discipline.ID))
			http.Error(w, "Discipline not found", http.StatusNotFound)
			return
		}
		a.log.Error("Failed to update discipline", zap.Error(er))
		http.Error(w, "Failed to update discipline", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(discipline); err != nil {
		a.log.Error("Failed to encode discipline", zap.Error(err))
		http.Error(w, "Failed to encode discipline", http.StatusInternalServerError)
		return
	}

	a.log.Info("Update discipline was successfully ended")
}

func (a *AdminHandlers) UpdateDocumentHandler(w http.ResponseWriter, r *http.Request) {
	var doc models.Document
	err := json.NewDecoder(r.Body).Decode(&doc)
	if err != nil {
		a.log.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	document, er := a.repo.UpdateDocument(r.Context(), doc)
	if er != nil {
		if errors.Is(er, errs.ErrDocNotFound) {
			a.log.Error("Document not found", zap.Int("ID:", document.ID))
			http.Error(w, "Document not found", http.StatusNotFound)
			return
		}
		a.log.Error("Failed to update document", zap.Error(er))
		http.Error(w, "Failed to update document", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(document); err != nil {
		a.log.Error("Failed to encode document", zap.Error(err))
		http.Error(w, "Failed to encode document", http.StatusInternalServerError)
		return
	}

	a.log.Info("Update document was successfully ended")
}

func (a *AdminHandlers) UploadDocumentHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		a.log.Error("Failed to parse data", zap.Error(err))
		http.Error(w, "Failed to parse data", http.StatusBadRequest)
		return
	}
	var document models.Document
	strdisciplineid := r.FormValue("discipline_id")
	discipline_id, err := strconv.Atoi(strdisciplineid)
	if err != nil {
		a.log.Error("Failed to convert discipline_id to int", zap.String("Operation: ", "Upload Document"), zap.Error(err))
		http.Error(w, "Failed to convert discipline_id to int", http.StatusInternalServerError)
		return
	}
	title := r.FormValue("title")
	document.DisciplineID = discipline_id
	document.Download_count = 0

	reader, header, err := r.FormFile("file")
	if err != nil {
		a.log.Error("Failed to get file", zap.Error(err))
		http.Error(w, "Failed to get file", http.StatusInternalServerError)
		return
	}
	ext := filepath.Ext(header.Filename) // тип файла
	uniqueName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	File_name := strings.Join([]string{title, ext}, "")
	document.Title = File_name
	file_path := fmt.Sprintf("/home/slava/uploads/%s", uniqueName)
	document.Filepath = file_path
	file, err := os.Create(file_path)
	if err != nil {
		a.log.Error("Failed to create file on server", zap.Error(err))
		http.Error(w, "Failed to create file on server", http.StatusInternalServerError)
		return
	}
	if _, err = io.Copy(file, reader); err != nil {
		a.log.Error("Failed to write file", zap.Error(err))
		http.Error(w, "Failed to write file on server", http.StatusInternalServerError)
		return
	}
	err = a.repo.UploadDocument(r.Context(), document)
	if err != nil {
		a.log.Error("Failed to add to db", zap.String("Operation: ", "UploadDocument"), zap.Error(err))
		http.Error(w, "Failed to add to db", http.StatusInternalServerError)
		return
	}
}

func (a *AdminHandlers) UploadDisciplineHandler(w http.ResponseWriter, r *http.Request) {
	var discipline models.CreateDiscipline
	if err := json.NewDecoder(r.Body).Decode(&discipline); err != nil {
		a.log.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if err := a.repo.UploadDiscipline(r.Context(), discipline); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			a.log.Error("Discipline already exists", zap.Error(err))
			http.Error(w, "Already exists", 409)
			return
		}
		a.log.Error("Failed to add to db", zap.String("Operation: ", "UploadDiscipline"), zap.Error(err))
		http.Error(w, "Failed to add to db", http.StatusInternalServerError)
		return
	}
}

func (a *AdminHandlers) GetDocumentHandler(w http.ResponseWriter, r *http.Request) {
	pathparts := strings.Split(r.URL.Path, "/")

	if len(pathparts) < 5 {
		a.log.Error("Bad request, small URL path")
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id := pathparts[4]
	fmt.Printf("id: %s /  // // / // ////", id)
	DocID, err := strconv.Atoi(id)
	fmt.Printf("Docid: %d", DocID)
	if err != nil {
		a.log.Error("Failed to convert DocID to int", zap.Error(err))
		http.Error(w, "Failed to convert DocID to int", http.StatusInternalServerError)
		return
	}
	document := a.repo.GetDocument(r.Context(), DocID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(document); err != nil {
		a.log.Error("Failed to encode document", zap.Error(err))
		http.Error(w, "Failed to encode document", http.StatusInternalServerError)
		return
	}
}

func (a *AdminHandlers) GET_Download_Count_Handler(w http.ResponseWriter, r *http.Request) {
	type count struct {
		Download_count int
	}
	var dwc count
	download_count := a.repo.GET_Download_Count(r.Context())
	dwc.Download_count = download_count
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dwc); err != nil {
		a.log.Error("Failed to get download count", zap.Error(err))
		http.Error(w, "Failed to get download count", http.StatusInternalServerError)
		return
	}
	a.log.Info("Getting Download count was successfully ended")
}

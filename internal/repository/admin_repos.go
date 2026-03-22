package repository

import (
	"context"
	"library/internal/errs"
	"library/internal/models"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type AdminRepo struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

func NewAdminRepo(pool *pgxpool.Pool, log *zap.Logger) *AdminRepo {
	return &AdminRepo{pool: pool, log: log}
}

func (a *AdminRepo) GetAllDocuments(ctx context.Context) ([]models.Document, error) {
	query := `
	SELECT * FROM documents
	`
	rows, err := a.pool.Query(ctx, query)
	if err != nil {
		a.log.Error("Failed to query documents", zap.String("operation: ", "GetAllDocuments"),
			zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	documents := make([]models.Document, 0)
	for rows.Next() {
		var document models.Document
		err := rows.Scan(&document.ID, &document.DisciplineID,
			&document.Title, &document.Filepath, &document.Download_count)
		if err != nil {
			a.log.Error("Failed to Scan document", zap.Error(err))
			return nil, err
		}
		documents = append(documents, document)
	}
	a.log.Info("Get documents was successfully ended")
	return documents, nil
}

func (a *AdminRepo) DeleteDocument(ctx context.Context, ID int) error {
	query_get_path := `
	SELECT file_path FROM documents 
	WHERE id = $1
	`
	row := a.pool.QueryRow(ctx, query_get_path, ID)
	var path string
	err := row.Scan(&path)
	if err != nil {
		a.log.Error("Failed to scan row", zap.String("Operation:", "Delete document"), zap.Error(err))
	}

	os.Remove(path)

	query := `
	DELETE FROM documents 
	WHERE id = $1
	`
	_, err = a.pool.Exec(ctx, query, ID)
	if err != nil {
		a.log.Error("Failed to Delete document", zap.Error(err))
	}
	a.log.Info("Successfully deleted document", zap.Int("With ID: ", ID))
	return nil
}

func (a *AdminRepo) DeleteDiscipline(ctx context.Context, DisciplineID int) error {
	query_get_path := `
	SELECT file_path FROM documents
	WHERE discipline_id = $1;
	`
	paths := make([]string, 0)
	rows, err := a.pool.Query(ctx, query_get_path, DisciplineID)
	if err != nil {
		a.log.Error("Failes to query document rows", zap.String("Operation: ", "Get paths"), zap.Error(err))
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var path string
		err := rows.Scan(&path)
		if err != nil {
			a.log.Error("Failed to scan row", zap.Error(err))
		}
		paths = append(paths, path)
	}

	query_delete := `
	DELETE FROM disciplines 
	WHERE id = $1
	`
	_, err = a.pool.Exec(ctx, query_delete, DisciplineID)
	if err != nil {
		a.log.Error("Failed to delete discipline", zap.Error(err))
	}

	for _, str := range paths {
		os.Remove(str)
	}

	a.log.Info("Delete discipline was successfully ended")
	return nil

}

func (a *AdminRepo) UploadDocument(ctx context.Context, document models.Document) error {
	query := `
	INSERT INTO documents
	(discipline_id, title, file_path)
	VALUES($1, $2, $3);
	`
	_, err := a.pool.Exec(ctx, query, document.DisciplineID, document.Title, document.Filepath)
	if err != nil {
		a.log.Error("Failed to Upload document", zap.Error(err))
	}
	return err
}

func (a *AdminRepo) UploadDiscipline(ctx context.Context, discipline models.CreateDiscipline) error {
	query := `
	INSERT INTO disciplines
	(title)
	VALUES($1);
	`
	_, err := a.pool.Exec(ctx, query, discipline.Title)
	if err != nil {
		a.log.Error("Failed to Upload discipline", zap.Error(err))
	}
	return err
}

func (a *AdminRepo) UpdateDiscipline(ctx context.Context, discipline models.Discipline) (models.Discipline, error) {
	query := `
	UPDATE disciplines
	SET title = $1 WHERE id = $2;
	`
	tag, err := a.pool.Exec(ctx, query, discipline.Title, discipline.ID)
	if err != nil {
		a.log.Error("Database error", zap.String("Operation: ", "UpdateDiscipline"), zap.Error(err))
	}
	if tag.RowsAffected() == 0 {
		a.log.Error("There is not discipline with that id")
		return discipline, errs.ErrDiscNotFound
	}
	query_get := `
	SELECT title FROM disciplines WHERE id = $1;
	`
	row := a.pool.QueryRow(ctx, query_get, discipline.ID)

	row.Scan(&discipline.Title)

	a.log.Info("Update discipline successfully ended")
	return discipline, nil
}

func (a *AdminRepo) UpdateDocument(ctx context.Context, document models.Document) (models.Document, error) {
	query := `
	UPDATE documents
	SET discipline_id = $1, title = $2, file_path = $3, 
	WHERE id = $4;
	`
	tag, err := a.pool.Exec(ctx, query, document.DisciplineID, document.Title, document.Filepath, document.ID)
	if err != nil {
		a.log.Error("Database error", zap.String("Operation: ", "Update Document"), zap.Error(err))
	}

	if tag.RowsAffected() == 0 {
		a.log.Error("There is no document with that id")
		return document, errs.ErrDocNotFound
	}

	query_get := `
	SELECT discipline_id, title, file_path WHERE id = $1;
	`
	row := a.pool.QueryRow(ctx, query_get, document.ID)
	row.Scan(&document.DisciplineID, &document.Title, &document.Filepath)
	a.log.Info("Update document was successfully ended")
	return document, nil
}

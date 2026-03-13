package repository

import (
	"context"
	"library/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DocumentRepository struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

func NewDocumentRepository(pool *pgxpool.Pool, log *zap.Logger) *DocumentRepository {
	return &DocumentRepository{pool: pool, log: log}
}

func (r *DocumentRepository) FindAllDisciplines(ctx context.Context) ([]models.Discipline, error) {
	query := `SELECT id, title FROM disciplines`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		r.log.Error("Failed to query disciplines", zap.String("operation", "FindAllDisciplines"),
			zap.Error(err))
		return nil, err
	}

	defer rows.Close()

	disciplines := make([]models.Discipline, 0)
	for rows.Next() {
		var discipline models.Discipline
		err := rows.Scan(&discipline.ID, &discipline.Title)
		if err != nil {
			r.log.Error("Failed to scan discipline row",
				zap.Error(err))
			return nil, err
		}
		disciplines = append(disciplines, discipline)
	}
	return disciplines, nil
}

func (r *DocumentRepository) FindDocumentsByDiscipline(ctx context.Context, DisciplineID int64) ([]models.Document, error) {
	query := `SELECT id, title, file_path FROM documents WHERE discipline_id = $1`

	rows, err := r.pool.Query(ctx, query, DisciplineID)
	if err != nil {
		r.log.Error("Failed to query documents", zap.String("operation", "FindDocumentsByDiscipline"),
			zap.String("id", string(DisciplineID)), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	documents := make([]models.Document, 0)
	for rows.Next() {
		var document models.Document
		err := rows.Scan(&document.ID, &document.Title, &document.Filepath)
		if err != nil {
			r.log.Error("Failed to scan document row", zap.Error(err))
			return nil, err
		}
		documents = append(documents, document)
	}

	return documents, nil
}

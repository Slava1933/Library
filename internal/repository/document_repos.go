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

func (r *DocumentRepository) FindByDiscipline(ctx context.Context, discipline string) ([]models.Document, error) {
	query := `SELECT * FROM documents WHERE discipline = $1`

	rows, err := r.pool.Query(ctx, query, discipline)
	if err != nil {
		r.log.Error("Failed to query documents", zap.String("operation", "FindByDiscipline"),
			zap.String("discipline", discipline), zap.Error(err))
		return nil, err
	}

	defer rows.Close()

	documents := make([]models.Document, 0)
	for rows.Next() {
		var document models.Document
		err := rows.Scan(&document.ID, &document.Discipline, &document.Filepath)
		if err != nil {
			r.log.Error("Failed to scan document row",
				zap.String("discipline", discipline),
				zap.Error(err),
			)
			return nil, err
		}
		documents = append(documents, document)
	}
	return documents, nil
}

func (r *DocumentRepository) FindByID(ctx context.Context, id int64) (models.Document, error) {
	query := `SELECT * FROM documents WHERE id = $1`

	row, err := r.pool.Query(ctx, query, id)
	if err != nil {
		r.log.Error("Failed to query documents", zap.String("operation", "FindByID"),
			zap.String("id", string(id)), zap.Error(err))
		return models.Document{}, err
	}
	defer row.Close()
	var document models.Document
	er := row.Scan(&document.ID, &document.Discipline, &document.Filepath)
	if er != nil {
		r.log.Error("Failed to scan document row",
			zap.String("id", string(id)),
			zap.Error(err),
		)
		return models.Document{}, err
	}
	return document, nil
}

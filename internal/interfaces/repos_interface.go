package interfaces

import (
	"context"
	"library/internal/models"
)

type Repository_interface interface {
	FindAllDisciplines(ctx context.Context) ([]models.Discipline, error)
	FindDocumentsByDiscipline(ctx context.Context, DisciplineID int64) ([]models.Document, error)
}

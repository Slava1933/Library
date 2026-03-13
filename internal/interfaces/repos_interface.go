package interfaces

import (
	"context"
	"library/internal/models"
)

type Repository_interface interface {
	FindAllDisciplines(ctx context.Context) ([]models.Discipline, error)
	FindDocumentsByDiscipline(ctx context.Context, DisciplineID int) ([]models.Document, error)
	FindDocument(ctx context.Context, ID int) (models.Document, error)
}

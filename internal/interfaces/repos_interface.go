package interfaces

import (
	"context"
	"library/internal/models"
)

type Repository_interface interface {
	FindByDiscipline(ctx context.Context, discipline string) ([]models.Document, error)
	FindByID(ctx context.Context, id int64) (models.Document, error)
}

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

type Admin_interface interface {
	GetAllDocuments(ctx context.Context) ([]models.Document, error)
	DeleteDocument(ctx context.Context, ID int) error
	DeleteDiscipline(ctx context.Context, DisciplineID int) error
	UploadDocument(ctx context.Context, document models.Document) error
	UploadDiscipline(ctx context.Context, discipline models.CreateDiscipline) error
	UpdateDiscipline(ctx context.Context, discipline models.Discipline) (models.Discipline, error)
	UpdateDocument(ctx context.Context, document models.Document) (models.Document, error)
	GetDocument(ctx context.Context, ID int) models.Document
	GET_Download_Count(ctx context.Context) int
	FilterDocumentsByDiscipline(ctx context.Context, DisciplineID int) ([]models.Document, error)
}

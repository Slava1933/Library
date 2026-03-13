package models

type Discipline struct {
	ID    int    `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
}

type Document struct {
	ID           int    `db:"id" json:"id"`
	DisciplineID int    `db:"discipline_id" json:"discipline_id"`
	Title        string `db:"title" json:"title"`
	Filepath     string `db:"file_path" json:"file_path"`
}

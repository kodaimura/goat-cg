package repository

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model/entity"
)


type ProjectRepository interface {
	Insert(p *entity.Project) error
    Update(id int, p *entity.Project) error
    
}


type projectRepository struct {
	db *sql.DB
}


func NewProjectRepository() ProjectRepository {
	db := db.GetDB()
	return &projectRepository{db}
}


func (rep *projectRepository) Insert(p *entity.Project) error {
	_, err := rep.db.Exec(
		`INSERT INTO projects (
			project_cd, 
			project_name
		 ) VALUES(?,?)`,
		p.ProjectCd, 
		p.ProjectName,
	)
	return err
}


func (rep *projectRepository) Update(id int, p *entity.Project) error {
	_, err := rep.db.Exec(
		`UPDATE projects 
		 SET project_name = ? 
		 WHERE project_id = ?`,
		p.ProjectName, 
		id,
	)
	return err
}
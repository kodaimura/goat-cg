package repository

import (
	"database/sql"

	"goat-cg/internal/shared/constant"
	"goat-cg/internal/core/db"
	"goat-cg/internal/model/entity"
)


type ProjectRepository interface {
	Insert(p *entity.Project) error
    Update(id int, p *entity.Project) error

    SelectByCd(cd string) (entity.Project, error)
    SelectByUserId(userId int) ([]entity.Project, error)
    SelectByCdAndUserId(cd string, userId int) (entity.Project, error)
    
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


func (rep *projectRepository) SelectByCd(cd string) (entity.Project, error) {
	var ret entity.Project
	err := rep.db.QueryRow(
		`SELECT 
			project_id,
			project_cd,
			project_name,
			create_at
		 FROM 
		 	projects
		 WHERE 
		 	project_cd = ?`, 
		 cd,
	).Scan(
		&ret.ProjectId, 
		&ret.ProjectCd, 
		&ret.ProjectName,
		&ret.CreateAt,
	)

	return ret, err
}


func (rep *projectRepository) SelectByUserId(userId int) ([]entity.Project, error){
	var ret []entity.Project
	rows, err := rep.db.Query(
		`SELECT 
			p.project_id,
			p.project_cd,
			p.project_name,
			p.create_at
		 FROM 
		 	projects p,
		 	users_projects up
		 WHERE 
		 	p.project_id = up.project_id
		 AND up.user_id = ?
		 AND up.state_cls = ?`, 
		 userId,
		 constant.STATE_CLS_JOIN,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := entity.Project{}
		err = rows.Scan(
			&p.ProjectId, 
			&p.ProjectCd, 
			&p.ProjectName,
			&p.CreateAt,
		)
		if err != nil {
			break
		}
		ret = append(ret, p)
	}

	return ret, err
}


func (rep *projectRepository) SelectByCdAndUserId(
	cd string,
	userId int,
) (entity.Project, error) {

	var ret entity.Project
	err := rep.db.QueryRow(
		`SELECT 
			p.project_id,
			p.project_name
		 FROM 
		 	projects p,
		 	users_projects up
		 WHERE 
		 	p.project_id = up.project_id
		 AND p.project_cd = ?
		 AND up.user_id = ?
		 AND up.state_cls = ?`, 
		 cd,
		 userId,
		 constant.STATE_CLS_JOIN,
	).Scan(
		&ret.ProjectId, 
		&ret.ProjectName,
	)

	return ret, err
}
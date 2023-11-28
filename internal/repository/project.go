package repository

import (
	"database/sql"

	"goat-cg/internal/shared/constant"
	"goat-cg/internal/core/db"
	"goat-cg/internal/model"
)


type ProjectRepository interface {
	Insert(p *model.Project) error
	Update(id int, p *model.Project) error

	GetByCd(cd string) (model.Project, error)
	GetByUserIdAndStateCls(
		userId int, state string,
	) ([]model.Project, error)
	GetByCdAndUserId(cd string, userId int) (model.Project, error)
	
}


type projectRepository struct {
	db *sql.DB
}


func NewProjectRepository() ProjectRepository {
	db := db.GetDB()
	return &projectRepository{db}
}


func (rep *projectRepository) Insert(p *model.Project) error {
	_, err := rep.db.Exec(
		`INSERT INTO project (
			project_cd, 
			project_name
		 ) VALUES(?,?)`,
		p.ProjectCd, 
		p.ProjectName,
	)
	return err
}


func (rep *projectRepository) Update(id int, p *model.Project) error {
	_, err := rep.db.Exec(
		`UPDATE project 
		 SET project_name = ? 
		 WHERE project_id = ?`,
		p.ProjectName, 
		id,
	)
	return err
}


func (rep *projectRepository) GetByCd(cd string) (model.Project, error) {
	var ret model.Project
	err := rep.db.QueryRow(
		`SELECT 
			project_id,
			project_cd,
			project_name,
			created_at 
		 FROM 
			 project
		 WHERE 
			 project_cd = ?`, 
		 cd,
	).Scan(
		&ret.ProjectId, 
		&ret.ProjectCd, 
		&ret.ProjectName,
		&ret.CreatedAt,
	)

	return ret, err
}


func (rep *projectRepository) GetByUserIdAndStateCls(
	userId int, state string,
) ([]model.Project, error){
	var ret []model.Project
	rows, err := rep.db.Query(
		`SELECT 
			p.project_id,
			p.project_cd,
			p.project_name,
			p.created_at 
		 FROM 
			 project p,
			 project_user pu
		 WHERE 
			 p.project_id = pu.project_id
		 AND pu.user_id = ?
		 AND pu.state_cls = ?`, 
		 userId,
		 state,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := model.Project{}
		err = rows.Scan(
			&p.ProjectId, 
			&p.ProjectCd, 
			&p.ProjectName,
			&p.CreatedAt,
		)
		if err != nil {
			break
		}
		ret = append(ret, p)
	}

	return ret, err
}


func (rep *projectRepository) GetByCdAndUserId(
	cd string,
	userId int,
) (model.Project, error) {

	var ret model.Project
	err := rep.db.QueryRow(
		`SELECT 
			p.project_id,
			p.project_name
		 FROM 
			 project p,
			 project_user pu
		 WHERE 
			 p.project_id = pu.project_id
		 AND p.project_cd = ?
		 AND pu.user_id = ?
		 AND pu.state_cls = ?`, 
		 cd,
		 userId,
		 constant.STATE_CLS_JOIN,
	).Scan(
		&ret.ProjectId, 
		&ret.ProjectName,
	)

	return ret, err
}
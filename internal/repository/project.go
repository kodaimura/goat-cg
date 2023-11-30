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

	GetByUserId(userId string) ([]model.Project, error)
	GetMemberProjects(userId int) ([]model.Project, error)
	GetByUniqueKey(username, projectName string) (model.Project, error)
	GetMemberProject(username, projectName string) (model.Project, error)
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
			project_name,
			project_memo,
			user_id,
			username 
		 ) VALUES(?,?,?,?)`,
		p.ProjectName, 
		p.ProjectMemo,
		p.UserId,
		p.Username,
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


func (rep *projectRepository) GetByUserId(userId int) ([]model.Project, error){
	var ret []model.Project
	rows, err := rep.db.Query(
		`SELECT 
			project_id,
			project_name,
			project_memo,
			created_at
			updated_at 
		 FROM 
			 project
		 WHERE user_id = ?`, 
		 userId,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := model.Project{}
		err = rows.Scan(
			&p.ProjectId, 
			&p.ProjectName,
			&p.ProjectMemo, 
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			break
		}
		ret = append(ret, p)
	}

	return ret, err
}


func (rep *projectRepository) GetMemberProjects(userId int) ([]model.Project, error){
	var ret []model.Project
	rows, err := rep.db.Query(
		`SELECT 
			p.project_id,
			p.project_name,
			p.project_memo,
			p.created_at
			p.updated_at 
		 FROM 
			 project p,
			 project_member pm
		 WHERE 
			 p.project_id = pm.project_id
		 AND pm.user_id = ?`, 
		 userId,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := model.Project{}
		err = rows.Scan(
			&p.ProjectId, 
			&p.ProjectName,
			&p.ProjectMemo, 
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			break
		}
		ret = append(ret, p)
	}

	return ret, err
}


func (rep *projectRepository) GetByUniqueKey(username, projectName string) (model.Project, error) {
	var ret model.Project
	err := rep.db.QueryRow(
		`SELECT 
			project_id,
			project_name,
			project_memo,
			user_id,
			username,
			created_at,
			updated_at 
		 FROM 
			 project
		 WHERE username = ?
		   AND project_name = ?`, 
		 username,
		 projectName,
	).Scan(
		&ret.ProjectId, 
		&ret.ProjectName, 
		&ret.ProjectMemo,
		&ret.UserId,
		&ret.Username,
		&ret.CreatedAt,
		&ret.UpdatedAt,
	)

	return ret, err
}


func (rep *projectRepository) GetMemberProject(username, projectName string) (model.Project, error) {
	var ret model.Project
	err := rep.db.QueryRow(
		`SELECT 
			p.project_id,
			p.project_name,
			p.project_memo,
			p.user_id,
			p.username,
			p.created_at
			p.updated_at 
		 FROM 
			 project p,
			 project_member pm
		 WHERE 
			 p.project_id = pm.project_id
		 AND pm.username = ?
		 AND p.project_name = ?`, 
		 username,
		 projectName,
	).Scan(
		&ret.ProjectId, 
		&ret.ProjectName, 
		&ret.ProjectMemo,
		&ret.UserId,
		&ret.Username,
		&ret.CreatedAt,
		&ret.UpdatedAt,
	)

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
			 project_member pu
		 WHERE 
			 p.project_id = pu.project_id
		 AND p.project_cd = ?
		 AND pu.user_id = ?
		 AND pu.user_status = ?`, 
		 cd,
		 userId,
		 constant.STATE_CLS_JOIN,
	).Scan(
		&ret.ProjectId, 
		&ret.ProjectName,
	)

	return ret, err
}
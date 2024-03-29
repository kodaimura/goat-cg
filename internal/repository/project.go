package repository

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model"
)


type ProjectRepository interface {
	Insert(p *model.Project) (int, error)
	Update(p *model.Project) error
	Delete(p *model.Project) error
	DeleteTx(p *model.Project, tx *sql.Tx) error
	
	GetById(projectId int) (model.Project, error)
	GetByUserId(userId int) ([]model.Project, error)
	GetMemberProjects(userId int) ([]model.Project, error)
	GetByUniqueKey(username, projectName string) (model.Project, error)
	GetMemberProject(userId int, ownername, projectName string) (model.Project, error)
}


type projectRepository struct {
	db *sql.DB
}


func NewProjectRepository() ProjectRepository {
	db := db.GetDB()
	return &projectRepository{db}
}


func (rep *projectRepository) Insert(p *model.Project) (int, error) {
	var projectId int

	err := rep.db.QueryRow(
		`INSERT INTO project (
			project_name,
			project_memo,
			user_id,
			username 
		 ) VALUES(?,?,?,?)
		 RETURNING project_id`,
		p.ProjectName, 
		p.ProjectMemo,
		p.UserId,
		p.Username,
	).Scan(
		&projectId,
	)

	return projectId, err
}


func (rep *projectRepository) Update(p *model.Project) error {
	_, err := rep.db.Exec(
		`UPDATE project 
		 SET project_name = ?,
		 	 project_memo = ? 
		 WHERE project_id = ?`,
		p.ProjectName, 
		p.ProjectMemo,
		p.ProjectId, 
	)

	return err
}


func (rep *projectRepository) Delete(p *model.Project) error {
	_, err := rep.db.Exec(
		`DELETE FROM project WHERE project_id = ?`, 
		p.ProjectId,
	)

	return err
}


func (rep *projectRepository) DeleteTx(p *model.Project, tx *sql.Tx) error {
	_, err := tx.Exec(
		`DELETE FROM project WHERE project_id = ?`, 
		p.ProjectId,
	)

	return err
}


func (rep *projectRepository) GetById(projectId int) (model.Project, error) {
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
		 WHERE project_id = ?`, 
		 projectId,
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


func (rep *projectRepository) GetByUserId(userId int) ([]model.Project, error){
	rows, err := rep.db.Query(
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
		 WHERE user_id = ?`, 
		 userId,
	)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	ret := []model.Project{}
	for rows.Next() {
		p := model.Project{}
		err = rows.Scan(
			&p.ProjectId, 
			&p.ProjectName,
			&p.ProjectMemo, 
			&p.UserId,
			&p.Username,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}

	return ret, nil
}


func (rep *projectRepository) GetMemberProjects(userId int) ([]model.Project, error){
	rows, err := rep.db.Query(
		`SELECT 
			p.project_id,
			p.project_name,
			p.project_memo,
			p.user_id,
			p.username,
			p.created_at,
			p.updated_at 
		 FROM 
			 project p,
			 project_member pm
		 WHERE 
			 p.project_id = pm.project_id
		 AND pm.user_id = ?`, 
		 userId,
	)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	ret := []model.Project{}
	for rows.Next() {
		p := model.Project{}
		err = rows.Scan(
			&p.ProjectId, 
			&p.ProjectName,
			&p.ProjectMemo, 
			&p.UserId,
			&p.Username,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}

	return ret, nil
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


func (rep *projectRepository) GetMemberProject(userId int, ownername, projectName string) (model.Project, error) {
	var ret model.Project
	err := rep.db.QueryRow(
		`SELECT 
			p.project_id,
			p.project_name,
			p.project_memo,
			p.user_id,
			p.username,
			p.created_at,
			p.updated_at 
		 FROM 
			 project p,
			 project_member pm
		 WHERE 
			 p.project_id = pm.project_id
		 AND pm.user_id = ?
		 AND p.username = ? 
		 AND p.project_name = ?`, 
		 userId,
		 ownername,
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
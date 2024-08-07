package repository

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model"
)


type ProjectRepository interface {
	Get(p *model.Project) ([]model.Project, error)
	GetOne(p *model.Project) (model.Project, error)
	Insert(p *model.Project, tx *sql.Tx) error
	Update(p *model.Project, tx *sql.Tx) error
	Delete(p *model.Project, tx *sql.Tx) error
	
	GetMemberProjects(userId int) ([]model.Project, error)
	GetMemberProject(userId int, ownername, projectName string) (model.Project, error)
}


type projectRepository struct {
	db *sql.DB
}


func NewProjectRepository() ProjectRepository {
	db := db.GetDB()
	return &projectRepository{db}
}


func (rep *projectRepository) Get(p *model.Project) ([]model.Project, error) {
	where, binds := db.BuildWhereClause(p)
	query := "SELECT * FROM project " + where
	rows, err := rep.db.Query(query, binds...)
	defer rows.Close()

	if err != nil {
		return []model.Project{}, err
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
			return []model.Project{}, err
		}
		ret = append(ret, p)
	}

	return ret, nil
}


func (rep *projectRepository) GetOne(p *model.Project) (model.Project, error) {
	var ret model.Project
	where, binds := db.BuildWhereClause(p)
	query := "SELECT * FROM project " + where

	err := rep.db.QueryRow(query, binds...).Scan(
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


func (rep *projectRepository) Insert(p *model.Project, tx *sql.Tx) error {
	cmd := 
	`INSERT INTO project (
		project_name,
		project_memo,
		user_id,
		username 
	 ) VALUES(?,?,?,?)`
	binds := []interface{}{p.ProjectName, p.ProjectMemo, p.UserId, p.Username}

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}


func (rep *projectRepository) Update(p *model.Project, tx *sql.Tx) error {
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


func (rep *projectRepository) Delete(p *model.Project, tx *sql.Tx) error {
	where, binds := db.BuildWhereClause(p)
	cmd := "DELETE FROM project " + where

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
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
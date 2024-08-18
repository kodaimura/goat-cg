package query

import (
	"database/sql"

	"goat-cg/internal/model"
	"goat-cg/internal/core/db"
)


type ProjectQuery interface {
	GetMemberProjects(userId int) ([]model.Project, error)
	GetMemberProject(userId int, ownername, projectName string) (model.Project, error)
}


type projectQuery struct {
	db *sql.DB
}


func NewProjectQuery() ProjectQuery {
	db := db.GetDB()
	return &projectQuery{db}
}


func (que *projectQuery) GetMemberProjects(userId int) ([]model.Project, error) {
	rows, err := que.db.Query(
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


func (que *projectQuery) GetMemberProject(userId int, ownername, projectName string) (model.Project, error) {
	var ret model.Project
	err := que.db.QueryRow(
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
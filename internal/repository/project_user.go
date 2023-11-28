package repository

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model"
)


type ProjectUserRepository interface {
	Select(userId, projectId int) (model.ProjectUser, error)
	Upsert(up *model.ProjectUser) error
	Delete(userId, projectId int) error
}


type projectUserRepository struct {
	db *sql.DB
}


func NewProjectUserRepository() ProjectUserRepository {
	db := db.GetDB()
	return &projectUserRepository{db}
}


func (rep *projectUserRepository) Select(userId, projectId int) (model.ProjectUser, error) {
	var ret model.ProjectUser

	err := rep.db.QueryRow(
		`SELECT
			user_id, 
			project_id, 
			state_cls, 
			role_cls
		 FROM project_user
		 WHERE user_id = ?
		  AND project_id = ?`,
		 userId,
		 projectId,
	).Scan(
		&ret.UserId, 
		&ret.ProjectId, 
		&ret.StateCls, 
		&ret.RoleCls,
	)

	return ret, err
}


func (rep *projectUserRepository) Upsert(up *model.ProjectUser) error {
	_, err := rep.db.Exec(
		`REPLACE INTO project_user (
			user_id, 
			project_id, 
			state_cls, 
			role_cls
		 ) 
		 VALUES(?,?,?,?)`,
		up.UserId, 
		up.ProjectId,
		up.StateCls,
		up.RoleCls,
	)

	return err
}


func (rep *projectUserRepository) Delete(userId, projectId int) error {
	_, err := rep.db.Exec(
		`DELETE FROM project_user
		 WHERE 
			 user_id = ?
		 AND project_id = ?`, 
		userId,
		projectId,
	)
	return err
}

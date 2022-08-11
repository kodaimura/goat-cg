package repository

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model/entity"
)


type UserProjectRepository interface {
	Select(userId, projectId int) (entity.UserProject, error)
    Upsert(up *entity.UserProject) error
    Delete(userId, projectId int) error
}


type userProjectRepository struct {
	db *sql.DB
}


func NewUserProjectRepository() UserProjectRepository {
	db := db.GetDB()
	return &userProjectRepository{db}
}


func (rep *userProjectRepository) Select(userId, projectId int) (entity.UserProject, error) {
	var ret entity.UserProject

	err := rep.db.QueryRow(
		`SELECT
			user_id, 
			project_id, 
			state_cls, 
			role_cls
		 FROM users_projects
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


func (rep *userProjectRepository) Upsert(up *entity.UserProject) error {
	_, err := rep.db.Exec(
		`REPLACE INTO users_projects (
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


func (rep *userProjectRepository) Delete(userId, projectId int) error {
	_, err := rep.db.Exec(
		`DELETE FROM users_projects
		 WHERE 
		 	user_id = ?
		 AND project_id = ?`, 
		userId,
		projectId,
	)
	return err
}
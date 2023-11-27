package dao

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model/entity"
)


type ProjectUserDao interface {
	Select(userId, projectId int) (entity.ProjectUser, error)
	Upsert(up *entity.ProjectUser) error
	Delete(userId, projectId int) error
}


type projectUserDao struct {
	db *sql.DB
}


func NewProjectUserDao() ProjectUserDao {
	db := db.GetDB()
	return &projectUserDao{db}
}


func (rep *projectUserDao) Select(userId, projectId int) (entity.ProjectUser, error) {
	var ret entity.ProjectUser

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


func (rep *projectUserDao) Upsert(up *entity.ProjectUser) error {
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


func (rep *projectUserDao) Delete(userId, projectId int) error {
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

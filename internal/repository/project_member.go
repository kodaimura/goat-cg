package repository

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model"
)


type ProjectMemberRepository interface {
	GetByPk(userId, projectId int) (model.ProjectMember, error)
	Upsert(up *model.ProjectMember) error
	Delete(userId, projectId int) error
}


type projectMemberRepository struct {
	db *sql.DB
}


func NewProjectMemberRepository() ProjectMemberRepository {
	db := db.GetDB()
	return &projectMemberRepository{db}
}


func (rep *projectMemberRepository) GetByPk(userId, projectId int) (model.ProjectMember, error) {
	var ret model.ProjectMember

	err := rep.db.QueryRow(
		`SELECT
			user_id, 
			project_id, 
			user_status, 
			user_role
		 FROM project_member
		 WHERE user_id = ?
		  AND project_id = ?`,
		 userId,
		 projectId,
	).Scan(
		&ret.UserId, 
		&ret.ProjectId, 
		&ret.UserStatus, 
		&ret.UserRole,
	)

	return ret, err
}


func (rep *projectMemberRepository) Upsert(up *model.ProjectMember) error {
	_, err := rep.db.Exec(
		`REPLACE INTO project_member (
			user_id, 
			project_id, 
			user_status, 
			user_role
		 ) 
		 VALUES(?,?,?,?)`,
		up.UserId, 
		up.ProjectId,
		up.UserStatus,
		up.UserRole,
	)

	return err
}


func (rep *projectMemberRepository) Delete(userId, projectId int) error {
	_, err := rep.db.Exec(
		`DELETE FROM project_member
		 WHERE 
			 user_id = ?
		 AND project_id = ?`, 
		userId,
		projectId,
	)
	return err
}

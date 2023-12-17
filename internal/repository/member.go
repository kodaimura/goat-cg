package repository

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model"
)


type MemberRepository interface {
	GetByPk(projectId, userId int) (model.Member, error)
	Insert(m *model.Member) error
	Upsert(m *model.Member) error
	Delete(m *model.Member) error
}


type memberRepository struct {
	db *sql.DB
}


func NewMemberRepository() MemberRepository {
	db := db.GetDB()
	return &memberRepository{db}
}


func (rep *memberRepository) GetByPk(projectId, userId int) (model.Member, error) {
	var ret model.Member

	err := rep.db.QueryRow(
		`SELECT
			project_id, 
			user_id, 
			user_status, 
			user_role
		 FROM project_member
		 WHERE user_id = ?
		  AND project_id = ?`,
		 userId,
		 projectId,
	).Scan(
		&ret.ProjectId, 
		&ret.UserId, 
		&ret.UserStatus, 
		&ret.UserRole,
	)

	return ret, err
}


func (rep *memberRepository) Insert(m *model.Member) error {
	_, err := rep.db.Exec(
		`INSERT INTO project_member (
			project_id,
			user_id, 
			user_status, 
			user_role
		 ) 
		 VALUES(?,?,?,?)`,
		m.ProjectId,
		m.UserId, 
		m.UserStatus,
		m.UserRole,
	)

	return err
}


func (rep *memberRepository) Upsert(m *model.Member) error {
	_, err := rep.db.Exec(
		`REPLACE INTO project_member (
			user_id, 
			project_id, 
			user_status, 
			user_role
		 ) 
		 VALUES(?,?,?,?)`,
		m.UserId, 
		m.ProjectId,
		m.UserStatus,
		m.UserRole,
	)

	return err
}


func (rep *memberRepository) Delete(m *model.Member) error {
	_, err := rep.db.Exec(
		`DELETE FROM project_member
		 WHERE user_id = ?
		   AND project_id = ?`, 
		m.UserId,
		m.ProjectId,
	)
	return err
}
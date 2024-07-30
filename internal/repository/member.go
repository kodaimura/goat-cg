package repository

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model"
)


type MemberRepository interface {
	Get(m *model.Member) ([]model.Member, error)
	GetOne(m *model.Member) (model.Member, error)
	Insert(m *model.Member, tx *sql.Tx) error
	Update(m *model.Member, tx *sql.Tx) error
	Delete(m *model.Member, tx *sql.Tx) error

	Upsert(m *model.Member, tx *sql.Tx) error
}


type memberRepository struct {
	db *sql.DB
}


func NewMemberRepository() MemberRepository {
	db := db.GetDB()
	return &memberRepository{db}
}


func (rep *memberRepository) Get(m *model.Member) ([]model.Member, error) {
	where, binds := db.BuildWhereClause(m)
	query := "SELECT * FROM project_member " + where

	rows, err := rep.db.Query(query, binds...)
	defer rows.Close()

	if err != nil {
		return []model.Member{}, err
	}

	ret := []model.Member{}
	for rows.Next() {
		m := model.Member{}
		err = rows.Scan(
			&m.ProjectId, 
			&m.UserId, 
			&m.UserStatus, 
			&m.UserRole,
			&m.CreatedAt, 
			&m.UpdatedAt,
		)
		if err != nil {
			return []model.Member{}, err
		}
		ret = append(ret, m)
	}

	return ret, err
}


func (rep *memberRepository) GetOne(m *model.Member) (model.Member, error) {
	var ret model.Member
	where, binds := db.BuildWhereClause(m)
	query := "SELECT * FROM project_member " + where

	err := rep.db.QueryRow(query, binds...).Scan(
		&ret.ProjectId, 
		&ret.UserId, 
		&ret.UserStatus, 
		&ret.UserRole,
		&ret.CreatedAt, 
		&ret.UpdatedAt,
	)

	return ret, err
}


func (rep *memberRepository) Insert(m *model.Member, tx *sql.Tx) error {
	cmd := 
	`INSERT INTO project_member (
		project_id,
		user_id, 
		user_status, 
		user_role
	 ) 
	 VALUES(?,?,?,?)`
	binds := []interface{}{m.ProjectId, m.UserId, m.UserStatus, m.UserRole}

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }

	return err
}


func (rep *memberRepository) Update(m *model.Member, tx *sql.Tx) error {
	cmd := 
	`UPDATE users 
	 SET user_status = ?,
	 	 user_role = ?
	 WHERE user_id = ?
	   AND project_id = ?`
	binds := []interface{}{m.UserStatus, m.UserRole, m.UserId, m.ProjectId}
	
	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}


func (rep *memberRepository) Delete(m *model.Member, tx *sql.Tx) error {
	cmd := 
	`DELETE FROM project_member 
	 WHERE user_id = ?
	   AND project_id = ?`
	binds := []interface{}{m.UserId, m.ProjectId}

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }

	return err
}


func (rep *memberRepository) Upsert(m *model.Member, tx *sql.Tx) error {
	cmd := 
	`REPLACE INTO project_member (
		project_id, 
		user_id, 
		user_status, 
		user_role
	 ) 
	 VALUES(?,?,?,?)`
	binds := []interface{}{m.ProjectId, m.UserId, m.UserStatus, m.UserRole}

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }

	return err
}
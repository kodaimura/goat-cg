package repository

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model"
)


type UserRepository interface {
	Get(u *model.User) ([]model.User, error)
	GetOne(u *model.User) (model.User, error)
	Insert(u *model.User, tx *sql.Tx) error
	Update(u *model.User, tx *sql.Tx) error
	Delete(u *model.User, tx *sql.Tx) error

	UpdateName(u *model.User, tx *sql.Tx) error
	UpdatePassword(u *model.User, tx *sql.Tx) error
	UpdateEmail(u *model.User, tx *sql.Tx) error
}


type userRepository struct {
	db *sql.DB
}

func NewUserRepository() UserRepository {
	db := db.GetDB()
	return &userRepository{db}
}

func (rep *userRepository) Get(u *model.User) ([]model.User, error) {
	where, binds := db.BuildWhereClause(u)
	query := "SELECT * FROM users " + where
	rows, err := rep.db.Query(query, binds...)
	defer rows.Close()

	if err != nil {
		return []model.User{}, err
	}

	ret := []model.User{}
	for rows.Next() {
		u := model.User{}
		err = rows.Scan(
			&u.UserId, 
			&u.Username, 
			&u.Password,
			&u.Email,
			&u.CreatedAt, 
			&u.UpdatedAt,
		)
		if err != nil {
			return []model.User{}, err
		}
		ret = append(ret, u)
	}

	return ret, nil
}


func (rep *userRepository) GetOne(u *model.User) (model.User, error) {
	var ret model.User
	where, binds := db.BuildWhereClause(u)
	query := "SELECT * FROM users " + where

	err := rep.db.QueryRow(query, binds...).Scan(
		&ret.UserId, 
		&ret.Username,  
		&ret.Password,
		&ret.Email,
		&ret.CreatedAt, 
		&ret.UpdatedAt,
	)

	return ret, err
}


func (rep *userRepository) Insert(u *model.User, tx *sql.Tx) error {
	cmd := 
	`INSERT INTO users (
		username, 
		password,
		email
	 ) VALUES(?,?,?)`
	binds := []interface{}{u.Username, u.Password, u.Email}

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}


func (rep *userRepository) Update(u *model.User, tx *sql.Tx) error {
	cmd := 
	`UPDATE users 
	 SET username = ?,
	     password = ?,
		 email = ?
	 WHERE user_id = ?`
	binds := []interface{}{u.Username, u.Password, u.Email, u.UserId}
	
	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}


func (rep *userRepository) Delete(u *model.User, tx *sql.Tx) error {
	cmd := "DELETE FROM users WHERE user_id = ?"
	binds := []interface{}{u.UserId}

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}


func (rep *userRepository) UpdateName(u *model.User, tx *sql.Tx) error {
	cmd := 
	`UPDATE users
	 SET username = ? 
	 WHERE user_id = ?`
	binds := []interface{}{u.Username, u.UserId}

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}


func (rep *userRepository) UpdatePassword(u *model.User, tx *sql.Tx) error {
	cmd := 
	`UPDATE users 
	 SET password = ? 
	 WHERE user_id = ?`
	binds := []interface{}{u.Password, u.UserId}

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }

	return err
}


func (rep *userRepository) UpdateEmail(u *model.User, tx *sql.Tx) error {
	cmd := 
	`UPDATE users 
	 SET email = ? 
	 WHERE user_id = ?`
	binds := []interface{}{u.Email, u.UserId}

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }

	return err
}
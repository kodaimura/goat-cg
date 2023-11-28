package repository

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model"
)


type UserRepository interface {
	Insert(u *model.User) error
	GetById(id int) (model.User, error)
	Update(id int, u *model.User) error
	Delete(id int) error
	
	/* 以降に追加 */
	GetByName(name string) (model.User, error)
	UpdatePassword(id int, password string) error
	UpdateName(id int, name string) error
}


type userRepository struct {
	db *sql.DB
}


func NewUserRepository() UserRepository {
	db := db.GetDB()
	return &userRepository{db}
}


func (rep *userRepository) GetById(id int) (model.User, error){
	var ret model.User

	err := rep.db.QueryRow(
		`SELECT 
			user_id, 
			user_name, 
			created_at , 
			updated_at 
		 FROM users 
		 WHERE user_id = ?`, 
		 id,
	).Scan(
		&ret.UserId, 
		&ret.UserName, 
		&ret.CreatedAt, 
		&ret.UpdatedAt,
	)

	return ret, err
}


func (rep *userRepository) Insert(u *model.User) error {
	_, err := rep.db.Exec(
		`INSERT INTO users (
			user_name, 
			password
		 ) VALUES(?,?)`,
		u.UserName, 
		u.Password,
	)
	return err
}


func (rep *userRepository) Update(id int, u *model.User) error {
	_, err := rep.db.Exec(
		`UPDATE users 
		 SET user_name = ? 
			  password = ?
		 WHERE user_id = ?`,
		u.UserName,
		u.Password, 
		id,
	)
	return err
}


func (rep *userRepository) Delete(id int) error {
	_, err := rep.db.Exec(
		`DELETE FROM users WHERE user_id = ?`, 
		id,
	)

	return err
}


func (rep *userRepository) UpdatePassword(id int, password string) error {
	_, err := rep.db.Exec(
		`UPDATE users 
		 SET password = ? 
		 WHERE user_id = ?`, 
		 password, 
		 id,
	)
	return err
}


func (rep *userRepository) UpdateName(id int, name string) error {
	_, err := rep.db.Exec(
		`UPDATE users
		 SET user_name = ? 
		 WHERE user_id = ?`, 
		name, 
		id,
	)
	return err
}



func (rep *userRepository) GetByName(name string) (model.User, error) {
	var ret model.User

	err := rep.db.QueryRow(
		`SELECT 
			user_id, 
			user_name, 
			password, 
			created_at , 
			updated_at 
		 FROM users 
		 WHERE user_name = ?`, 
		 name,
	).Scan(
		&ret.UserId, 
		&ret.UserName, 
		&ret.Password, 
		&ret.CreatedAt, 
		&ret.UpdatedAt,
	)

	return ret, err
}
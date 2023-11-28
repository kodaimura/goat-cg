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
			username, 
			created_at , 
			updated_at 
		 FROM users 
		 WHERE user_id = ?`, 
		 id,
	).Scan(
		&ret.UserId, 
		&ret.Username, 
		&ret.CreatedAt, 
		&ret.UpdatedAt,
	)

	return ret, err
}


func (rep *userRepository) Insert(u *model.User) error {
	_, err := rep.db.Exec(
		`INSERT INTO users (
			username, 
			password
		 ) VALUES(?,?)`,
		u.Username, 
		u.Password,
	)
	return err
}


func (rep *userRepository) Update(id int, u *model.User) error {
	_, err := rep.db.Exec(
		`UPDATE users 
		 SET username = ? 
			  password = ?
		 WHERE user_id = ?`,
		u.Username,
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
		 SET username = ? 
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
			username, 
			password, 
			created_at , 
			updated_at 
		 FROM users 
		 WHERE username = ?`, 
		 name,
	).Scan(
		&ret.UserId, 
		&ret.Username, 
		&ret.Password, 
		&ret.CreatedAt, 
		&ret.UpdatedAt,
	)

	return ret, err
}
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
	
	GetByName(name string) (model.User, error)
	GetByEmail(email string) (model.User, error)
	UpdatePassword(id int, password string) error
	UpdateEmail(id int, email string) error
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
			email,
			created_at , 
			updated_at 
		 FROM users 
		 WHERE user_id = ?`, 
		 id,
	).Scan(
		&ret.UserId, 
		&ret.Username, 
		&ret.Email, 
		&ret.CreatedAt, 
		&ret.UpdatedAt,
	)

	return ret, err
}


func (rep *userRepository) Insert(u *model.User) error {
	_, err := rep.db.Exec(
		`INSERT INTO users (
			username, 
			password,
			email
		 ) VALUES(?,?, ?)`,
		u.Username, 
		u.Password,
		u.Email
	)
	return err
}


func (rep *userRepository) Update(id int, u *model.User) error {
	_, err := rep.db.Exec(
		`UPDATE users 
		 SET username = ? 
			 password = ?
			 emial = ?
		 WHERE user_id = ?`,
		u.Username,
		u.Password,
		u.Email, 
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


func (rep *userRepository) UpdateEmail(id int, email string) error {
	_, err := rep.db.Exec(
		`UPDATE users
		 SET email = ? 
		 WHERE user_id = ?`, 
		email, 
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
			email,
			created_at , 
			updated_at 
		 FROM users 
		 WHERE username = ?`, 
		 name,
	).Scan(
		&ret.UserId, 
		&ret.Username, 
		&ret.Password, 
		&ret.Email, 
		&ret.CreatedAt, 
		&ret.UpdatedAt,
	)

	return ret, err
}


func (rep *userRepository) GetByEmail(email string) (model.User, error) {
	var ret model.User

	err := rep.db.QueryRow(
		`SELECT 
			user_id, 
			username, 
			password, 
			email,
			created_at , 
			updated_at 
		 FROM users 
		 WHERE email = ?`, 
		 email,
	).Scan(
		&ret.UserId, 
		&ret.Username, 
		&ret.Password, 
		&ret.Email, 
		&ret.CreatedAt, 
		&ret.UpdatedAt,
	)

	return ret, err
}
package repository

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model"
)


type UserRepository interface {
	Insert(u *model.User) (int, error)
	Get() ([]model.User, error)
	GetById(id int) (model.User, error)
	GetByName(name string) (model.User, error)
	GetByEmail(email string) (model.User, error)
	Update(u *model.User) error
	UpdateEmail(u *model.User) error
	UpdatePassword(u *model.User) error
	Delete(u *model.User) error
}


type userRepository struct {
	db *sql.DB
}


func NewUserRepository() UserRepository {
	db := db.GetDB()
	return &userRepository{db}
}


func (rep *userRepository) Insert(u *model.User) (int, error) {
	var userId int

	err := rep.db.QueryRow(
		`INSERT INTO users (
			username, 
			password,
			email
		 ) VALUES(?,?, ?)`,
		u.Username, 
		u.Password,
		u.Email,
	).Scan(
		&userId,
	)

	return userId, err
}


func (ur *userRepository) Get() ([]model.User, error) {
	rows, err := ur.db.Query(
		`SELECT 
			id, 
			username, 
			email,
			created_at, 
			updated_at 
		 FROM users`,
	)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	ret := []model.User{}
	for rows.Next() {
		u := model.User{}
		err = rows.Scan(
			&u.UserId, 
			&u.Username,
			&u.Email,
			&u.CreatedAt, 
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		ret = append(ret, u)
	}

	return ret, nil
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


func (rep *userRepository) Update(u *model.User) error {
	_, err := rep.db.Exec(
		`UPDATE users 
		 SET username = ? 
			 password = ?
			 email = ?
		 WHERE user_id = ?`,
		u.Username,
		u.Password,
		u.Email, 
		u.UserId,
	)
	return err
}


func (ur *userRepository) UpdatePassword(u *model.User) error {
	_, err := ur.db.Exec(
		`UPDATE users 
		 SET password = ? 
		 WHERE user_id = ?`, 
		 u.Password, 
		 u.UserId,
	)
	return err
}


func (ur *userRepository) UpdateEmail(u *model.User) error {
	_, err := ur.db.Exec(
		`UPDATE users
		 SET email = ? 
		 WHERE user_id = ?`, 
		u.Email, 
		u.UserId,
	)
	return err
}


func (rep *userRepository) Delete(u *model.User) error {
	_, err := rep.db.Exec(
		`DELETE FROM users WHERE user_id = ?`, 
		u.UserId,
	)

	return err
}



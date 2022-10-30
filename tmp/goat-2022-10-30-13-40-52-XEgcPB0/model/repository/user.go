package repository


import (
	"database/sql"

	"xxxxx/internal/core/db"
	"xxxxx/internal/model/entity"
)


type UserRepository interface {
	Insert(u *entity.User) error
	Select(id int) (entity.User, error)
	Update(id int, u *entity.User) error
	Delete(id int) error
	SelectAll() ([]entity.User, error)
}


type userRepository struct {
	db *sql.DB
}


func NewUserRepository() UserRepository {
	db := db.GetDB()
	return &userRepository{db}
}


func (rep *userRepository) Insert(u *entity.User) error {
	_, err := rep.db.Exec(
		`INSERT INTO user (
			test_column
		 ) VALUES(?)`,
		u.TestColumn,
	)

	return err
}


func (rep *userRepository) Select(id int) (entity.User, error) {
	var ret entity.User

	err := rep.db.QueryRow(
		`SELECT
			user_id,
			test_column,
			create_at,
			update_at
		 FROM user
		 WHERE user_id = ?`,
		id,
	).Scan(
		&ret.UserId,
		&ret.TestColumn,
		&ret.CreateAt,
		&ret.UpdateAt,
	)

	return ret, err
}


func (rep *userRepository) Update(id int, u *entity.User) error {
	_, err := rep.db.Exec(
		`UPDATE user
		 SET
			test_column = ?
		 WHERE user_id = ?`,
		u.TestColumn,
		id,
	)

	return err
}


func (rep *userRepository) Delete(id int) error {
	_, err := rep.db.Exec(
		`DELETE FROM user
		 WHERE user_id = ?`,
		id,
	)

	return err
}


func (rep *userRepository) SelectAll() ([]entity.User, error) {
	var ret []entity.User

	rows, err := rep.db.Query(
		`SELECT
			user_id,
			test_column,
			create_at,
			update_at
		 FROM user`,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		u := entity.User{}
		err = rows.Scan(
			&u.UserId,
			&u.TestColumn,
			&u.CreateAt,
			&u.UpdateAt,
		)
		if err != nil {
			break
		}
		ret = append(ret, u)
	}

	return ret, err
}
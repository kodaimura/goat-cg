package queryservice

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model/entity"
)


type UserQueryService interface {
	QueryUsers() ([]entity.User, error)
	QueryUser(id int) (entity.User, error)
	QueryUserByName(name string) (entity.User, error)
}

type userQueryService struct {
	db *sql.DB
}

func NewUserQueryService() UserQueryService {
	db := db.GetDB()
	return &userQueryService{db}
}


func (qs *userQueryService) QueryUsers() ([]entity.User, error){
	var ret []entity.User

	rows, err := qs.db.Query(
		`SELECT 
			user_id,
			user_name, 
			create_at, 
			update_at 
		 FROM users`,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		u := entity.User{}
		err = rows.Scan(
			&u.UserId, 
			&u.UserName, 
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


func (qs *userQueryService) QueryUser(id int) (entity.User, error) {
	var ret entity.User

	err := qs.db.QueryRow(
		`SELECT 
			user_id, 
			user_name, 
			create_at, 
			update_at 
		 FROM users 
		 WHERE user_id = ?`, 
		 id,
	).Scan(
		&ret.UserId, 
		&ret.UserName, 
		&ret.CreateAt, 
		&ret.UpdateAt,
	)

	return ret, err
}


func (qs *userQueryService) QueryUserByName(name string) (entity.User, error) {
	var ret entity.User

	err := qs.db.QueryRow(
		`SELECT 
			user_id, 
			user_name, 
			password,
			create_at, 
			update_at 
		 FROM users 
		 WHERE user_name = ?`, 
		 name,
	).Scan(
		&ret.UserId, 
		&ret.UserName, 
		&ret.Password, 
		&ret.CreateAt, 
		&ret.UpdateAt,
	)

	return ret, err
}

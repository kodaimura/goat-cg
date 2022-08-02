package queryservice

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model/entity"
)


type UserProjectQueryService interface {
	QueryUserProjectsByUserId(id int) ([]entity.UserProject, error)
	QueryUserProjectsByProjectId(id int) ([]entity.UserProject, error)
	QueryUserProject(userId, projectId int) (entity.UserProject, error)
	
}

type userProjectQueryService struct {
	db *sql.DB
}

func NewUserProjectQueryService() UserProjectQueryService {
	db := db.GetDB()
	return &userProjectQueryService{db}
}


func (qs *userProjectQueryService) QueryUserProjectsByUserId(
	id int,
) ([]entity.UserProject, error) {
	
	var ret []entity.UserProject

	rows, err := qs.db.Query(
		`SELECT 
			user_id, 
			project_id, 
			state_cls, 
			role_cls
		 FROM 
		 	users_projects
		 WHERE
		 	user_id = ?`,
		 id,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		u := entity.UserProject{}
		err = rows.Scan(
			&u.UserId, 
			&u.ProjectId, 
			&u.StateCls, 
			&u.RoleCls,
		)
		if err != nil {
			break
		}
		ret = append(ret, u)
	}

	return ret, err
}


func (qs *userProjectQueryService) QueryUserProjectsByProjectId(
	id int,
) ([]entity.UserProject, error) {

	var ret []entity.UserProject

	rows, err := qs.db.Query(
		`SELECT 
			user_id, 
			project_id, 
			state_cls, 
			role_cls
		 FROM 
		 	users_projects
		 WHERE
		 	project_id = ?`,
		 id,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		u := entity.UserProject{}
		err = rows.Scan(
			&u.UserId, 
			&u.ProjectId, 
			&u.StateCls, 
			&u.RoleCls,
		)
		if err != nil {
			break
		}
		ret = append(ret, u)
	}

	return ret, err
}


func (qs *userProjectQueryService) QueryUserProject(
	userId int, 
	projectId int,
) (entity.UserProject, error) {
	
	var ret entity.UserProject

	err := qs.db.QueryRow(
		`SELECT 
			user_id, 
			project_id, 
			state_cls, 
			role_cls,
			create_at,
			update_at
		 FROM 
		 	users_projects
		 WHERE 
		 	user_id = ?
		 AND project_id = ?`, 
		 userId,
		 projectId,
	).Scan(
		&ret.UserId, 
		&ret.ProjectId, 
		&ret.StateCls, 
		&ret.RoleCls,
		&ret.UpdateAt,
	)

	return ret, err
}
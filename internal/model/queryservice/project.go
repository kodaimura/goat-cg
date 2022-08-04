package queryservice

import (
	"database/sql"

	"goat-cg/internal/shared/constant"
	"goat-cg/internal/core/db"
	"goat-cg/internal/model/entity"
)


type ProjectQueryService interface {
	QueryProjectByCd(cd string) (entity.Project, error)
	QueryProjectsByUserId(userId int) ([]entity.Project, error)
	QueryProjectByCdAndUserId(cd string, userId int) (entity.Project, error)
}

type projectQueryService struct {
	db *sql.DB
}

func NewProjectQueryService() ProjectQueryService {
	db := db.GetDB()
	return &projectQueryService{db}
}


func (qs *projectQueryService) QueryProjectByCd(
	cd string,
) (entity.Project, error) {

	var ret entity.Project
	err := qs.db.QueryRow(
		`SELECT 
			project_id,
			project_cd,
			project_name
		 FROM 
		 	projects
		 WHERE 
		 	project_cd = ?`, 
		 cd,
	).Scan(
		&ret.ProjectId, 
		&ret.ProjectCd, 
		&ret.ProjectName,
	)

	return ret, err
}


func (qs *projectQueryService) QueryProjectsByUserId(
	userId int,
) ([]entity.Project, error){

	var ret []entity.Project
	rows, err := qs.db.Query(
		`SELECT 
			p.project_id,
			p.project_cd,
			p.project_name
		 FROM 
		 	projects p,
		 	users_projects up
		 WHERE 
		 	p.project_id = up.project_id
		 AND up.user_id = ?
		 AND up.state_cls = ?`, 
		 userId,
		 constant.STATE_CLS_JOIN,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := entity.Project{}
		err = rows.Scan(
			&p.ProjectId, 
			&p.ProjectCd, 
			&p.ProjectName,
		)
		if err != nil {
			break
		}
		ret = append(ret, p)
	}

	return ret, err
}


func (qs *projectQueryService) QueryProjectByCdAndUserId(
	cd string,
	userId int,
) (entity.Project, error) {

	var ret entity.Project
	err := qs.db.QueryRow(
		`SELECT 
			p.project_id,
			p.project_name
		 FROM 
		 	projects p,
		 	users_projects up
		 WHERE 
		 	p.project_id = up.project_id
		 AND p.project_cd = ?
		 AND up.user_id = ?
		 AND up.state_cls = ?`, 
		 cd,
		 userId,
		 constant.STATE_CLS_JOIN,
	).Scan(
		&ret.ProjectId, 
		&ret.ProjectName,
	)

	return ret, err
}
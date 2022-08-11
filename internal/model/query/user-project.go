package query

import (
	"database/sql"

	"goat-cg/internal/shared/dto"
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/core/db"
)


type UserProjectQuery interface {
	QueryJoinRequests(userId int) ([]dto.QueOutJoinRequest, error)
}


type userProjectQuery struct {
	db *sql.DB
}


func NewUserProjectQuery() UserProjectQuery {
	db := db.GetDB()
	return &userProjectQuery{db}
}


func (que *userProjectQuery)QueryJoinRequests(userId int) ([]dto.QueOutJoinRequest, error){

	var ret []dto.QueOutJoinRequest
	rows, err := que.db.Query(
		`SELECT
			u.user_id,
			u.user_name,
			p.project_id,
			p.project_cd,
			p.project_name,
			up2.update_at
		 FROM 
		 	users_projects up1,
		 	users_projects up2,
		 	users u,
		 	projects p
		 WHERE up1.user_id = ?
		  AND up1.role_cls in (?, ?)
		  AND up2.project_id = up1.project_id
		  AND up2.state_cls = ?
		  AND u.user_id = up2.user_id
		  AND p.project_id = up2.project_id `, 
		 userId,
		 constant.ROLE_CLS_ADMIN,
		 constant.ROLE_CLS_OWNER,
		 constant.STATE_CLS_REQUEST,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		x := dto.QueOutJoinRequest{}
		err = rows.Scan(
			&x.UserId,
			&x.UserName,
			&x.ProjectId,
			&x.ProjectCd,
			&x.ProjectName,
			&x.UpdateAt,
		)
		if err != nil {
			break
		}
		ret = append(ret, x)
	}

	return ret, err
}
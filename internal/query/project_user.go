package query

import (
	"database/sql"

	"goat-cg/internal/shared/dto"
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/core/db"
)


type ProjectUserQuery interface {
	QueryJoinRequests(userId int) ([]dto.QueOutJoinRequest, error)
}


type projectUserQuery struct {
	db *sql.DB
}


func NewProjectUserQuery() ProjectUserQuery {
	db := db.GetDB()
	return &projectUserQuery{db}
}


func (que *projectUserQuery)QueryJoinRequests(userId int) ([]dto.QueOutJoinRequest, error){

	var ret []dto.QueOutJoinRequest
	rows, err := que.db.Query(
		`SELECT
			u.user_id,
			u.username,
			p.project_id,
			p.project_cd,
			p.project_name,
			pu2.updated_at
		 FROM 
			 project_user pu1,
			 project_user pu2,
			 users u,
			 project p
		 WHERE pu1.user_id = ?
		  AND pu1.role_cls in (?, ?)
		  AND pu2.project_id = pu1.project_id
		  AND pu2.state_cls = ?
		  AND u.user_id = pu2.user_id
		  AND p.project_id = pu2.project_id `, 
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
			&x.Username,
			&x.ProjectId,
			&x.ProjectCd,
			&x.ProjectName,
			&x.UpdatedAt,
		)
		if err != nil {
			break
		}
		ret = append(ret, x)
	}

	return ret, err
}
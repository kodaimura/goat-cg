package query
/*
import (
	"database/sql"

	"goat-cg/internal/shared/dto"
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/core/db"
)


type ProjectMemberQuery interface {
	QueryJoinRequests(userId int) ([]dto.QueOutJoinRequest, error)
}


type projectMemberQuery struct {
	db *sql.DB
}


func NewProjectMemberQuery() ProjectMemberQuery {
	db := db.GetDB()
	return &projectMemberQuery{db}
}


func (que *projectMemberQuery)QueryJoinRequests(userId int) ([]dto.QueOutJoinRequest, error){

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
			 project_member pu1,
			 project_member pu2,
			 users u,
			 project p
		 WHERE pu1.user_id = ?
		  AND pu1.user_role in (?, ?)
		  AND pu2.project_id = pu1.project_id
		  AND pu2.user_status = ?
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
*/
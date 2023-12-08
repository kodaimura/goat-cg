package query

import (
	"database/sql"

	"goat-cg/internal/shared/dto"
	"goat-cg/internal/core/db"
)


type ProjectMemberQuery interface {
	GetProjectMember(projectId int) ([]dto.ProjectMember, error)
}


type projectMemberQuery struct {
	db *sql.DB
}


func NewProjectMemberQuery() ProjectMemberQuery {
	db := db.GetDB()
	return &projectMemberQuery{db}
}


func (que *projectMemberQuery)GetProjectMember(projectId int) ([]dto.ProjectMember, error){

	var ret []dto.ProjectMember
	rows, err := que.db.Query(
		`SELECT
			pm.project_id,
			pm.user_id,
			u.username,
			u.email,
			pm.user_status,
			pm.user_role,
			pm.created_at,
			pm.updated_at
		 FROM 
			 project_member pm,
			 users u
		 WHERE pm.project_id = ?
		  AND u.user_id = pm.user_id`, 
		 projectId,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		x := dto.ProjectMember{}
		err = rows.Scan(
			&x.ProjectId,
			&x.UserId,
			&x.Username,
			&x.Email,
			&x.UserStatus,
			&x.UserRole,
			&x.CreatedAt,
			&x.UpdatedAt,
		)
		if err != nil {
			break
		}
		ret = append(ret, x)
	}

	return ret, err
}
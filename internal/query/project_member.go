package query

import (
	"database/sql"

	"goat-cg/internal/dto"
	"goat-cg/internal/core/db"
)


type ProjectMemberQuery interface {
	GetProjectMembers(projectId int) ([]dto.ProjectMember, error)
	GetProjectMember(projectId, userId int) (dto.ProjectMember, error)
}


type projectMemberQuery struct {
	db *sql.DB
}


func NewProjectMemberQuery() ProjectMemberQuery {
	db := db.GetDB()
	return &projectMemberQuery{db}
}


func (que *projectMemberQuery)GetProjectMembers(projectId int) ([]dto.ProjectMember, error) {
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
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	ret := []dto.ProjectMember{}
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
			return nil, err
		}
		ret = append(ret, x)
	}

	return ret, nil
}


func (que *projectMemberQuery)GetProjectMember(projectId, userId int) (dto.ProjectMember, error) {
	var ret dto.ProjectMember
	err := que.db.QueryRow(
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
		  AND u.user_id = pm.user_id
		  AND u.user_id = ? `,
		 projectId,
		 userId,
	).Scan(
		&ret.ProjectId,
		&ret.UserId,
		&ret.Username,
		&ret.Email,
		&ret.UserStatus,
		&ret.UserRole,
		&ret.CreatedAt,
		&ret.UpdatedAt,
	)

	return ret, err
}
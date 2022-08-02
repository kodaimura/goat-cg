package service

import (
	"goat-cg/internal/model/queryservice"
)


// GetProjectId() Return value
/*----------------------------------------*/
const GET_PROJECT_ID_NOT_FOUND_INT = -1
// 正常時: プロジェクトID
/*----------------------------------------*/

func GetProjectId(userId int, projectCd string) int {
	pQue := queryservice.NewProjectQueryService()
	project, err := pQue.QueryProjectByCdAndUserId(projectCd, userId)

	if err != nil {
		return GET_PROJECT_ID_NOT_FOUND_INT
	}

	return project.ProjectId
}
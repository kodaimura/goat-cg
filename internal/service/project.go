package service

import (	
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model/entity"
	"goat-cg/internal/model/repository"
)


type ProjectService interface {
	GetProjectId(userId int, projectCd string) int 
	GetProjects(userId int) ([]entity.Project, error)
	GetProjectsPendingApproval(userId int) ([]entity.Project, error)
	GetProjectByCd(projectCd string) entity.Project
	CreateProject(userId int, projectCd, projectName string) int
	JoinRequest(userId, projectId int) int
	CancelJoinRequest(userId, projectId int) int
}


type projectService struct {
	pRep repository.ProjectRepository
	upRep repository.UserProjectRepository
}


func NewProjectService() ProjectService {
	pRep := repository.NewProjectRepository()
	upRep := repository.NewUserProjectRepository()
	return &projectService{pRep, upRep}
}

// GetProjectId() Return value
/*----------------------------------------*/
const GET_PROJECT_ID_NOT_FOUND_INT = -1
// 正常時: プロジェクトID
/*----------------------------------------*/

func (serv *projectService) GetProjectId(
	userId int, 
	projectCd string,
) int {
	project, err := serv.pRep.SelectByCdAndUserId(projectCd, userId)

	if err != nil {
		return GET_PROJECT_ID_NOT_FOUND_INT
	}

	return project.ProjectId
}


func (serv *projectService) GetProjects(
	userId int,
) ([]entity.Project, error) {
	projects, err := serv.pRep.SelectByUserIdAndStateCls(
		userId, constant.STATE_CLS_JOIN,
	)

	if err != nil {
		logger.LogError(err.Error())
	}

	return projects, err
}


func (serv *projectService) GetProjectsPendingApproval(
	userId int,
) ([]entity.Project, error) {
	projects, err := serv.pRep.SelectByUserIdAndStateCls(
		userId, constant.STATE_CLS_REQUEST,
	)

	if err != nil {
		logger.LogError(err.Error())
	}

	return projects, err
}


func (serv *projectService) GetProjectByCd(
	projectCd string,
) entity.Project {
	project, _ := serv.pRep.SelectByCd(projectCd)

	return project
}


// CreateProject() Return value
/*----------------------------------------*/
const CREATE_PROJECT_SUCCESS_INT = 0
const CREATE_PROJECT_CONFLICT_INT = 1
const CREATE_PROJECT_ERROR_INT = 2
/*----------------------------------------*/

func (serv *projectService) CreateProject(
	userId int, 
	projectCd string, 
	projectName string,
) int {
	_, err := serv.pRep.SelectByCd(projectCd)
	if err == nil {
		return CREATE_PROJECT_CONFLICT_INT
	}

	var p entity.Project
	p.ProjectCd = projectCd
	p.ProjectName = projectName
	err = serv.pRep.Insert(&p)

	if err != nil {
		logger.LogError(err.Error())
		return CREATE_PROJECT_ERROR_INT
	}

	project, err := serv.pRep.SelectByCd(projectCd)

	if err != nil {
		logger.LogError(err.Error())
		return CREATE_PROJECT_ERROR_INT
	}

	var up entity.UserProject
	up.UserId = userId
	up.ProjectId = project.ProjectId
	up.StateCls = constant.STATE_CLS_JOIN
	up.RoleCls = constant.ROLE_CLS_OWNER
	err = serv.upRep.Upsert(&up)

	if err != nil {
		logger.LogError(err.Error())
		return CREATE_PROJECT_ERROR_INT
	}

	return CREATE_PROJECT_SUCCESS_INT
}


// JoinRequest() Return value
/*----------------------------------------*/
const JOIN_REQUEST_SUCCESS_INT = 0
const JOIN_REQUEST_ALREADY_INT= 1
const JOIN_REQUEST_ERROR_INT = 2
/*----------------------------------------*/

func (serv *projectService) JoinRequest(
	userId int, projectId int,
) int {
	up0, err := serv.upRep.Select(userId, projectId)

	if err == nil && up0.StateCls == constant.STATE_CLS_JOIN {
		return JOIN_REQUEST_ALREADY_INT
	} 

	var up entity.UserProject
	up.UserId = userId
	up.ProjectId = projectId
	up.StateCls = constant.STATE_CLS_REQUEST
	up.RoleCls = constant.ROLE_CLS_NOMAL

	err = serv.upRep.Upsert(&up)

	if err != nil {
		logger.LogError(err.Error())
		return JOIN_REQUEST_ERROR_INT
	}

	return JOIN_REQUEST_SUCCESS_INT

}


// CancelJoinRequest() Return value
/*----------------------------------------*/
const CANCEL_JOIN_REQUEST_SUCCESS_INT = 0
const CANCEL_JOIN_REQUEST_ERROR_INT= 1
/*----------------------------------------*/

func (serv *projectService) CancelJoinRequest(
	userId int, projectId int,
) int {
	err := serv.upRep.Delete(userId, projectId)

	if err != nil {
		logger.LogError(err.Error())
		return CANCEL_JOIN_REQUEST_ERROR_INT
	}

	return CANCEL_JOIN_REQUEST_SUCCESS_INT

}


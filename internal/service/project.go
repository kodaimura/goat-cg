package service

import (	
	"goat-cg/internal/constant"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model/entity"
	"goat-cg/internal/model/repository"
	"goat-cg/internal/model/queryservice"
)


type ProjectService interface {
	GetProjects(userId int) ([]entity.Project, error)
	CreateProject(userId int, projectCd, projectName string) int
}


type projectService struct {
	pRep repository.ProjectRepository
	upRep repository.UserProjectRepository
	pQue queryservice.ProjectQueryService
}


func NewProjectService() ProjectService {
	pRep := repository.NewProjectRepository()
	upRep := repository.NewUserProjectRepository()
	pQue := queryservice.NewProjectQueryService()
	return &projectService{pRep, upRep, pQue}
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
	_, err := serv.pQue.QueryProjectByCd(projectCd)
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

	project, err := serv.pQue.QueryProjectByCd(projectCd)

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


func (serv *projectService) GetProjects(
	userId int,
) ([]entity.Project, error) {
	projects, err := serv.pQue.QueryProjectsByUserId(userId)

	if err != nil {
		logger.LogError(err.Error())
	}

	return projects, err
}
package service

import (	
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model"
	"goat-cg/internal/repository"
)


type ProjectService interface {
	GetProjectId(userId int, projectCd string) int 
	GetProjects(userId int) ([]model.Project, error)
	GetMemberProjects(userId int) ([]model.Project, error)
	GetProjectByCd(projectCd string) model.Project
	CreateProject(userId int, username, projectName, projectMemo string) int
}


type projectService struct {
	projectRepository repository.ProjectRepository
	projectMemberRepository repository.ProjectMemberRepository
}


func NewProjectService() ProjectService {
	projectRepository := repository.NewProjectRepository()
	projectMemberRepository := repository.NewProjectMemberRepository()
	return &projectService{projectRepository, projectMemberRepository}
}


/*----------------------------------------*/
const GET_PROJECT_ID_NOT_FOUND_INT = -1
// 正常時: プロジェクトID
/*----------------------------------------*/

// GetProjectId get projectId by userId and projectCd.
// if not return -1.  
func (serv *projectService) GetProjectId(
	userId int, 
	projectCd string,
) int {
	project, err := serv.projectRepository.GetByCdAndUserId(projectCd, userId)

	if err != nil {
		return GET_PROJECT_ID_NOT_FOUND_INT
	}

	return project.ProjectId
}


// ログインユーザのプロジェクを取得
func (serv *projectService) GetProjects(userId int) ([]model.Project, error) {
	projects, err := serv.projectRepository.GetByUserId(userId)

	if err != nil {
		logger.Error(err.Error())
	}

	return projects, err
}


//参画しているプロジェクトを取得
func (serv *projectService) GetMemberProjects(userId int) ([]model.Project, error) {
	projects, err := serv.projectRepository.GetMemberProjects(userId)

	if err != nil {
		logger.Error(err.Error())
	}

	return projects, err
}


func (serv *projectService) GetProjectByCd(
	projectCd string,
) model.Project {
	project, _ := serv.projectRepository.GetByCd(projectCd)

	return project
}


/*----------------------------------------*/
const CREATE_PROJECT_SUCCESS_INT = 0
const CREATE_PROJECT_CONFLICT_INT = 1
const CREATE_PROJECT_ERROR_INT = 2
/*----------------------------------------*/

// CreateProject create new Project.
func (serv *projectService) CreateProject(
	userId int, 
	username string,
	projectName string, 
	projectMemo string,
) int {
	_, err := serv.projectRepository.GetByUniqueKey(username, projectName)
	if err == nil {
		return CREATE_PROJECT_CONFLICT_INT
	}

	var p model.Project
	p.ProjectCd = projectCd
	p.ProjectName = projectName
	err = serv.projectRepository.Insert(&p)

	if err != nil {
		logger.Error(err.Error())
		return CREATE_PROJECT_ERROR_INT
	}

	project, err := serv.projectRepository.GetByCd(projectCd)

	if err != nil {
		logger.Error(err.Error())
		return CREATE_PROJECT_ERROR_INT
	}

	var up model.ProjectMember
	up.UserId = userId
	up.ProjectId = project.ProjectId
	up.UserStatus = constant.STATE_CLS_JOIN
	up.UserRole = constant.ROLE_CLS_OWNER
	err = serv.projectMemberRepository.Upsert(&up)

	if err != nil {
		logger.Error(err.Error())
		return CREATE_PROJECT_ERROR_INT
	}

	return CREATE_PROJECT_SUCCESS_INT
}
package service

import (	
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model/entity"
	"goat-cg/internal/model/dao"
)


type ProjectService interface {
	GetProjectId(userId int, projectCd string) int 
	GetProjects(userId int) ([]entity.Project, error)
	GetProjectsPendingApproval(userId int) ([]entity.Project, error)
	GetProjectByCd(projectCd string) entity.Project
	CreateProject(userId int, projectCd, projectName string) int
}


type projectService struct {
	pDao dao.ProjectDao
	upDao dao.ProjectUserDao
}


func NewProjectService() ProjectService {
	pDao := dao.NewProjectDao()
	upDao := dao.NewProjectUserDao()
	return &projectService{pDao, upDao}
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
	project, err := serv.pDao.SelectByCdAndUserId(projectCd, userId)

	if err != nil {
		return GET_PROJECT_ID_NOT_FOUND_INT
	}

	return project.ProjectId
}


// GetProjects get projects: the state user join.
func (serv *projectService) GetProjects(
	userId int,
) ([]entity.Project, error) {
	projects, err := serv.pDao.SelectByUserIdAndStateCls(
		userId, constant.STATE_CLS_JOIN,
	)

	if err != nil {
		logger.LogError(err.Error())
	}

	return projects, err
}


// GetProjects get projects: the state user are applying for joinrequest.
func (serv *projectService) GetProjectsPendingApproval(
	userId int,
) ([]entity.Project, error) {
	projects, err := serv.pDao.SelectByUserIdAndStateCls(
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
	project, _ := serv.pDao.SelectByCd(projectCd)

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
	projectCd string, 
	projectName string,
) int {
	_, err := serv.pDao.SelectByCd(projectCd)
	if err == nil {
		return CREATE_PROJECT_CONFLICT_INT
	}

	var p entity.Project
	p.ProjectCd = projectCd
	p.ProjectName = projectName
	err = serv.pDao.Insert(&p)

	if err != nil {
		logger.LogError(err.Error())
		return CREATE_PROJECT_ERROR_INT
	}

	project, err := serv.pDao.SelectByCd(projectCd)

	if err != nil {
		logger.LogError(err.Error())
		return CREATE_PROJECT_ERROR_INT
	}

	var up entity.ProjectUser
	up.UserId = userId
	up.ProjectId = project.ProjectId
	up.StateCls = constant.STATE_CLS_JOIN
	up.RoleCls = constant.ROLE_CLS_OWNER
	err = serv.upDao.Upsert(&up)

	if err != nil {
		logger.LogError(err.Error())
		return CREATE_PROJECT_ERROR_INT
	}

	return CREATE_PROJECT_SUCCESS_INT
}
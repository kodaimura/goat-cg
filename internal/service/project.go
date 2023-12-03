package service

import (
	"goat-cg/internal/core/logger"
	"goat-cg/internal/core/errs"
	"goat-cg/internal/model"
	"goat-cg/internal/repository"
)


type ProjectService interface {
	GetProject(projectId int) (model.Project, error)
	GetProjects(userId int) ([]model.Project, error)
	GetMemberProjects(userId int) ([]model.Project, error)
	CreateProject(userId int, username, projectName, projectMemo string) error
	UpdateProject(username string, projectId int, projectName, projectMemo string) error
}


type projectService struct {
	projectRepository repository.ProjectRepository
	//projectMemberRepository repository.ProjectMemberRepository
}


func NewProjectService() ProjectService {
	projectRepository := repository.NewProjectRepository()
	//projectMemberRepository := repository.NewProjectMemberRepository()
	return &projectService{projectRepository}//, projectMemberRepository}
}


func (serv *projectService) GetProject(projectId int) (model.Project, error) {
	project, err := serv.projectRepository.GetById(projectId)

	if err != nil {
		logger.Error(err.Error())
	}

	return project, err
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


func (serv *projectService) CreateProject(userId int, username, projectName, projectMemo string) error {
	_, err := serv.projectRepository.GetByUniqueKey(username, projectName)
	if err == nil {
		return errs.NewUniqueConstraintError("project_name")
	}

	var p model.Project
	p.ProjectName = projectName
	p.ProjectMemo = projectMemo
	p.UserId = userId
	p.Username = username
	err = serv.projectRepository.Insert(&p)

	if err != nil {
		logger.Error(err.Error())
	}

	return err
}


func (serv *projectService) UpdateProject(username string, projectId int, projectName, projectMemo string) error {
	project, err := serv.projectRepository.GetByUniqueKey(username, projectName)
	if err == nil && project.ProjectId != projectId {
		return errs.NewUniqueConstraintError("project_name")
	}

	var p model.Project
	p.ProjectId = projectId
	p.ProjectName = projectName
	p.ProjectMemo = projectMemo
	err = serv.projectRepository.Update(&p)

	if err != nil {
		logger.Error(err.Error())
	}

	return err
}
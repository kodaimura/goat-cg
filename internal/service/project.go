package service

import (
	"database/sql"

	"goat-cg/internal/core/logger"
	"goat-cg/internal/core/errs"
	"goat-cg/internal/core/db"
	"goat-cg/internal/model"
	"goat-cg/internal/repository"
)


type ProjectService interface {
	GetProject(projectId int) (model.Project, error)
	GetProjects(userId int) ([]model.Project, error)
	GetMemberProjects(userId int) ([]model.Project, error)
	CreateProject(userId int, username, projectName, projectMemo string) error
	UpdateProject(username string, projectId int, projectName, projectMemo string) error
	DeleteProject(projectId int) error
}


type projectService struct {
	projectRepository repository.ProjectRepository
	tableRepository repository.TableRepository
	columnRepository repository.ColumnRepository
}


func NewProjectService() ProjectService {
	projectRepository := repository.NewProjectRepository()
	tableRepository := repository.NewTableRepository()
	columnRepository := repository.NewColumnRepository()
	return &projectService{projectRepository, tableRepository, columnRepository}
}


func (srv *projectService) GetProject(projectId int) (model.Project, error) {
	project, err := srv.projectRepository.GetOne(&model.Project{ProjectId: projectId})

	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug(err.Error())
		} else {
			logger.Error(err.Error())
		}
	}

	return project, err
}


// ログインユーザのプロジェクを取得
func (srv *projectService) GetProjects(userId int) ([]model.Project, error) {
	projects, err := srv.projectRepository.Get(&model.Project{UserId: userId})

	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug(err.Error())
		} else {
			logger.Error(err.Error())
		}
	}

	return projects, err
}


//参画しているプロジェクトを取得
func (srv *projectService) GetMemberProjects(userId int) ([]model.Project, error) {
	projects, err := srv.projectRepository.GetMemberProjects(userId)

	if err != nil {
		logger.Error(err.Error())
	}

	return projects, err
}


func (srv *projectService) CreateProject(userId int, username, projectName, projectMemo string) error {
	_, err := srv.projectRepository.GetOne(&model.Project{Username: username, ProjectName: projectName})
	if err == nil {
		return errs.NewUniqueConstraintError("project_name")
	}

	var p model.Project
	p.ProjectName = projectName
	p.ProjectMemo = projectMemo
	p.UserId = userId
	p.Username = username

	if err = srv.projectRepository.Insert(&p, nil); err != nil {
		logger.Error(err.Error())
	}

	return err
}


func (srv *projectService) UpdateProject(username string, projectId int, projectName, projectMemo string) error {
	project, err := srv.projectRepository.GetOne(&model.Project{Username: username, ProjectName: projectName})
	if err == nil && project.ProjectId != projectId {
		return errs.NewUniqueConstraintError("project_name")
	}

	p, err := srv.projectRepository.GetOne(&model.Project{ProjectId: projectId})
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	p.ProjectName = projectName
	p.ProjectMemo = projectMemo
	if err = srv.projectRepository.Update(&p, nil); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}


func (srv *projectService) DeleteProject(projectId int) error {
	tables, err := srv.tableRepository.Get(&model.Table{ProjectId: projectId})
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	tx, err := db.GetDB().Begin()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	if err = srv.projectRepository.Delete(&model.Project{ProjectId: projectId}, tx); err != nil {
		tx.Rollback()
		logger.Error(err.Error())
		return err
	}
	
	if err = srv.tableRepository.Delete(&model.Table{ProjectId: projectId}, tx); err != nil {
		tx.Rollback()
		logger.Error(err.Error())
		return err
	}

	for _, table := range tables {
		if err = srv.columnRepository.DeleteByTableIdTx(table.TableId, tx); err != nil {
			tx.Rollback()
			logger.Error(err.Error())
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
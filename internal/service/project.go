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


func (serv *projectService) GetProject(projectId int) (model.Project, error) {
	project, err := serv.projectRepository.GetById(projectId)

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
func (serv *projectService) GetProjects(userId int) ([]model.Project, error) {
	projects, err := serv.projectRepository.GetByUserId(userId)

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

	if err = serv.projectRepository.Insert(&p); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
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

	if err = serv.projectRepository.Update(&p); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}


func (serv *projectService) DeleteProject(projectId int) error {
	var p model.Project
	p.ProjectId= projectId

	tables, err := serv.tableRepository.GetByProjectId(projectId)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	tx, err := db.GetDB().Begin()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	if err = serv.projectRepository.DeleteTx(&p, tx); err != nil {
		tx.Rollback()
		logger.Error(err.Error())
		return err
	}
	
	if err = serv.tableRepository.DeleteByProjectIdTx(projectId, tx); err != nil {
		tx.Rollback()
		logger.Error(err.Error())
		return err
	}

	for _, table := range tables {
		if err = serv.columnRepository.DeleteByTableIdTx(table.TableId, tx); err != nil {
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
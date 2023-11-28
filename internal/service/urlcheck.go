package service

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/repository"
)


type UrlCheckService interface {
	CheckProjectCdAndGetProjectId(c *gin.Context) int 
	CheckTableIdAndGetTableId(c *gin.Context, projectId int) int
	CheckColumnIdAndGetColumnId(c *gin.Context, tableId int) int
}


type urlCheckServiceice struct {
	projectRepository repository.ProjectRepository
	tableRepository repository.TableRepository
	columnRepository repository.ColumnRepository
}


func NewUrlCheckService() UrlCheckService {
	projectRepository := repository.NewProjectRepository()
	tableRepository := repository.NewTableRepository()
	columnRepository := repository.NewColumnRepository()
	return &urlCheckServiceice{projectRepository, tableRepository, columnRepository}
}


// CheckProjectCdAndGetProjectId check url parameter (/:project_cd).
// if user is accessible to the project return projectId.
func (serv *urlCheckServiceice) CheckProjectCdAndGetProjectId(
	c *gin.Context,
) int {
	userId := jwt.GetUserId(c)
	projectCd := c.Param("project_cd")

	project, err := serv.projectRepository.GetByCdAndUserId(projectCd, userId)

	if err != nil {
		logger.Error(err.Error())
		c.Redirect(303, "/projects")
		c.Abort()
		return -1
	}

	return project.ProjectId
}


// CheckTableIdAndGetTableId check url parameter (/:table_id).
// if table related to the project return tableId.
func (serv *urlCheckServiceice) CheckTableIdAndGetTableId(
	c *gin.Context, 
	projectId int,
) int {
	tableId, err := strconv.Atoi(c.Param("table_id"))
	if err != nil {
		logger.Error(err.Error())
		c.Redirect(303, "/projects")
		c.Abort()
		return -1
	}

	table, err := serv.tableRepository.GetById(tableId)

	if err != nil || table.ProjectId != projectId {
		if err != nil {
			logger.Error(err.Error())
		} else {
			logger.Error(fmt.Sprintf("table.ProjectId:%d/projectId:%d", table.ProjectId, projectId))
		}
		c.Redirect(303, "/projects")
		c.Abort()
		return -1
	}

	return tableId
}


// CheckColumnIdAndGetColumnId check url parameter (/:column_id).
// if column related to the table return columnId.
func (serv *urlCheckServiceice) CheckColumnIdAndGetColumnId(
	c *gin.Context, 
	tableId int,
) int {
	columnId, err := strconv.Atoi(c.Param("column_id"))
	if err != nil {
		logger.Error(err.Error())
		c.Redirect(303, "/projects")
		c.Abort()
		return -1
	}

	column, err := serv.columnRepository.GetById(columnId)

	if err != nil || column.TableId != tableId{
		if err != nil {
			logger.Error(err.Error())
		} else {
			logger.Error(fmt.Sprintf("column.TableId:%d/tableId:%d", column.TableId, tableId))
		}
		
		c.Redirect(303, "/projects")
		c.Abort()
		return -1
	}

	return columnId
}
package service

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model/repository"
)


type UrlCheckService interface {
	CheckProjectCdAndGetProjectId(c *gin.Context) int 
	CheckTableIdAndGetTableId(c *gin.Context, projectId int) int
	CheckColumnIdAndGetColumnId(c *gin.Context, tableId int) int
}


type urlCheckService struct {
	pRep repository.ProjectRepository
	tRep repository.TableRepository
	cRep repository.ColumnRepository
}


func NewUrlCheckService() UrlCheckService {
	pRep := repository.NewProjectRepository()
	tRep := repository.NewTableRepository()
	cRep := repository.NewColumnRepository()
	return &urlCheckService{pRep, tRep, cRep}
}


// CheckProjectCdAndGetProjectId check url parameter (/:project_cd).
// if user is accessible to the project return projectId.
func (serv *urlCheckService) CheckProjectCdAndGetProjectId(
	c *gin.Context,
) int {
	userId := jwt.GetUserId(c)
	projectCd := c.Param("project_cd")

	project, err := serv.pRep.SelectByCdAndUserId(projectCd, userId)

	if err != nil {
		logger.LogError(err.Error())
		c.Redirect(303, "/projects")
		c.Abort()
		return -1
	}

	return project.ProjectId
}


// CheckTableIdAndGetTableId check url parameter (/:table_id).
// if table related to the project return tableId.
func (serv *urlCheckService) CheckTableIdAndGetTableId(
	c *gin.Context, 
	projectId int,
) int {
	tableId, err := strconv.Atoi(c.Param("table_id"))
	if err != nil {
		logger.LogError(err.Error())
		c.Redirect(303, "/projects")
		c.Abort()
		return -1
	}

	table, err := serv.tRep.Select(tableId)

	if err != nil || table.ProjectId != projectId {
		if err != nil {
			logger.LogError(err.Error())
		} else {
			logger.LogError(fmt.Sprintf("table.ProjectId:%d/projectId:%d", table.ProjectId, projectId))
		}
		c.Redirect(303, "/projects")
		c.Abort()
		return -1
	}

	return tableId
}


// CheckColumnIdAndGetColumnId check url parameter (/:column_id).
// if column related to the table return columnId.
func (serv *urlCheckService) CheckColumnIdAndGetColumnId(
	c *gin.Context, 
	tableId int,
) int {
	columnId, err := strconv.Atoi(c.Param("column_id"))
	if err != nil {
		logger.LogError(err.Error())
		c.Redirect(303, "/projects")
		c.Abort()
		return -1
	}

	column, err := serv.cRep.Select(columnId)

	if err != nil || column.TableId != tableId{
		if err != nil {
			logger.LogError(err.Error())
		} else {
			logger.LogError(fmt.Sprintf("column.TableId:%d/tableId:%d", column.TableId, tableId))
		}
		
		c.Redirect(303, "/projects")
		c.Abort()
		return -1
	}

	return columnId
}
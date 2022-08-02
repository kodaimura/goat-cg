package service

import (
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model/entity"
	"goat-cg/internal/model/repository"
	"goat-cg/internal/model/queryservice"
)


type TableService interface {
	GetTables(projectId int) ([]entity.Table, error)
	CreateTable(
		userId, projectId int,
		tableName, tableNameLogical string,
	) int	
}


type tableService struct {
	tRep repository.TableRepository
	tQue queryservice.TableQueryService
}


func NewTableService() TableService {
	tRep := repository.NewTableRepository()
	tQue := queryservice.NewTableQueryService()

	return &tableService{tRep, tQue}
}


func (serv *tableService) GetTables(
	projectId int,
) ([]entity.Table, error) {
	tables, err := serv.tQue.QueryTables(projectId)

	if err != nil {
		logger.LogError(err.Error())
	}

	return tables, err
}


// CreateTable() Return value
/*----------------------------------------*/
const CREATE_TABLE_SUCCESS_INT = 0
const CREATE_TABLE_ERROR_INT = 1
/*----------------------------------------*/

func (serv *tableService) CreateTable(
	userId, projectId int,
	tableName, tableNameLogical string, 
) int {

	var t entity.Table
	t.ProjectId = projectId
	t.TableName = tableName
	t.TableNameLogical = tableNameLogical
	t.CreateUserId = userId
	t.UpdateUserId = userId
	err := serv.tRep.Insert(&t)

	if err != nil {
		logger.LogError(err.Error())
		return CREATE_TABLE_ERROR_INT
	}

	return CREATE_TABLE_SUCCESS_INT
}
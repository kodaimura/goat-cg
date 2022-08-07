package service

import (
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model/entity"
	"goat-cg/internal/model/repository"
)


type TableService interface {
	GetTables(projectId int) ([]entity.Table, error)
	GetTable(tableId int) (entity.Table, error)
	CreateTable(
		projectId, userId int,
		tableName, tableNameLogical string,
	) int
	UpdateTable(
		tableId, userId int,
		tableName, tableNameLogical string,
	) int	
}


type tableService struct {
	tRep repository.TableRepository
}


func NewTableService() TableService {
	tRep := repository.NewTableRepository()

	return &tableService{tRep}
}


func (serv *tableService) GetTables(
	projectId int,
) ([]entity.Table, error) {
	tables, err := serv.tRep.SelectByProjectId(projectId)

	if err != nil {
		logger.LogError(err.Error())
	}

	return tables, err
}


func (serv *tableService) GetTable(tableId int) (entity.Table, error) {
	table, err := serv.tRep.Select(tableId)

	if err != nil {
		logger.LogError(err.Error())
	}

	return table, err
}


// CreateTable() Return value
/*----------------------------------------*/
const CREATE_TABLE_SUCCESS_INT = 0
const CREATE_TABLE_ERROR_INT = 1
/*----------------------------------------*/

func (serv *tableService) CreateTable(
	projectId, userId int,
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


// UpdateTable() Return value
/*----------------------------------------*/
const UPDATE_TABLE_SUCCESS_INT = 0
const UPDATE_TABLE_ERROR_INT = 1
/*----------------------------------------*/

func (serv *tableService) UpdateTable(
	tableId, userId int,
	tableName, tableNameLogical string, 
) int {

	var t entity.Table
	t.TableName = tableName
	t.TableNameLogical = tableNameLogical
	t.UpdateUserId = userId
	err := serv.tRep.Update(tableId, &t)

	if err != nil {
		logger.LogError(err.Error())
		return UPDATE_TABLE_ERROR_INT
	}

	return UPDATE_TABLE_SUCCESS_INT
}
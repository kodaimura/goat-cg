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
		delFlg int,
	) int
	DeleteTable(tableId int) int 
}


type tableService struct {
	tRep repository.TableRepository
	cRep repository.ColumnRepository
}


func NewTableService() TableService {
	tRep := repository.NewTableRepository()
	cRep := repository.NewColumnRepository()

	return &tableService{tRep, cRep}
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
	delFlg int,
) int {

	var t entity.Table
	t.TableName = tableName
	t.TableNameLogical = tableNameLogical
	t.UpdateUserId = userId
	t.DelFlg = delFlg
	err := serv.tRep.Update(tableId, &t)

	if err != nil {
		logger.LogError(err.Error())
		return UPDATE_TABLE_ERROR_INT
	}

	return UPDATE_TABLE_SUCCESS_INT
}


// DeleteTable() Return value
/*----------------------------------------*/
const DELETE_TABLE_SUCCESS_INT = 0
const DELETE_TABLE_ERROR_INT = 1
/*----------------------------------------*/

func (serv *tableService) DeleteTable(tableId int) int {
	err := serv.tRep.Delete(tableId)

	if err != nil {
		logger.LogError(err.Error())
		return DELETE_TABLE_ERROR_INT
	}

	err = serv.cRep.DeleteByTableId(tableId)

	if err != nil {
		logger.LogError(err.Error())
		return DELETE_TABLE_ERROR_INT
	}

	return DELETE_TABLE_SUCCESS_INT
}
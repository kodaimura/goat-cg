package service

import (
	"goat-cg/internal/shared/dto"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model/entity"
	"goat-cg/internal/model/repository"
	"goat-cg/internal/model/query"
)


type TableService interface {
	GetTables(projectId int) ([]entity.Table, error)
	GetTable(tableId int) (entity.Table, error)
	CreateTable(
		projectId, userId int,
		tableName, tableNameLogical string,
	) int
	UpdateTable(
		projectId, tableId, userId int,
		tableName, tableNameLogical string,
		delFlg int,
	) int
	DeleteTable(tableId int) int 
	GetTableLog(tableId int) ([]dto.QueOutTableLog, error)
}


type tableService struct {
	tRep repository.TableRepository
	cRep repository.ColumnRepository
	tQue query.TableQuery
}


func NewTableService() TableService {
	tRep := repository.NewTableRepository()
	cRep := repository.NewColumnRepository()
	tQue := query.NewTableQuery()

	return &tableService{tRep, cRep, tQue}
}


// GetTables get tables by projeectId.
func (serv *tableService) GetTables(
	projectId int,
) ([]entity.Table, error) {
	tables, err := serv.tRep.SelectByProjectId(projectId)

	if err != nil {
		logger.LogError(err.Error())
	}

	return tables, err
}


// GetTable get table by tableId.
func (serv *tableService) GetTable(tableId int) (entity.Table, error) {
	table, err := serv.tRep.Select(tableId)

	if err != nil {
		logger.LogError(err.Error())
	}

	return table, err
}


/*----------------------------------------*/
const CREATE_TABLE_SUCCESS_INT = 0
const CREATE_TABLE_CONFLICT_INT = 1
const CREATE_TABLE_ERROR_INT = 2
/*----------------------------------------*/

// CreateTable create new Table.
func (serv *tableService) CreateTable(
	projectId, userId int,
	tableName, tableNameLogical string, 
) int {

	_, err := serv.tRep.SelectByNameAndProjectId(tableName, projectId)
	if err == nil {
		return CREATE_TABLE_CONFLICT_INT
	}

	var t entity.Table
	t.ProjectId = projectId
	t.TableName = tableName
	t.TableNameLogical = tableNameLogical
	t.CreateUserId = userId
	t.UpdateUserId = userId
	err = serv.tRep.Insert(&t)

	if err != nil {
		logger.LogError(err.Error())
		return CREATE_TABLE_ERROR_INT
	}

	return CREATE_TABLE_SUCCESS_INT
}


/*----------------------------------------*/
const UPDATE_TABLE_SUCCESS_INT = 0
const UPDATE_TABLE_CONFLICT_INT = 1
const UPDATE_TABLE_ERROR_INT = 2
/*----------------------------------------*/

// UpdateTable update Table by tableId.
// contains logical delete. 
func (serv *tableService) UpdateTable(
	projectId, tableId, userId int,
	tableName, tableNameLogical string, 
	delFlg int,
) int {

	t0, err := serv.tRep.SelectByNameAndProjectId(tableName, projectId)
	if err == nil && t0.TableId != tableId{
		return UPDATE_TABLE_CONFLICT_INT
	}

	var t entity.Table
	t.TableName = tableName
	t.TableNameLogical = tableNameLogical
	t.UpdateUserId = userId
	t.DelFlg = delFlg
	err = serv.tRep.Update(tableId, &t)

	if err != nil {
		logger.LogError(err.Error())
		return UPDATE_TABLE_ERROR_INT
	}

	return UPDATE_TABLE_SUCCESS_INT
}


/*----------------------------------------*/
const DELETE_TABLE_SUCCESS_INT = 0
const DELETE_TABLE_ERROR_INT = 1
/*----------------------------------------*/

// DeleteTable delete Table by tableId.
// (physical delete)
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


// GetTableLog get Table chenge log.
func (serv *tableService) GetTableLog(tableId int) ([]dto.QueOutTableLog, error) {
	tableLog, err := serv.tQue.QueryTableLog(tableId)

	if err != nil {
		logger.LogError(err.Error())
	}

	return tableLog, err
}
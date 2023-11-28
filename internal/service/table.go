package service

import (
	"goat-cg/internal/shared/dto"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model"
	"goat-cg/internal/repository"
	"goat-cg/internal/query"
)


type TableService interface {
	GetTables(projectId int) ([]model.Table, error)
	GetTable(tableId int) (model.Table, error)
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
	tableRepository repository.TableRepository
	columnRepository repository.ColumnRepository
	tableQuery query.TableQuery
}


func NewTableService() TableService {
	tableRepository := repository.NewTableRepository()
	columnRepository := repository.NewColumnRepository()
	tableQuery := query.NewTableQuery()

	return &tableService{tableRepository, columnRepository, tableQuery}
}


// GetTables get tables by projeectId.
func (serv *tableService) GetTables(
	projectId int,
) ([]model.Table, error) {
	tables, err := serv.tableRepository.GetByProjectId(projectId)

	if err != nil {
		logger.Error(err.Error())
	}

	return tables, err
}


// GetTable get table by tableId.
func (serv *tableService) GetTable(tableId int) (model.Table, error) {
	table, err := serv.tableRepository.GetById(tableId)

	if err != nil {
		logger.Error(err.Error())
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

	_, err := serv.tableRepository.GetByNameAndProjectId(tableName, projectId)
	if err == nil {
		return CREATE_TABLE_CONFLICT_INT
	}

	var t model.Table
	t.ProjectId = projectId
	t.TableName = tableName
	t.TableNameLogical = tableNameLogical
	t.CreateUserId = userId
	t.UpdateUserId = userId
	err = serv.tableRepository.Insert(&t)

	if err != nil {
		logger.Error(err.Error())
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

	t0, err := serv.tableRepository.GetByNameAndProjectId(tableName, projectId)
	if err == nil && t0.TableId != tableId{
		return UPDATE_TABLE_CONFLICT_INT
	}

	var t model.Table
	t.TableName = tableName
	t.TableNameLogical = tableNameLogical
	t.UpdateUserId = userId
	t.DelFlg = delFlg
	err = serv.tableRepository.Update(tableId, &t)

	if err != nil {
		logger.Error(err.Error())
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
	err := serv.tableRepository.Delete(tableId)

	if err != nil {
		logger.Error(err.Error())
		return DELETE_TABLE_ERROR_INT
	}

	err = serv.columnRepository.DeleteByTableId(tableId)

	if err != nil {
		logger.Error(err.Error())
		return DELETE_TABLE_ERROR_INT
	}

	return DELETE_TABLE_SUCCESS_INT
}


// GetTableLog get Table chenge log.
func (serv *tableService) GetTableLog(tableId int) ([]dto.QueOutTableLog, error) {
	tableLog, err := serv.tableQuery.QueryTableLog(tableId)

	if err != nil {
		logger.Error(err.Error())
	}

	return tableLog, err
}
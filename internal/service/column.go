package service

import (
	"goat-cg/internal/shared/dto"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model"
	"goat-cg/internal/repository"
	"goat-cg/internal/query"
)


type ColumnService interface {
	GetColumn(columnId int) (model.Column, error)
	GetColumns(tableId int) ([]model.Column, error)
	CreateColumn(in dto.ServInCreateColumn) int
	UpdateColumn(sin dto.ServInCreateColumn) int
	DeleteColumn(columnId int) int
	GetColumnLog(columnId int) ([]dto.QueOutColumnLog, error)
}


type columnService struct {
	columnRepository repository.ColumnRepository
	tableRepository repository.TableRepository
	columnQuery query.ColumnQuery
}


func NewColumnService() ColumnService {
	columnRepository := repository.NewColumnRepository()
	tableRepository := repository.NewTableRepository()
	columnQuery := query.NewColumnQuery()
	return &columnService{columnRepository, tableRepository, columnQuery}
}


// GetColumn get Column record by columnId.
func (serv *columnService) GetColumn(columnId int) (model.Column, error) {
	column, err := serv.columnRepository.GetById(columnId)

	if err != nil {
		logger.Error(err.Error())
	}

	return column, err
}


// GetColumn get Column records by tableId.
func (serv *columnService) GetColumns(tableId int) ([]model.Column, error) {
	columns, err := serv.columnRepository.GetByTableId(tableId)

	if err != nil {
		logger.Error(err.Error())
	}

	return columns, err
}


/*----------------------------------------*/
const CREATE_COLUMN_SUCCESS_INT = 0
const CREATE_COLUMN_CONFLICT_INT = 1
const CREATE_COLUMN_ERROR_INT = 2
/*----------------------------------------*/

// CreateColumn create new Column record.
func (serv *columnService) CreateColumn(sin dto.ServInCreateColumn) int {
	_, err := serv.columnRepository.GetByNameAndTableId(sin.ColumnName, sin.TableId)
	if err == nil {
		return CREATE_COLUMN_CONFLICT_INT
	}
	
	column := sin.ToColumn()
	err = serv.columnRepository.Insert(&column)

	if err != nil {
		logger.Error(err.Error())
		return CREATE_COLUMN_ERROR_INT
	}

	return CREATE_COLUMN_SUCCESS_INT
}


/*----------------------------------------*/
const UPDATE_COLUMN_SUCCESS_INT = 0
const UPDATE_COLUMN_CONFLICT_INT = 1
const UPDATE_COLUMN_ERROR_INT = 2
/*----------------------------------------*/

// UpdateColumn update Column record by columnId.
func (serv *columnService) UpdateColumn(sin dto.ServInCreateColumn) int {
	col, err := serv.columnRepository.GetByNameAndTableId(sin.ColumnName, sin.TableId)
	
	if err == nil && col.ColumnId != sin.ColumnId {
		return UPDATE_COLUMN_CONFLICT_INT
	}
	
	column := sin.ToColumn()
	err = serv.columnRepository.Update(&column)

	if err != nil {
		logger.Error(err.Error())
		return UPDATE_COLUMN_ERROR_INT
	}

	return UPDATE_COLUMN_SUCCESS_INT
}


/*----------------------------------------*/
const DELETE_COLUMN_SUCCESS_INT = 0
const DELETE_COLUMN_ERROR_INT = 1
/*----------------------------------------*/

// DeleteColumn delete Column record by columnId.
// (physical delete)
func (serv *columnService) DeleteColumn(columnId int) int {
	_, err := serv.columnRepository.GetById(columnId)

	if err != nil {
		logger.Error(err.Error())
		return DELETE_COLUMN_ERROR_INT
	}

	err = serv.columnRepository.Delete(columnId)

	if err != nil {
		logger.Error(err.Error())
		return DELETE_COLUMN_ERROR_INT
	}

	return DELETE_COLUMN_SUCCESS_INT
}


// GetColumnLog get Column chenge log.
func (serv *columnService) GetColumnLog(columnId int) ([]dto.QueOutColumnLog, error) {
	columnLog, err := serv.columnQuery.QueryColumnLog(columnId)

	if err != nil {
		logger.Error(err.Error())
	}

	return columnLog, err
}
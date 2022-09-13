package service

import (
	"goat-cg/internal/shared/dto"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model/entity"
	"goat-cg/internal/model/repository"
	"goat-cg/internal/model/query"
)


type ColumnService interface {
	GetColumn(columnId int) (entity.Column, error)
	GetColumns(tableId int) ([]entity.Column, error)
	CreateColumn(in dto.ServInCreateColumn) int
	UpdateColumn(columnId int, sin dto.ServInCreateColumn) int
	DeleteColumn(columnId int) int
	GetColumnLog(columnId int) ([]dto.QueOutColumnLog, error)
}


type columnService struct {
	cRep repository.ColumnRepository
	tRep repository.TableRepository
	cQue query.ColumnQuery
}


func NewColumnService() ColumnService {
	cRep := repository.NewColumnRepository()
	tRep := repository.NewTableRepository()
	cQue := query.NewColumnQuery()
	return &columnService{cRep, tRep, cQue}
}


// GetColumn get Column record by columnId.
func (serv *columnService) GetColumn(columnId int) (entity.Column, error) {
	column, err := serv.cRep.Select(columnId)

	if err != nil {
		logger.LogError(err.Error())
	}

	return column, err
}


// GetColumn get Column records by tableId.
func (serv *columnService) GetColumns(tableId int) ([]entity.Column, error) {
	columns, err := serv.cRep.SelectByTableId(tableId)

	if err != nil {
		logger.LogError(err.Error())
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
	_, err := serv.cRep.SelectByNameAndTableId(sin.ColumnName, sin.TableId)
	if err == nil {
		return CREATE_COLUMN_CONFLICT_INT
	}
	
	column := sin.ToColumn()
	err = serv.cRep.Insert(&column)

	if err != nil {
		logger.LogError(err.Error())
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
func (serv *columnService) UpdateColumn(
	columnId int, sin dto.ServInCreateColumn,
) int {
	col, err := serv.cRep.SelectByNameAndTableId(sin.ColumnName, sin.TableId)
	
	if err == nil && col.ColumnId != columnId {
		return UPDATE_COLUMN_CONFLICT_INT
	}
	
	column := sin.ToColumn()
	err = serv.cRep.Update(columnId, &column)

	if err != nil {
		logger.LogError(err.Error())
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
	_, err := serv.cRep.Select(columnId)

	if err != nil {
		logger.LogError(err.Error())
		return DELETE_COLUMN_ERROR_INT
	}

	err = serv.cRep.Delete(columnId)

	if err != nil {
		logger.LogError(err.Error())
		return DELETE_COLUMN_ERROR_INT
	}

	return DELETE_COLUMN_SUCCESS_INT
}


// GetColumnLog get Column chenge log.
func (serv *columnService) GetColumnLog(columnId int) ([]dto.QueOutColumnLog, error) {
	columnLog, err := serv.cQue.QueryColumnLog(columnId)

	if err != nil {
		logger.LogError(err.Error())
	}

	return columnLog, err
}
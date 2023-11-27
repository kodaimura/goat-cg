package service

import (
	"goat-cg/internal/shared/dto"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model"
	"goat-cg/internal/dao"
	"goat-cg/internal/query"
)


type ColumnService interface {
	GetColumn(columnId int) (model.Column, error)
	GetColumns(tableId int) ([]model.Column, error)
	CreateColumn(in dto.ServInCreateColumn) int
	UpdateColumn(columnId int, sin dto.ServInCreateColumn) int
	DeleteColumn(columnId int) int
	GetColumnLog(columnId int) ([]dto.QueOutColumnLog, error)
}


type columnService struct {
	cDao dao.ColumnDao
	tDao dao.TableDao
	cQue query.ColumnQuery
}


func NewColumnService() ColumnService {
	cDao := dao.NewColumnDao()
	tDao := dao.NewTableDao()
	cQue := query.NewColumnQuery()
	return &columnService{cDao, tDao, cQue}
}


// GetColumn get Column record by columnId.
func (serv *columnService) GetColumn(columnId int) (model.Column, error) {
	column, err := serv.cDao.Select(columnId)

	if err != nil {
		logger.Error(err.Error())
	}

	return column, err
}


// GetColumn get Column records by tableId.
func (serv *columnService) GetColumns(tableId int) ([]model.Column, error) {
	columns, err := serv.cDao.SelectByTableId(tableId)

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
	_, err := serv.cDao.SelectByNameAndTableId(sin.ColumnName, sin.TableId)
	if err == nil {
		return CREATE_COLUMN_CONFLICT_INT
	}
	
	column := sin.ToColumn()
	err = serv.cDao.Insert(&column)

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
func (serv *columnService) UpdateColumn(
	columnId int, sin dto.ServInCreateColumn,
) int {
	col, err := serv.cDao.SelectByNameAndTableId(sin.ColumnName, sin.TableId)
	
	if err == nil && col.ColumnId != columnId {
		return UPDATE_COLUMN_CONFLICT_INT
	}
	
	column := sin.ToColumn()
	err = serv.cDao.Update(columnId, &column)

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
	_, err := serv.cDao.Select(columnId)

	if err != nil {
		logger.Error(err.Error())
		return DELETE_COLUMN_ERROR_INT
	}

	err = serv.cDao.Delete(columnId)

	if err != nil {
		logger.Error(err.Error())
		return DELETE_COLUMN_ERROR_INT
	}

	return DELETE_COLUMN_SUCCESS_INT
}


// GetColumnLog get Column chenge log.
func (serv *columnService) GetColumnLog(columnId int) ([]dto.QueOutColumnLog, error) {
	columnLog, err := serv.cQue.QueryColumnLog(columnId)

	if err != nil {
		logger.Error(err.Error())
	}

	return columnLog, err
}
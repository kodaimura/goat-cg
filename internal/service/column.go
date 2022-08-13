package service

import (
	"fmt"

	"goat-cg/internal/shared/dto"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model/entity"
	"goat-cg/internal/model/repository"
)


type ColumnService interface {
	GetColumn(columnId int) (entity.Column, error)
	GetColumns(tableId int) ([]entity.Column, error)
	CreateColumn(in dto.ServInCreateColumn) int
	UpdateColumn(columnId int, sin dto.ServInCreateColumn) int
	DeleteColumn(columnId int) int

	updateTableLastLog(tableId int, action, columnName string)
}


type columnService struct {
	cRep repository.ColumnRepository
	tRep repository.TableRepository
}


func NewColumnService() ColumnService {
	cRep := repository.NewColumnRepository()
	tRep := repository.NewTableRepository()
	return &columnService{cRep, tRep}
}


func (serv *columnService) GetColumn(
	columnId int,
) (entity.Column, error) {
	column, err := serv.cRep.Select(columnId)

	if err != nil {
		logger.LogError(err.Error())
	}

	return column, err
}


func (serv *columnService) GetColumns(
	tableId int,
) ([]entity.Column, error) {
	columns, err := serv.cRep.SelectByTableId(tableId)

	if err != nil {
		logger.LogError(err.Error())
	}

	return columns, err
}


// CreateColumn() Return value
/*----------------------------------------*/
const CREATE_COLUMN_SUCCESS_INT = 0
const CREATE_COLUMN_CONFLICT_INT = 1
const CREATE_COLUMN_ERROR_INT = 2
/*----------------------------------------*/

func (serv *columnService) CreateColumn(
	sin dto.ServInCreateColumn,
) int {
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

	serv.updateTableLastLog(sin.TableId, "create", sin.ColumnName)

	return CREATE_COLUMN_SUCCESS_INT
}


// UpdateColumn() Return value
/*----------------------------------------*/
const UPDATE_COLUMN_SUCCESS_INT = 0
const UPDATE_COLUMN_CONFLICT_INT = 1
const UPDATE_COLUMN_ERROR_INT = 2
/*----------------------------------------*/

func (serv *columnService) UpdateColumn(
	columnId int,
	sin dto.ServInCreateColumn,
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

	serv.updateTableLastLog(sin.TableId, "update", sin.ColumnName)

	return UPDATE_COLUMN_SUCCESS_INT
}


// DeleteColumn() Return value
/*----------------------------------------*/
const DELETE_COLUMN_SUCCESS_INT = 0
const DELETE_COLUMN_ERROR_INT = 1
/*----------------------------------------*/

func (serv *columnService) DeleteColumn(columnId int) int {
	col, err := serv.cRep.Select(columnId)

	if err != nil {
		logger.LogError(err.Error())
		return DELETE_COLUMN_ERROR_INT
	}

	err = serv.cRep.Delete(columnId)

	if err != nil {
		logger.LogError(err.Error())
		return DELETE_COLUMN_ERROR_INT
	}

	serv.updateTableLastLog(col.TableId, "delete", col.ColumnName)

	return DELETE_COLUMN_SUCCESS_INT
}


func (serv *columnService) updateTableLastLog(tableId int, action, columnName string) {
	msg := fmt.Sprintf("%s: %s", action, columnName)

	err := serv.tRep.UpdateLastLog(tableId, msg)

	if err != nil {
		logger.LogError(err.Error())
	}
}
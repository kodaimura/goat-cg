package service

import (
	"goat-cg/internal/shared/dto"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/core/errs"
	"goat-cg/internal/model"
	"goat-cg/internal/repository"
	"goat-cg/internal/query"
)


type ColumnService interface {
	GetColumn(columnId int) (model.Column, error)
	GetColumns(tableId int) ([]model.Column, error)
	CreateColumn(in dto.ServInCreateColumn) error
	UpdateColumn(sin dto.ServInCreateColumn) error
	DeleteColumn(columnId int) error
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


// CreateColumn create new Column record.
func (serv *columnService) CreateColumn(sin dto.ServInCreateColumn) error {
	_, err := serv.columnRepository.GetByUniqueKey(sin.ColumnName, sin.TableId)
	if err == nil {
		return errs.NewUniqueConstraintError("column_name")
	}
	
	column := sin.ToColumn()
	err = serv.columnRepository.Insert(&column)

	if err != nil {
		logger.Error(err.Error())
	}

	return err
}


// UpdateColumn update Column record by columnId.
func (serv *columnService) UpdateColumn(sin dto.ServInCreateColumn) error {
	col, err := serv.columnRepository.GetByUniqueKey(sin.ColumnName, sin.TableId)

	if err == nil && col.ColumnId != sin.ColumnId {
		return errs.NewUniqueConstraintError("column_name")
	}
	
	column := sin.ToColumn()
	err = serv.columnRepository.Update(&column)

	if err != nil {
		logger.Error(err.Error())
	}

	return err
}


// DeleteColumn delete Column record by columnId.
// (physical delete)
func (serv *columnService) DeleteColumn(columnId int) error {
	err := serv.columnRepository.Delete(columnId)

	if err != nil {
		logger.Error(err.Error())
	}

	return err
}


// GetColumnLog get Column chenge log.
func (serv *columnService) GetColumnLog(columnId int) ([]dto.QueOutColumnLog, error) {
	columnLog, err := serv.columnQuery.QueryColumnLog(columnId)

	if err != nil {
		logger.Error(err.Error())
	}

	return columnLog, err
}
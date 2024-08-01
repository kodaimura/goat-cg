package service

import (
	"database/sql"

	"goat-cg/internal/dto"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/core/errs"
	"goat-cg/internal/model"
	"goat-cg/internal/repository"
	"goat-cg/internal/query"
)


type ColumnService interface {
	GetColumn(columnId int) (model.Column, error)
	GetColumns(tableId int) ([]model.Column, error)
	CreateColumn(in dto.CreateColumn) error
	UpdateColumn(sin dto.CreateColumn) error
	DeleteColumn(columnId int) error
	GetColumnLog(columnId int) ([]dto.ColumnLog, error)
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
func (srv *columnService) GetColumn(columnId int) (model.Column, error) {
	column, err := srv.columnRepository.GetById(columnId)

	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug(err.Error())
		} else {
			logger.Error(err.Error())
		}
	}

	return column, err
}


// GetColumn get Column records by tableId.
func (srv *columnService) GetColumns(tableId int) ([]model.Column, error) {
	columns, err := srv.columnRepository.GetByTableId(tableId)

	if err != nil {
		logger.Error(err.Error())
	}

	return columns, err
}


// CreateColumn create new Column record.
func (srv *columnService) CreateColumn(sin dto.CreateColumn) error {
	_, err := srv.columnRepository.GetByUniqueKey(sin.ColumnName, sin.TableId)
	if err == nil {
		return errs.NewUniqueConstraintError("column_name")
	}
	
	column := sin.ToColumn()
	
	_, err = srv.columnRepository.Insert(&column)
	if err != nil {
		logger.Error(err.Error())
	}

	return err
}


// UpdateColumn update Column record by columnId.
func (srv *columnService) UpdateColumn(sin dto.CreateColumn) error {
	col, err := srv.columnRepository.GetByUniqueKey(sin.ColumnName, sin.TableId)

	if err == nil && col.ColumnId != sin.ColumnId {
		return errs.NewUniqueConstraintError("column_name")
	}
	
	column := sin.ToColumn()

	if err = srv.columnRepository.Update(&column); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}


// DeleteColumn delete Column record by columnId.
// (physical delete)
func (srv *columnService) DeleteColumn(columnId int) error {
	var c model.Column
	c.ColumnId= columnId

	if err := srv.columnRepository.Delete(&c); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}


// GetColumnLog get Column chenge log.
func (srv *columnService) GetColumnLog(columnId int) ([]dto.ColumnLog, error) {
	columnLog, err := srv.columnQuery.GetColumnLog(columnId)

	if err != nil {
		logger.Error(err.Error())
	}

	return columnLog, err
}
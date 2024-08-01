package service

import (
	"database/sql"

	"goat-cg/internal/dto"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/core/errs"
	"goat-cg/internal/model"
	"goat-cg/internal/repository"
	"goat-cg/internal/core/db"
	"goat-cg/internal/query"
)


type TableService interface {
	GetTables(projectId int) ([]model.Table, error)
	GetTable(tableId int) (model.Table, error)
	CreateTable(projectId, userId int, tableName, tableNameLogical string) error
	UpdateTable(projectId, tableId, userId int, tableName, tableNameLogical string, delFlg int) error
	DeleteTable(tableId int) error 
	GetTableLog(tableId int) ([]dto.TableLog, error)
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
func (srv *tableService) GetTables(projectId int) ([]model.Table, error) {
	tables, err := srv.tableRepository.GetByProjectId(projectId)

	if err != nil {
		logger.Error(err.Error())
	}

	return tables, err
}


// GetTable get table by tableId.
func (srv *tableService) GetTable(tableId int) (model.Table, error) {
	table, err := srv.tableRepository.GetById(tableId)

	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug(err.Error())
		} else {
			logger.Error(err.Error())
		}
	}

	return table, err
}

// CreateTable create new Table.
func (srv *tableService) CreateTable(projectId, userId int, tableName, tableNameLogical string) error {
	_, err := srv.tableRepository.GetByUniqueKey(tableName, projectId)
	if err == nil {
		return errs.NewUniqueConstraintError("table_name")
	}

	var t model.Table
	t.ProjectId = projectId
	t.TableName = tableName
	t.TableNameLogical = tableNameLogical
	t.CreateUserId = userId
	t.UpdateUserId = userId

	_, err = srv.tableRepository.Insert(&t);
	if err != nil {
		logger.Error(err.Error())
	}

	return err
}


// UpdateTable update Table by tableId.
// contains logical delete. 
func (srv *tableService) UpdateTable(projectId, tableId, userId int, tableName, tableNameLogical string, delFlg int) error {
	table, err := srv.tableRepository.GetByUniqueKey(tableName, projectId)
	if err == nil && table.TableId != tableId {
		return errs.NewUniqueConstraintError("table_name")
	}

	var t model.Table
	t.TableId = tableId
	t.TableName = tableName
	t.TableNameLogical = tableNameLogical
	t.UpdateUserId = userId
	t.DelFlg = delFlg

	if err = srv.tableRepository.Update(&t); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}


// DeleteTable delete Table by tableId.
// (physical delete)
func (srv *tableService) DeleteTable(tableId int) error {
	var t model.Table
	t.TableId= tableId

	tx, err := db.GetDB().Begin()
	if err != nil {
		tx.Rollback()
		logger.Error(err.Error())
		return err
	}

	if err = srv.tableRepository.DeleteTx(&t, tx); err != nil {
		tx.Rollback()
		logger.Error(err.Error())
		return err
	}

	if err = srv.columnRepository.DeleteByTableIdTx(tableId, tx); err != nil {
		tx.Rollback()
		logger.Error(err.Error())
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}


// GetTableLog get Table chenge log.
func (srv *tableService) GetTableLog(tableId int) ([]dto.TableLog, error) {
	tableLog, err := srv.tableQuery.GetTableLog(tableId)

	if err != nil {
		logger.Error(err.Error())
	}

	return tableLog, err
}
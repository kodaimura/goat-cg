package service

import (
	"goat-cg/internal/shared/dto"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model/entity"
	"goat-cg/internal/model/repository"
	"goat-cg/internal/model/queryservice"
)


type ColumnService interface {
	GetColumns(tableId int) ([]entity.Column, error)
	CreateColumn(in dto.ServInCreateColumn) int
}


type columnService struct {
	cRep repository.ColumnRepository
	cQue queryservice.ColumnQueryService
}


func NewColumnService() ColumnService {
	cRep := repository.NewColumnRepository()
	cQue := queryservice.NewColumnQueryService()

	return &columnService{cRep, cQue}
}


func (serv *columnService) GetColumns(
	tableId int,
) ([]entity.Column, error) {
	columns, err := serv.cQue.QueryColumns(tableId)

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
	_, err := serv.cQue.QueryColumnByNameAndTableId(sin.ColumnName, sin.TableId)
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
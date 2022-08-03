package service

import (
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model/entity"
	"goat-cg/internal/model/repository"
	"goat-cg/internal/model/queryservice"
)


type ColumnService interface {
	GetColumns(tableId int) ([]entity.Column, error)
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
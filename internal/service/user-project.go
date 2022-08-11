package service

import (	
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/shared/dto"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model/entity"
	"goat-cg/internal/model/repository"
	"goat-cg/internal/model/query"
)


type UserProjectService interface {
	JoinRequest(userId, projectId int) int
	CancelJoinRequest(userId, projectId int) int
	PermitJoinRequest(userId, projectId int) int
	GetJoinRequests(userId int) ([]dto.QueOutJoinRequest, error)
}


type userProjectService struct {
	upRep repository.UserProjectRepository
	upQue query.UserProjectQuery
}


func NewUserProjectService() UserProjectService {
	upRep := repository.NewUserProjectRepository()
	upQue := query.NewUserProjectQuery()
	return &userProjectService{upRep, upQue}
}



// JoinRequest() Return value
/*----------------------------------------*/
const JOIN_REQUEST_SUCCESS_INT = 0
const JOIN_REQUEST_ALREADY_INT= 1
const JOIN_REQUEST_ERROR_INT = 2
/*----------------------------------------*/

func (serv *userProjectService) JoinRequest(
	userId int, projectId int,
) int {
	up0, err := serv.upRep.Select(userId, projectId)

	if err == nil && up0.StateCls == constant.STATE_CLS_JOIN {
		return JOIN_REQUEST_ALREADY_INT
	} 

	var up entity.UserProject
	up.UserId = userId
	up.ProjectId = projectId
	up.StateCls = constant.STATE_CLS_REQUEST
	up.RoleCls = constant.ROLE_CLS_NOMAL

	err = serv.upRep.Upsert(&up)

	if err != nil {
		logger.LogError(err.Error())
		return JOIN_REQUEST_ERROR_INT
	}

	return JOIN_REQUEST_SUCCESS_INT

}


// CancelJoinRequest() Return value
/*----------------------------------------*/
const CANCEL_JOIN_REQUEST_SUCCESS_INT = 0
const CANCEL_JOIN_REQUEST_ERROR_INT= 1
/*----------------------------------------*/

func (serv *userProjectService) CancelJoinRequest(
	userId int, projectId int,
) int {
	err := serv.upRep.Delete(userId, projectId)

	if err != nil {
		logger.LogError(err.Error())
		return CANCEL_JOIN_REQUEST_ERROR_INT
	}

	return CANCEL_JOIN_REQUEST_SUCCESS_INT

}


// PermitJoinRequest() Return value
/*----------------------------------------*/
const PERMIT_JOIN_REQUEST_SUCCESS_INT = 0
const PERMIT_JOIN_REQUEST_ERROR_INT= 1
/*----------------------------------------*/

func (serv *userProjectService) PermitJoinRequest(
	userId int, projectId int,
) int {

	var up entity.UserProject
	up.UserId = userId
	up.ProjectId = projectId
	up.StateCls = constant.STATE_CLS_JOIN
	up.RoleCls = constant.ROLE_CLS_NOMAL

	err := serv.upRep.Upsert(&up)

	if err != nil {
		logger.LogError(err.Error())
		return PERMIT_JOIN_REQUEST_ERROR_INT
	}

	return PERMIT_JOIN_REQUEST_SUCCESS_INT

}


func (serv *userProjectService) GetJoinRequests(
	userId int,
) ([]dto.QueOutJoinRequest, error) {
	jrs, err := serv.upQue.QueryJoinRequests(userId)

	if err != nil {
		logger.LogError(err.Error())
	}

	return jrs, err

}
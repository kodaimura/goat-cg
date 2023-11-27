package service

import (	
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/shared/dto"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model"
	"goat-cg/internal/dao"
	"goat-cg/internal/query"
)


type ProjectUserService interface {
	JoinRequest(userId, projectId int) int
	CancelJoinRequest(userId, projectId int) int
	PermitJoinRequest(userId, projectId int) int
	GetJoinRequests(userId int) ([]dto.QueOutJoinRequest, error)
}


type projectUserService struct {
	upDao dao.ProjectUserDao
	upQue query.ProjectUserQuery
}


func NewProjectUserService() ProjectUserService {
	upDao := dao.NewProjectUserDao()
	upQue := query.NewProjectUserQuery()
	return &projectUserService{upDao, upQue}
}


/*----------------------------------------*/
const JOIN_REQUEST_SUCCESS_INT = 0
const JOIN_REQUEST_ALREADY_INT= 1
const JOIN_REQUEST_ERROR_INT = 2
/*----------------------------------------*/

// JoinRequest
func (serv *projectUserService) JoinRequest(
	userId int, projectId int,
) int {
	up0, err := serv.upDao.Select(userId, projectId)

	if err == nil && up0.StateCls == constant.STATE_CLS_JOIN {
		return JOIN_REQUEST_ALREADY_INT
	} 

	var up model.ProjectUser
	up.UserId = userId
	up.ProjectId = projectId
	up.StateCls = constant.STATE_CLS_REQUEST
	up.RoleCls = constant.ROLE_CLS_NOMAL

	err = serv.upDao.Upsert(&up)

	if err != nil {
		logger.Error(err.Error())
		return JOIN_REQUEST_ERROR_INT
	}

	return JOIN_REQUEST_SUCCESS_INT

}


/*----------------------------------------*/
const CANCEL_JOIN_REQUEST_SUCCESS_INT = 0
const CANCEL_JOIN_REQUEST_ERROR_INT= 1
/*----------------------------------------*/

// CancelJoinRequest
func (serv *projectUserService) CancelJoinRequest(
	userId int, projectId int,
) int {
	err := serv.upDao.Delete(userId, projectId)

	if err != nil {
		logger.Error(err.Error())
		return CANCEL_JOIN_REQUEST_ERROR_INT
	}

	return CANCEL_JOIN_REQUEST_SUCCESS_INT

}


/*----------------------------------------*/
const PERMIT_JOIN_REQUEST_SUCCESS_INT = 0
const PERMIT_JOIN_REQUEST_ERROR_INT= 1
/*----------------------------------------*/

// PermitJoinRequest
func (serv *projectUserService) PermitJoinRequest(
	userId int, projectId int,
) int {

	var up model.ProjectUser
	up.UserId = userId
	up.ProjectId = projectId
	up.StateCls = constant.STATE_CLS_JOIN
	up.RoleCls = constant.ROLE_CLS_NOMAL

	err := serv.upDao.Upsert(&up)

	if err != nil {
		logger.Error(err.Error())
		return PERMIT_JOIN_REQUEST_ERROR_INT
	}

	return PERMIT_JOIN_REQUEST_SUCCESS_INT

}


// GetJoinRequests get requests to join project.
// that login user have the authority to permit.
func (serv *projectUserService) GetJoinRequests(
	userId int,
) ([]dto.QueOutJoinRequest, error) {
	jrs, err := serv.upQue.QueryJoinRequests(userId)

	if err != nil {
		logger.Error(err.Error())
	}

	return jrs, err

}
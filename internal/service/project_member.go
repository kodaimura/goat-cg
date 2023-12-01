package service

//import (	
//	"goat-cg/internal/shared/constant"
//	"goat-cg/internal/shared/dto"
//	"goat-cg/internal/core/logger"
//	"goat-cg/internal/model"
//	"goat-cg/internal/repository"
//	"goat-cg/internal/query"
//)
//
//
//type ProjectMemberService interface {
//	JoinRequest(userId, projectId int) int
//	CancelJoinRequest(userId, projectId int) int
//	PermitJoinRequest(userId, projectId int) int
//	GetJoinRequests(userId int) ([]dto.QueOutJoinRequest, error)
//}
//
//
//type projectMemberService struct {
//	projectMemberRepository repository.ProjectMemberRepository
//	projectMemberQuery query.ProjectMemberQuery
//}
//
//
//func NewProjectMemberService() ProjectMemberService {
//	projectMemberRepository := repository.NewProjectMemberRepository()
//	projectMemberQuery := query.NewProjectMemberQuery()
//	return &projectMemberService{projectMemberRepository, projectMemberQuery}
//}

///*----------------------------------------*/
//const JOIN_REQUEST_SUCCESS_INT = 0
//const JOIN_REQUEST_ALREADY_INT= 1
//const JOIN_REQUEST_ERROR_INT = 2
///*----------------------------------------*/
//
//// JoinRequest
//func (serv *projectMemberService) JoinRequest(
//	userId int, projectId int,
//) int {
//	up0, err := serv.projectMemberRepository.GetByPk(userId, projectId)
//
//	if err == nil && up0.UserStatus == constant.STATE_CLS_JOIN {
//		return JOIN_REQUEST_ALREADY_INT
//	} 
//
//	var up model.ProjectMember
//	up.UserId = userId
//	up.ProjectId = projectId
//	up.UserStatus = constant.STATE_CLS_REQUEST
//	up.UserRole = constant.ROLE_CLS_NOMAL
//
//	err = serv.projectMemberRepository.Upsert(&up)
//
//	if err != nil {
//		logger.Error(err.Error())
//		return JOIN_REQUEST_ERROR_INT
//	}
//
//	return JOIN_REQUEST_SUCCESS_INT
//
//}
//
//
///*----------------------------------------*/
//const CANCEL_JOIN_REQUEST_SUCCESS_INT = 0
//const CANCEL_JOIN_REQUEST_ERROR_INT= 1
///*----------------------------------------*/
//
//// CancelJoinRequest
//func (serv *projectMemberService) CancelJoinRequest(
//	userId int, projectId int,
//) int {
//	err := serv.projectMemberRepository.Delete(userId, projectId)
//
//	if err != nil {
//		logger.Error(err.Error())
//		return CANCEL_JOIN_REQUEST_ERROR_INT
//	}
//
//	return CANCEL_JOIN_REQUEST_SUCCESS_INT
//
//}
//
//
///*----------------------------------------*/
//const PERMIT_JOIN_REQUEST_SUCCESS_INT = 0
//const PERMIT_JOIN_REQUEST_ERROR_INT= 1
///*----------------------------------------*/
//
//// PermitJoinRequest
//func (serv *projectMemberService) PermitJoinRequest(
//	userId int, projectId int,
//) int {
//
//	var up model.ProjectMember
//	up.UserId = userId
//	up.ProjectId = projectId
//	up.UserStatus = constant.STATE_CLS_JOIN
//	up.UserRole = constant.ROLE_CLS_NOMAL
//
//	err := serv.projectMemberRepository.Upsert(&up)
//
//	if err != nil {
//		logger.Error(err.Error())
//		return PERMIT_JOIN_REQUEST_ERROR_INT
//	}
//
//	return PERMIT_JOIN_REQUEST_SUCCESS_INT
//
//}
//
//
//// GetJoinRequests get requests to join project.
//// that login user have the authority to permit.
//func (serv *projectMemberService) GetJoinRequests(
//	userId int,
//) ([]dto.QueOutJoinRequest, error) {
//	jrs, err := serv.projectMemberQuery.QueryJoinRequests(userId)
//
//	if err != nil {
//		logger.Error(err.Error())
//	}
//
//	return jrs, err
//
//}
//
//
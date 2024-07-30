package service

import (
	"goat-cg/internal/core/errs"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model"
	"goat-cg/internal/dto"
	"goat-cg/internal/repository"
	"goat-cg/internal/query"
)


type MemberService interface {
	Invite(projectId int, email string) error
	GetMembers(projectId int) ([]dto.ProjectMember, error)
	GetMember(projectId, userId int) (dto.ProjectMember, error)
	DeleteMember(projectId, userId int) error
}


type memberService struct {
	memberRepository repository.MemberRepository
	userRepository repository.UserRepository
	projectMemberQuery query.ProjectMemberQuery
}


func NewMemberService() MemberService {
	memberRepository := repository.NewMemberRepository()
	userRepository := repository.NewUserRepository()
	projectMemberQuery := query.NewProjectMemberQuery()
	return &memberService{memberRepository, userRepository, projectMemberQuery}
}

func (srv *memberService) Invite(projectId int, email string) error {
	user, err := srv.userRepository.GetOne(&model.User{Email: email})
	if err != nil {
		return errs.NewNotFoundError()
	}

	_, err = srv.memberRepository.GetOne(&model.Member{ProjectId: projectId, UserId: user.UserId})
	if err == nil {
		return errs.NewAlreadyRegisteredError()
	}

	var m model.Member
	m.ProjectId = projectId
	m.UserId = user.UserId
	m.UserStatus = "0"
	m.UserRole = "0"

	if err = srv.memberRepository.Insert(&m, nil); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (srv *memberService) GetMembers(projectId int) ([]dto.ProjectMember, error) {
	return srv.projectMemberQuery.GetProjectMembers(projectId)
}

func (srv *memberService) GetMember(projectId, userId int) (dto.ProjectMember, error) {
	return srv.projectMemberQuery.GetProjectMember(projectId, userId)
}

func (srv *memberService) DeleteMember(projectId, userId int) error {
	var m model.Member
	m.ProjectId = projectId
	m.UserId = userId

	if err := srv.memberRepository.Delete(&m, nil); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}
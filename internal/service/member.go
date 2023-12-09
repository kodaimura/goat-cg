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

func (serv *memberService) Invite(projectId int, email string) error {
	user, err := serv.userRepository.GetByEmail(email)
	if err != nil {
		return errs.NewNotFoundError()
	}

	_, err = serv.memberRepository.GetByPk(projectId, user.UserId)
	if err == nil {
		return errs.NewAlreadyRegisteredError()
	}

	var m model.Member
	m.ProjectId = projectId
	m.UserId = user.UserId
	m.UserStatus = "0"
	m.UserRole = "0"

	if err = serv.memberRepository.Insert(&m); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (serv *memberService) GetMembers(projectId int) ([]dto.ProjectMember, error) {
	return serv.projectMemberQuery.GetProjectMember(projectId)
}
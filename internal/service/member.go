package service

import (
	"goat-cg/internal/model"
	"goat-cg/internal/repository"
	"goat-cg/internal/core/errs"
	"goat-cg/internal/core/logger"
)


type MemberService interface {
	Invite(projectId int, email string) error
}


type memberService struct {
	memberRepository repository.MemberRepository
	userRepository repository.UserRepository
}


func NewMemberService() MemberService {
	memberRepository := repository.NewMemberRepository()
	userRepository := repository.NewUserRepository()
	return &memberService{memberRepository, userRepository}
}

func (serv *memberService) Invite(projectId int, email string) error {
	user, err := serv.userRepository.GetByEmail(email)
	if err != nil {
		return errs.NewNotFoundError()
	}

	_, err = serv.memberRepository.GetByPk(projectId, user.UserId)
	if err != nil {
		return errs.NewNotFoundError()
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
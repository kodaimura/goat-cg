package service

import (
	"goat-cg/internal/repository"
)


type MemberService interface {
}


type memberService struct {
	memberRepository repository.MemberRepository
}


func NewMemberService() MemberService {
	memberRepository := repository.NewMemberRepository()
	return &memberService{memberRepository}
}

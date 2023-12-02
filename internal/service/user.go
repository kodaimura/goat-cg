package service

import (
	"golang.org/x/crypto/bcrypt"

	"goat-cg/internal/core/jwt"
	
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model"
	"goat-cg/internal/repository"
)


type UserService interface {
	Signup(username, password string) int
	Login(username, password string) int
	GenerateJWT(userId int) string
	GetProfile(userId int) (model.User, error)
	ChangeUsername(userId int, username string) int
	ChangePassword(userId int, password string) int
	DeleteUser(userId int) int
}


type userService struct {
	userRepository repository.UserRepository
}


func NewUserService() UserService {
	userRepository := repository.NewUserRepository()
	return &userService{userRepository}
}


// Signup() Return value
/*----------------------------------------*/
const SIGNUP_SUCCESS_INT = 0
const SIGNUP_CONFLICT_INT = 1
const SIGNUP_ERROR_INT = 2
/*----------------------------------------*/

func (serv *userService) Signup(username, password string) int {
	_, err := serv.userRepository.GetByName(username)

	if err == nil {
		return SIGNUP_CONFLICT_INT
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		logger.Error(err.Error())
		return SIGNUP_ERROR_INT
	}

	var user model.User
	user.Username = username
	user.Password = string(hashed)

	err = serv.userRepository.Insert(&user)

	if err != nil {
		logger.Error(err.Error())
		return SIGNUP_ERROR_INT
	}

	return SIGNUP_SUCCESS_INT
}


// Login() Return value
/*----------------------------------------*/
const LOGIN_FAILURE_INT = -1
// 正常時: ユーザ識別ID
/*----------------------------------------*/

func (serv *userService) Login(username, password string) int {
	user, err := serv.userRepository.GetByName(username)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return LOGIN_FAILURE_INT
	}

	return user.UserId
}


// GenerateJWT() Return value
/*----------------------------------------*/
const GENERATE_JWT_FAILURE_STR = ""
// 正常時: jwt文字列
/*----------------------------------------*/

func (serv *userService) GenerateJWT(userId int) string {
	user, err := serv.userRepository.GetById(userId)
	
	if err != nil {
		logger.Error(err.Error())
		return GENERATE_JWT_FAILURE_STR
	}

	var cc jwt.CustomClaims
	cc.UserId = user.UserId
	cc.Username = user.Username
	jwtStr, err := jwt.GenerateJWT(cc)

	if err != nil {
		logger.Error(err.Error())
		return GENERATE_JWT_FAILURE_STR
	}

	return jwtStr
}


func (serv *userService) GetProfile(userId int) (model.User, error) {
	user, err := serv.userRepository.GetById(userId)

	if err != nil {
		logger.Error(err.Error())
	}

	return user, err
}


// ChangeUsername() Return value
/*----------------------------------------*/
const CHANGE_EMAIL_SUCCESS_INT = 0
const CHANGE_EMAIL_FAILURE_INT = 1
/*----------------------------------------*/
func (serv *userService) ChangeEmail(userId int, email string) int {
	err := serv.userRepository.UpdateEmail(userId, email)

	if err != nil {
		logger.Error(err.Error())
		return CHANGE_EMAIL_FAILURE_INT
	}

	return CHANGE_EMAIL_SUCCESS_INT
}


// ChangePassword() Return value
/*----------------------------------------*/
const CHANGE_PASSWORD_SUCCESS_INT = 0
const CHANGE_PASSWORD_FAILURE_INT = 1
/*----------------------------------------*/
func (serv *userService) ChangePassword(userId int, password string) int {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		logger.Error(err.Error())
		return CHANGE_PASSWORD_FAILURE_INT
	}

	err = serv.userRepository.UpdatePassword(userId, string(hashed))
	
	if err != nil {
		logger.Error(err.Error())
		return CHANGE_PASSWORD_FAILURE_INT
	}

	return CHANGE_PASSWORD_SUCCESS_INT
}


// DeleteUser() Return value
/*----------------------------------------*/
const DELETE_USER_SUCCESS_INT = 0
const DELETE_USER_FAILURE_INT = 1
/*----------------------------------------*/
func (serv *userService) DeleteUser(userId int) int {
	err := serv.userRepository.Delete(userId)

	if err != nil {
		logger.Error(err.Error())
		return DELETE_USER_FAILURE_INT
	}

	return DELETE_USER_SUCCESS_INT
}

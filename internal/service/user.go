package service

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/core/errs"
	"goat-cg/internal/model"
	"goat-cg/internal/repository"
)


type UserService interface {
	Signup(username, password, email string) error
	Login(username, password string) (model.User, error)
	GenerateJWT(id int) (string, error)
	GetProfile(id int) (model.User, error)
	UpdateEmail(id int, email string) error
	UpdatePassword(id int, password string) error
	DeleteUser(id int) error
}


type userService struct {
	userRepository repository.UserRepository
}

func NewUserService() UserService {
	return &userService{
		userRepository: repository.NewUserRepository(),
	}
}


func (us *userService) Signup(username, password, email string) error {
	_, err := us.userRepository.GetByName(username)

	if err == nil {
		return errs.NewUniqueConstraintError("username")
	}

	_, err = us.userRepository.GetByEmail(email)

	if err == nil {
		return errs.NewUniqueConstraintError("email")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	var user model.User
	user.Username = username
	user.Password = string(hashed)
	user.Email = email

	if err = us.userRepository.Insert(&user); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}


func (us *userService) Login(username, password string) (model.User, error) {
	user, err := us.userRepository.GetByName(username)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug(err.Error())
		} else {
			logger.Error(err.Error())
		}
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.Error(err.Error())
	}

	return user, err
}


func (us *userService) GenerateJWT(id int) (string, error) {
	user, err := us.userRepository.GetById(id)
	
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	var cc jwt.CustomClaims
	cc.UserId = user.UserId
	cc.Username = user.Username
	cc.Email = user.Email
	jwtStr, err := jwt.GenerateJWT(cc)

	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	return jwtStr, nil
}


func (us *userService) GetProfile(id int) (model.User, error) {
	user, err := us.userRepository.GetById(id)

	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug(err.Error())
		} else {
			logger.Error(err.Error())
		}
	}

	return user, err
}


func (us *userService) UpdateEmail(id int, email string) error {
	u, err := us.userRepository.GetByEmail(email)

	if err == nil && u.UserId != id{
		return errs.NewUniqueConstraintError("email")
	}

	var user model.User
	user.UserId = id
	user.Email = email

	if err = us.userRepository.UpdateEmail(&user); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}


func (us *userService) UpdatePassword(id int, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	var user model.User
	user.UserId = id
	user.Password = string(hashed)
	
	if err = us.userRepository.UpdatePassword(&user); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}


func (us *userService) DeleteUser(id int) error {
	var user model.User
	user.UserId = id

	if err := us.userRepository.Delete(&user); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

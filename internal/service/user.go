package service

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/core/errs"
	"goat-cg/internal/model"
	"goat-cg/internal/dto"
	"goat-cg/internal/repository"
)


type UserService interface {
	Signup(name, password, email string) error
	Login(name, password string) (dto.User, error)
	GenerateJWT(id int) (string, error)
	GetProfile(id int) (dto.User, error)
	UpdateName(id int, name string) error
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


func (srv *userService) toUserDTO(user model.User) dto.User {
	return dto.User{
		UserId:    user.UserId,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}


func (srv *userService) Signup(name, password, email string) error {
	_, err := srv.userRepository.GetOne(&model.User{Username: name})
	if err == nil {
		return errs.NewUniqueConstraintError("username")
	}
	_, err = srv.userRepository.GetOne(&model.User{Email: email})
	if err == nil {
		return errs.NewUniqueConstraintError("email")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	var user model.User
	user.Username = name
	user.Password = string(hashed)
	user.Email = email

	err = srv.userRepository.Insert(&user, nil);
	if err != nil {
		logger.Error(err.Error())
	}

	return err
}


func (srv *userService) Login(name, password string) (dto.User, error) {
	user, err := srv.userRepository.GetOne(&model.User{Username: name})
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug(err.Error())
		} else {
			logger.Error(err.Error())
		}
		return dto.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.Error(err.Error())
	}

	return srv.toUserDTO(user), err
}


func (srv *userService) GetProfile(id int) (dto.User, error) {
	user, err := srv.userRepository.GetOne(&model.User{UserId: id})
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug(err.Error())
		} else {
			logger.Error(err.Error())
		}
	}

	return srv.toUserDTO(user), err
}


func (srv *userService) GenerateJWT(id int) (string, error) {
	user, err := srv.userRepository.GetOne(&model.User{UserId: id})
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


func (srv *userService) UpdateName(id int, name string) error {
	u, err := srv.userRepository.GetOne(&model.User{Username: name})
	if err == nil && u.UserId != id{
		return errs.NewUniqueConstraintError("username")
	}

	var user model.User
	user.UserId = id
	user.Username = name

	if err = srv.userRepository.UpdateName(&user, nil); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}


func (srv *userService) UpdatePassword(id int, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	var user model.User
	user.UserId = id
	user.Password = string(hashed)
	
	if err = srv.userRepository.UpdatePassword(&user, nil); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}


func (srv *userService) UpdateEmail(id int, email string) error {
	u, err := srv.userRepository.GetOne(&model.User{Email: email})
	if err == nil && u.UserId != id{
		return errs.NewUniqueConstraintError("email")
	}

	var user model.User
	user.UserId = id
	user.Email = email

	if err = srv.userRepository.UpdateEmail(&user, nil); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}


func (srv *userService) DeleteUser(id int) error {
	var user model.User
	user.UserId = id

	if err := srv.userRepository.Delete(&user, nil); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}
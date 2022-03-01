package service

/*
Service ではドメインの仕様をチェックする
*/
import (
	"22dojo-online/pkg/domain/entity"
	"22dojo-online/pkg/domain/repository"
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

type UserServiceInterface interface {
	SelectUserByPrimaryKey(string) (*entity.User, error)
	SelectUserByAuthToken(string) (*entity.User, error)
	InsertUser(user *entity.User) error
	UpdateUserByPrimaryKey(user *entity.User) error
}

type UserService struct {
	UserRepository repository.UserRepository
}

func (us *UserService) InsertUser(user *entity.User) error {
	if user.Name == "" {
		validate := validator.New()
		err := validate.Struct(user)
		return fmt.Errorf("user name is empty. err = %w", err)
	}
	if err := us.UserRepository.InsertUser(user); err != nil {
		return fmt.Errorf("failed to insert user. err = %w", err)
	}
	return nil
}

func (us *UserService) SelectUserByPrimaryKey(userID string) (*entity.User, error) {
	user, err := us.UserRepository.SelectUserByPrimaryKey(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query a user by primary key. err = %w", err)
	}
	return user, nil
}

func (us *UserService) SelectUserByAuthToken(authToken string) (*entity.User, error) {
	user, err := us.UserRepository.SelectUserByPrimaryKey(authToken)
	if err != nil {
		return nil, fmt.Errorf("failed to query a user by primary key. err = %w", err)
	}
	return user, nil
}

func (us *UserService) UpdateUserByPrimaryKey(user *entity.User) error {
	if err := us.UserRepository.UpdateUserByPrimaryKey(user); err != nil {
		return fmt.Errorf("failed to update user. err = %w", err)
	}
	return nil
}

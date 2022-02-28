package usecase

import (
	"22dojo-online/pkg/domain"
	"22dojo-online/pkg/interfaces/database"
	"fmt"
)

type UserInteractor struct {
	UserRepository UserRepository
}

func NewUserInteractor(sqlHandler database.SQLHandler) *UserInteractor {
	return &UserInteractor{
		UserRepository: &database.UserRepository{
			SQLHandler: sqlHandler,
		},
	}
}

func (interactor *UserInteractor) SelectUserByPrimaryKey(userID string) (*domain.User, error) {
	user, err := interactor.UserRepository.SelectUserByPrimaryKey(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query a user by primary key. err = %w", err)
	}
	return user, nil
}

func (interactor *UserInteractor) SelectUserByAuthToken(authToken string) (*domain.User, error) {
	user, err := interactor.UserRepository.SelectUserByAuthToken(authToken)
	if err != nil {
		return nil, fmt.Errorf("failed to query a user by auth token. err = %w", err)
	}
	return user, nil
}

func (interactor *UserInteractor) InsertUser(user *domain.User) error {
	if err := interactor.UserRepository.InsertUser(user); err != nil {
		return fmt.Errorf("failed to insert user. err = %w", err)
	}
	return nil
}

func (interactor *UserInteractor) UpdateUserByPrimaryKey(user *domain.User) error {
	if err := interactor.UserRepository.UpdateUserByPrimaryKey(user); err != nil {
		return fmt.Errorf("failed to update user. err = %w", err)
	}
	return nil
}

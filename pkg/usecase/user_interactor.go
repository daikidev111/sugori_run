package usecase

import (
	"22dojo-online/pkg/domain"
	"22dojo-online/pkg/interfaces/database"
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

func (interactor *UserInteractor) SelectUserByPrimaryKey(userID string) (user domain.User, err error) {
	user, err = interactor.UserRepository.SelectUserByPrimaryKey(userID)
	return
}

func (interactor *UserInteractor) SelectUserByAuthToken(authToken string) (user domain.User, err error) {
	user, err = interactor.UserRepository.SelectUserByAuthToken(authToken)
	return
}

func (interactor *UserInteractor) InsertUser(user domain.User) (err error) {
	err = interactor.UserRepository.InsertUser(user)
	return
}

func (interactor *UserInteractor) UpdateUserByPrimaryKey(user domain.User) (err error) {
	err = interactor.UserRepository.UpdateUserByPrimaryKey(user)
	return
}

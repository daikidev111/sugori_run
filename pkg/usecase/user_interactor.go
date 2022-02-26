package usecase

import (
	"22dojo-online/pkg/domain"
	"22dojo-online/pkg/interfaces/database"
)

type UserInteractor struct {
	UserRepository UserRepository
}

func NewUserInteractor(sqlHandler database.SqlHandler) *UserInteractor {
	return &UserInteractor{
		UserRepository: &database.UserRepository{
			SqlHandler: sqlHandler,
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

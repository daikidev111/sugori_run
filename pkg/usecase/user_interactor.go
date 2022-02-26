package usecase

import "22dojo-online/pkg/domain"

type UserInteractor struct {
	UserRepository UserRepository
}

func (interactor *UserInteractor) SelectUserByPrimaryKey(userID string) (user domain.User, err error) {
	user, err = interactor.UserRepository.SelectUserByPrimaryKey(userID)
	return
}

func (interactor *UserInteractor) SelectUserByAuthToken(authToken string) (user domain.User, err error) {
	user, err = interactor.UserRepository.SelectUserByAuthToken(authToken)
	return
}

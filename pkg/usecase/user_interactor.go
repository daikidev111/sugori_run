package usecase

import "22dojo-online/pkg/domain"

type UserInteractor struct {
	UserRepository UserRepository
}

func (interactor *UserInteractor) SelectUserByUserID(userID string) (user domain.User, err error) {
	user, err = interactor.UserRepository.SelectUserByUserID(userID)
	return
}

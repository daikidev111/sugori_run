package usecase

import "22dojo-online/pkg/domain"

type UserRepository interface {
	SelectUserByUserID(string) (domain.User, error)
}

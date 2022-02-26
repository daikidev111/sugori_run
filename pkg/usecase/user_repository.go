package usecase

import "22dojo-online/pkg/domain"

type UserRepository interface {
	SelectUserByPrimaryKey(string) (domain.User, error)
	SelectUserByAuthToken(string) (domain.User, error)
}

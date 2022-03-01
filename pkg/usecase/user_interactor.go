package usecase

import (
	"fmt"

	"22dojo-online/pkg/domain/entity"
	"22dojo-online/pkg/domain/service"
	"22dojo-online/pkg/driver/mysql/database"
)

// type UserInteractor struct {
// 	UserRepository repository.UserRepository
// }

type UserInteractor struct {
	UserService service.UserServiceInterface
}

func NewUserInteractor(sqlHandler database.SQLHandler) *UserInteractor {
	return &UserInteractor{
		UserService: &service.UserService{
			UserRepository: &database.UserRepository{
				SQLHandler: sqlHandler,
			},
		},
	}
}

//TODO: change the func name to get user
func (interactor *UserInteractor) SelectUserByPrimaryKey(userID string) (*entity.User, error) {
	user, err := interactor.SelectUserByPrimaryKey(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query a user by primary key. err = %w", err)
	}
	return user, nil
}

func (interactor *UserInteractor) SelectUserByAuthToken(authToken string) (*entity.User, error) {
	user, err := interactor.SelectUserByAuthToken(authToken)
	if err != nil {
		return nil, fmt.Errorf("failed to query a user by auth token. err = %w", err)
	}
	return user, nil
}

func (interactor *UserInteractor) InsertUser(user *entity.User) error {
	if err := interactor.InsertUser(user); err != nil {
		return fmt.Errorf("failed to insert user. err = %w", err)
	}
	return nil
}

func (interactor *UserInteractor) UpdateUserByPrimaryKey(user *entity.User) error {
	if err := interactor.UpdateUserByPrimaryKey(user); err != nil {
		return fmt.Errorf("failed to update user. err = %w", err)
	}
	return nil
}

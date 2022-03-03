package mock

import (
	"22dojo-online/pkg/domain/entity"
	"22dojo-online/pkg/domain/repository"
	"22dojo-online/pkg/domain/service"
	"fmt"
	"testing"
)

type userServiceFaker struct {
	UserRepository repository.UserRepository
}

func (usf *userServiceFaker) SelectUserByPrimaryKeyF(userID string) (*entity.User, error) {
	user, err := usf.UserRepository.SelectUserByPrimaryKey(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query a user by primary key. err = %w", err)
	}
	return user, nil
}

// test
func TestSelectUserByPrimaryKeyService(t *testing.T) {
	table := []struct {
		testName     string
		id           string
		nameToChange string
		user         entity.User
		err          error
	}{
		{
			"FIRST TEST CASE: InsertUser from pkg/interfaces/database/user_repository.go",
			"78164dcf-6b7c-45e4-862a-2a0f6735a449",
			"whattt",
			entity.User{
				ID:        "78164dcf-6b7c-45e4-862a-2a0f6735a449",
				AuthToken: "b187b9e0-08e6-42dd-a9b3-a900b137983c",
				Name:      "whatt",
				HighScore: 100,
				Coin:      10000,
			},
			nil,
		},
	}
	mockService := new(userServiceFaker)
	serviceF := service.NewUserService(mockService.UserRepository)
	t.Logf("service %s", b.ID)

	// for _, tt := range table {
	// 	t.Run(tt.testName, func(t *testing.T) {
	// 		b, _ := serviceF.SelectUserByPrimaryKey(tt.id)
	// 		if b == nil {
	// 			return
	// 		}
	// 		// assert.Equal(t, tt.err, err)
	// 		// assert.Equal(t, tt.user, b)
	// 	})
	// }
}

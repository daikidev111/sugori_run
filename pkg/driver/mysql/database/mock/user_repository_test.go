package mock

import (
	"database/sql"
	"regexp"
	"testing"

	"22dojo-online/pkg/domain/entity"
	driver "22dojo-online/pkg/driver/mysql"
	"22dojo-online/pkg/driver/mysql/database"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func DummySQLHandler(db *sql.DB) database.SQLHandler {
	SQLHandler := new(driver.SQLHandlerImpl)
	SQLHandler.Conn = db
	return SQLHandler
}

// go test -v pkg/driver/mysql/database/mock/user_repository_test.go
func TestSelectUserByPrimaryKey(t *testing.T) {
	table := []struct {
		testName string
		id       string
		user     entity.User
		err      error
	}{
		{
			"FIRST TEST CASE: SelectUserByPrimaryKey from pkg/interfaces/database/user_repository.go",
			"78164dcf-6b7c-45e4-862a-2a0f6735a449",
			entity.User{
				ID:        "78164dcf-6b7c-45e4-862a-2a0f6735a449",
				AuthToken: "b187b9e0-08e6-42dd-a9b3-a900b137983c",
				Name:      "whatt",
				HighScore: 100,
				Coin:      47800,
			},
			nil,
		},
		{ // second test case
			"SECOND TEST CASE: SelectUserByPrimaryKey from pkg/interfaces/database/user_repository.go",
			"78164dcf-6b7c-45e4-862a-2a0f6735a449",
			entity.User{
				ID:        "78164dcf-6b7c-45e4-862a-2a0f6735a449",
				AuthToken: "b187b9e0-08e6-42dd-a9b3-a900b137983c",
				Name:      "whatt",
				HighScore: 100,
				Coin:      -1000,
			},
			nil,
		},
	}

	/*   prepare   */
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error("sqlmock not work")
	}
	defer db.Close()

	repo := &database.UserRepository{
		SQLHandler: DummySQLHandler(db),
	}

	query := "SELECT `id`, `auth_token`, `name`, `high_score`, `coin` FROM `user` WHERE `id`= ?"

	for _, tt := range table {
		t.Run(tt.testName, func(t *testing.T) {
			b := &tt.user
			rows := sqlmock.NewRows([]string{
				"id", "auth_token", "name", "high_score", "coin",
			}).AddRow(b.ID, b.AuthToken, b.Name, b.HighScore, b.Coin)
			mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(tt.id).WillReturnRows(rows)
			// ExpectQuery expects Query() or QueryRow() to be called with expectedSQL query.
			// the *ExpectedQuery allows to mock database response.

			got, err := repo.SelectUserByPrimaryKey(b.ID)

			assert.Equal(t, tt.err, err)
			assert.Equal(t, b, got)
		})
	}
}

// func TestSelectUserByAuthToken(t *testing.T) {
// 	// table for test
// 	table := []struct {
// 		testName string
// 		id       string
// 		user     entity.User
// 		err      error
// 	}{
// 		{
// 			"FIRST TEST CASE: SelectUserByPrimaryKey from pkg/interfaces/database/user_repository.go",
// 			"78164dcf-6b7c-45e4-862a-2a0f6735a449",
// 			entity.User{
// 				ID:        "78164dcf-6b7c-45e4-862a-2a0f6735a449",
// 				AuthToken: "b187b9e0-08e6-42dd-a9b3-a900b137983c",
// 				Name:      "whatt",
// 				HighScore: 100,
// 				Coin:      47800,
// 			},
// 			nil,
// 		},
// 	}

// 	/*   prepare   */
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Error("sqlmock not work")
// 	}
// 	defer db.Close()

// 	repo := &database.UserRepository{
// 		SQLHandler: DummySQLHandler(db),
// 	}

// 	query := "SELECT `id`, `auth_token`, `name`, `high_score`, `coin` FROM `user` WHERE `auth_token`= ?"

// 	for _, tt := range table {
// 		t.Run(tt.testName, func(t *testing.T) {
// 			b := &tt.user
// 			rows := sqlmock.NewRows([]string{
// 				"id", "auth_token", "name", "high_score", "coin",
// 			}).AddRow(b.ID, b.AuthToken, b.Name, b.HighScore, b.Coin)
// 			mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(tt.id).WillReturnRows(rows)
// 			// ExpectQuery expects Query() or QueryRow() to be called with expectedSQL query.
// 			// the *ExpectedQuery allows to mock database response.

// 			got, err := repo.SelectUserByAuthToken(b.AuthToken)

// 			assert.Equal(t, tt.err, err)
// 			assert.Equal(t, b, got)
// 		})
// 	}
// }

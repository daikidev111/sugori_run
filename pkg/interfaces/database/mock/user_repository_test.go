package mock

//nolint: gofmt, goimports// this is why
import (
	"22dojo-online/pkg/domain"
	"22dojo-online/pkg/infrastructure"
	"22dojo-online/pkg/interfaces/database"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func DummySQLHandler(db *sql.DB) database.SQLHandler {
	SQLHandler := new(infrastructure.SQLHandler)
	SQLHandler.Conn = db
	return SQLHandler
}

func TestSelectUserByPrimaryKey(t *testing.T) {
	// table for test
	table := []struct {
		testName string
		id       string
		user     domain.User
		err      error
	}{
		{
			"Testing SelectUserByPrimaryKey from pkg/interfaces/database/user_repository.go",
			"78164dcf-6b7c-45e4-862a-2a0f6735a449",
			domain.User{
				ID:        "78164dcf-6b7c-45e4-862a-2a0f6735a449",
				AuthToken: "b187b9e0-08e6-42dd-a9b3-a900b137983c",
				Name:      "whatt",
				HighScore: 100,
				Coin:      47800,
			},
			nil,
		},
		// { // second test case
		// 	"Testing SelectUserByPrimaryKey from pkg/interfaces/database/user_repository.go",
		// 	"78164dcf-6b7c-45e4-862a-2a0f6735a449",
		// 	domain.User{
		// 		ID:        "78164dcf-6b7c-45e4-862a-2a0f6735a449",
		// 		AuthToken: "b187b9e0-08e6-42dd-a9b3-a900b137983c",
		// 		Name:      "whatt",
		// 		HighScore: 100,
		// 		Coin:      -1000,
		// 	},
		// 	nil,
		// },
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

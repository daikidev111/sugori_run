package unit

// import (
// 	"22dojo-online/pkg/domain/entity"
// 	driver "22dojo-online/pkg/driver/mysql"
// 	"22dojo-online/pkg/driver/mysql/database"
// 	"database/sql"
// 	"fmt"
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/google/go-cmp/cmp"
// )

// func DummySQLHandler(db *sql.DB) database.SQLHandler {
// 	SQLHandler := new(driver.SQLHandlerImpl)
// 	SQLHandler.Conn = db
// 	return SQLHandler
// }

// func TestSelectUserByPrimaryKey(t *testing.T) {
// 	// table for test
// 	t.Parallel()

// 	type args struct {
// 		userID  string
// 		want    entity.User
// 		wantErr bool
// 	}

// 	tests := map[string]args{
// 		"success to calculate": {
// 			userID: "78164dcf-6b7c-45e4-862a-2a0f6735a449",
// 			want: entity.User{
// 				ID:        "78164dcf-6b7c-45e4-862a-2a0f6735a449",
// 				AuthToken: "b187b9e0-08e6-42dd-a9b3-a900b137983c",
// 				Name:      "whatt",
// 				HighScore: 100,
// 				Coin:      47800,
// 			},
// 			wantErr: false,
// 		},
// 	}

// 	db, _, err := sqlmock.New()
// 	if err != nil {
// 		t.Error("sqlmock not work")
// 	}
// 	defer db.Close()

// 	repo := &database.UserRepository{
// 		SQLHandler: DummySQLHandler(db),
// 	}

// 	for n, tt := range tests {
// 		tt := tt
// 		n := n
// 		t.Run(n, func(t *testing.T) {
// 			t.Parallel() // 逐一になってしまうから並行処理
// 			got, err := repo.SelectUserByPrimaryKey(tt.userID)
// 			if (err != nil) != tt.wantErr {
// 				fmt.Errorf("SelectUserByPrimaryKey(userID) error = %s", err.Error())
// 				return
// 			}
// 			if diff := cmp.Diff(tt.want, got); diff != "" { // TestCaseを
// 				fmt.Errorf("SelectUserByPrimaryKey(userID) mistmatch (-want +got): \n%s", diff)
// 			}
// 		})

// 	}
// }

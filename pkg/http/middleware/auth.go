package middleware

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"22dojo-online/pkg/db"
	"22dojo-online/pkg/dcontext"
	"22dojo-online/pkg/domain"
	"22dojo-online/pkg/http/response"
)

type Auth struct {
	db *sql.DB
}

func NewAuth(db *sql.DB) *Auth {
	return &Auth{
		db: db,
	}
}

// ここがSQLに依存しているのでここをどう対応するか考えないといけない
func (auth *Auth) Authenticate(nextFunc http.HandlerFunc) http.HandlerFunc {

	// Authenticate ユーザ認証を行ってContextへユーザID情報を保存する
	// func Authenticate(nextFunc http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		// リクエストヘッダからx-token(認証トークン)を取得
		token := request.Header.Get("x-token")
		if token == "" {
			log.Println("x-token is empty")
			return
		}

		row := db.Conn.QueryRow("SELECT `id`, `auth_token`, `name`, `high_score`, `coin` FROM `user` WHERE `auth_token`=?", token)

		user := domain.User{}
		err := row.Scan(&user.ID, &user.AuthToken, &user.Name, &user.HighScore, &user.Coin)
		if err != nil {
			if err == sql.ErrNoRows {
				return
			}
			log.Println(err)
			return
		}

		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Invalid token")
			return
		}
		// if user == nil {
		// 	log.Printf("user not found. token=%s", token)
		// 	response.BadRequest(writer, "Invalid token")
		// 	return
		// }

		// ユーザIDをContextへ保存して以降の処理に利用する
		ctx = dcontext.SetUserID(ctx, user.ID)

		// 次の処理
		nextFunc(writer, request.WithContext(ctx))
	}
}

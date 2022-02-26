package middleware

import (
	"22dojo-online/pkg/dcontext"
	"22dojo-online/pkg/http/response"
	"22dojo-online/pkg/usecase"
	"context"
	"log"
	"net/http"
)

type Auth struct {
	UserInteractor *usecase.UserInteractor
}

func NewAuth(UserInteractor *usecase.UserInteractor) *Auth {
	return &Auth{
		UserInteractor: UserInteractor,
	}
}

func (auth *Auth) Authenticate(nextFunc http.HandlerFunc) http.HandlerFunc {
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

		user, err := auth.UserInteractor.SelectUserByAuthToken(token)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
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

package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"22dojo-online/pkg/db"
	"22dojo-online/pkg/dcontext"
	"22dojo-online/pkg/http/response"
	"22dojo-online/pkg/server/model"
)

type gameFinishRequest struct {
	Score int32 `json:"score"`
}

type gameFinishResponse struct {
	Coin int32 `json:"coin"`
}

// HandleGameFinishPost GameFinish時のリクエストとレスポンスの処理
func HandleGameFinishPost() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Requestのスコアのdecode
		var requestBody gameFinishRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
			response.BadRequest(writer, "Bad Request")
			return
		}

		// scoreが1以下の場合の考慮
		if requestBody.Score < 1 { // 0でも処理をしなくて良い
			log.Println("Negative score is invalid")
			response.BadRequest(writer, "Negative score is invalid")
			return
		}

		// int32へキャスト
		score := requestBody.Score

		// ユーザー認証(middleware)からのユーザーIDの取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println("userID is empty")
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// pkg/db/conn.goのトランザクション用のwrapper関数の呼び出し
		err := db.Transact(ctx, db.Conn, func(tx *sql.Tx) error {
			user, err := model.SelectUserByPrimaryKeyWithLock(userID, tx)
			if err != nil {
				log.Println(err)
				return err
			}

			// if score > user.HighScore {
			// 	err = model.UpdateScoreByPrimaryKeyWithLock(userID, tx, score)
			// 	if err != nil {
			// 		log.Println(err)
			// 		return err
			// 	}
			// }

			coin := user.Coin + score // コインの計算方法
			// err = model.UpdateCoinByPrimaryKeyWithLock(userID, tx, coin)
			// if err != nil {
			// 	log.Println(err)
			// 	return err
			// }

			err = model.UpdateCoinAndScoreByPrimaryKeyWithLock(userID, tx, score, coin, user.HighScore)
			if err != nil {
				log.Println(err)
				return err
			}

			return nil
		})
		if err != nil { // トランザクションが失敗した場合
			log.Println("DB Transaction failed: ")
			log.Println(err)
			response.BadRequest(writer, "Internal Server Error")
			return
		}

		// responseとしてコインの値を返す
		response.Success(writer, &gameFinishResponse{
			Coin: score, // 取得コインはスコアと同様にしている
		})
	}
}

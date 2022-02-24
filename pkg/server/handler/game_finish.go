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
	Score int `json:"score"`
}

type gameFinishResponse struct {
	Coin int32 `json:"coin"`
}

// HandleGameFinshPost GameFinish時のリクエストとレスポンスの処理
func HandleGameFinshPost() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var requestBody gameFinishRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
			response.BadRequest(writer, "Bad Request")
			return
		}

		if requestBody.Score < 0 {
			log.Println("Negative score is invalid")
			response.BadRequest(writer, "Negative score is invalid")
		}

		score := int32(requestBody.Score)

		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println("userID is empty")
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		err := db.Transact(db.Conn, func(tx *sql.Tx) error {
			user, err := model.SelectUserByPrimaryKeyWithLock(userID, tx)
			if err != nil {
				log.Println(err)
				return err
			}

			if score > user.HighScore {
				err = model.UpdateScoreByPrimaryKeyWithLock(userID, tx, score)
				if err != nil {
					log.Println(err)
					return err
				}
			}

			coin := user.Coin + score
			err = model.UpdateCoinByPrimaryKeyWithLock(userID, tx, coin)
			if err != nil {
				log.Println(err)
				return err
			}

			return nil
		})
		if err != nil {
			log.Println("DB Transaction failed")
			response.BadRequest(writer, "Internal Server Error")
			return
		}

		response.Success(writer, &gameFinishResponse{
			Coin: score,
		})
	}
}
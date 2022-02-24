package handler

import (
	"encoding/json"
	"log"
	"net/http"

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

		err := model.UpdateCoinAndScoreByPrimaryKeyTx(ctx, userID, score)
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

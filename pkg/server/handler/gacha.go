package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"22dojo-online/pkg/dcontext"
	"22dojo-online/pkg/http/response"
)

type gachaResponse struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int32  `json:"rarity"`
	isNew        bool   `json:"isNew"`
}

// type gachaListResponse struct {
// 	Result []*gachaResponse `json:"gachaList"`
// }

type gachaListResponse struct {
	Result int `json:"gachaList"`
}

type gachaRequest struct {
	Times int `json:"times"`
}

// HandleCollectionGet
func HandleGachaPost() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var requestBody gachaRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
			response.BadRequest(writer, "Bad Request")
			return
		}

		log.Println(requestBody.Times)

		// Contextから認証済みのユーザIDを取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println("userID is empty")
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		response.Success(writer, &gachaListResponse{
			Result: 1,
		})
	}
}

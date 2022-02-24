package handler

import (
	"log"
	"net/http"
	"strconv"

	"22dojo-online/pkg/http/response"
	"22dojo-online/pkg/server/model"
)

type rankingGetResponse struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	Rank     int32  `json:"rank"`
	Score    int32  `json:"score"`
}

type collectionRankingResponse struct {
	Ranks []*rankingGetResponse `json:"ranks"`
}

// HandleRankingGet ランキング更新
func HandleRankingGet() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		key := request.URL.Query().Get("start")

		startKey, err := strconv.Atoi(key)
		if err != nil {
			log.Println("Failed to convert the start key to int data type: Check Atoi in line 32")
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		if startKey < 1 {
			log.Println("start key cannot be less than 1")
			response.BadRequest(writer, "Bad Request")
			return
		}

		userRankings, err := model.SelectUsersFromRankingStart(startKey)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		userRankingsArr := make([]*rankingGetResponse, 0, len(userRankings))

		starKeyCounter, err := strconv.ParseInt(key, 10, 32)
		if err != nil {
			log.Println("Failed to convert the start key to int64 data type: Check Atoi in line 48")
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		startKeyCounter := int32(starKeyCounter)

		for _, user := range userRankings {
			userRankingsArr = append(userRankingsArr, &rankingGetResponse{
				UserID:   user.UserID,
				UserName: user.UserName,
				Score:    user.HighScore,
				Rank:     startKeyCounter,
			})

			startKeyCounter++
		}

		response.Success(writer, &collectionRankingResponse{
			Ranks: userRankingsArr,
		})
	}
}

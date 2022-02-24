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
		keys, ok := request.URL.Query()["start"]
		if !ok || len(keys[0]) < 1 { // [start] < 1 is invalid condition
			log.Println("URL param 'start' is missing")
			return
		}

		startKey, err := strconv.Atoi(keys[0])
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

		startKeyCounter := startKey
		for _, user := range userRankings {
			r := rankingGetResponse{
				UserID:   user.UserID,
				UserName: user.UserName,
				Score:    user.HighScore,
				Rank:     int32(startKeyCounter),
			}
			userRankingsArr = append(userRankingsArr, &r)

			startKeyCounter++
		}

		response.Success(writer, &collectionRankingResponse{
			Ranks: userRankingsArr,
		})
	}
}

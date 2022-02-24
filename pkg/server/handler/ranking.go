package handler

import (
	"22dojo-online/pkg/http/response"
	"22dojo-online/pkg/server/model"
	"log"
	"net/http"
	"strconv"
)

// {"ranks":[{"userId":"731269d9-b7b2-4931-b2e1-ad600226d5f1","userName":"ばったー","rank":31,"score":1477}]
type RankingGetResponse struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	Rank     int32  `json:"rank"`
	Score    int32  `json:"score"`
}

type CollectionRankingResponse struct {
	Ranks []RankingGetResponse `json:"ranks"`
}

// HandleRankingGet ランキング更新
func HandleRankingGet() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		keys, ok := request.URL.Query()["start"]
		if !ok || len(keys[0]) < 1 { // [start] < 1 is invalid condition
			log.Println("URL param 'start' is missing")
			return
		}

		startKey, _ := strconv.Atoi(keys[0]) // TODO: error handling

		userRankings, err := model.SelectUsersFromRankingStart(startKey)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		if userRankings == nil {
			log.Println("userRankings is empty")
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		userRankingsArr := make([]RankingGetResponse, 0, len(userRankings))

		startKeyCounter := startKey
		for i := range userRankings {
			r := RankingGetResponse{}

			r.UserID = userRankings[i].UserID
			r.UserName = userRankings[i].UserName
			r.Score = userRankings[i].HighScore
			r.Rank = int32(startKeyCounter)

			userRankingsArr = append(userRankingsArr, r)
			startKeyCounter++
		}

		response.Success(writer, &CollectionRankingResponse{
			Ranks: userRankingsArr,
		})

	}
}

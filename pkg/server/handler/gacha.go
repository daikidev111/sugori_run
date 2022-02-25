package handler

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"22dojo-online/pkg/constant"
	"22dojo-online/pkg/dcontext"
	"22dojo-online/pkg/http/response"
	"22dojo-online/pkg/server/model"
)

type gachaResponse struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int32  `json:"rarity"`
	IsNew        bool   `json:"isNew"`
}

type gachaListResponse struct {
	Result []*gachaResponse `json:"gachaList"`
}

// type gachaListResponse struct {
// 	Result int `json:"gachaList"`
// }

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

		if requestBody.Times < 1 {
			log.Println("Times cannot be less than 1")
			response.BadRequest(writer, "Bad Request")
			return
		}

		// Contextから認証済みのユーザIDを取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println("userID is empty")
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// トランザクションをここから開始する

		user, err := model.SelectUserByPrimaryKey(userID)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		times := int32(requestBody.Times)
		if user.Coin < constant.GachaCoinConsumption*times {
			log.Println("The amount of coins the user has is less than the consumption of coin for the gacha draw")
			response.BadRequest(writer, "Invalid gacha draw")
			return
		}

		collectionItems, err := model.SelectAllCollectionItems()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		if len(collectionItems) == 0 {
			log.Println("A collection of items is not found")
			response.BadRequest(writer, "A collection of items is not found.")
			return
		}

		userCollectionItems, err := model.SelectUserCollectionItemByUserID(userID)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		gachaProbabilities, err := model.SelectAllCollectionItemProbability()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		var SumOfRatio int
		for _, gachaProb := range gachaProbabilities {
			SumOfRatio += int(gachaProb.Ratio) //TODO:  cast を治す -> 構造体のratioをイントにしても良いかも？？
		}

		// ここからループの実装？？ times 範囲内での
		// 乱数の取得
		randInt := rand.Intn(SumOfRatio)

		var targetRatio int
		var targetCollectionID string

		for _, gachaProb := range gachaProbabilities {
			targetRatio += int(gachaProb.Ratio)
			// ガチャで排出確率に基づいたコレクションアイテムの取得
			if targetRatio > randInt {
				targetCollectionID = gachaProb.CollectionID
				break
			}
		}

		/*
			すでに所持しているかをIDを突き合わせて判定
			（新しく獲得したアイテムはisNewがtrue,既に持っているアイテムはisNewがfalseとなります。）重複はなし
		*/

		userCollectionItemsMap := make(map[string]bool, len(userCollectionItems))
		for i := range userCollectionItems {
			userCollectionItemsMap[userCollectionItems[i].CollectionID] = true
		}

		itemCollectionMap := make(map[string]*model.CollectionItem, len(collectionItems))
		for i, collectionItem := range collectionItems {
			itemCollectionMap[collectionItem.ID] = collectionItems[i]
		}

		// user collectionの中で獲得されたアイテムの存在確認
		if !userCollectionItemsMap[targetCollectionID] { // もしuser collectionのアイテムの重複していない場合はcollectionitemuserに格納する bulk insertでの実装？？
			err = model.InsertUserCollectionItemByUserID(userID, targetCollectionID)
			if err != nil {
				log.Println("Failed to insert the new item into the user's collection", err)
				response.InternalServerError(writer, "Internal Server Error")
				return
			}
		}

		// 共通処理は、responseに格納
		gachaCollectionList := make([]*gachaResponse, 0, times)
		gachaCollectionList = append(gachaCollectionList, &gachaResponse{
			CollectionID: targetCollectionID,
			Name:         itemCollectionMap[targetCollectionID].Name,
			Rarity:       itemCollectionMap[targetCollectionID].Rarity,
			IsNew:        !userCollectionItemsMap[targetCollectionID], // 新しく獲得したアイテムはisNewがtrue,既に持っているアイテムはisNewがfalse
		})

		//コイン消費（コインをマイナスにしてアップデート処理）
		user.Coin -= constant.GachaCoinConsumption
		err = model.UpdateCoinByPrimaryKey(userID, user.Coin)
		if err != nil {
			log.Println("Failed to update the user's coin", err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		response.Success(writer, &gachaListResponse{
			Result: gachaCollectionList,
		})
	}
}

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

		// トランザクションをここから開始する

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

		UserCollcetionItemsArr := make([]*model.UserCollectionItem, 0, times)                // times -1かどうかの確認　Gachaで取得された新アイテムの格納
		gachaCollectionList := make([]*gachaResponse, 0, times)                              // Response用のスライス
		itemCollectionMap := make(map[string]*model.CollectionItem, len(collectionItems))    // itemCollectionのマップ
		userCollectionItemsMap := make(map[string]bool, len(userCollectionItems)+int(times)) // user_collection_itemsのマップ

		// すでに所持しているかをIDを突き合わせて判定するためのマップ(動的）
		for _, userCollectionItem := range userCollectionItems {
			userCollectionItemsMap[userCollectionItem.CollectionID] = true
		}

		// IDのキーと格アイテムの情報が入っている静的マップの生成
		for i, collectionItem := range collectionItems {
			itemCollectionMap[collectionItem.ID] = collectionItems[i]
		}

		var SumOfRatio int
		for _, gachaProb := range gachaProbabilities {
			SumOfRatio += int(gachaProb.Ratio) //TODO:  cast を治す -> 構造体のratioをイントにしても良いかも？？
		}

		// =========== start of the loop ===================
		for i := 0; i < requestBody.Times; i++ {

			randInt := rand.Intn(SumOfRatio) // 乱数の取得

			var targetRatio int           // Ratioから取得される値を足す際に必要
			var targetCollectionID string // 乱数を越した際のcollection ID(ガチャの引きアイテム)

			// ガチャで排出確率に基づいたコレクションアイテムの取得
			for _, gachaProb := range gachaProbabilities {
				targetRatio += int(gachaProb.Ratio) // TODO: cast直し
				if targetRatio > randInt {
					targetCollectionID = gachaProb.CollectionID
					break
				}
			}

			// 共通処理: responseに格納
			gachaCollectionList = append(gachaCollectionList, &gachaResponse{
				CollectionID: targetCollectionID,
				Name:         itemCollectionMap[targetCollectionID].Name,
				Rarity:       itemCollectionMap[targetCollectionID].Rarity,
				IsNew:        !userCollectionItemsMap[targetCollectionID], // 新しく獲得したアイテムはisNewがtrue,既に持っているアイテムはisNewがfalse
			})

			// 新アイテムはUserCollcetionItemsArrに格納(bulk insert時に必要となる)
			if !userCollectionItemsMap[targetCollectionID] {
				UserCollcetionItemsArr = append(UserCollcetionItemsArr, &model.UserCollectionItem{ // TODO: fix typo UserCollcetionItemsArr
					UserID:       userID,
					CollectionID: targetCollectionID,
				})
				userCollectionItemsMap[targetCollectionID] = true // 格納後はtrueに変換する
			}
		}
		// ========== end of the loop ==============

		// bulk insertの開始
		if len(UserCollcetionItemsArr) > 0 {
			err = model.InsertUserCollectionItemsByUserID(UserCollcetionItemsArr)
			if err != nil {
				log.Println("Failed to insert the new item into the user's collection", err)
				response.InternalServerError(writer, "Internal Server Error")
				return
			}
		}

		//コイン消費（コインをマイナスにしてアップデート処理）
		user.Coin -= constant.GachaCoinConsumption * times
		err = model.UpdateCoinByPrimaryKey(userID, user.Coin)
		if err != nil {
			log.Println("Failed to update the user's coin", err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// responseを返す
		response.Success(writer, &gachaListResponse{
			Result: gachaCollectionList,
		})
	}
}

// transaction -> select user by primary key, select all user collection items, update coin, inseruser collections,

package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"22dojo-online/pkg/constant"
	"22dojo-online/pkg/db"
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
	Result []*gachaResponse `json:"results"`
}

type gachaRequest struct {
	Times int `json:"times"`
}

//nolint: gocyclo // this is why
func HandleGachaPost() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var requestBody gachaRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
			response.BadRequest(writer, "Bad Request")
			return
		}

		if requestBody.Times < 1 || requestBody.Times > constant.GachaUpperLimit {
			log.Println("Times cannot be less than 1 or greater than 10")
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

		gachaProbabilities, err := model.SelectAllCollectionItemProbability()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		collectionItems, err := model.SelectAllCollectionItems()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		if len(collectionItems) == 0 {
			log.Println("A collection of items is empty")
			response.BadRequest(writer, "A collection of items is not found.")
			return
		}

		gachaCollectionList := make([]*gachaResponse, 0, requestBody.Times)               // Response用のスライス
		itemCollectionMap := make(map[string]*model.CollectionItem, len(collectionItems)) // itemCollectionのマップ

		// IDのキーと各アイテムの情報が入っている静的マップの生成
		for i, collectionItem := range collectionItems {
			itemCollectionMap[collectionItem.ID] = collectionItems[i]
		}

		// トランザクション開始
		err = db.Transact(ctx, db.Conn, func(tx *sql.Tx) error {
			user, err := model.SelectUserByPrimaryKeyWithLock(userID, tx)
			if err != nil {
				log.Println("SelectUserByPrimaryKeyWithLock: Failed to fetch a row of the user's information")
				response.InternalServerError(writer, "Internal Server Error")
				return err
			}

			times := int32(requestBody.Times)
			totalCoin := constant.GachaCoinConsumption * times
			if user.Coin < totalCoin {
				log.Println("Not enough coins")
				response.BadRequest(writer, "Invalid gacha draw")
				return err
			}

			userCollectionItems, err := model.SelectUserCollectionItemByUserID(userID)
			if err != nil {
				log.Println("SelectUserCollectionItemByUserIDWithLock: Failed to fetch rows from the user collection item")
				log.Println(err)
				response.InternalServerError(writer, "Internal Server Error")
				return err
			}

			userCollectionItemsArr := make([]*model.UserCollectionItem, 0, times)                // Gachaで取得された新アイテムの格納
			userCollectionItemsMap := make(map[string]bool, len(userCollectionItems)+int(times)) // user_collection_itemsのマップ

			// すでに所持しているかをIDを突き合わせて判定するためのマップ(動的）
			for _, userCollectionItem := range userCollectionItems {
				userCollectionItemsMap[userCollectionItem.CollectionID] = true
			}

			var sumOfRatio int
			for _, gachaProb := range gachaProbabilities {
				sumOfRatio += gachaProb.Ratio
			}

			// =========== start of the loop ===================
			for i := 0; i < int(times); i++ {
				//nolint: gosec // this is why
				randInt := rand.Intn(sumOfRatio) // 乱数の取得

				var targetRatio int                            // Ratioから取得される値を足す際に必要
				var targetCollectionID string                  // 乱数を越した際のcollection ID(ガチャの引きアイテム)
				for _, gachaProb := range gachaProbabilities { // ガチャで排出確率に基づいたコレクションアイテムの取得
					targetRatio += gachaProb.Ratio
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
					userCollectionItemsArr = append(userCollectionItemsArr, &model.UserCollectionItem{
						UserID:       userID,
						CollectionID: targetCollectionID,
					})
					userCollectionItemsMap[targetCollectionID] = true // 格納後はtrueに変換する
				}
			}
			// ========== end of the loop ==============

			// コイン消費（コインをマイナスにしてアップデート処理）
			user.Coin -= totalCoin
			if err := model.UpdateCoinAndScoreByPrimaryKeyWithLock(tx, userID, user.HighScore, user.Coin); err != nil {
				log.Println("UpdateCoinAndScoreByPrimaryKeyWithLock: Failed to update the user's coin", err)
				response.InternalServerError(writer, "Internal Server Error")
				return err
			}

			// bulk insertの開始
			if len(userCollectionItemsArr) > 0 {
				err := model.InsertUserCollectionItemsByUserIDWithLock(tx, userCollectionItemsArr)
				if err != nil {
					log.Println("InsertUserCollectionItemsByUserIDWithLock: Failed to insert the new item(s) into the user's collection", err)
					response.InternalServerError(writer, "Internal Server Error")
					return err
				}
			}
			return nil
		})
		if err != nil { // トランザクションが失敗した場合
			log.Println("DB Transaction failed or a query with the lock failed to execute. Rollback the transaction.")
			log.Println(err)
			response.BadRequest(writer, "Internal Server Error")
			return
		}

		// responseを返す
		response.Success(writer, &gachaListResponse{
			Result: gachaCollectionList,
		})
	}
}

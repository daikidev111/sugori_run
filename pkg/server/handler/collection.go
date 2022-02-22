package handler

import (
	"log"
	"net/http"

	"22dojo-online/pkg/dcontext"
	"22dojo-online/pkg/http/response"
	"22dojo-online/pkg/server/model"
)

type collectionGetResponse struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int32  `json:"rarity"`
	HasItem      bool   `json:"hasItem"`
}

type collectionListResponse struct {
	Collections []collectionGetResponse `json:"collections"`
}

// HandleSettingGet ゲーム設定情報取得処理
func HandleCollectionGet() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Contextから認証済みのユーザIDを取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println("userID is empty")
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		var collectionItems []*model.CollectionItem
		var err error
		collectionItems, err = model.SelectAllCollectionItems()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		if collectionItems == nil {
			log.Println("A collection of items is not found")
			response.BadRequest(writer, "A collection of items is not found.")
			return
		}

		var userCollectionItems []*model.UserCollectionItem
		userCollectionItems, err = model.SelectUserCollectionItemByUserID(userID)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		collectionItemArr := make([]collectionGetResponse, 0, len(collectionItems))

		for i := range collectionItems {
			c := collectionGetResponse{}
			if len(userCollectionItems) > i {
				if collectionItems[i].ID == userCollectionItems[i].CollectionID {
					c.HasItem = true
					c.CollectionID = collectionItems[i].ID
					c.Name = collectionItems[i].Name
					c.Rarity = collectionItems[i].Rarity
				} else {
					c.HasItem = false
					c.CollectionID = collectionItems[i].ID
					c.Name = collectionItems[i].Name
					c.Rarity = collectionItems[i].Rarity
				}
			} else {
				c.HasItem = false
				c.CollectionID = collectionItems[i].ID
				c.Name = collectionItems[i].Name
				c.Rarity = collectionItems[i].Rarity
			}

			collectionItemArr = append(collectionItemArr, c)
		}

		response.Success(writer, &collectionListResponse{
			Collections: collectionItemArr,
		})
	}
}

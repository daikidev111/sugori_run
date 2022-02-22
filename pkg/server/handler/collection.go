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

// HandleCollectionGet
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

		userCollectionItems, err := model.SelectUserCollectionItemByUserID(userID)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		collectionItemArr := make([]collectionGetResponse, 0, len(collectionItems))

		// valueはstructではなくboolにした
		userCollectionItemsMap := make(map[string]bool, len(userCollectionItems))
		for i := range userCollectionItems {
			userCollectionItemsMap[userCollectionItems[i].CollectionID] = false
		}

		for i := range collectionItems {
			c := collectionGetResponse{}

			c.CollectionID = collectionItems[i].ID
			c.Name = collectionItems[i].Name
			c.Rarity = collectionItems[i].Rarity
			_, c.HasItem = userCollectionItemsMap[collectionItems[i].ID]

			collectionItemArr = append(collectionItemArr, c)
		}

		response.Success(writer, &collectionListResponse{
			Collections: collectionItemArr,
		})
	}
}

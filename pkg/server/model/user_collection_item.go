package model

import (
	"database/sql"
	"log"

	"22dojo-online/pkg/db"
)

type UserCollectionItem struct {
	UserID       string
	CollectionID string
}

func SelectUserCollectionItemByUserID(userID string) ([]*UserCollectionItem, error) {
	rows, err := db.Conn.Query("SELECT user_id, collection_item_id FROM user_collection_item where user_id = ? ORDER BY collection_item_id;", userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	return convertToUserCollectionItems(rows)
}

// convertToUserCollectionItemsでrowデータをUserCollectionItemデータへ変換する
func convertToUserCollectionItems(rows *sql.Rows) ([]*UserCollectionItem, error) {
	var userCollectionItems []*UserCollectionItem

	for rows.Next() {
		userCollectionItem := &UserCollectionItem{}
		err := rows.Scan(&userCollectionItem.UserID, &userCollectionItem.CollectionID)
		if err != nil {
			return nil, err
		}
		userCollectionItems = append(userCollectionItems, userCollectionItem)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userCollectionItems, nil
}

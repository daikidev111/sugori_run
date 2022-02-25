package model

import (
	"database/sql"
	"fmt"
	"strings"

	"22dojo-online/pkg/db"
)

type UserCollectionItem struct {
	UserID       string
	CollectionID string
}

func SelectUserCollectionItemByUserID(userID string) ([]*UserCollectionItem, error) {
	rows, err := db.Conn.Query("SELECT user_id, collection_item_id FROM user_collection_item where user_id = ?;", userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return convertToUserCollectionItems(rows)
}

func SelectUserCollectionItemByUserIDWithLock(tx *sql.Tx, userID string) ([]*UserCollectionItem, error) {
	rows, err := tx.Query("SELECT user_id, collection_item_id FROM user_collection_item where user_id = ? FOR UPDATE;", userID)
	if err != nil {
		return nil, err
	}

	return convertToUserCollectionItems(rows)
}

//TODO: ADD TRANSACTION !!!!!!!
func InsertUserCollectionItemsByUserIDWithLock(tx *sql.Tx, userCollectionItems []*UserCollectionItem) error {
	valueStrings := make([]string, 0, len(userCollectionItems))
	valueArgs := make([]interface{}, 0, len(userCollectionItems)*3)
	for _, userCollectionItem := range userCollectionItems {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, userCollectionItem.UserID)
		valueArgs = append(valueArgs, userCollectionItem.CollectionID)
	}
	stmt := fmt.Sprintf("INSERT INTO user_collection_item (user_id, collection_item_id) VALUES %s",
		strings.Join(valueStrings, ","))
	_, err := tx.Exec(stmt, valueArgs...) // call Statusln with a variable number of arguments
	// maybe remove statusln and directly access to the index of the arrays
	return err
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

package model

import (
	"database/sql"
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

func InsertUserCollectionItemsByUserIDWithLock(tx *sql.Tx, userCollectionItems []*UserCollectionItem) error {
	valueStrings := make([]string, 0, len(userCollectionItems))
	valueArgs := make([]interface{}, 0, len(userCollectionItems)*2)
	for _, userCollectionItem := range userCollectionItems {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, userCollectionItem.UserID, userCollectionItem.CollectionID)
	}

	var sb strings.Builder
	sb.WriteString("INSERT INTO user_collection_item (user_id, collection_item_id) VALUES")
	sb.WriteString(strings.Join(valueStrings, ","))
	stmt := sb.String()

	_, err := tx.Exec(stmt, valueArgs...)
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

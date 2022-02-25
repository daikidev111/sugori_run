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

//TODO: ADD TRANSACTION !!!!!!!
func InsertUserCollectionItemsByUserID(userCollectionItems []*UserCollectionItem) error {
	valueStrings := make([]string, 0, len(userCollectionItems))
	valueArgs := make([]interface{}, 0, len(userCollectionItems)*3)
	for _, userCollectionItem := range userCollectionItems {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, userCollectionItem.UserID)
		valueArgs = append(valueArgs, userCollectionItem.CollectionID)
	}
	stmt := fmt.Sprintf("INSERT INTO user_collection_item (user_id, collection_item_id) VALUES %s",
		strings.Join(valueStrings, ","))
	_, err := db.Conn.Exec(stmt, valueArgs...) // call Statusln with a variable number of arguments
	// maybe remove statusln and directly access to the index of the arrays
	return err

}

// func BulkInsert(unsavedRows []*gachaResponse) error {
// 	valueStrings := make([]string, 0, len(unsavedRows))
// 	valueArgs := make([]interface{}, 0, len(unsavedRows)*3)
// 	for _, post := range unsavedRows {
// 		valueStrings = append(valueStrings, "(?, ?, ?)")
// 		valueArgs = append(valueArgs, post.Column1)
// 		valueArgs = append(valueArgs, post.Column2)
// 		valueArgs = append(valueArgs, post.Column3)
// 	}
// 	stmt := fmt.Sprintf("INSERT INTO my_sample_table (column1, column2, column3) VALUES %s",
// 		strings.Join(valueStrings, ","))
// 	_, err := db.Exec(stmt, valueArgs...)
// 	return err
// }

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

package model

import (
	"database/sql"
	"log"

	"22dojo-online/pkg/db"
)

type CollectionItem struct {
	ID     string
	Name   string
	Rarity int32
}

func SelectAllCollectionItems() ([]*CollectionItem, error) {
	// Obtain all the items exisiting in the collectionItem table
	// Obtain all the items that the user, specified with user id from the user_collectionItem
	// Conduct a loop to find the one that the user has, which will be marked as True for hasItem or else will be marked as False
	rows, err := db.Conn.Query("SELECT id, name, rarity FROM collection_item ORDER BY id;")

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	return convertToCollectionItem(rows)
}

// convertToCollectionItemでrowデータをCollectionItemデータへ変換する
func convertToCollectionItem(rows *sql.Rows) ([]*CollectionItem, error) {
	var CollectionItems []*CollectionItem

	for rows.Next() {
		collectionItem := &CollectionItem{}
		err := rows.Scan(&collectionItem.ID, &collectionItem.Name, &collectionItem.Rarity)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		CollectionItems = append(CollectionItems, collectionItem)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return CollectionItems, nil
}

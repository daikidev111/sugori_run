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
	// Obtain all the items existing in the collectionItem table
	// Obtain all the items that the user, specified by user id from the user_collectionItem
	// Conduct a loop to find the one that the user has, which will be marked as True for hasItem or else will be marked as False

	rows_count, err := db.Conn.Query("SELECT id, name, rarity FROM collection_item ORDER BY id;")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	count := 0
	for rows_count.Next() {
		count += 1
	}

	defer rows_count.Close()

	rows, err := db.Conn.Query("SELECT id, name, rarity FROM collection_item ORDER BY id;")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	return convertToCollectionItem(rows, count)
}

// convertToCollectionItemでrowデータをCollectionItemデータへ変換する
func convertToCollectionItem(rows *sql.Rows, num_rows int) ([]*CollectionItem, error) {
	CollectionItems := make([]*CollectionItem, 0, num_rows)

	for rows.Next() {
		collectionItem := &CollectionItem{}
		err := rows.Scan(&collectionItem.ID, &collectionItem.Name, &collectionItem.Rarity)
		if err != nil {
			return nil, err
		}
		CollectionItems = append(CollectionItems, collectionItem)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return CollectionItems, nil
}

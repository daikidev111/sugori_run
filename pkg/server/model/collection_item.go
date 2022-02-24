package model

import (
	"database/sql"

	"22dojo-online/pkg/db"
)

type CollectionItem struct {
	ID     string
	Name   string
	Rarity int32
}

func SelectAllCollectionItems() ([]*CollectionItem, error) {
	rows, err := db.Conn.Query("SELECT id, name, rarity FROM collection_item;")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return convertToCollectionItem(rows)
}

// convertToCollectionItemでrowデータをCollectionItemデータへ変換する
func convertToCollectionItem(rows *sql.Rows) ([]*CollectionItem, error) {
	CollectionItems := make([]*CollectionItem, 0)

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

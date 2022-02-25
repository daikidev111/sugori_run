package model

import (
	"database/sql"

	"22dojo-online/pkg/db"
)

type GachaProbability struct {
	CollectionID string
	Name         string
	Rarity       int32
	Ratio        int
}

func SelectAllCollectionItemProbability() ([]*GachaProbability, error) {
	rows, err := db.Conn.Query("SELECT c.id, c.name, c.rarity, g.ratio FROM gacha_probability g JOIN collection_item c ON g.collection_item_id = c.id;")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return convertToGachaProbability(rows)
}

// convertToCollectionItemでrowデータをCollectionItemデータへ変換する
func convertToGachaProbability(rows *sql.Rows) ([]*GachaProbability, error) {
	GachaProbabilityList := make([]*GachaProbability, 0)

	for rows.Next() {
		gachaProbability := &GachaProbability{}
		err := rows.Scan(&gachaProbability.CollectionID, &gachaProbability.Name, &gachaProbability.Rarity, &gachaProbability.Ratio)
		if err != nil {
			return nil, err
		}
		GachaProbabilityList = append(GachaProbabilityList, gachaProbability)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return GachaProbabilityList, nil
}

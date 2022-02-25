package model

import (
	"database/sql"
	"log"

	"22dojo-online/pkg/constant"
	"22dojo-online/pkg/db"
)

// User userテーブルデータ
type User struct {
	ID        string
	AuthToken string
	Name      string
	HighScore int32
	Coin      int32
}

type UserRanking struct {
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
	HighScore int32  `json:"score"`
}

// InsertUser データベースをレコードを登録する
func InsertUser(record *User) error {
	// TODO: usersテーブルへのレコードの登録を行うSQLを入力する
	_, err := db.Conn.Exec(
		"INSERT INTO `user` (`id`, `auth_token`, `name`, `high_score`, `coin`) VALUES (?, ?, ?, ?, ?);",
		record.ID, record.AuthToken, record.Name, record.HighScore, record.Coin)
	return err
}

// SelectUserByAuthToken auth_tokenを条件にレコードを取得する
func SelectUserByAuthToken(authToken string) (*User, error) {
	row := db.Conn.QueryRow("SELECT `id`, `auth_token`, `name`, `high_score`, `coin` FROM `user` WHERE `auth_token`=?", authToken)
	return convertToUser(row)
}

// SelectUserByPrimaryKey 主キーを条件にレコードを取得する
func SelectUserByPrimaryKey(userID string) (*User, error) {
	row := db.Conn.QueryRow("SELECT `id`, `auth_token`, `name`, `high_score`, `coin` FROM `user` WHERE `id`=?", userID)
	return convertToUser(row)
}

// UpdateUserByPrimaryKey 主キーを条件にレコードを更新する
func UpdateUserByPrimaryKey(record *User) error {
	_, err := db.Conn.Exec(
		"UPDATE user SET name = ? WHERE id = ?",
		record.Name, record.ID)
	return err
}

// SelectUsersFromRankingStart ランキングの表示をスコア降順で同率の場合はIDの昇順
func SelectUsersFromRankingStart(start int) ([]*UserRanking, error) {
	rows, err := db.Conn.Query("SELECT `id`, `name`, `high_score` FROM `user` ORDER BY `high_score` DESC, `id` ASC LIMIT ?, ?;", start-1, constant.RankingNum)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return convertToUserForRanking(rows)
}

// SelectUserByPrimaryKeyWithLock 主キーを条件にユーザーを出力(排他制御あり)
func SelectUserByPrimaryKeyWithLock(userID string, tx *sql.Tx) (*User, error) {
	row := tx.QueryRow("SELECT `id`, `auth_token`, `name`, `high_score`, `coin` FROM `user` WHERE `id` = ? FOR UPDATE", userID)
	return convertToUser(row)
}

// UpdateCoinAndScoreByPrimaryKeyWithLock 主キーを条件にスコアとコインの更新処理(排他制御あり)
func UpdateCoinAndScoreByPrimaryKeyWithLock(tx *sql.Tx, userID string, highScore, coin int32) error {
	_, err := tx.Exec(
		"UPDATE user SET coin = ?, high_score = ? WHERE id = ?",
		coin, highScore, userID)
	if err != nil {
		return err
	}

	return nil
}

// convertToUserForRanking rowデータをUserRankingsデータへ変換する
func convertToUserForRanking(rows *sql.Rows) ([]*UserRanking, error) {
	UserRankings := make([]*UserRanking, 0)
	for rows.Next() {
		UserRanking := &UserRanking{}
		err := rows.Scan(&UserRanking.UserID, &UserRanking.UserName, &UserRanking.HighScore)
		if err != nil {
			return nil, err
		}
		UserRankings = append(UserRankings, UserRanking)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return UserRankings, nil
}

// convertToUser rowデータをUserデータへ変換する
func convertToUser(row *sql.Row) (*User, error) {
	user := User{}
	err := row.Scan(&user.ID, &user.AuthToken, &user.Name, &user.HighScore, &user.Coin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println(err)
		return nil, err
	}
	return &user, nil
}

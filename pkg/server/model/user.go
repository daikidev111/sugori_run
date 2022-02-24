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
	// TODO: auth_tokenを条件にSELECTを行うSQLを第1引数に入力する
	row := db.Conn.QueryRow("SELECT `id`, `auth_token`, `name`, `high_score`, `coin` FROM `user` WHERE `auth_token`=?", authToken)
	return convertToUser(row)
}

// SelectUserByPrimaryKey 主キーを条件にレコードを取得する
func SelectUserByPrimaryKey(userID string) (*User, error) {
	// TODO: idを条件にSELECTを行うSQLを第1引数に入力する
	row := db.Conn.QueryRow("SELECT `id`, `auth_token`, `name`, `high_score`, `coin` FROM `user` WHERE `id`=?", userID)
	return convertToUser(row)
}

// UpdateUserByPrimaryKey 主キーを条件にレコードを更新する
func UpdateUserByPrimaryKey(record *User) error {
	// TODO: idを条件に指定した値でnameカラムの値を更新するSQLを入力する
	_, err := db.Conn.Exec(
		"UPDATE user SET name = ? WHERE id = ?",
		record.Name, record.ID)
	return err
}

func SelectUsersFromRankingStart(start int) ([]*UserRanking, error) {
	rows, err := db.Conn.Query("SELECT `id`, `name`, `high_score` FROM `user` ORDER BY `high_score` DESC, `id` ASC LIMIT ?, ?;", start, constant.RankingNum)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return convertToUserForRanking(rows)
}

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

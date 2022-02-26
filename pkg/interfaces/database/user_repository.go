package database

import (
	"22dojo-online/pkg/domain"
)

type UserRepository struct {
	SqlHandler
}

// SelectUserByPrimaryKey auth_tokenを条件にレコードを取得する
func (repo *UserRepository) SelectUserByPrimaryKey(userID string) (user domain.User, err error) {
	row, err := repo.Query("SELECT `auth_token`, `name`, `high_score`, `coin` FROM `user` WHERE `id`= ?", userID)
	if err != nil {
		return
	}

	defer row.Close()

	var authToken string
	var name string
	var highScore int32
	var coin int32

	row.Next()
	if err = row.Scan(&authToken, &name, &highScore, &coin); err != nil {
		return
	}

	user.ID = userID
	user.AuthToken = authToken
	user.Name = name
	user.HighScore = highScore
	user.Coin = coin
	return
}

// SelectUserByAuthToken auth_tokenを条件にレコードを取得する
func (repo *UserRepository) SelectUserByAuthToken(authToken string) (user domain.User, err error) {
	row, err := repo.Query("SELECT `id`, `name`, `high_score`, `coin` FROM `user` WHERE `auth_token`=?", authToken)
	if err != nil {
		return
	}

	defer row.Close()

	var userID string
	var name string
	var highScore int32
	var coin int32

	row.Next()
	if err = row.Scan(&userID, &name, &highScore, &coin); err != nil {
		return
	}

	user.ID = userID
	user.Name = name
	user.HighScore = highScore
	user.Coin = coin

	return
}

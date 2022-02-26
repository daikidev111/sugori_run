package database

import "22dojo-online/pkg/domain"

type UserRepository struct {
	SqlHandler
}

// SelectUserByAuthToken auth_tokenを条件にレコードを取得する
func (repo *UserRepository) SelectUserByUserID(userID string) (user domain.User, err error) {
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

package database

import (
	"22dojo-online/pkg/domain/entity"
	"log"
)

type UserRepository struct {
	SQLHandler
}

// SelectUserByPrimaryKey auth_tokenを条件にレコードを取得する
func (repo *UserRepository) SelectUserByPrimaryKey(userID string) (user *entity.User, err error) {
	user = &entity.User{}
	row := repo.QueryRow("SELECT `id`, `auth_token`, `name`, `high_score`, `coin` FROM `user` WHERE `id`= ?", userID)
	if err != nil {
		log.Println(err, "Query Row error")
		return
	}

	if err := row.Scan(&user.ID, &user.AuthToken, &user.Name, &user.HighScore, &user.Coin); err != nil {
		return nil, err
	}

	return user, nil
}

// SelectUserByAuthToken auth_tokenを条件にレコードを取得する
func (repo *UserRepository) SelectUserByAuthToken(authToken string) (user *entity.User, err error) {
	user = &entity.User{}
	row := repo.QueryRow("SELECT `id`, `auth_token`, `name`, `high_score`, `coin` FROM `user` WHERE `auth_token`=?", authToken)
	if err != nil {
		log.Println(err, "Query Row error")
		return
	}

	if err := row.Scan(&user.ID, &user.AuthToken, &user.Name, &user.HighScore, &user.Coin); err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) InsertUser(user *entity.User) error {
	_, err := repo.Execute(
		"INSERT INTO `user` (`id`, `auth_token`, `name`, `high_score`, `coin`) VALUES (?, ?, ?, ?, ?);",
		user.ID, user.AuthToken, user.Name, user.HighScore, user.Coin)
	return err
}

func (repo *UserRepository) UpdateUserByPrimaryKey(user *entity.User) error {
	_, err := repo.Execute(
		"UPDATE user SET name = ? WHERE id = ?",
		user.Name, user.ID)
	return err
}

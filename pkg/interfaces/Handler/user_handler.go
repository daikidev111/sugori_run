package handler

import (
	"22dojo-online/pkg/domain"
	"22dojo-online/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
)

type Handler interface {
	GetUserHandler(writer http.ResponseWriter, response *domain.User)
}

func New() Handler {
	return &UserHandler{}
}

type UserHandler struct {
}

func (uh *UserHandler) GetUserHandler(writer http.ResponseWriter, response *domain.User) {
	if response == nil {
		return
	}
	body := &domain.UserGetResponse{
		ID:        response.ID,
		Name:      response.Name,
		HighScore: response.HighScore,
		Coin:      response.Coin,
	}

	data, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		utils.InternalServerError(writer, "marshal error")
		return
	}
	if _, err := writer.Write(data); err != nil {
		log.Println(err)
	}
}

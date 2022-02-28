package handler

// import (
// 	"22dojo-online/pkg/domain"
// 	"22dojo-online/pkg/utils"
// 	"encoding/json"
// 	"log"
// 	"net/http"
// )

// type UserGetResponse struct {
// 	ID        string `json:"id"`
// 	Name      string `json:"name"`
// 	HighScore int32  `json:"highScore"`
// 	Coin      int32  `json:"coin"`
// }

// type UserCreateResponse struct {
// 	Token string `json:"token"`
// }

// type Handler interface {
// 	GetUserHandler(writer http.ResponseWriter, response *domain.User)
// 	CreateUserHandler(writer http.ResponseWriter, response string)
// }

// func New() Handler {
// 	return &UserHandler{}
// }

// type UserHandler struct {
// }

// func (uh *UserHandler) GetUserHandler(writer http.ResponseWriter, response *domain.User) {
// 	if response == nil {
// 		return
// 	}
// 	body := &UserGetResponse{
// 		ID:        response.ID,
// 		Name:      response.Name,
// 		HighScore: response.HighScore,
// 		Coin:      response.Coin,
// 	}

// 	data, err := json.Marshal(body)
// 	if err != nil {
// 		log.Println(err)
// 		utils.InternalServerError(writer, "marshal error")
// 		return
// 	}
// 	if _, err := writer.Write(data); err != nil {
// 		log.Println(err)
// 	}
// }

// func (uh *UserHandler) CreateUserHandler(writer http.ResponseWriter, response string) {
// 	if response == "" {
// 		log.Println("Auth token is empty")
// 		return
// 	}
// 	body := &UserCreateResponse{
// 		Token: response,
// 	}

// 	data, err := json.Marshal(body)
// 	if err != nil {
// 		log.Println(err)
// 		utils.InternalServerError(writer, "marshal error")
// 		return
// 	}
// 	if _, err := writer.Write(data); err != nil {
// 		log.Println(err)
// 	}
// }

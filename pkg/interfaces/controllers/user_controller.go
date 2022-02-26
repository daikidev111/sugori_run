package controllers

import (
	"22dojo-online/pkg/dcontext"
	"22dojo-online/pkg/domain"
	"22dojo-online/pkg/http/response"
	"22dojo-online/pkg/interfaces/database"
	"22dojo-online/pkg/usecase"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type UserController struct {
	Interactor usecase.UserInteractor
}

type userGetResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	HighScore int32  `json:"highScore"`
	Coin      int32  `json:"coin"`
}

type userCreateRequest struct {
	Name string `json:"name"`
}

type userCreateResponse struct {
	Token string `json:"token"`
}

type userUpdateRequest struct {
	Name string `json:"name"`
}

func NewUserController(sqlHandler database.SQLHandler) *UserController {
	return &UserController{
		Interactor: usecase.UserInteractor{
			UserRepository: &database.UserRepository{
				SQLHandler: sqlHandler,
			},
		},
	}
}

func (controller *UserController) GetUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println("userID is empty")
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		user, err := controller.Interactor.SelectUserByPrimaryKey(userID)

		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		response.Success(writer, &userGetResponse{
			ID:        user.ID,
			Name:      user.Name,
			HighScore: user.HighScore,
			Coin:      user.Coin,
		})
	}
}

func (controller *UserController) InsertUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// リクエストBodyから更新後情報を取得
		var requestBody userCreateRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
			response.BadRequest(writer, "Bad Request")
			return
		}

		// UUIDでユーザIDを生成する
		userID, err := uuid.NewRandom()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// UUIDで認証トークンを生成する
		authToken, err := uuid.NewRandom()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// データベースにユーザデータを登録する
		err = controller.Interactor.InsertUser(domain.User{
			ID:        userID.String(),
			AuthToken: authToken.String(),
			Name:      requestBody.Name,
			HighScore: 0,
			Coin:      0,
		})
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// 生成した認証トークンを返却
		response.Success(writer, &userCreateResponse{Token: authToken.String()})
	}
}

func (controller *UserController) UpdateUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// リクエストBodyから更新後情報を取得
		var requestBody userUpdateRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
			response.BadRequest(writer, "Bad Request")
			return
		}

		// Contextから認証済みのユーザIDを取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println("userID is empty")
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		user, err := controller.Interactor.SelectUserByPrimaryKey(userID)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		user.Name = requestBody.Name

		// if user == nil {
		// 	log.Println("user not found")
		// 	response.BadRequest(writer, fmt.Sprintf("user not found. userID=%s", userID))
		// 	return
		// }

		err = controller.Interactor.UpdateUserByPrimaryKey(user)
		if err != nil {
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		response.Success(writer, nil)
	}
}

package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"22dojo-online/pkg/dcontext"
	"22dojo-online/pkg/domain"
	handler "22dojo-online/pkg/interfaces/Handler"
	"22dojo-online/pkg/interfaces/database"
	"22dojo-online/pkg/usecase"
	"22dojo-online/pkg/utils"

	"github.com/google/uuid"
)

type UserController struct {
	Interactor usecase.UserInteractor
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
			utils.InternalServerError(writer, "Internal Server Error")
			return
		}
		var user *domain.User
		user, err := controller.Interactor.SelectUserByPrimaryKey(userID)

		if err != nil {
			log.Println(err)
			utils.InternalServerError(writer, "Internal Server Error")
			return
		}
		if user == nil {
			log.Println("user not found")
			utils.BadRequest(writer, fmt.Sprintf("user not found. userID=%s", userID))
			return
		}

		userHandler := handler.New()
		userHandler.GetUserHandler(writer, user)
	}
}

func (controller *UserController) InsertUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// リクエストBodyから更新後情報を取得
		var requestBody domain.UserCreateRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
			utils.BadRequest(writer, "Bad Request")
			return
		}

		// UUIDでユーザIDを生成する
		userID, err := uuid.NewRandom()
		if err != nil {
			log.Println(err)
			utils.InternalServerError(writer, "Internal Server Error")
			return
		}

		// UUIDで認証トークンを生成する
		authToken, err := uuid.NewRandom()
		if err != nil {
			log.Println(err)
			utils.InternalServerError(writer, "Internal Server Error")
			return
		}

		// データベースにユーザデータを登録する
		err = controller.Interactor.InsertUser(&domain.User{
			ID:        userID.String(),
			AuthToken: authToken.String(),
			Name:      requestBody.Name,
			HighScore: 0,
			Coin:      0,
		})

		if err != nil {
			log.Println(err)
			utils.InternalServerError(writer, "Internal Server Error")
			return
		}

		userHandler := handler.New()
		userHandler.CreateUserHandler(writer, authToken.String())
	}
}

func (controller *UserController) UpdateUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// リクエストBodyから更新後情報を取得
		var requestBody domain.UserUpdateRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
			utils.BadRequest(writer, "Bad Request")
			return
		}

		// Contextから認証済みのユーザIDを取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println("userID is empty")
			utils.InternalServerError(writer, "Internal Server Error")
			return
		}

		var user *domain.User
		var err error
		user, err = controller.Interactor.SelectUserByPrimaryKey(userID)
		if err != nil {
			log.Println(err)
			utils.InternalServerError(writer, "Internal Server Error")
			return
		}
		if user == nil {
			log.Println("user not found")
			utils.BadRequest(writer, fmt.Sprintf("user not found. userID=%s", userID))
			return
		}
		user.Name = requestBody.Name

		err = controller.Interactor.UpdateUserByPrimaryKey(user)
		if err != nil {
			utils.InternalServerError(writer, "Internal Server Error")
			return
		}
	}
}

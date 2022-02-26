package server

import (
	"database/sql"
	"log"
	"net/http"

	"22dojo-online/pkg/http/middleware"
	"22dojo-online/pkg/infrastructure"
	"22dojo-online/pkg/interfaces/controllers"
	"22dojo-online/pkg/server/handler"
)

// Serve HTTPサーバを起動する
func Serve(addr string) {
	var db *sql.DB
	m := middleware.NewAuth(db)

	userController := controllers.NewUserController(infrastructure.NewSqlHandler())

	/* ===== URLマッピングを行う ===== */
	http.HandleFunc("/setting/get", get(controllers.GetSetting()))
	http.HandleFunc("/user/get",
		get(m.Authenticate(userController.GetUser()))) // middleware.Authenticateでhandler funcを囲む
	// http.HandleFunc("/user/get",
	// 	get(middleware.Authenticate(handler.HandleUserGet()))) // middleware.Authenticateでhandler funcを囲む

	http.HandleFunc("/user/create", post(handler.HandleUserCreate()))

	http.HandleFunc("/user/update",
		post(handler.HandleUserUpdate()))

	/* ===== サーバの起動 ===== */
	log.Println("Server running...")
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Listen and serve failed. %+v", err)
	}
}

// get GETリクエストを処理する
func get(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodGet)
}

// post POSTリクエストを処理する
func post(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodPost)
}

// httpMethod 指定したHTTPメソッドでAPIの処理を実行する
func httpMethod(apiFunc http.HandlerFunc, method string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// CORS対応
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,Accept,Origin,x-token")

		// プリフライトリクエストは処理を通さない
		if request.Method == http.MethodOptions {
			return
		}
		// 指定のHTTPメソッドでない場合はエラー
		if request.Method != method {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			if _, err := writer.Write([]byte("Method Not Allowed")); err != nil {
				log.Println(err)
			}
			return
		}

		// 共通のレスポンスヘッダを設定
		writer.Header().Add("Content-Type", "application/json")
		apiFunc(writer, request)
	}
}

package server

import (
	"log"
	"net/http"

	"22dojo-online/pkg/http/middleware"
	"22dojo-online/pkg/server/handler"
)

// Serve HTTPサーバを起動する
func Serve(addr string) {
	/* ===== URLマッピングを行う ===== */
	http.HandleFunc("/setting/get", get(handler.HandleSettingGet()))
	http.HandleFunc("/user/create", post(handler.HandleUserCreate()))

	// middlewareは 22dojo-online/pkg/http/middleware パッケージを利用する
	// middleware を利用することでauth_tokenありきのoperationができるようになる
	http.HandleFunc("/user/get",
		get(middleware.Authenticate(handler.HandleUserGet()))) // middleware.Authenticateでhandler funcを囲む
	http.HandleFunc("/user/update",
		post(middleware.Authenticate(handler.HandleUserUpdate())))

	// Collection List の表示の際のhttp handler
	http.HandleFunc("/collection/list",
		get(middleware.Authenticate(handler.HandleCollectionGet())))

	http.HandleFunc("/ranking/list",
		get(middleware.Authenticate(handler.HandleRankingGet())))

	http.HandleFunc("/gacha/draw",
		post(middleware.Authenticate(handler.HandleGachaPost())))

	http.HandleFunc("/game/finish",
		post(middleware.Authenticate(handler.HandleGameFinishPost())))

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

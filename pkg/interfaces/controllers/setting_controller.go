package controllers

import (
	"22dojo-online/pkg/interfaces/handler"
	"net/http"
)

// GetSetting ゲーム設定情報取得処理
func GetSetting() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		settingHandler := handler.NewSettingHandler()
		settingHandler.GetSettingHandler(writer)
	}
}

type settingGetResponse struct {
	GachaCoinConsumption int32 `json:"gachaCoinConsumption"`
}

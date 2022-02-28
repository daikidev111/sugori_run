package controllers

import (
	"22dojo-online/pkg/constant"
	"22dojo-online/pkg/errors"
	"encoding/json"
	"log"
	"net/http"
)

// GetSetting ゲーム設定情報取得処理
func GetSetting() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// settingHandler := handler.NewSettingHandler()
		// settingHandler.GetSettingHandler(writer, user)
		response := &settingGetResponse{
			GachaCoinConsumption: constant.GachaCoinConsumption,
		}
		if response == nil {
			return
		}

		data, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
			errors.InternalServerError(writer, "marshal error")
			return
		}
		if _, err := writer.Write(data); err != nil {
			log.Println(err)
		}
	}
}

type settingGetResponse struct {
	GachaCoinConsumption int32 `json:"gachaCoinConsumption"`
}

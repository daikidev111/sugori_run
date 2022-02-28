package handler

// import (
// 	"22dojo-online/pkg/constant"
// 	"22dojo-online/pkg/domain"
// 	"22dojo-online/pkg/utils"
// 	"encoding/json"
// 	"log"
// 	"net/http"
// )

// type SettingHandler interface {
// 	GetSettingHandler(writer http.ResponseWriter)
// }

// func NewSettingHandler() SettingHandler {
// 	return &SettingHandle{}
// }

// type SettingHandle struct {
// }

// func (uh *SettingHandle) GetSettingHandler(writer http.ResponseWriter) {
// 	body := &domain.SettingGetResponse{
// 		GachaCoinConsumption: constant.GachaCoinConsumption,
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

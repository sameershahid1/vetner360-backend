package controller

import (
	"encoding/json"
	"net/http"
	"vetner360-backend/model"
	data_type "vetner360-backend/utils/type"
)

func GetRoles(response http.ResponseWriter, request *http.Request) {
	requestResponse := data_type.Response[model.Role]{Status: true, Message: "Hi, admin"}
	jsonData, _ := json.Marshal(requestResponse)
	response.WriteHeader(http.StatusOK)
	response.Write(jsonData)
}

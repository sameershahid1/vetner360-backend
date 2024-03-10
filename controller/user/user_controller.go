package userController

import (
	"encoding/json"
	"net/http"
)

type postResposne struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	var response postResposne = postResposne{Status: true, Message: "Successfully "}
	jsonData := json.Marshal(response)

	w.Write(jsonData)
}

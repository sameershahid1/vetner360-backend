package controller

import (
	"log"
	"net/http"
)

func SocketHandling(response http.ResponseWriter, request *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

}

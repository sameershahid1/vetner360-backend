package custom_middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func ValidateJsonFormat(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		var body bytes.Buffer

		_, err := io.Copy(&body, request.Body)
		if err != nil {
			http.Error(response, fmt.Sprintf("Error reading request body: %v", err), http.StatusBadRequest)
			return
		}

		request.Body = ioutil.NopCloser(&body)
		if !json.Valid(body.Bytes()) {
			http.Error(response, "Invalid json format", http.StatusBadRequest)
			return
		}

	})
}

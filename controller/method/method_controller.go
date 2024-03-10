package methodController

import "net/http"

func RouteDoesExists(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("route does not exist"))
}

func MethodNotExists(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(405)
	w.Write([]byte("method is not valid"))
}

package router

import "net/http"

func authTestHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You are authenticated!"))
}

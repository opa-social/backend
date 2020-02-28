package router

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/opa-social/backend/internal/firebase"
)

type Router struct {
	address  string
	port     uint
	firebase *firebase.Controller
	server   *http.Server
}

func Setup(address string, port uint, firebase *firebase.Controller) Router {
	return Router{
		address:  address,
		port:     port,
		firebase: firebase,
	}
}

func (r *Router) Create() {
	router := mux.NewRouter()

	// Use firebase authentication middleware.
	router.Use(r.firebase.AuthMiddleware)

	r.server = &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%d", r.address, r.port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

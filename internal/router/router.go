package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/opa-social/backend/internal/firebase"
)

// Router stores internal fields required to run HTTP router.
type Router struct {
	address  string
	port     uint
	firebase *firebase.Controller
	server   *http.Server
}

// Setup creates a Router with minimum required values.
func Setup(address string, port uint, firebase *firebase.Controller) Router {
	return Router{
		address:  address,
		port:     port,
		firebase: firebase,
	}
}

// Create intializes the router, path handlers, middleware, etc.
// It also creates the HTTP server struct.
func (r *Router) Create() {
	router := mux.NewRouter()

	// Use firebase authentication middleware.
	router.Use(r.firebase.AuthMiddleware)

	// Routes.
	router.HandleFunc("/test", authTestHandler).Methods(http.MethodGet)

	r.server = &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%d", r.address, r.port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

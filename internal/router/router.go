package router

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
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

// Serve takes the default server configuration and runs it in the main execution
// context. It also optionally sets up the server configuration if it hasn't been
// setup already.
func (r *Router) Serve() {
	if r.server == nil {
		r.Create()
	}

	go func() {
		log.Println("Now serving...")
		log.Fatal(r.server.ListenAndServe())
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	<-c

	log.Println("\nServer shutting down. Goodbye...")
}

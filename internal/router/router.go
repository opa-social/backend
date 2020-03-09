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
	router.HandleFunc("/event/new", r.newEventHandler).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json;charset=utf-8")
	router.HandleFunc("/event/{event:[A-Za-z]{5}}/join", r.joinEventHandler).
		Methods(http.MethodPost)
	router.HandleFunc("/event/{event:[A-Za-z]{5}}/matches", r.getEventMatches).
		Methods(http.MethodGet)

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
		log.Printf("Now serving on %s:%d...\n", r.address, r.port)
		log.Fatal(r.server.ListenAndServe())
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	<-c

	fmt.Print("\n")
	log.Println("Server shutting down. Goodbye...")
}

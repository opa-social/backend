package router

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/opa-social/backend/internal/firebase"
)

func (router *Router) newEventHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Got request to make new event")
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Unable to read body from %s with error: \"%s\"", r.RemoteAddr, err)
		http.Error(w, "Bad request", http.StatusBadRequest)

		return
	}

	request := &firebase.EventRequest{}
	err = json.Unmarshal(body, request)
	if err != nil {
		log.Printf("Unable to read body from %s with error: \"%s\"", r.RemoteAddr, err)
		http.Error(w, "Malformed JSON request", http.StatusBadRequest)

		return
	}

	response, err := router.firebase.CreateNewEvent(r.Header.Get("X-OPA-UID"), request)
	if err != nil {
		http.Error(w, "Could not form response", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (router *Router) joinEventHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("User %s joined event %s", r.Header.Get("X-OPA-UID"), mux.Vars(r)["event"])

	w.WriteHeader(http.StatusAccepted)
}

func (router *Router) getEventMatches(w http.ResponseWriter, r *http.Request) {
	uid := r.Header.Get("X-OPA-UID")
	eventID := mux.Vars(r)["event"]

	log.Printf("User %s wants matches for event %s", uid, eventID)
}

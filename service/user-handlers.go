package service

import (
	"encoding/json"
	"net/http"

	"github.com/clockworksoul/cog2/data/rest"
	"github.com/clockworksoul/cog2/dataaccess"
	"github.com/gorilla/mux"
)

// handleGetUsers handles "GET /v2/user"
func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := dataaccess.UserList()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

// handleGetUser handles "GET /v2/user/{username}"
func handleGetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	exists, err := dataaccess.UserExists(params["username"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "No such user", http.StatusNotFound)
		return
	}

	user, err := dataaccess.UserGet(params["username"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// handlePostUser handles "POST /v2/user/{username}"
func handlePostUser(w http.ResponseWriter, r *http.Request) {
	var user rest.User
	var err error

	params := mux.Vars(r)

	_ = json.NewDecoder(r.Body).Decode(&user)

	user.Username = params["username"]

	exists, err := dataaccess.UserExists(user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// NOTE: Should we just make "update" create users that don't exist?
	if exists {
		err = dataaccess.UserUpdate(user)
	} else {
		err = dataaccess.UserCreate(user)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handlePostUser handles "DELETE /v2/user/{username}"
func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := dataaccess.UserDelete(params["username"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addUserMethodsToRouter(router *mux.Router) {
	router.HandleFunc("/v2/user", handleGetUsers).Methods("GET")
	router.HandleFunc("/v2/user/{username}", handleGetUser).Methods("GET")
	router.HandleFunc("/v2/user/{username}", handlePostUser).Methods("POST")
	router.HandleFunc("/v2/user/{username}", handleDeleteUser).Methods("DELETE")
}
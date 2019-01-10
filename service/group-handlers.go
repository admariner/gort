package service

import (
	"encoding/json"
	"net/http"

	"github.com/clockworksoul/cog2/dal"
	"github.com/clockworksoul/cog2/data/rest"
	"github.com/gorilla/mux"
)

// handleGetGroup handles "GET /v2/group/{groupname}"
func handleGetGroup(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	exists, err := dal.GroupExists(params["groupname"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "No such group", http.StatusNotFound)
		return
	}

	group, err := dal.GroupGet(params["groupname"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(group)
}

// handleGetGroups handles "GET /v2/group"
func handleGetGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := dal.GroupList()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(groups)
}

// handleGetGroupMembers handles "GET /v2/group/{groupname}/member"
func handleGetGroupMembers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	groupname := params["groupname"]

	exists, err := dal.GroupExists(groupname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "no such group", http.StatusNotFound)
		return
	}

	group, err := dal.GroupGet(groupname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(group.Users)
}

// handlePostGroup handles "POST /v2/group/{groupname}"
func handlePostGroup(w http.ResponseWriter, r *http.Request) {
	var group rest.Group
	var err error

	params := mux.Vars(r)

	_ = json.NewDecoder(r.Body).Decode(&group)

	group.Name = params["groupname"]

	exists, err := dal.GroupExists(group.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// NOTE: Should we just make "update" create groups that don't exist?
	if exists {
		err = dal.GroupUpdate(group)
	} else {
		err = dal.GroupCreate(group)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleDeleteGroupMember handles "DELETE "/v2/group/{groupname}/member/{username}""
func handleDeleteGroupMember(w http.ResponseWriter, r *http.Request) {
	var exists bool
	var err error

	params := mux.Vars(r)
	groupname := params["groupname"]
	username := params["username"]

	exists, err = dal.GroupExists(groupname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "no such group", http.StatusNotFound)
		return
	}

	exists, err = dal.UserExists(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "no such user", http.StatusNotFound)
		return
	}

	err = dal.GroupRemoveUser(groupname, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handlePostGroupMember handles "POST "/v2/group/{groupname}/member/{username}""
func handlePostGroupMember(w http.ResponseWriter, r *http.Request) {
	var exists bool
	var err error

	params := mux.Vars(r)
	groupname := params["groupname"]
	username := params["username"]

	exists, err = dal.GroupExists(groupname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "no such group", http.StatusNotFound)
		return
	}

	exists, err = dal.UserExists(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "no such user", http.StatusNotFound)
		return
	}

	err = dal.GroupAddUser(groupname, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handlePostGroup handles "DELETE /v2/group/{groupname}"
func handleDeleteGroup(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := dal.GroupDelete(params["groupname"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addGroupMethodsToRouter(router *mux.Router) {
	// Basic group methods
	router.HandleFunc("/v2/group", handleGetGroups).Methods("GET")
	router.HandleFunc("/v2/group/{groupname}", handleGetGroup).Methods("GET")
	router.HandleFunc("/v2/group/{groupname}", handlePostGroup).Methods("POST")
	router.HandleFunc("/v2/group/{groupname}", handleDeleteGroup).Methods("DELETE")

	// Group user membership
	router.HandleFunc("/v2/group/{groupname}/member", handleGetGroupMembers).Methods("GET")
	router.HandleFunc("/v2/group/{groupname}/member/{username}", handleDeleteGroupMember).Methods("DELETE")
	router.HandleFunc("/v2/group/{groupname}/member/{username}", handlePostGroupMember).Methods("POST")
}
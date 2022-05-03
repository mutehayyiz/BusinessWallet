package controller

import (
	"BusinessWallet/auth"
	"BusinessWallet/storage"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type User struct{}

func (u User) Contact(w http.ResponseWriter, r *http.Request) {
	contactId := mux.Vars(r)["id"]
	claims := r.Context().Value("claims").(*auth.Claims)

	cId, err := strconv.ParseInt(contactId, 10, 64)

	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if r.Method == http.MethodGet {
		err = storage.User.AddContact(claims.Id, int(cId), 0)
	} else {
		err = storage.User.DeleteContact(claims.Id, int(cId), 0)
	}

	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ReturnResponse(w, http.StatusOK, "ok")
}

/*
func (u User) Contacts(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*auth.Claims)

	contacts, err := storage.User.Contacts(claims.Id)
}
*/

/*
func (u User) AddContact(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	eventId := mux.Vars(r)["eventId"]
	claims := r.Context().Value("claims").(*auth.Claims)

	err := storage.User.AddContact(claims.Id, userId, claims)
}

*/

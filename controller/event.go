package controller

import (
	"BusinessWallet/auth"
	"BusinessWallet/model"
	"BusinessWallet/storage"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Event struct{}

func (e Event) Create(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*auth.Claims)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ReturnError(w, http.StatusBadRequest, err.Error())
		return
	}

	var create model.CreateEventRequest
	if err = json.Unmarshal(body, &create); err != nil {
		ReturnError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = create.Validate(); err != nil {
		ReturnError(w, http.StatusBadRequest, err.Error())
		return
	}

	event, err := storage.Event.Create(&create, claims.Id)
	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ReturnResponse(w, http.StatusOK, event)
}

func (e Event) Get(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	event, err := storage.Event.Get(eventID)
	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ReturnResponse(w, http.StatusOK, event)
}

func (e Event) Attend(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	claims := r.Context().Value("claims").(*auth.Claims)

	eId, err := strconv.ParseInt(eventID, 10, 64)
	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = storage.Event.Attend(int(eId), claims.Id)
	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ReturnResponse(w, http.StatusOK, "ok")
}

func (e Event) Leave(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	claims := r.Context().Value("claims").(*auth.Claims)

	err := storage.Event.Leave(eventID, claims.Id)
	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ReturnResponse(w, http.StatusOK, "ok")
}

func (e Event) Delete(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	claims := r.Context().Value("claims").(*auth.Claims)

	err := storage.Event.Delete(eventID, claims.Id)
	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ReturnResponse(w, http.StatusOK, "ok")
}

func (e Event) Past(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*auth.Claims)

	events, err := storage.Event.Past(claims.Id)
	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ReturnResponse(w, http.StatusOK, events)
}

func (e Event) Active(w http.ResponseWriter, r *http.Request) {
	events, err := storage.Event.Active()
	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ReturnResponse(w, http.StatusOK, events)
}

func (e Event) Now(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*auth.Claims)

	events, err := storage.Event.Now(claims.Id)
	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ReturnResponse(w, http.StatusOK, events)
}

func (e Event) Together(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*auth.Claims)
	contactID := mux.Vars(r)["userId"]

	cId, err := strconv.ParseInt(contactID, 10, 64)
	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	events, err := storage.Event.Together(claims.Id, int(cId))
	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ReturnResponse(w, http.StatusOK, events)
}

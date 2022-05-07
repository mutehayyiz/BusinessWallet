package controller

import (
	"BusinessWallet/auth"
	"BusinessWallet/storage"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type User struct{}

func (u User) AddOrDeleteContact(w http.ResponseWriter, r *http.Request) {
	contactId := mux.Vars(r)["id"]
	claims := r.Context().Value("claims").(*auth.Claims)

	cId, err := strconv.ParseInt(contactId, 10, 64)
	if err != nil {
		logrus.Error(err)
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	switch r.Method {
	case http.MethodPost:
		logrus.Info("adding contact with id ", cId)
		err = storage.User.AddContact(claims.Id, int(cId), 0)
	case http.MethodDelete:
		logrus.Info("deleting contact with id ", cId)
		err = storage.User.DeleteContact(claims.Id, int(cId), 0)
	}

	if err != nil {
		logrus.WithError(err).Error("storage error ")
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("success contact with id %d", cId)
	ReturnResponse(w, http.StatusOK, "ok")
}

func (u User) Me(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*auth.Claims)

	user, err := storage.User.GetUser(claims.Id)

	if err != nil {
		logrus.WithError(err).Error("storage error ")
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("me, %+v", user)
	ReturnResponse(w, http.StatusOK, user)
}

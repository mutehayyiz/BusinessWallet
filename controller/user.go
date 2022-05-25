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

func (u User) GetContact(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	cId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logrus.Error(err)
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}
	user, err := storage.User.GetUser(int(cId))

	if err != nil {
		logrus.WithError(err).Error("storage error ")
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	mask := struct {
		Name     string `json:"name"`
		Surname  string `json:"surname"`
		Email    string `json:"email" gorm:"index:idx_email,unique"`
		Phone    string `json:"phone" gorm:"index:idx_phone,unique"`
		Linkedin string `json:"linkedin"`
		Company  string `json:"company"`
	}{
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Phone:    user.Phone,
		Linkedin: user.Linkedin,
		Company:  user.Company,
	}

	logrus.Infof("user info, %+v", mask)
	ReturnResponse(w, http.StatusOK, mask)
}

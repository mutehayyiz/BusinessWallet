package controller

import (
	"BusinessWallet/auth"
	"BusinessWallet/model"
	"BusinessWallet/storage"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type Auth struct{}

func (Auth) Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.WithError(err).Error("bad request ")
		ReturnError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(string(body))

	var register model.RegisterRequest
	if err = json.Unmarshal(body, &register); err != nil {
		logrus.WithError(err).Error("unmarshal error ")
		ReturnError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = register.Validate(); err != nil {
		logrus.WithError(err).Error("validation error ")

		ReturnError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := storage.User.Register(&register)
	if err != nil {
		logrus.WithError(err).Error("storage error ")
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := auth.CreateToken(int(user.ID), user.Name)
	if err != nil {
		logrus.WithError(err).Error("create token error ")
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := model.LoginResponse{
		UserData: *user,
		Token:    token,
	}

	logrus.Infof("registered in , %+v", user)
	ReturnResponse(w, http.StatusOK, resp)
}

func (Auth) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.WithError(err).Error("bad request ")
		ReturnError(w, http.StatusBadRequest, err.Error())
		return
	}

	var login model.LoginRequest
	if err = json.Unmarshal(body, &login); err != nil {
		logrus.WithError(err).Error("storage error ")
		ReturnError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = login.Validate(); err != nil {
		logrus.WithError(err).Error("validation error ")
		ReturnError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := storage.User.Login(&login)
	if err != nil {
		logrus.WithError(err).Error("storage error ")
		ReturnError(w, http.StatusBadRequest, "email or password is wrong!")
		return
	}

	token, err := auth.CreateToken(int(user.ID), user.Name)
	if err != nil {
		logrus.WithError(err).Error("crete token error: ")
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := model.LoginResponse{
		UserData: *user,
		Token:    token,
	}

	logrus.Infof("logged in %+v", user)
	ReturnResponse(w, http.StatusOK, resp)
}

func (a Auth) VerifyToken(w http.ResponseWriter, r *http.Request) {
	logrus.Info("token verified")
	ReturnResponse(w, http.StatusOK, "ok")
	return
}

/*
func (Auth) RefreshToken(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value("claims").(*auth.Claims)
	token, err := auth.CreateToken(claims.Credentials)
	if err != nil {
		ReturnError(w, http.StatusInternalServerError, "Token couldn't created")
		return
	}

	data := map[string]string{
		"data":  "Token Refreshed",
		"token": token,
	}
	ReturnResponse(w, 200, data)
}

*/

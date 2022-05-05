package main

import (
	"BusinessWallet/config"
	"BusinessWallet/controller"
	"BusinessWallet/middleware"
	"BusinessWallet/storage"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	err := config.Global.Load("config.json")
	if err != nil {
		logrus.WithError(err).Fatal("config error")
		os.Exit(1)
	}

	err = storage.Connect(config.Global.Storage)
	if err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}

	err = storage.Delete()
	if err != nil {
		logrus.Error(err)
	}

	storage.Seed()

	router := GenerateRoutes()

	logrus.Info("listening on port: ", config.Global.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Global.Port), router)
	if err != nil {
		panic(err)
	}
}

func GenerateRoutes() *mux.Router {
	r := mux.NewRouter()

	auth := controller.Auth{}

	r.HandleFunc("/register", auth.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", auth.Login).Methods(http.MethodPost)
	r.HandleFunc("/verify_token", middleware.Auth(auth.VerifyToken)).Methods(http.MethodPost)

	r.HandleFunc("/delete", controller.DeleteAll).Methods(http.MethodDelete)
	r.HandleFunc("/seed", controller.Seed).Methods(http.MethodGet)

	api := r.PathPrefix("/api").Subrouter()

	event := controller.Event{}

	api.HandleFunc("/event", middleware.Auth(event.Create)).Methods(http.MethodPost)
	api.HandleFunc("/event/past", middleware.Auth(event.Past)).Methods(http.MethodGet)
	api.HandleFunc("/event/active", middleware.Auth(event.Active)).Methods(http.MethodGet)
	api.HandleFunc("/event/current", middleware.Auth(event.Now)).Methods(http.MethodGet)
	api.HandleFunc("/event/{id}", middleware.Auth(event.Get)).Methods(http.MethodGet)
	api.HandleFunc("/event/{id}/attend", middleware.Auth(event.Attend)).Methods(http.MethodPost)
	api.HandleFunc("/event/{id}/leave", middleware.Auth(event.Leave)).Methods(http.MethodPost)
	api.HandleFunc("/event/{id}/delete", middleware.Auth(event.Delete)).Methods(http.MethodDelete)
	api.HandleFunc("/event/{userId}/together", middleware.Auth(event.Together)).Methods(http.MethodGet)

	user := controller.User{}

	api.HandleFunc("/user/contact/{id}", middleware.Auth(user.Contact)).Methods(http.MethodGet)
	api.HandleFunc("/user/contact/{id}", middleware.Auth(user.Contact)).Methods(http.MethodDelete)

	return r
}

package api

import (
	"BusinessWallet/config"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var server = controller.Server{}

func Run() {
	// load config.conf
	err := config.Global.Load("config.conf")
	if err != nil {
		logrus.WithError(err).Fatal("config error")
		os.Exit(1)
	}

	server.ConnectDB(config.Global.Storage)
	server.GenerateRoutes()

	logrus.Info(fmt.Sprintf("server listening on port: %d", config.Global.Port))

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", config.Global.Host, config.Global.Port), server.Router)
	if err != nil {
		logrus.WithError(err).Fatal("server error")
		os.Exit(1)
	}
}

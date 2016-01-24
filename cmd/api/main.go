package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gophergala2016/blogalert/cmd/api/controllers"
	"github.com/gophergala2016/blogalert/repositories/rethink"
)

func main() {
	config, err := OpenConfig("config.yml")
	if err != nil {
		logrus.WithError(err).Fatal("Error loading config file")
	}

	logrus.Info("Loaded config")

	repo, err := rethink.NewRepo(config.RethinkDB)
	if err != nil {
		logrus.WithError(err).Fatal("Error connecting to repository")
	}

	logrus.Info("Connected to repository")

	tokenValidator := controllers.NewTokenValidator(config.Token.ClientID)
	updateController := controllers.NewUpdateController(repo, tokenValidator)
	readController := controllers.NewReadController(repo, tokenValidator)
	readAllController := controllers.NewReadAllController(repo, tokenValidator)
	subscribeController := controllers.NewSubscribeController(repo, tokenValidator)
	unsubscribeController := controllers.NewUnsubscribeController(repo, tokenValidator)

	http.Handle("/updates", updateController)
	http.Handle("/read", readController)
	http.Handle("/readall", readAllController)
	http.Handle("/subscribe", subscribeController)
	http.Handle("/unsubscribe", unsubscribeController)

	logrus.Info("Starting server")
	http.ListenAndServe(config.Server.Listen, nil)
}

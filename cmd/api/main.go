package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
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

	tokenValidator := NewTokenValidator(config)
	updateController := NewUpdateController(repo, tokenValidator)
	readController := NewReadController(repo, tokenValidator)

	http.Handle("/updates", updateController)
	http.Handle("/read", readController)

	logrus.Info("Starting server")
	http.ListenAndServe(config.Server.Listen, nil)
}

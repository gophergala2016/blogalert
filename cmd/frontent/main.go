//go:generate rice embed-go
package main

import (
	"html/template"
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/Sirupsen/logrus"
)

var index *template.Template
var fileServer http.Handler

var context struct {
	Config *Config
}

func main() {
	box, err := rice.FindBox("assets")
	if err != nil {
		logrus.WithError(err).Fatal("Error loading assets")
	}

	str, err := box.String("index.tpl.html")
	if err != nil {
		logrus.WithError(err).Fatal("Error loading index")
	}
	tmpl, err := template.New("name").Parse(str)
	if err != nil {
		logrus.WithError(err).Fatal("Error parsing index")
	}

	index = tmpl

	config, err := OpenConfig("config.yml")
	if err != nil {
		logrus.WithError(err).Fatal("Error loading config file")
	}
	logrus.Info("Loaded config")

	context.Config = config

	logrus.Info("Starting server")

	fileServer = http.FileServer(box.HTTPBox())

	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(config.Server.Listen, nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		index.Execute(w, &context)
	} else {
		fileServer.ServeHTTP(w, r)
	}
}

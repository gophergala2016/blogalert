package main

import (
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gophergala2016/blogalert"
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

	wp := blogalert.NewWorkerPool(runtime.NumCPU())
	extractor := blogalert.NewExtractor(repo, wp, logrus.StandardLogger())

	ticker := time.NewTicker(config.Refresh)
	logrus.Info("Starting crawl loop")
	for {
		<-ticker.C
		logrus.Info("Crawl started")

		blogs, err := repo.GetAllBlogs()
		if err != nil {
			logrus.WithError(err).Error("Error getting blog list")
			continue
		}

		for _, blog := range blogs {
			extractor.Crawl(blog)
			wp.Wait()
		}

		logrus.Info("Crawl finished")
	}
}

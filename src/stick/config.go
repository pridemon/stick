package main

import (
	"os"
	"time"

	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

var (
	cluster *elastic.Client
)

func NewWorker(queue chan string) *Worker {
	commitInterval := GetenvInt("STICK_COMMIT_INTERVAL", 5)
	return &Worker{
		CommitInterval: time.Duration(commitInterval) * time.Second,
		CommitAmount:   GetenvInt("STICK_COMMIT_AMOUNT", 500),

		queue:   queue,
		elastic: cluster.Bulk(),
	}
}

func NewQueue() chan string {
	return make(chan string, GetenvInt("STICK_INTERNAL_BUFFER", 1000))
}

func NewWorkers(queue chan string) []*Worker {
	cnt := GetenvInt("STICK_WORKERS", 5)
	workers := make([]*Worker, 0, cnt)

	for i := 0; i < cnt; i++ {
		workers = append(workers, NewWorker(queue))
	}

	return workers
}

func init() {
	var err error
	url := os.Getenv("STICK_ELASTIC_URL")
	healthInterval := GetenvInt("STICK_ELASTIC_HEALTHCHECK_INTERVAL", 10)
	maxRetries := GetenvInt("STICK_ELASTIC_MAX_RETRIES", 5)

	// TODO: config for healthcheck and retries
	cluster, err = elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(time.Duration(healthInterval)*time.Second),
		elastic.SetMaxRetries(maxRetries),
	)
	if err != nil {
		log.WithError(err).WithField("url", url).Fatal("can't connect to elastic")
	}
}

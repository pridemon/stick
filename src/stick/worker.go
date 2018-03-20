package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

type Worker struct {
	CommitInterval time.Duration
	CommitAmount   int

	queue   chan string
	elastic *elastic.BulkService
}

type Message map[string]interface{}

func (w *Worker) Start() {
	tick := time.Tick(w.CommitInterval)

	for {
		select {
		case line, ok := <-w.queue:
			if !ok {
				w.Flush()
				return // queue closed, we can exit
			}

			w.AddMsg(line)

		case <-tick:
			w.Flush() // send bulk query to elasticsearch every tick event
		}
	}
}

func (w *Worker) AddMsg(line string) {
	var msg Message

	err := json.Unmarshal([]byte(line), &msg)
	if err != nil {
		log.WithError(err).WithField("line", line).Error("can't parse line as json")
		return
	}

	msg["@timestamp"] = time.Now().UTC().Format(time.RFC3339)
	indexName := time.Now().UTC().Format("stick-2006.01.02")

	req := elastic.NewBulkIndexRequest().
		Index(indexName).
		Type("stick").
		Doc(msg)

	// add one more request to bulk query, send them all if bulk query has more than CommitAmount requests
	if w.elastic.Add(req).NumberOfActions() >= w.CommitAmount {
		w.Flush()
	}
}

func (w *Worker) Flush() {
	log.WithField("count", w.elastic.NumberOfActions()).Debug("flushing to elastic search")

	w.elastic.Do(context.Background()) // sends all accumulated messages to elasticsearch
}

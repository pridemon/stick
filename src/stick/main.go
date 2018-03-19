package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/olivere/elastic"
)

func main() {
	client, err := elastic.NewClient(
		elastic.SetURL("http://10.0.0.12:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetMaxRetries(5),
	)

	if err != nil {
		panic(err)
	}

	t := time.Now()
	indexName := t.Format("stick-2006.01.02")
	exists, err := client.IndexExists(indexName).Do(context.Background())

	if err != nil {
		panic(err)
	}

	if !exists {
		// mapping :=   https://godoc.org/github.com/olivere/elastic  (example)

		createIndex, err := client.CreateIndex(indexName).Do(context.Background())
		if err != nil {
			panic(err)
		}

		if !createIndex.Acknowledged {
			log.Fatalln("index not acknowledged")
		}
	}

	stamp := time.Now().UTC().Format(time.RFC3339)
	message := `{"@timestamp": "` + stamp + `", "test": "me"}`
	resIndex, err := client.Index().
		Index(indexName).
		Type("doc").
		BodyString(message).
		Do(context.Background())

	if err != nil {
		panic(err)
	}

	fmt.Printf("Indexed doc %s to index %s, type %s\n", resIndex.Id, resIndex.Index, resIndex.Type)

}

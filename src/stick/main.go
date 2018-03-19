package main

import (
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()

	queue := NewQueue()

	consumer := &Consumer{
		queue: queue,
	}

	for _, worker := range NewWorkers(queue) {
		go worker.Start()
	}

	app.Post("/", consumer.Handle)

	_ = app.Run(
		iris.Addr(":8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithoutBanner,
		iris.WithoutVersionChecker,
	)

}

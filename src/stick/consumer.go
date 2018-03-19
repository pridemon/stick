package main

import (
	"bufio"
	"io"

	"github.com/kataras/iris/context"
	log "github.com/sirupsen/logrus"
)

type Consumer struct {
	queue chan string
}

func (c *Consumer) Handle(ctx context.Context) {
	reader := bufio.NewReader(ctx.Request().Body)
	for {
		line, err := reader.ReadString('\n')

		if err != nil && err != io.EOF {
			log.WithError(err).Error("can't read line from POST body")
			return
		}

		c.queue <- line

		if err == io.EOF {
			return
		}
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/davyj0nes/worker-example/metrics"

	"github.com/davyj0nes/worker-example/config"
	"github.com/davyj0nes/worker-example/publisher"
	"github.com/davyj0nes/worker-example/srv"
	"github.com/davyj0nes/worker-example/wrkr"
)

func init() {
	metrics.Initialise()
}

func main() {
	c := config.New()
	fmt.Printf("%+v", c)

	startPublisher(c)
	startWorker(c)
	startServer()
}

func startPublisher(config config.Configuration) {
	log.Println("starting publisher...")

	params := publisher.Params{
		Endpoint: config.AWS.Endpoint,
		Region:   config.AWS.Region,
		TopicARN: config.AWS.TopicARN,
	}

	go publisher.Go(params)
}

func startWorker(config config.Configuration) {
	log.Println("starting worker...")

	params := wrkr.Params{
		Endpoint:  config.AWS.Endpoint,
		QueueName: config.AWS.QueueName,
		Region:    config.AWS.Region,
	}

	go wrkr.Go(params)
}

func startServer() {
	log.Println("starting server...")

	mux := srv.New()
	http.ListenAndServe(":8080", mux)

}

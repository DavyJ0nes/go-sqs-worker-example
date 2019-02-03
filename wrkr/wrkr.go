package wrkr

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	awsconfig "github.com/NYTimes/gizmo/config/aws"
	"github.com/NYTimes/gizmo/pubsub"
	"github.com/NYTimes/gizmo/pubsub/aws"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	messagesProcessed prometheus.Counter
	processingLatency prometheus.Histogram
	messagesInFlight  prometheus.Gauge
)

func init() {
	messagesProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "messages_processed_total",
		Help: "Total Count of messages processed from SQS",
	})

	processingLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "messages_processing_latency",
		Help:    "Time spent processing messages",
		Buckets: prometheus.DefBuckets,
	})

	messagesInFlight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "messages_in_flight",
		Help: "current amount of messages being processed",
	})

	prometheus.MustRegister(messagesProcessed)
	prometheus.MustRegister(processingLatency)
	prometheus.MustRegister(messagesInFlight)
}

type Params struct {
	Region    string
	Endpoint  string
	QueueName string
}

// Go initialises and runs a message worker
func Go(params Params) {
	log.Println("starting worker...")

	awsConfig := awsconfig.Config{
		Region:      params.Region,
		EndpointURL: &params.Endpoint,
	}

	queueURL := fmt.Sprintf("%s/%s", params.Endpoint, params.QueueName)

	config := aws.SQSConfig{
		Config:   awsConfig,
		QueueURL: queueURL,
	}

	log.Println(config)
	sub, err := aws.NewSubscriber(config)
	if err != nil {
		panic(err)
	}

	errc := make(chan error, 1)
	go handleError(errc)

	msgs := sub.Start()
	defer sub.Stop()

	p := Processor{}

	for msg := range msgs {
		go p.Process(msg)
	}
}

type Processor struct{}

func (p Processor) Process(msg pubsub.SubscriberMessage) {
	start := time.Now()
	messagesInFlight.Inc()
	messagesProcessed.Inc()

	delay()

	log.Println("---:", string(msg.Message()))
	msg.Done()

	messagesInFlight.Dec()
	processingLatency.Observe(time.Since(start).Seconds())
}

func handleError(errc <-chan error) {
	if err := <-errc; err != nil {
		panic(err)
	}
}

func delay() {
	i := rand.Intn(10000)
	time.Sleep(time.Duration(i) * time.Millisecond)
}

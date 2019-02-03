package publisher

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/NYTimes/gizmo/pubsub"
	"github.com/prometheus/client_golang/prometheus"

	awsconfig "github.com/NYTimes/gizmo/config/aws"
	"github.com/NYTimes/gizmo/pubsub/aws"
	"github.com/Pallinder/go-randomdata"
)

var (
	messagesPublished *prometheus.CounterVec
)

func init() {
	messagesPublished = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "messages_published_total",
		Help: "Total Count of messagesPublished published to SNS",
	}, []string{"topic"})

	prometheus.MustRegister(messagesPublished)
}

type Params struct {
	Region   string
	Endpoint string
	TopicARN string
}

// Go initialises and starts publishing messagesPublished to SNS
func Go(params Params) {
	awsConfig := awsconfig.Config{
		Region:      params.Region,
		EndpointURL: &params.Endpoint,
	}

	config := aws.SNSConfig{
		Config: awsConfig,
		Topic:  params.TopicARN,
	}

	pub, err := aws.NewPublisher(config)
	if err != nil {
		panic(err)
	}

	startPublishing(pub, params.TopicARN)
}

func startPublishing(pub pubsub.Publisher, topic string) {
	ctx := context.Background()
	key := "test"

	for {
		delay()

		body := randomdata.Paragraph()

		if err := pub.PublishRaw(ctx, key, []byte(body)); err != nil {
			log.Println(err)
		}

		messagesPublished.WithLabelValues(topic).Inc()
	}
}

func delay() {
	i := rand.Intn(1000)
	time.Sleep(time.Duration(i) * time.Millisecond)
}

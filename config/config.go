package config

import (
	"log"

	"github.com/spf13/viper"
)

// Configuration defines the apps config
type Configuration struct {
	AWS AWSConfiguration
}

// New creates new configuration
func New() Configuration {
	var configuration Configuration

	viper.AutomaticEnv()
	viper.BindEnv("aws.Region", "AWS_REGION")
	viper.BindEnv("aws.Endpoint", "AWS_ENDPOINT")
	viper.BindEnv("aws.AccountID", "AWS_ACCOUNT_ID")
	viper.BindEnv("aws.QueueName", "AWS_QUEUE_NAME")
	viper.BindEnv("aws.TopicARN", "AWS_TOPIC_ARN")

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return configuration
}

package common

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/rs/zerolog/log"
)

// EventMessaging is an abstraction for message buses used for decoupling components.
type EventMessaging interface {
	Publish(event *Event) error
}

// NewEventMessaging is a factory method for an `EventMessaging` based on the `MessageConfig` in the parameter
func NewEventMessaging(c MessagingConfig) (EventMessaging, error) {
	if strings.ToLower(c.Type) == snsMessagingType {
		return newSNSMessaging(c)
	}
	return newLoggerMessaging(c), nil
}

// snsMessaging is the AWS SNS implementation of `EventMessaging`.
type snsMessaging struct {
	client *sns.Client
	topic  string
}

// newSNSMessaging creates an `snsMessaging` from the config provided as parameter
func newSNSMessaging(c MessagingConfig) (*snsMessaging, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			c.AwsKey,
			c.AwsSecret, "")),
		config.WithRegion(c.AwsRegion),
	)
	if err != nil {
		return nil, err
	}
	return &snsMessaging{
		client: sns.NewFromConfig(awsConfig),
		topic:  c.SnsTopicArn,
	}, nil
}

// Publish creates a new SNS event from the parameter `event` and sends it on the topic of `snsMessaging`
func (s snsMessaging) Publish(event *Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = s.client.Publish(context.TODO(), &sns.PublishInput{
		TopicArn: aws.String(s.topic),
		Message:  aws.String(string(payload)),
	})
	return err
}

// loggerMessaging is a "NOOP" messaging implementation, writing debug logs
type loggerMessaging struct {
	topic string
}

// newLoggerMessaging creates a `loggerMessaging` from the config provided as parameter
func newLoggerMessaging(c MessagingConfig) *loggerMessaging {
	return &loggerMessaging{
		topic: c.SnsTopicArn,
	}
}

// Publish creates a new debug log entry from the parameter `event`.
func (s loggerMessaging) Publish(event *Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	log.Debug().RawJSON("payload", payload).Str("topic", s.topic).Msg("Emitted event")
	return nil
}

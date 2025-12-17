package consumer

import "github.com/IBM/sarama"

type Option func(config *sarama.Config)

func WithAllowAutoTopicCreation(allow bool) func(config *sarama.Config) {
	return func(config *sarama.Config) {
		config.Metadata.AllowAutoTopicCreation = allow
	}
}

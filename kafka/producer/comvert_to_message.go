package producer

import (
	"github.com/IBM/sarama"
)

func (i *impl) convertToMessage(topic string, data []byte, option MessageOption) (sarama.ProducerMessage, error) {
	// Validation
	if topic == "" {
		return sarama.ProducerMessage{}, ErrEmptyTopic
	}
	if len(data) == 0 {
		return sarama.ProducerMessage{}, ErrNoData
	}
	if option.Key == "" {
		option.Key = i.uuid.String()
	}

	// Create sarama producer message
	message := sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(option.Key),
		Value: sarama.ByteEncoder(data),
	}

	// Add headers to message
	if count := len(option.Headers); count > 0 {
		message.Headers = make([]sarama.RecordHeader, 0, count)
		for k, v := range option.Headers {
			message.Headers = append(message.Headers, sarama.RecordHeader{
				Key:   []byte(k),
				Value: []byte(v),
			})
		}
	}

	return message, nil
}

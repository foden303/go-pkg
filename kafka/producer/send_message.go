package producer

import (
	"context"
)

func (i *impl) SendMessage(ctx context.Context, topic string, data []byte, option MessageOption) (int32, int64, error) {
	// Create producer message
	message, err := i.convertToMessage(topic, data, option)
	if err != nil {
		return 0, 0, err
	}

	// Publish to Kafka
	partition, offset, err := i.syncProducer.SendMessage(&message)
	if err != nil {
		return 0, 0, err
	}
	return partition, offset, nil
}

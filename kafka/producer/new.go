package producer

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type impl struct {
	client       sarama.Client
	syncProducer sarama.SyncProducer
	uuid         uuid.UUID
}

func New(cfg Config, opts ...Option) Producer {
	// Default config
	samaraCfg := sarama.NewConfig()
	samaraCfg.Producer.Return.Successes = true // Mandatory setting

	// Override values
	samaraCfg.RackID = cfg.RackID
	samaraCfg.ClientID = cfg.ClientID

	// Options
	for _, opt := range opts {
		opt(samaraCfg)
	}

	// Create a client
	client, err := sarama.NewClient(cfg.Brokers, samaraCfg)
	if err != nil {
		log.Fatalf("unable to create kafka client. Err: %v", err)
	}

	// Create kafka producer
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		log.Fatalf("unable to create kafka producer client. Err: %v", err)
	}

	return &impl{
		client:       client,
		syncProducer: producer,
		uuid:         uuid.New(),
	}
}

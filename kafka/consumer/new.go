package consumer

import (
	"go-pkg/logger"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

type impl struct {
	client   sarama.Client
	consumer sarama.ConsumerGroup
	logger   logger.Logger
	topics   []string
	handler  sarama.ConsumerGroupHandler
	ready    chan bool
	wg       *sync.WaitGroup
}

func New(cfg Config, logger logger.Logger, opts ...Option) Consumer {
	// Default config
	samaraCfg := sarama.NewConfig()
	samaraCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	samaraCfg.Consumer.Return.Errors = true

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

	consumer, err := sarama.NewConsumerGroupFromClient(cfg.GroupID, client)
	if err != nil {
		log.Fatalf("unable to create kafka consumer group. Err: %v", err)
	}

	lg := logger.With(map[string]interface{}{
		"kafka.group.id":  cfg.GroupID,
		"kafka.client.id": cfg.ClientID,
	})
	wg := sync.WaitGroup{}
	return &impl{
		client:   client,
		consumer: consumer,
		logger:   lg,
		handler:  newHandler(lg, cfg.Handler, &wg),
		topics:   cfg.Topics,
		wg:       &wg,
		ready:    make(chan bool),
	}
}

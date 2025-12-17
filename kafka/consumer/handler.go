package consumer

import (
	"context"
	"fmt"
	"go-pkg/logger"
	"runtime/debug"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
)

type implHandler struct {
	logger  logger.Logger
	handler ConsumeHandler
	wg      *sync.WaitGroup
}

func newHandler(logger logger.Logger, handler ConsumeHandler, wg *sync.WaitGroup) sarama.ConsumerGroupHandler {
	return &implHandler{
		logger:  logger,
		handler: handler,
		wg:      wg,
	}
}

func (h *implHandler) Setup(session sarama.ConsumerGroupSession) error {
	lg := h.logger.With(map[string]interface{}{
		"consumer.member.id": session.MemberID(),
	})
	lg.Infof("Consumer Ready. GenerationID: [%d], Member ID: [%s], Partition Allocation: [%v]", session.GenerationID(), session.MemberID(), session.Claims())
	return nil
}

func (h *implHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	h.logger.Infof("Cleaning up...")
	return nil
}

func (h *implHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.consume(session, msg)
	}
	return nil
}

func (h *implHandler) consume(session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage) {
	h.logger.Infof("[START] Consume Handler: Topic[%s], Partition[%d]", message.Topic, message.Partition)

	// Recovery panic error
	start := time.Now()
	defer func(logger logger.Logger) {
		if rcv := recover(); rcv != nil {
			err := errors.WithStack(fmt.Errorf("panic err: %s", rcv))
			logger.Errorf(err, "Caught PANIC. Stack Trace: %s", debug.Stack())
		}
	}(h.logger)

	// Create a context message
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Read message key
	var messageKey string
	if message.Key != nil {
		messageKey = string(message.Key)
	}

	// Create a message
	msg := Message{
		ID: MessageID{
			Topic:     message.Topic,
			Partition: message.Partition,
			Offset:    message.Offset,
			Key:       messageKey,
		},
		Value:   message.Value,
		Headers: make(map[string]string, len(message.Headers)),
	}
	if message.Headers != nil {
		for _, r := range message.Headers {
			msg.Headers[string(r.Key[:])] = string(r.Value[:])
		}
	}

	// TODO: backoff retry
	if err := h.handler(ctx, msg); err != nil {
		h.logger.Errorf(err, "failed to handle message. Err: %v", err)
		return
	}

	// Commit message at
	h.commitMessageOffset(ctx, session, msg.ID)

	h.logger.Infof("[END] Consume Handler: Topic[%s], Partition[%d], Offset[%d], Duration[%dms]", message.Topic, message.Partition, message.Offset, time.Since(start).Milliseconds())
}

func (h *implHandler) commitMessageOffset(ctx context.Context, session sarama.ConsumerGroupSession, messageID MessageID) {
	offsetToCommit := messageID.Offset + 1 // Should always commit next offset as best practice
	session.MarkOffset(messageID.Topic, messageID.Partition, offsetToCommit, "")
}

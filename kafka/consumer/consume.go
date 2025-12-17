package consumer

import (
	"context"
	"fmt"

	pkgerrors "github.com/pkg/errors"
)

// Consume consumes messages in a loop
func (i *impl) Consume(ctx context.Context) error {
	consumeErr := make(chan error, 1)
	go func() {
		for {
			if ctx.Err() != nil {
				return
			}
			if err := i.consumer.Consume(ctx, i.topics, i.handler); err != nil {
				consumeErr <- pkgerrors.WithStack(fmt.Errorf("consuming failed. err: %w", err))
				return
			}
		}
	}()

	select {
	case err := <-consumeErr:
		return err
	case <-ctx.Done():
		i.logger.Infof("Closing consumer....")
		if err := i.consumer.Close(); err != nil {
			return pkgerrors.Wrap(err, "could not stop consumer")
		}
		if !i.client.Closed() {
			if err := i.client.Close(); err != nil {
				return pkgerrors.Wrap(err, "could not stop consumer client")
			}
		}
		i.logger.Infof("Consumer closed")
		return nil
	}
}

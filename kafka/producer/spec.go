package producer

import "context"

type Producer interface {
	Shutdown()
	SendMessage(ctx context.Context, topic string, data []byte, option MessageOption) (int32, int64, error)
}

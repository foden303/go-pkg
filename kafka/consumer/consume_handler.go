package consumer

import "context"

type ConsumeHandler func(ctx context.Context, message Message) error

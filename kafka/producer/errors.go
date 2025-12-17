package producer

import "errors"

var (
	ErrEmptyTopic = errors.New("topic is empty")
	ErrNoData     = errors.New("data isn't provided")
)

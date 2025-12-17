package consumer

type Message struct {
	ID      MessageID
	Value   []byte
	Headers map[string]string
}

type MessageID struct {
	Key       string
	Topic     string
	Partition int32
	Offset    int64
}

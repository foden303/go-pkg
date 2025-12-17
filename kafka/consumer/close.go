package consumer

func (i *impl) Close() error {
	return i.consumer.Close()
}

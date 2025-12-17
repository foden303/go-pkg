package producer

func (i *impl) Shutdown() {
	i.syncProducer.Close()
}

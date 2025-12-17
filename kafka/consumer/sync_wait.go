package consumer

func (i *impl) SyncWait() {
	i.wg.Wait()
}

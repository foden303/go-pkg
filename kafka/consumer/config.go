package consumer

type Config struct {
	GroupID  string
	RackID   string
	ClientID string
	Brokers  []string
	Topics   []string
	Handler  ConsumeHandler
}

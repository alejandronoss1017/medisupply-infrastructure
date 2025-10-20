package driven

// Publisher defines the interface for publishing messages
type Publisher interface {
	Publish(exchange, routingKey string, body []byte) error
	Close() error
}

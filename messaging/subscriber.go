package messaging

type Subscriber interface {
	Consume() (<-chan Message, error)
	Close() error
}

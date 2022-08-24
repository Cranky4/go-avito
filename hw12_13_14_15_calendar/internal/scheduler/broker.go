package internalscheduler

type Message struct {
	Text, Topic string
}

type Producer interface {
	// Connect(addr string) error
	Send(message Message) error
}

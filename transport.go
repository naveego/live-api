package live

type Transport interface {

	// Name should return the name of the transport
	Name() string

	// WriteMessage writes a message to the transport
	WriteMessage(message Message) error

	// ReadMessage reads a message from the transport
	ReadMessage() (Message, error)

	// Close closes the transport
	Close() error
}

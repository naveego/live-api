package live

import (
	"encoding/json"
	"fmt"
)

const (
	MessageTypeHello   = "HELLO"
	MessageTypePing    = "PING"
	MessageTypePong    = "PONG"
	MessageTypeMessage = "MESSAGE"
	MessageTypeGoodbye = "GOODBYE"
)

type MessageType string

type Message struct {
	Type          MessageType
	ContentType   string
	ContentLength int32
	Content       []byte
}

func NewJSONMessage(data interface{}) (Message, error) {

	var msg Message

	buf, err := json.Marshal(data)
	if err != nil {
		return msg, err
	}

	msg = Message{
		Type:          MessageTypeMessage,
		ContentType:   "application/json",
		ContentLength: int32(len(buf)),
		Content:       buf,
	}

	return msg, nil

}

func NewPingMessage() Message {
	return Message{
		Type:          MessageTypePing,
		ContentLength: int32(0),
	}
}

func NewPongMessage() Message {
	return Message{
		Type:          MessageTypePong,
		ContentLength: int32(0),
	}
}

func NewHelloMessage(clientID, host string) Message {

	msgContent := Hello{
		ClientID: clientID,
		Host:     host,
	}

	buf, _ := json.Marshal(msgContent)

	return Message{
		Type:          MessageTypeHello,
		ContentType:   "application/json",
		ContentLength: int32(len(buf)),
		Content:       buf,
	}
}

func NewGoodbyeMessage() Message {
	return Message{
		Type:          MessageTypeGoodbye,
		ContentLength: int32(0),
	}
}

func (msg Message) ReadJSON(ptr interface{}) error {
	if msg.ContentType != "application/json" {
		return fmt.Errorf("Message Error: Expected content type to be 'application/json' but was '%s'", msg.ContentType)
	}

	err := json.Unmarshal(msg.Content, ptr)
	if err != nil {
		return err
	}

	return nil
}

func hasContent(messageType MessageType) bool {
	return messageType == MessageTypeMessage || messageType == MessageTypeHello
}

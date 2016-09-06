package live

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

type wsTransport struct {
	wsConn   *websocket.Conn
	incoming chan *Message
}

func NewWebSocketTransport(conn *websocket.Conn) Transport {

	p := &wsTransport{
		wsConn:   conn,
		incoming: make(chan *Message),
	}

	conn.SetPingHandler(p.ping)
	conn.SetPongHandler(p.pong)

	go p.readLoop()

	return p
}

func (p *wsTransport) Name() string {
	return "ws"
}

func (p *wsTransport) WriteMessage(message Message) error {
	var err error
	switch message.Type {
	case MessageTypeMessage:
		if message.ContentType == "application/json" {
			err = p.wsConn.WriteMessage(websocket.TextMessage, message.Content)
		} else {
			err = p.wsConn.WriteMessage(websocket.BinaryMessage, message.Content)
		}
	case MessageTypeHello:
		err = p.wsConn.WriteMessage(websocket.TextMessage, message.Content)
	case MessageTypeGoodbye:
		err = p.wsConn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(5*time.Second))
	case MessageTypePing:
		err = p.wsConn.WriteMessage(websocket.PingMessage, nil)
	case MessageTypePong:
		err = p.wsConn.WriteMessage(websocket.PongMessage, nil)
	}

	return err
}

func (p *wsTransport) ReadMessage() (Message, error) {
	msg := <-p.incoming
	return *msg, nil
}

func (p *wsTransport) Close() error {
	return p.wsConn.Close()
}

func (p *wsTransport) readLoop() {

	defer func() {
		if err := recover(); err != nil {
		}
	}()

	for {
		var msg Message

		mt, message, err := p.wsConn.ReadMessage()
		if err != nil {

			if _, ok := err.(*websocket.CloseError); ok {
				msg.Type = MessageTypeGoodbye
			}

		}

		msg.ContentLength = int32(len(message))
		msg.Content = message

		switch mt {
		case websocket.TextMessage: // TextMessage
			msg.Type = MessageTypeMessage
			msg.ContentType = "application/json"
		case websocket.BinaryMessage: // BinaryMessage
			msg.Type = MessageTypeMessage
			msg.ContentType = "application/octet-stream"
		case websocket.PingMessage: // PingMessage
			logrus.Debug("PING")
			msg.Type = MessageTypePing
		case websocket.PongMessage: // PongMessage
			logrus.Debug("PONG")
			msg.Type = MessageTypePong

		}

		p.incoming <- &msg
	}
}

func (p *wsTransport) ping(data string) error {
	pingMsg := NewPingMessage()
	p.incoming <- &pingMsg
	return nil
}

func (p *wsTransport) pong(data string) error {
	pongMsg := NewPongMessage()
	p.incoming <- &pongMsg
	return nil
}

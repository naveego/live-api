package live

import (
	"github.com/gorilla/websocket"
)

func NewWebSocketClient(addr, id, host string) (*Client, error) {
	wsConn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}

	transport := NewWebSocketTransport(wsConn)
	cli, err := NewClientWithTransport(transport, id, host)
	return cli, err
}

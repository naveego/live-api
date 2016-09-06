package live

import (
	"net"
	"time"

	"github.com/Sirupsen/logrus"
)

type Client struct {
	transport Transport
	ticker    *time.Ticker
	incoming  chan Message
	errors    chan error
	closing   chan struct{}
}

func NewTCPClient(addr, id, host string) (*Client, error) {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return &Client{}, err
	}

	transport := NewTCPTransport(conn)
	cli, err := NewClientWithTransport(transport, id, host)
	return cli, err
}

func NewClientWithTransport(transport Transport, id, host string) (*Client, error) {

	cli := &Client{
		transport: transport,
		ticker:    time.NewTicker(5 * time.Second),
		incoming:  make(chan Message),
		errors:    make(chan error),
		closing:   make(chan struct{}),
	}

	helloMsg := NewHelloMessage(id, host)
	err := transport.WriteMessage(helloMsg)
	if err != nil {
		return cli, err
	}

	go withRecover(cli.read)
	go withRecover(cli.heartbeat)

	return cli, nil

}

func (cli *Client) Incoming() <-chan Message {
	return cli.incoming
}

func (cli *Client) Errors() <-chan error {
	return cli.errors
}

func (cli *Client) Close() {
	close(cli.closing)
	cli.transport.Close()
}

func (cli *Client) read() {
	for {
		msg, err := cli.transport.ReadMessage()
		if err != nil {
			cli.errors <- err
		}

		switch msg.Type {
		case MessageTypePing:
			logrus.Debug("PING")
		case MessageTypePong:
			logrus.Debug("PONG")
		default:
			cli.incoming <- msg
		}
	}
}

func (cli *Client) heartbeat() {
	for {
		select {
		case <-cli.ticker.C:
			cli.sendPing()
		case <-cli.closing:
			cli.ticker.Stop()
			return
		}
	}
}

func (cli *Client) sendPing() {
	pingMsg := NewPingMessage()
	cli.transport.WriteMessage(pingMsg)
}

func withRecover(f func()) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warn(err)
		}
	}()

	f()
}

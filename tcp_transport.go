package live

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
)

type tcpTransport struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

// NewTCPTransport creates a new tcp Transport for sending/receiving
// messages.
func NewTCPTransport(connection net.Conn) Transport {
	return &tcpTransport{
		conn:   connection,
		reader: bufio.NewReader(connection),
		writer: bufio.NewWriter(connection),
	}
}

func (p *tcpTransport) Name() string {
	return "tcp"
}

func (p *tcpTransport) WriteMessage(message Message) error {
	err := p.writeType(message.Type)
	if err != nil {
		return err
	}

	if !hasContent(message.Type) {
		p.writer.Flush()
		return nil
	}

	err = p.writeContentType(message.ContentType)
	if err != nil {
		return err
	}

	err = p.writeContent(message)
	if err != nil {
		return err
	}

	return p.writer.Flush()
}

func (p *tcpTransport) ReadMessage() (msg Message, err error) {

	msg.Type, err = p.readType()
	if err != nil {
		return
	}

	if !hasContent(msg.Type) {
		return
	}

	msg.ContentType, err = p.readContentType()
	if err != nil {
		return
	}

	err = p.readContent(&msg)
	return

}

func (p *tcpTransport) Close() error {
	p.conn.Close()
	return nil
}

func (p *tcpTransport) readType() (MessageType, error) {
	var t uint16
	err := binary.Read(p.reader, binary.BigEndian, &t)
	if err != nil {
		return "", nil
	}

	switch t {
	case 0:
		return MessageTypePing, nil
	case 1:
		return MessageTypePong, nil
	case 2:
		return MessageTypeHello, nil
	case 3:
		return MessageTypeGoodbye, nil
	case 10:
		return MessageTypeMessage, nil
	default:
		return "", fmt.Errorf("Protocol Error: Unknown message type")
	}
}

func (p *tcpTransport) readContentType() (string, error) {
	var ctLen uint16
	err := binary.Read(p.reader, binary.BigEndian, &ctLen)
	if err != nil {
		return "", err
	}

	buf := make([]byte, ctLen)
	_, err = p.reader.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func (p *tcpTransport) readContent(message *Message) error {
	var cLen int32
	err := binary.Read(p.reader, binary.BigEndian, &cLen)
	if err != nil {
		return err
	}

	message.Content = make([]byte, cLen)
	_, err = p.reader.Read(message.Content)
	if err != nil {
		return err
	}

	return nil
}

func (p *tcpTransport) writeType(t MessageType) (err error) {
	switch t {
	case MessageTypePing:
		err = binary.Write(p.writer, binary.BigEndian, uint16(0))
	case MessageTypePong:
		err = binary.Write(p.writer, binary.BigEndian, uint16(1))
	case MessageTypeHello:
		err = binary.Write(p.writer, binary.BigEndian, uint16(2))
	case MessageTypeGoodbye:
		err = binary.Write(p.writer, binary.BigEndian, uint16(3))
	case MessageTypeMessage:
		err = binary.Write(p.writer, binary.BigEndian, uint16(10))
	default:
		err = fmt.Errorf("Protocol Error: Unknown message type '%s'", t)
	}

	return
}

func (p *tcpTransport) writeContentType(cType string) (err error) {
	buf := []byte(cType)
	bufLen := len(buf)

	err = binary.Write(p.writer, binary.BigEndian, uint16(bufLen))
	if err != nil {
		return
	}

	_, err = p.writer.Write(buf)
	return err
}

func (p *tcpTransport) writeContent(message Message) (err error) {

	err = binary.Write(p.writer, binary.BigEndian, message.ContentLength)
	if err != nil {
		return
	}

	_, err = p.writer.Write(message.Content)
	return
}

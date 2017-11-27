package main

import (
	"fmt"
	"net"
)

type Connection struct {
	ip         net.IP
	port       uint16
	conn       net.TCPConn
	infoHash   string
	peerId     string
	choke      bool
	interested bool
}

func (c *Connection) Handshake() {
	zeroBytes := []byte("\x00\x00\x00\x00\x00\x00\x00\x00")

	var hello []byte
	hello = append(hello, "\x13BitTorrent protocol"...)
	hello = append(hello, zeroBytes...)
	hello = append(hello, []byte(c.infoHash)...)
	hello = append(hello, c.peerId...)

	c.Send(hello)
}

func (c *Connection) Bitfield() string {
	return ""
}

func (c *Connection) Choke() string {
	return ""
}

func (c *Connection) Unchoke() string {
	return ""
}

func (c *Connection) Connect() {
	addr := net.TCPAddr{
		IP:   c.ip,
		Port: int(c.port),
	}
	laddr := net.TCPAddr{
		Port: 50005,
	}
	conn, err := net.DialTCP("tcp", &laddr, &addr)
	check(err)
	c.conn = *conn
}

func (c *Connection) Send(message []byte) {
	_, err := c.conn.Write(message)
	if err != nil {
		println("Connection failed:", err.Error())
	}
}

func (c *Connection) Receive() []byte {
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := c.conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	return buf
}

func (c *Connection) Listen(ch chan []byte) []byte {
	for {
		response := c.Receive()
		ch <- response
	}
}

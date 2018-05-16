package main

import (
	"fmt"
	"github.com/libp2p/go-reuseport"
	"net"
)

type Connection struct {
	ip         net.IP
	port       uint16
	conn       net.Conn
	infoHash   string
	peerId     string
	choke      bool
	interested bool
	laddr      net.TCPAddr
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
	conn, err := reuseport.Dial("tcp", "192.168.1.15:50005", c.ip.String()+":"+fmt.Sprint(c.port))
	check(err)
	c.conn = conn
}

func (c *Connection) Send(message []byte) {
	fmt.Println("Sending something to " + c.Ip().String())
	fmt.Println(message)
	_, err := c.conn.Write(message)
	if err != nil {
		println("Connection failed:", err.Error())
	}
}

func (c *Connection) Receive() []byte {
	fmt.Println("Receiving something from " + c.Ip().String())
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := c.conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	return buf
}

func (c *Connection) Close() {
	c.conn.Close()
}

func (c *Connection) Ip() net.IP {
	return c.ip
}

func (c *Connection) Port() uint16 {
	return c.port
}

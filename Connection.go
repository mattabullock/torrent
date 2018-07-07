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
	peerHave   []bool
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

func (c *Connection) Choke() {
	c.Open()
	c.choke = true
}

func (c *Connection) Unchoke() {
	c.choke = false
}

func (c *Connection) Interested() {
	c.interested = true
}

func (c *Connection) Uninterested() {
	c.interested = false
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

func (c *Connection) handshake() {
	c.Handshake()
	c.Receive()
	// check correct handshake received
}

func (c *Connection) bitfield(f File) {
	bitfield := []byte{5}
	for i := 0; i < f.numPieces; i += 8 {
		currentByte := 0
		for j := 8; j > 0; j-- {
			if f.havePieces[i] {
				currentByte *= "1" << j
			}
		}
		bitfield[i/8+1] = currentByte
	}

	c.send(bitfield)
}

func (c *Connection) handleRequest(message []byte) {
	switch message[0] {
	case "\x01":
		c.Choke()
		break
	case "\x02":
		c.Unchoke()
		break
	case "\x03":
		c.Interested()
		break
	case "\x04":
		c.Uninterested()
		break
	case "\x05":
		break
	default:
		break
	}
}

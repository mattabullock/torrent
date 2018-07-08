package main

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"github.com/libp2p/go-reuseport"
	"github.com/mattabullock/bencode-go"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("USAGE: torrent <file>")
		return
	}

	args := os.Args[1:]

	// Get data from metadata file
	data, err := ioutil.ReadFile(args[0])
	check(err)

	ann, file := ReadMetadata(data)
	infoHash := sha1sum(data[239 : len(data)-1])
	peerId := "-TR2840-0p8s3d54k2co"
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	trackerURL := ann.GenerateTrackerURL()

	for {
		resp, err := netClient.Get(trackerURL)

		check(err)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err.Error())
		}

		trackerResponse := bencode.Decode(body).(map[string]interface{})
		if val, ok := trackerResponse["failure reason"]; ok {
			panic(val.(string))
		}

		peers := []byte(trackerResponse["peers"].(string))
		connections := make(map[string]Connection)
		//go Listen()

		for i := 0; i < len(peers); i += 6 {
			ip := net.IPv4(peers[i], peers[i+1], peers[i+2], peers[i+3])
			port := binary.BigEndian.Uint16([]byte{peers[i+4], peers[i+5]})

			if _, ok := connections[ip.String()]; ok {
				continue
			}

			conn := Connection{
				ip:         ip,
				port:       port,
				infoHash:   infoHash,
				peerId:     peerId,
				choke:      true,
				interested: false,
			}

			connections[ip.String()] = conn

			//go handleConnection(conn)
			handleConnection(conn, file)
		}
		time.Sleep(60 * time.Second)
	}

	for {

	}
}

func handleConnection(conn Connection, file File) {
	conn.Connect()
	defer conn.Close()
	conn.Handshake()
	conn.Bitfield(file)
	for {
		message := conn.Receive()
		conn.handleRequest(message)
	}
}

func sha1sum(s []byte) string {
	h := sha1.New()
	h.Write(s)
	sha1sum := h.Sum(nil)
	return string(sha1sum)
}

func ReadMetadata(data []byte) (Announce, File) {
	infoHash := url.QueryEscape(sha1sum(data[239 : len(data)-1]))
	metadata := bencode.Decode(data).(map[string]interface{})
	announceURL, ok := metadata["announce"].(string)
	if !ok {
		panic(ok)
	}

	info := metadata["info"].(map[string]interface{})
	pieces := info["pieces"].(string)
	pieceLength := info["piece length"].(uint64)
	length := info["length"].(uint64)
	numPieces := length / pieceLength

	file := File{
		length:      length,
		pieceLength: pieceLength,
		numPieces:   numPieces,
		pieces:      pieces,
		havePieces:  make([]bool, numPieces),
	}

	ann := Announce{
		url:           announceURL,
		infoHash:      infoHash,
		peerId:        url.QueryEscape("-TR2840-0p8s3d54k2co"),
		port:          "50005",
		uploaded:      0,
		downloaded:    0,
		left:          699400192,
		numwant:       80,
		key:           "5c179003",
		compact:       "1",
		supportcrypto: true,
		event:         "started",
	}

	return ann, file
}

func Listen() []byte {
	l, err := reuseport.Listen("tcp", "192.168.1.15:50005")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	Log("Listening on " + l.Addr().String())
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		Log("localAddr: " + conn.LocalAddr().String() + "->" + conn.RemoteAddr().String())
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	// Close the connection when you're done with it.
	conn.Close()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Log(message string) {
	currentTime := time.Now()
	message = currentTime.Format(time.RFC3339Nano) + " - " + message
	fmt.Println(message)
}

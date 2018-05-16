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

	ann := ReadMetadata(data)
	infoHash := sha1sum(data[239 : len(data)-1])
	peerId := "-TR2840-0p8s3d54k2co"
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	trackerURL := ann.GenerateTrackerURL()
	resp, err := netClient.Get(trackerURL)

	check(err)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	trackerResponse := bencode.Decode(body).(map[string]interface{})
	fmt.Println([]byte(trackerResponse["peers"].(string)))
	if val, ok := trackerResponse["failure reason"]; ok {
		panic(val.(string))
	}

	peers := []byte(trackerResponse["peers"].(string))
	go Listen()

	for i := 0; i < len(peers); i += 6 {
		ip := net.IPv4(peers[i], peers[i+1], peers[i+2], peers[i+3])
		port := binary.BigEndian.Uint16([]byte{peers[i+4], peers[i+5]})

		conn := Connection{
			ip:         ip,
			port:       port,
			infoHash:   infoHash,
			peerId:     peerId,
			choke:      true,
			interested: false,
		}

		conn.Connect()
		conn.Handshake()
		message := conn.Receive()
		fmt.Println(message)
		conn.Close()
		//ch := make(chan []byte)
		//go handleResponse(ch)
	}
}

func handleResponse(ch chan []byte) {
	for i := range ch {
		fmt.Printf("%x\n", i)
	}
}

func sha1sum(s []byte) string {
	h := sha1.New()
	h.Write(s)
	sha1sum := h.Sum(nil)
	return string(sha1sum)
}

func ReadMetadata(data []byte) Announce {
	infoHash := url.QueryEscape(sha1sum(data[239 : len(data)-1]))
	metadata := bencode.Decode(data).(map[string]interface{})
	announceURL, ok := metadata["announce"].(string)
	if !ok {
		panic(ok)
	}

	var have map[string]bool
	have = make(map[string]bool)
	pieces := metadata["info"].(map[string]interface{})["pieces"].(string)
	for i := 0; i < len(pieces); i += 20 {
		have[pieces[i:i+19]] = false
	}

	//info := bencode.Encode(metadata["info"].(map[string]interface{}))
	//fmt.Println(string(info))
	//infoHash := url.PathEscape(sha1sum(info))
	//fmt.Println(infoHash)

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

	return ann
}

func Listen() []byte {
	l, err := reuseport.Listen("tcp", "192.168.1.15:50005")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + l.Addr().String())
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		fmt.Println("localAddr: " + conn.LocalAddr().String() + "->" + conn.RemoteAddr().String())
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
	// Send a response back to person contacting us.
	fmt.Println("request!")
	// Close the connection when you're done with it.
	conn.Close()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

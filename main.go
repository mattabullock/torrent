package main

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"github.com/mattabullock/bencode-go"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

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

	//info := bencode.Encode(metadata["info"].(map[string]interface{}))
	//fmt.Println(string(info))
	//infoHash := url.PathEscape(sha1sum(info))
	//fmt.Println(infoHash)

	ann := Announce{
		url:           announceURL,
		infoHash:      infoHash,
		peerId:        "-TR2840-0p8s3d54k2co",
		port:          "56026",
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

func GenerateHandshake(data []byte) []byte {
	infoHash := sha1sum(data[239 : len(data)-1])
	zeroBytes := []byte("\x00\x00\x00\x00\x00\x00\x00\x00")
	peerId := []byte("-TR2840-0p8s3d54k2co")
	var hello []byte
	hello = append(hello, "\x13BitTorrent protocol"...)
	hello = append(hello, zeroBytes...)
	hello = append(hello, []byte(infoHash)...)
	hello = append(hello, peerId...)

	return hello
}

func main() {

	args := os.Args[1:]

	// Get data from metadata file
	data, err := ioutil.ReadFile(args[0])
	check(err)

	ann := ReadMetadata(data)

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
	if val, ok := trackerResponse["failure reason"]; ok {
		panic(val.(string))
	}

	peers := []byte(trackerResponse["peers"].(string))

	//for i := 0; i < len(peers); i += 6 {
	i := 0
	ip := net.IPv4(peers[i], peers[i+1], peers[i+2], peers[i+3])
	port := binary.BigEndian.Uint16([]byte{peers[i+4], peers[i+5]})
	fmt.Println(ip)
	addr := net.TCPAddr{
		IP:   ip,
		Port: int(port),
	}
	laddr := net.TCPAddr{
		Port: 50005,
	}

	conn, err := net.DialTCP("tcp", &laddr, &addr)
	check(err)
	hello := GenerateHandshake(data)
	_, err = conn.Write(hello)
	if err != nil {
		println("Write to server failed:", err.Error())
	}
	result, err := ioutil.ReadAll(conn)
	check(err)

	fmt.Println(string(result))
	//}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

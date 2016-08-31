package main

import (
	"bencode-go"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func main() {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	ann := Announce{
		url:           "http://torrent.ubuntu.com:6969/announce",
		infoHash:      "%90%28%9f%d3M%fc%1c%f8%f3%16%a2h%ad%d85L%853DX",
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

	trackerURL := ann.GenerateTrackerURL()
	fmt.Println(trackerURL)
	resp, err := netClient.Get(trackerURL)

	check(err)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	stuff := bencode.Decode(body)
	cast := stuff.(map[string]interface{})
	peers := []byte(cast["peers"].(string))

	fmt.Println(len(peers))
	for i := 0; i < len(peers); i += 6 {
		//fmt.Println(peers[i])
		ip := net.IPv4(peers[i], peers[i+1], peers[i+2], peers[i+3])
		port := binary.BigEndian.Uint16([]byte{peers[i+4], peers[i+5]})
		fmt.Println(ip, ":", port)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

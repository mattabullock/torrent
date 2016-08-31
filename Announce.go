package main

import (
	"bytes"
	"strconv"
)

type Announce struct {
	url           string
	infoHash      string
	peerId        string
	port          string
	uploaded      uint64
	downloaded    uint64
	left          uint64
	numwant       uint64
	key           string
	compact       string
	supportcrypto bool
	event         string
}

func (a *Announce) GenerateTrackerURL() string {
	baseURL := a.Url() + "?"
	fullURL := bytes.NewBufferString(baseURL)

	fullURL.WriteString("info_hash=" + a.InfoHash())
	fullURL.WriteString("&peer_id=" + a.PeerId())
	fullURL.WriteString("&port=" + a.Port())
	fullURL.WriteString("&uploaded=" + string(a.Uploaded()))
	fullURL.WriteString("&downloaded=" + string(a.Downloaded()))
	fullURL.WriteString("&left=" + string(a.Left()))
	fullURL.WriteString("&numwant=" + string(a.Numwant()))
	fullURL.WriteString("&key=" + a.Key())
	fullURL.WriteString("&compact=" + a.Compact())
	fullURL.WriteString("&supportcrypto=" + a.Supportcrypto())
	fullURL.WriteString("&event=" + a.Event())

	return fullURL.String()
}

func (a *Announce) Url() string {
	return a.url
}

func (a *Announce) InfoHash() string {
	return a.infoHash
}

func (a *Announce) PeerId() string {
	return a.peerId
}

func (a *Announce) Port() string {
	return a.port
}

func (a *Announce) Uploaded() string {
	return strconv.FormatUint(a.uploaded, 10)
}

func (a *Announce) Downloaded() string {
	return strconv.FormatUint(a.downloaded, 10)
}

func (a *Announce) Left() string {
	return strconv.FormatUint(a.left, 10)
}

func (a *Announce) Numwant() string {
	return strconv.FormatUint(a.numwant, 10)
}

func (a *Announce) Key() string {
	return a.key
}

func (a *Announce) Compact() string {
	return a.compact
}

func (a *Announce) Supportcrypto() string {
	if a.supportcrypto {
		return "1"
	}
	return "0"
}

func (a *Announce) Event() string {
	return a.event
}

package main

type File struct {
	length          uint64
	pieceLength     uint64
	numPieces       uint64
	pieces          string
	data            string
	havePieces      []bool
	requestedPieces []bool
}

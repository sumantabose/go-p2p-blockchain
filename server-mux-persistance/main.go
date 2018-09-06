package main

import (
	"log"
	"time"
)

func init() { // Idea from https://appliedgo.net/networking/
	log.SetFlags(log.Lshortfile)
	readFlags() // in defs.go
}

func main() {
	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{0, t.String(), 0, calculateHash(genesisBlock), ""}

	Blockchain = append(Blockchain, genesisBlock)

	p2pInit()

	log.Fatal(muxServer(*listenF)) // function is in mux.go
}


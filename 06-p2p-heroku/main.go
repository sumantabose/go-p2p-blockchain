package main

import (
	"log"
)

func init() { // Idea from https://appliedgo.net/networking/
	log.SetFlags(log.Lshortfile)
	readFlags() // in defs.go
	registerGOB() // in blockchain.go
}

func main() {
	// p2pInit() // Initialize P2P Network from Bootstrapper
	log.Fatal(muxServer()) // function is in mux.go
}


package main

import (
	"log"
	"sync"
	"flag"
	"math/rand"
	"time"
	gonet "net"
	host "github.com/libp2p/go-libp2p-host"
)

///// FLAG & VARIABLES

	//var listenF *int
	//var target *string
	var secio *bool
	var verbose *bool
	var seed *int64

	var ha host.Host

/////

// Block represents each 'item' in the blockchain
type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
}

// Blockchain is a series of validated Blocks
var Blockchain []Block

// Message takes incoming JSON payload for writing heart rate
type Message struct {
	BPM int
}

type newTarget_json struct {
	NewTarget string
}

var mutex = &sync.Mutex{}



////////  HELPER FUNCTIONS

func readFlags() {
	// Parse options from the command line
	//listenF = flag.Int("l", 0, "wait for incoming connections")
	//target = flag.String("d", "", "target peer to dial")
	secio = flag.Bool("secio", true, "enable secio")	
	verbose = flag.Bool("v", false, "enable verbose")
	seed = flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()
}

func  GetMyIP() string {
	var MyIP string

	conn, err := gonet.Dial("udp", "8.8.8.8:80")
	if err != nil {
	   log.Fatalln(err)
	} else {
		localAddr := conn.LocalAddr().(*gonet.UDPAddr)
		MyIP = localAddr.IP.String()
	}
	return MyIP
}

func genRandInt(n int) int {
    myRandSource := rand.NewSource(time.Now().UnixNano())
   	myRand := rand.New(myRandSource)
   	val := myRand.Intn(n)
   	return val
}
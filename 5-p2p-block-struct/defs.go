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

	var secio *bool
	var verbose *bool
	var seed *int64
	var dataDir *string // data directory prefix where the gob files are stored

	var ha host.Host

/////

// Raw Material Transaction (Type 1)
type RawMaterialTransaction struct {
	Name string
}

// Delivery Transaction (Type 2)
type DeliveryTransaction struct {
	Age int
}

// Block represents each 'item' in the blockchain
type Block struct {
	Index     int
	Timestamp string
	TxnType int
	TxnPayload interface{}
	Comment	string
	Proposer string
	PrevHash  string
	Hash	string
}

// Blockchain is a series of validated Blocks
var Blockchain []Block

// Standard Input takes incoming string from terminal for `comment` in block (Type 0)
type StdInput struct {
	Comment string
}

var mutex = &sync.Mutex{}

////////  HELPER FUNCTIONS

func readFlags() {
	// Parse options from the command line
	secio = flag.Bool("secio", true, "enable secio")	
	verbose = flag.Bool("v", false, "enable verbose")
	seed = flag.Int64("seed", 0, "set random seed for id generation")
	dataDir = flag.String("data", "data", "pathname of data directory")
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
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
	var bootstrapperAddr *string

	var ha host.Host

	var PeerStart = false

/////

// Raw Material Transaction (Type 1)
type RawMaterialTransaction struct {
	SerialNo int
	ProductCode string
	ProductName string
	ProductBatchNo string
	Quantity int
	RawMaterial RawMaterial
}
type RawMaterial struct {
	RawMaterialBatchNo string
	RawMaterialsID string
	RawMaterialName string
	RawMaterialQuantity float32
	RawMaterialMeasurementUnit string
}

// Delivery Transaction (Type 2)
type DeliveryTransaction struct {
	SerialNo int
	RecGenerator string
	ShipmentID string
	Timestamp string	
	Longitude string
	Latitude string
	ShippedFromCompID string
	ShippedToCompID string
	LocationID string
	DeliveryStatus string 	
	DeliveryType string
	Product Product
	Document Document
}
type Product struct {
	ProductCode string
	ProductName string
	ProductBatch ProductBatch
}
type ProductBatch struct {
	ProductBatchNo string
	ProductBatchQuantity int
}
type Document struct {
	DocumentURL string
	DocumentType string
	DocumentHash string
	DocumentSign string
}

// Block represents each 'item' in the blockchain
type Block struct {
	Index		int
	Timestamp	string
	TxnType		int // 0 for StdInput, 1 for RawMaterialTransaction, 2 for DeliveryTransaction
	TxnPayload	interface{}
	Comment		string
	Proposer	string
	PrevHash 	string
	ThisHash	string
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
	verbose = flag.Bool("v", true, "enable verbose")
	seed = flag.Int64("seed", 0, "set random seed for id generation")
	dataDir = flag.String("data", "data", "pathname of data directory")
	bootstrapperAddr = flag.String("b", "heroku", "Address of bootstrapper")
	flag.Parse()

	if *bootstrapperAddr == "heroku" {
		*bootstrapperAddr = "https://ntu-blockchain-bootstrapper.herokuapp.com/"
	} else if *bootstrapperAddr == "local" {
		*bootstrapperAddr = "http://localhost:" + bootstrapperPort + "/"
	} else {
		*bootstrapperAddr = "http://" + *bootstrapperAddr + ":" + bootstrapperPort + "/"
	}
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
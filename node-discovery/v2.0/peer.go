/* README

Written by Sumanta Bose, 11 Sept 2018

*/

package main

import (
	"os"
    "log"
    "sync"
    "time"
    "flag"
    "bytes"
    "syscall"
    "net/http"
    "math/rand"
    "os/signal"
    "io/ioutil"
    "encoding/json"

    "github.com/davecgh/go-spew/spew"
)

///// GLOBAL CONSTS & VARIABLES 

const (
    bootstrapperPort = "51000"
)

type Peer struct {
    PeerAddress string `json:"PeerAddress"`
}

type PeerProfile struct { // connections of one peer
    ThisPeer Peer `json:"ThisPeer"` // any node
    Neighbors []Peer `json:"Neighbors"` // edges to that node
    Status bool `json:"Status"` // Status: Alive or Dead
    Connected bool `json:"Connected"` // If a node is connected or not [To be used later]
}

var peerProfile PeerProfile // used to enroll THIS peer | connectP2PNet() & enrollP2PNet)()
var PeerGraph = make(map[string]PeerProfile) // Key = Node.PeerAddress; Value.Neighbors = Edges
var graphMutex sync.RWMutex
var verbose *bool

///// LIST OF FUNCTIONS

func init() {
	log.SetFlags(log.Lshortfile)
    go func() {
    	signalChan := make(chan os.Signal)
    	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
        <- signalChan
        log.Println("Received Interrupt. Exiting now.")
        //ExitProtocol()
        os.Exit(1)
    }()
    verbose = flag.Bool("v", false, "enable verbose")
    flag.Parse()
}

func main() {
	queryP2PList()
	connectP2PNet()
	enrollP2PNet()
	queryP2PList()
	select{}
}

func queryP2PList() { // Query the list of peers in the P2P Network from the Bootstrapper
	log.Println("Querying list of peers")
	
	response, err := http.Get("http://localhost:" + bootstrapperPort + "/query-p2p-list")
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}

	graphMutex.RLock()
		json.Unmarshal(responseData, &PeerGraph)
		if *verbose { log.Println("PeerGraph = ", PeerGraph) ; spew.Dump(PeerGraph) }
	graphMutex.RUnlock()
}

func enrollP2PNet() { // Enroll to the P2P Network by adding THIS peer with Bootstrapper
	log.Println("Enrolling in P2P network")

	jsonValue, err := json.Marshal(peerProfile)
	if err != nil {
		log.Println(err)
		return
	}

	url := "http://localhost:" + bootstrapperPort + "/enroll-p2p-net"
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()

    if *verbose { log.Println("response Status:", response.Status) }
    if *verbose { log.Println("response Headers:", response.Header) }
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Println(err)
        return
    }
    log.Println("response Body:", string(body))
}

///// LIST OF HELPER FUNCTIONS

func connectP2PNet() {
	peerProfile.ThisPeer = Peer {PeerAddress : genRandString(15)}

	if len(PeerGraph) == 0 { // first node in the network
		// do nothing
	} else {
		log.Println("Connecting to P2P network")

		if *verbose { log.Println("len(PeerGraph) = ", len(PeerGraph)) }
		// make connection with peers[choice]

		choice := genRandInt(len(PeerGraph))
		log.Println("choice = ", choice)
		peers := make([]string, 0, len(PeerGraph))
		for p, _ := range PeerGraph {
		    peers = append(peers, p)
		}
		peerProfile.Neighbors = append(peerProfile.Neighbors, Peer {PeerAddress : peers[choice]})
		if *verbose { log.Println("peers[choice] = ", peers[choice]) }
	}
}

///// LIST OF MISC FUNCTIONS

func genRandString(n int) string { // generate Random String of length 'n'
    var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    val := make([]rune, n)
    for i := range val {
        myRandSource := rand.NewSource(time.Now().UnixNano())
        myRand := rand.New(myRandSource)
        val[i] = letterRunes[myRand.Intn(len(letterRunes))]
    }
    return string(val)
}

func genRandInt(n int) int {
    myRandSource := rand.NewSource(time.Now().UnixNano())
   	myRand := rand.New(myRandSource)
   	val := myRand.Intn(n)
   	return val
}
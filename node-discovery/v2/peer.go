/* README

Written by Sumanta Bose, 7 Sept 2018

*/

package main

import (
	"os"
    "log"
    "sync"
    "time"
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

type PeerLinks struct { // connections of one peer
    peer Peer `json:"peer"` // any node
    links []Peer `json:"links"` // edges to that node
}

type PeerGraph struct { // undirected p2p graph
    peers []*Peer `json:"peers"` // nodes
    links map[Peer][]*Peer `json:"links"` // edges
}

var peerGraph PeerGraph
var thisPeer PeerLinks
var graphMutex sync.RWMutex

/////

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
}

func main() {
	queryP2PList()
	connectP2PNet()
	joinP2PNet()
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

	graphMutex.Lock()
		json.Unmarshal(responseData, &peerGraph)
	graphMutex.Unlock()
	graphMutex.RLock()
		spew.Dump(peerGraph)
	graphMutex.RUnlock()
}

func joinP2PNet() { // Join the P2P Network by adding THIS peer's address to list of peers with Bootstrapper
	log.Println("Joining P2P network")

	log.Println(thisPeer)

	jsonValue, err := json.Marshal(thisPeer)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(jsonValue)

	url := "http://localhost:" + bootstrapperPort + "/join-p2p-net"
	// response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// defer response.Body.Close()

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    response, err := client.Do(req)
    if err != nil {
        log.Println(err)
        return
    }
    defer response.Body.Close()

    log.Println("response Status:", response.Status)
    log.Println("response Headers:", response.Header)
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Println(err)
        return
    }
    log.Println("response Body:", string(body))
}

///// LIST OF HELPER FUNCTIONS

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

/////


func connectP2PNet() {

	thisPeer.peer = Peer {PeerAddress : genRandString(15)}
	if len(peerGraph.peers) == 0 { // first node
		thisPeer.links = append(thisPeer.links, Peer {PeerAddress : "NULL"})
	} else {
		choice := genRandInt(len(peerGraph.peers))
		thisPeer.links = append(thisPeer.links, *peerGraph.peers[choice])
	}

	log.Println(thisPeer)
}








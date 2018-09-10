/* README

Written by Sumanta Bose, 10 Sept 2018

*/

package main

import (
	"os"
    "log"
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

// var thisPeer Peer

type PeerProfile struct {
    ThisPeer Peer `json:"ThisPeer"`
    Status bool `json:"Status"`
}

var peer PeerProfile
var ListOfPeers []Peer

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
	peer = PeerProfile {
		ThisPeer : Peer {PeerAddress : genRandString(15)},
		Status : true,
	}
	log.Println(peer)
	queryP2PList()
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

	json.Unmarshal(responseData, &ListOfPeers)
	spew.Dump(ListOfPeers)
}

func joinP2PNet() { // Join the P2P Network by adding THIS peer's address to list of peers with Bootstrapper
	log.Println("Joining P2P network")

	jsonValue, err := json.Marshal(peer)
	if err != nil {
		log.Println(err)
		return
	}

	url := "http://localhost:" + bootstrapperPort + "/join-p2p-net"
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
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
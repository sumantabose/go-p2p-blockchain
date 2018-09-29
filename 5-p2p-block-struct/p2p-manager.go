package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"os"
	"time"
	"bytes"
	"os/signal"
    "syscall"
    "sync"
	"net/http"
	"io/ioutil"
    "encoding/json"

	golog "github.com/ipfs/go-log"
	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	gologging "github.com/whyrusleeping/go-logging"
	"github.com/davecgh/go-spew/spew"

)

///// GLOBAL CONSTS & VARIABLES 

const (
    bootstrapperPort = "51000"
)

var MaxPeerPort int

type Peer struct {
    PeerAddress string `json:"PeerAddress"`
}

type PeerProfile struct { // connections of one peer
    ThisPeer Peer `json:"ThisPeer"` // any node
    PeerPort int `json:"PeerPort"` // port of peer
    Neighbors []Peer `json:"Neighbors"` // edges to that node
    Status bool `json:"Status"` // Status: Alive or Dead
    Connected bool `json:"Connected"` // If a node is connected or not [To be used later]
}

var thisPeerFullAddr string
var peerProfile PeerProfile // used to enroll THIS peer | connectP2PNet() & enrollP2PNet)()
var PeerGraph = make(map[string]PeerProfile) // Key = Node.PeerAddress; Value.Neighbors = Edges
var graphMutex sync.RWMutex

///// LIST OF FUNCTIONS

func p2pInit() {
	// LibP2P code uses golog to log messages. They log with different string IDs (i.e. "swarm")
	// We can control the verbosity level for all loggers with:
	golog.SetAllLoggers(gologging.INFO) // Change to DEBUG for extra info

	requestPort() // request THIS peer's port from bootstrapper
	if *verbose { log.Println("PeerIP = ", GetMyIP()) }
	queryP2PGraph() // query graph of peers in the P2P Network

	// Make a host that listens on the given multiaddress
	makeBasicHost(peerProfile.PeerPort, *secio, *seed)
	ha.SetStreamHandler("/p2p/1.0.0", handleStream)

	log.Println("Peerstore().Peers() before connecting =", ha.Peerstore().Peers())
	connectP2PNet()
	enrollP2PNet()
	log.Println("Peerstore().Peers() after connecting =", ha.Peerstore().Peers())

	go func() {
		for {
			time.Sleep(5 * time.Second)
			fmt.Println("\nList of peers:")
			for i, _ := range ha.Peerstore().Peers() {
				log.Println("-->", ha.Peerstore().Peers()[i].Pretty())
			}
		}
	}()
}

func requestPort() { // Requesting PeerPort
	log.Println("Requesting PeerPort from Bootstrapper", *bootstrapperAddr) //http://" + *bootstrapperIP + ":" + bootstrapperPort)

	response, err := http.Get(*bootstrapperAddr + "port-request")
	if err != nil {
		log.Println(err)
		log.Fatalln("PANIC: Unable to requestPort() from bootstrapper. Bootstrapper may be down.")
		// return
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		log.Fatalln("PANIC: Unable to requestPort() from bootstrapper. Bootstrapper may be down.")
		// return
	}

	json.Unmarshal(responseData, &peerProfile.PeerPort)
	if *verbose { log.Println("PeerPort = ", peerProfile.PeerPort) }

	if peerProfile.PeerPort == 0 {
		log.Println("PANIC: Exiting Program. PeerPort = 0. Bootstrapper may be down.")
		os.Exit(1)
	}
}

func queryP2PGraph() { // Query the graph of peers in the P2P Network from the Bootstrapper
	log.Println("Querying graph of peers from Bootstrapper", *bootstrapperAddr)
	
	response, err := http.Get(*bootstrapperAddr + "query-p2p-graph")
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

func connectP2PNet() {
	peerProfile.ThisPeer = Peer {PeerAddress : thisPeerFullAddr}

	if len(PeerGraph) == 0 { // first node in the network
		log.Println("I'm first peer. Creating Genesis Block.")
		Blockchain = append(Blockchain, generateGenesisBlock())
		save2File()
		spew.Dump(Blockchain)
		log.Println("I'm first peer. Listening for connections.")
	} else {
		log.Println("Connecting to P2P network")

		if *verbose { log.Println("Cardinality of PeerGraph = ", len(PeerGraph)) }

		// make connection with peers[choice]
		choice := genRandInt(len(PeerGraph))
		log.Println("Connecting choice = ", choice)

		peers := make([]string, 0, len(PeerGraph))
		for p, _ := range PeerGraph {
		    peers = append(peers, p)
		}
		log.Println("Connecting to", peers[choice])
		connect2Target(peers[choice])
		peerProfile.Neighbors = append(peerProfile.Neighbors, Peer {PeerAddress : peers[choice]})
		if *verbose { log.Println("peers[choice] = ", peers[choice]) }
	}
}

func enrollP2PNet() { // Enroll to the P2P Network by adding THIS peer with Bootstrapper
	log.Println("Enrolling in P2P network at Bootstrapper", *bootstrapperAddr)

	jsonValue, err := json.Marshal(peerProfile)
	if err != nil {
		log.Println(err)
		return
	}

	url := *bootstrapperAddr + "enroll-p2p-net"
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

func connect2Target(newTarget string) {
	log.Println("Attempting to connect to", newTarget)
	// The following code extracts target's peer ID from the
	// given multiaddress
	ipfsaddr, err := ma.NewMultiaddr(newTarget)
	if err != nil {
		log.Fatalln(err)
	}
	if *verbose { log.Printf("ipfsaddr = ", ipfsaddr) }
	if *verbose { log.Printf("Target = ", newTarget) }

	pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
	if err != nil {
		log.Fatalln(err)
	}
	if *verbose { log.Printf("pid = ", pid) }
	if *verbose { log.Printf("ma.P_IPFS = ", ma.P_IPFS) }

	peerid, err := peer.IDB58Decode(pid)
	if err != nil {
		log.Fatalln(err)
	}
	if *verbose { log.Println("peerid = ", peerid) }

	// Decapsulate the /ipfs/<peerID> part from the target
	// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, _ := ma.NewMultiaddr(
		fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)
	if *verbose { log.Printf("targetPeerAddr = ", targetPeerAddr) }
	if *verbose { log.Printf("targetAddr = ", targetAddr) }

	// We have a peer ID and a targetAddr so we add it to the peerstore
	// so LibP2P knows how to contact it
	ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

	log.Println("opening stream to", newTarget)
	// make a new stream from host B to host A
	// it should be handled on host A by the handler we set above because
	// we use the same /p2p/1.0.0 protocol
	s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
	if err != nil {
		log.Fatalln(err)
	}
	// Create a buffered stream so that read and writes are non blocking.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	// Create a thread to read and write data.
	go p2pWriteData(rw)
	go p2pReadData(rw)

	ch := make(chan os.Signal)
    signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-ch
        log.Println("Received Interrupt. Exiting now.")
        cleanup(rw)
        os.Exit(1)
    }()
	//select {} // hang forever
}

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress. It will use secio if secio is true.
func makeBasicHost(listenPort int, secio bool, randseed int64) { //(host.Host, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}
	if *verbose { log.Printf("r = ", r) }

	// Generate a key pair for this host. We will use it to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		//return nil, err
		log.Fatal(err)
	}
	if *verbose { log.Printf("priv = ", priv) }

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/"+GetMyIP()+"/tcp/%d", listenPort)),
		libp2p.Identity(priv),
	}
	if *verbose { log.Printf("opts = ", opts) }

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		//return nil, err
		log.Fatal(err)
	}
	if *verbose {
		log.Printf("basicHost = ", basicHost)
		log.Printf("basicHost.ID() = ", basicHost.ID())
		log.Printf("basicHost.ID().Pretty() = ", basicHost.ID().Pretty())
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))
	if *verbose { log.Printf("hostAddr = ", hostAddr) }

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := basicHost.Addrs()[0]
	if *verbose { log.Printf("addr = ", addr) }
	fullAddr := addr.Encapsulate(hostAddr)
	log.Printf("My fullAddr = %s\n", fullAddr)
	thisPeerFullAddr = fullAddr.String()

	//return basicHost, nil
	ha = basicHost // ha defined in defs.go
	if *verbose { log.Printf("basicHost = ", ha) }
}

func cleanup(rw *bufio.ReadWriter) {
	fmt.Println("cleanup")
    mutex.Lock()
    rw.WriteString("Exit\n")
    rw.Flush()
    mutex.Unlock()
}
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"os/signal"
    "syscall"

	golog "github.com/ipfs/go-log"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	gologging "github.com/whyrusleeping/go-logging"
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

	// LibP2P code uses golog to log messages. They log with different
	// string IDs (i.e. "swarm"). We can control the verbosity level for
	// all loggers with:
	golog.SetAllLoggers(gologging.INFO) // Change to DEBUG for extra info

	// // Parse options from the command line
	// listenF = flag.Int("l", 0, "wait for incoming connections")
	// target = flag.String("d", "", "target peer to dial")
	// secio = flag.Bool("secio", false, "enable secio")	
	// verbose = flag.Bool("verbose", false, "enable verbose")
	// seed = flag.Int64("seed", 0, "set random seed for id generation")
	// flag.Parse()

	if *listenF == 0 {
		log.Fatal("Please provide a port to bind on with -l")
	}

	// Make a host that listens on the given multiaddress
	ha, err := makeBasicHost(*listenF, *secio, *seed)
	if err != nil {
		log.Fatal(err)
	}
	if *verbose { log.Printf("ha = ", ha) }

	if *target == "" {
		log.Println("listening for connections")
		// Set a stream handler on host A. /p2p/1.0.0 is
		// a user-defined protocol name.
		ha.SetStreamHandler("/p2p/1.0.0", handleStream)

		//select {} // hang forever
		/**** This is where the listener code ends ****/
	} else {
		ha.SetStreamHandler("/p2p/1.0.0", handleStream)

		// The following code extracts target's peer ID from the
		// given multiaddress
		ipfsaddr, err := ma.NewMultiaddr(*target)
		if err != nil {
			log.Fatalln(err)
		}
		if *verbose { log.Printf("ipfsaddr = ", ipfsaddr) }
		if *verbose { log.Printf("*target = ", *target) }

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

		log.Println("opening stream")
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
		go writeData(rw)
		go readData(rw)

		ch := make(chan os.Signal)
	    signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	    go func() {
	        <-ch
	        cleanup(rw)
	        os.Exit(1)
	    }()

		//select {} // hang forever

	}

	log.Fatal(muxServer(*listenF)) // function is in mux.go
}


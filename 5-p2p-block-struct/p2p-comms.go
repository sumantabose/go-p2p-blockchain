package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"os/signal"
    "syscall"
    "runtime"
	"encoding/gob"

	"github.com/davecgh/go-spew/spew"
	net "github.com/libp2p/go-libp2p-net"

)

func handleStream(s net.Stream) {

	log.Println("Got a new stream!")
	log.Println("New list of peers =", ha.Peerstore().Peers())

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go p2pReadData(rw)
	go p2pWriteData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).

	ch := make(chan os.Signal)
    signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-ch
        log.Println("Received Interrupt. Exiting now.")
        cleanup(rw)
        os.Exit(1)
    }()
}

func p2pReadData(rw *bufio.ReadWriter) {

	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			//log.Fatal(err)
			log.Println(err)
		}

		if str == "" {
			return
		}
		if str != "Exit\n" { // Old condition was if str != "\n"

			chain := make([]Block, 0)
			if err := json.Unmarshal([]byte(str), &chain); err != nil {
				log.Fatal(err)
			}

			mutex.Lock()
			if len(chain) >= len(Blockchain) {
				Blockchain = chain

				save2File()

				bytes, err := json.MarshalIndent(Blockchain, "", "  ")
				if err != nil {

					log.Fatal(err)
				}
				// Green console color: 	\x1b[32m
				// Reset console color: 	\x1b[0m
				// fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
				if len(Blockchain) > LastRcvdBlockchainLen {
					// Green console color: 	\x1b[32m
					// Reset console color: 	\x1b[0m
					fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
					LastRcvdBlockchainLen = len(Blockchain)
				}
			}
			mutex.Unlock()
		}
	}
}

func p2pWriteData(rw *bufio.ReadWriter) {

	go func() {
		for {
			time.Sleep(5 * time.Second)
			mutex.Lock()
			bytes, err := json.Marshal(Blockchain)
			if err != nil {
				log.Println(err)
			}
			mutex.Unlock()

			// mutex.Lock()
			// rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
			// rw.Flush()
			// mutex.Unlock()

			mutex.Lock()
			rw.WriteString(string(bytes)+"\n")
			rw.Flush()
			if len(Blockchain) > LastSentBlockchainLen {
				fmt.Sprintf("%s\n", string(bytes))
				LastSentBlockchainLen = len(Blockchain)
			}
			mutex.Unlock()

		}
	}()

	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if sendData != "\n" {
			sendData = strings.Replace(sendData, "\n", "", -1) + " (From terminal)"
			newBlock := generateBlock(Blockchain[len(Blockchain)-1], sendData, "", 0)

			if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				mutex.Lock()
				Blockchain = append(Blockchain, newBlock)
				mutex.Unlock()
			}

			bytes, err := json.Marshal(Blockchain)
			if err != nil {
				log.Println(err)
			}

			spew.Dump(Blockchain)

			// mutex.Lock()
			// rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
			// rw.Flush()
			// mutex.Unlock()

			mutex.Lock()
			rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
			rw.Flush()
			LastSentBlockchainLen = len(Blockchain)
			mutex.Unlock()
		}
	}
}

func save2File() {
	gobCheck(writeGob())
}

func writeGob() error {

	dataDirFull := *dataDir + strconv.Itoa(peerProfile.PeerPort)

    if _, err := os.Stat(dataDirFull); os.IsNotExist(err) { // if dataDirFull does not exist
    	log.Println("`", dataDirFull, "` does not exist. Creating directory.")
    	os.Mkdir(dataDirFull, 0755) // https://stackoverflow.com/questions/14249467/os-mkdir-and-os-mkdirall-permission-value
	}

    filePath := dataDirFull + "/blockchain-" + strconv.Itoa(len(Blockchain)-1) + ".gob"
    file, err := os.Create(filePath)
    if err == nil {
        encoder := gob.NewEncoder(file)

		// The problem is that in the old code, there was an error when executing
		// encoder.Encode(Blockchain) but since we didn’t check for error, we didn’t realise that.
		// The output is blank because Blockchain fails to get encoded properly.
		// If we add error checking, we’d see something like:
		// encode error:gob: type not registered for interface: main.<Type>

		// If we add one line to our previous code before encoder.Encode(Blockchain) we could get the expected output,
		// gob.Register(Type{})
		// But what is this Type?

		// gob.Register(??)
		// READ: https://stackoverflow.com/questions/14121422/de-and-encode-interface-with-gob
		// READ: https://play.golang.org/p/xt4zNyPZ2W

		// ALTERNATE: Consider encoding and saving in JSON
        err = encoder.Encode(Blockchain)
        if err != nil {
			log.Println("GOB encode error:", err)
		}
    }
    file.Close()
    return err
}

func gobCheck(e error) { // Inspired from http://www.robotamer.com/code/go/gotamer/gob.html
    if e != nil {
        _, file, line, _ := runtime.Caller(1)
        log.Println(line, "\t", file, "\n", e)
        os.Exit(1)
    }
}
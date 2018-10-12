package main

import (
	"encoding/json"
	"time"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"net/http"
	"io"
	"os"
	"log"
	// "strconv"
)

type NewTarget struct { // Used to force a P2P connection through MUX [May dissable this feature later]
	Target string
}

// web server
func muxServer() error {
	mux := makeMuxRouter()
	//httpPort := os.Getenv("PORT")
	// log.Println("HTTP Server Listening on port :", peerProfile.PeerPort+1500) // peerProfile.PeerPort in peer-manager.go
	log.Println("HTTP MUX server listening on " + GetMyIP() + ":" + os.Getenv("PORT")) // listenPort is determined on the go
	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// create handlers
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/start", handleStart).Methods("GET")
	muxRouter.HandleFunc("/blockchain", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/raw", handleRawMaterialTxnWriteBlock).Methods("POST")
	muxRouter.HandleFunc("/del", handleDeliveryTxnWriteBlock).Methods("POST")
	muxRouter.HandleFunc("/comment", handleCommentWriteBlock).Methods("POST")
	muxRouter.HandleFunc("/connect", handleConnect).Methods("POST")
	muxRouter.HandleFunc("/query/{type}/{field}/{value}", handleQuery).Methods("GET")
	return muxRouter
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	if PeerStart == false {
		p2pInit() // Initialize P2P Network from Bootstrapper
		PeerStart = true
	} else {
		respondWithJSON(w, r, http.StatusBadRequest, "Bad request. Peer is already up.")
	}
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
    txnType, _ := params["type"]
    field, _ := params["field"]
    value, _ := params["value"]

    TxnArray := query(txnType, field, value)
	bytes, err := json.MarshalIndent(TxnArray, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	spew.Dump(TxnArray)
	io.WriteString(w, string(bytes))
}

// write blockchain when we receive an http request
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	mutex.Unlock()
	
	io.WriteString(w, string(bytes))
}

// takes JSON payload as an input for comment
func handleCommentWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var t StdInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	mutex.Lock()
	newBlock := generateBlock(Blockchain[len(Blockchain)-1], t.Comment, "", 0)
	mutex.Unlock()

	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		mutex.Lock()
		Blockchain = append(Blockchain, newBlock)
		save2File()
		mutex.Unlock()
		spew.Dump(Blockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

// takes JSON payload as an input for Raw Material Transaction
func handleRawMaterialTxnWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rmTxn RawMaterialTransaction

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rmTxn); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	mutex.Lock()
	newBlock := generateBlock(Blockchain[len(Blockchain)-1], "Raw Material Transaction (Type 1)", rmTxn, 1)
	mutex.Unlock()

	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		mutex.Lock()
		Blockchain = append(Blockchain, newBlock)
		save2File()
		mutex.Unlock()
		spew.Dump(Blockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

// takes JSON payload as an input for Delivery Transaction
func handleDeliveryTxnWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var delTxn DeliveryTransaction

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&delTxn); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	mutex.Lock()
	newBlock := generateBlock(Blockchain[len(Blockchain)-1], "Delivery Transaction (Type 2)", delTxn, 2)
	mutex.Unlock()

	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		mutex.Lock()
		Blockchain = append(Blockchain, newBlock)
		save2File()
		mutex.Unlock()
		spew.Dump(Blockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func handleConnect(w http.ResponseWriter, r *http.Request) { // May dissable this feature later
	w.Header().Set("Content-Type", "application/json")
	var m NewTarget

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		log.Println(err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	log.Println("MUX NewTarget =", m.Target)
	connect2Target(m.Target)
	respondWithJSON(w, r, http.StatusCreated, m.Target)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
package main

import (
	"encoding/json"
	"time"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"net/http"
	"io"
	"log"
	"strconv"
)


// web server
func muxServer() error {
	mux := makeMuxRouter()
	//httpPort := os.Getenv("PORT")
	log.Println("HTTP Server Listening on port :", peerProfile.PeerPort) // peerProfile.PeerPort in peer-manager.go
	s := &http.Server{
		Addr:           ":" + strconv.Itoa(peerProfile.PeerPort),
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
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	muxRouter.HandleFunc("/connect", handleConnect).Methods("POST")
	return muxRouter
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

// takes JSON payload as an input for heart rate (BPM)
func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	mutex.Lock()
	newBlock := generateBlock(Blockchain[len(Blockchain)-1], m.BPM)
	mutex.Unlock()

	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		mutex.Lock()
		Blockchain = append(Blockchain, newBlock)
		mutex.Unlock()
		spew.Dump(Blockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func handleConnect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m newTarget_json

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		log.Println(err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	log.Println("mux NewTarget =", m.NewTarget)
	connect2Target(m.NewTarget)
	respondWithJSON(w, r, http.StatusCreated, m.NewTarget)
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

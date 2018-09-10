/* README

Written by Sumanta Bose, 7 Sept 2018

*/

package main

import (
    "io"
    "log"
    "sync"
    "time"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"
    "github.com/davecgh/go-spew/spew"
)

///// GLOBAL CONSTS & VARIABLES 

const (
    listenPort = "51000"
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
var graphMutex sync.RWMutex

/////

///// LIST OF DRIVER FUNCTIONS

func init() {
    log.SetFlags(log.Lshortfile)
}

func main() {
    log.Fatal(launchMUXServer())
}

///// LIST OF MUX FUNCTIONS

func launchMUXServer() error { // launch MUX server
    mux := makeMUXRouter()
    log.Println("HTTP MUX server listening on port:", listenPort) // listenPort is a global const
    s := &http.Server{
        Addr:           ":" + listenPort,
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

func makeMUXRouter() http.Handler { // create handlers
    muxRouter := mux.NewRouter()
    muxRouter.HandleFunc("/query-p2p-list", handleQuery).Methods("GET")
    muxRouter.HandleFunc("/join-p2p-net", handleJoin).Methods("POST")
    return muxRouter
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
    log.Println("handleQuery() API called")
    graphMutex.RLock()
    bytes, err := json.MarshalIndent(peerGraph, "", "  ")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    io.WriteString(w, string(bytes))
    spew.Dump(bytes)
    graphMutex.RUnlock()
}

func handleJoin(w http.ResponseWriter, r *http.Request) {
    log.Println("handleJoin() API called")
    w.Header().Set("Content-Type", "application/json")
    //var peer Peer
    var incomingPeer PeerLinks

    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&incomingPeer); err != nil {
        log.Println(err)
        respondWithJSON(w, r, http.StatusBadRequest, r.Body)
        return
    }
    defer r.Body.Close()

    log.Println(incomingPeer)
    spew.Dump(incomingPeer)

    //ListOfPeers = append(ListOfPeers, peer)
    log.Println("Join request from:", incomingPeer, "successful")
    respondWithJSON(w, r, http.StatusCreated, incomingPeer)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")

    response, err := json.MarshalIndent(payload, "", "  ")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("HTTP 500: Internal Server Error"))
        return
    }
    w.WriteHeader(code)
    w.Write(response)
}

///// LIST OF HELPER FUNCTIONS


func inspect(incomingPeer PeerLinks) {
    spew.Dump(incomingPeer)


}









/* README

Written by Sumanta Bose, 11 Sept 2018

*/

package main

import (
    "io"
    "log"
    "time"
    "sync"
    "flag"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"
    "github.com/davecgh/go-spew/spew"
)

///// GLOBAL CONSTS, STRUCTS, VARS & METHODS 

const (
    listenPort = "51000"
)

type Peer struct {
    PeerAddress string `json:"PeerAddress"`
}

func (p Peer) Addr() string {
    return p.PeerAddress
}

type PeerProfile struct { // connections of one peer
    ThisPeer Peer `json:"ThisPeer"` // any node
    Neighbors []Peer `json:"Neighbors"` // edges to that node
    Status bool `json:"Status"` // Status: Alive or Dead
    Connected bool `json:"Connected"` // If a node is connected or not [To be used later]
}

var PeerGraph = make(map[string]PeerProfile) // Key = Node.PeerAddress; Value.Neighbors = Edges
var graphMutex sync.RWMutex
var verbose *bool

///// LIST OF FUNCTIONS

func init() {
    log.SetFlags(log.Lshortfile)
    verbose = flag.Bool("v", false, "enable verbose")
    flag.Parse()
}

func main() {
    log.Fatal(launchMUXServer())
}

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
    muxRouter.HandleFunc("/enroll-p2p-net", handleEnroll).Methods("POST")
    return muxRouter
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
    log.Println("handleQuery() API called")
    graphMutex.RLock()
    defer graphMutex.RUnlock() // until the end of the handleQuery()
    bytes, err := json.Marshal(PeerGraph) // MarshalIndent(PeerGraph, "", "  ")
    if err != nil {
        log.Println(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    io.WriteString(w, string(bytes))
    if *verbose { log.Println("PeerGraph = ", PeerGraph) ; spew.Dump(PeerGraph) }
}

func handleEnroll(w http.ResponseWriter, r *http.Request) {
    log.Println("handleEnroll() API called")
    w.Header().Set("Content-Type", "application/json")
    var incomingPeer PeerProfile

    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&incomingPeer); err != nil {
        log.Println(err)
        respondWithJSON(w, r, http.StatusBadRequest, r.Body)
        return
    }
    defer r.Body.Close()

    _ = updatePeerGraph(incomingPeer)
    log.Println("Enroll request from:", incomingPeer.ThisPeer, "successful")
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

func updatePeerGraph(inPeer PeerProfile) error {
    if *verbose { log.Println("incomingPeer = ", inPeer) ; spew.Dump(PeerGraph) }
    
    // Update PeerGraph
    graphMutex.Lock()
        if *verbose { log.Println("PeerGraph before update = ", PeerGraph) }
        PeerGraph[inPeer.ThisPeer.Addr()] = inPeer
        for _, neighbor := range inPeer.Neighbors {
            profile := PeerGraph[neighbor.Addr()]
            profile.Neighbors = append(profile.Neighbors, inPeer.ThisPeer)
            PeerGraph[neighbor.Addr()] = profile
        }
        if *verbose { log.Println("PeerGraph after update = ", PeerGraph) ; spew.Dump(PeerGraph)}
    graphMutex.Unlock()
    return nil
}
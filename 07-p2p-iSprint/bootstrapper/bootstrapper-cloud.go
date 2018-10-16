/* README

Written by Sumanta Bose, 29 Sept 2018
Modified on 2 Oct 2018
Deployed on: https://blockchain-bootstrapper.herokuapp.com/

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
    gonet "net"
    "os"

    "github.com/gorilla/mux"
    "github.com/davecgh/go-spew/spew"
)

///// GLOBAL CONSTS, STRUCTS, VARS & METHODS

// const (
//     listenPort = "51000"
// )

var MaxPeerPort int

type Peer struct {
    PeerAddress string `json:"PeerAddress"`
}

func (p Peer) Addr() string {
    return p.PeerAddress
}

type PeerProfile struct { // connections of one peer
    ThisPeer Peer `json:"ThisPeer"` // any node
    PeerPort int `json:"PeerPort"` // port of peer
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
    MaxPeerPort = 4999 // starting peer port
}

func main() {
    log.Fatal(launchMUXServer())
}

func launchMUXServer() error { // launch MUX server
    mux := makeMUXRouter()
    log.Println("HTTP MUX server listening on " + GetMyIP() + ":" + os.Getenv("PORT")) // listenPort is determined on the go
    s := &http.Server{
        Addr:           ":"+os.Getenv("PORT"),
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
    muxRouter.HandleFunc("/", handleHome).Methods("GET")
    muxRouter.HandleFunc("/query-p2p-graph", handleQuery).Methods("GET")
    muxRouter.HandleFunc("/enroll-p2p-net", handleEnroll).Methods("POST")
    muxRouter.HandleFunc("/port-request", handlePortReq).Methods("GET")
    return muxRouter
}

func handleHome(w http.ResponseWriter, r *http.Request) {
    log.Println("handleHome() API called")
    io.WriteString(w, "This is NTU blockchain bootstrapper. Developed by Sumanta Bose and Mayank Raikwar.")
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

func handlePortReq(w http.ResponseWriter, r *http.Request) {
    log.Println("handlePortReq() API called")
    MaxPeerPort = MaxPeerPort + 1
    bytes, err := json.Marshal(MaxPeerPort)
    if err != nil {
        log.Println(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    io.WriteString(w, string(bytes))
    if *verbose { log.Println("MaxPeerPort = ", MaxPeerPort) ; spew.Dump(MaxPeerPort) }
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

func  GetMyIP() string {
    var MyIP string

    conn, err := gonet.Dial("udp", "8.8.8.8:80")
    if err != nil {
       log.Fatalln(err)
    } else {
        localAddr := conn.LocalAddr().(*gonet.UDPAddr)
        MyIP = localAddr.IP.String()
    }
    return MyIP
}

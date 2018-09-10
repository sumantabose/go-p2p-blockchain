/* README

Written by Sumanta Bose, 10 Sept 2018

*/

package main

import (
    "io"
    "log"
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

type PeerProfile struct {
    ThisPeer Peer `json:"ThisPeer"`
    Status bool `json:"Status"`
}

var ListOfPeers []Peer

/////

///// LIST OF FUNCTIONS

func init() {
    log.SetFlags(log.Lshortfile)
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
    muxRouter.HandleFunc("/join-p2p-net", handleJoin).Methods("POST")
    return muxRouter
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
    log.Println("handleQuery() API called")
    bytes, err := json.MarshalIndent(ListOfPeers, "", "  ")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    io.WriteString(w, string(bytes))
    spew.Dump(bytes)
}

func handleJoin(w http.ResponseWriter, r *http.Request) {
    log.Println("handleJoin() API called")
    w.Header().Set("Content-Type", "application/json")
    var peer PeerProfile

    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&peer); err != nil {
        log.Println(err)
        respondWithJSON(w, r, http.StatusBadRequest, r.Body)
        return
    }
    defer r.Body.Close()

    log.Println(peer)

    ListOfPeers = append(ListOfPeers, peer.ThisPeer)
    log.Println("Join request from:", peer.ThisPeer, "successful")
    respondWithJSON(w, r, http.StatusCreated, peer)
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
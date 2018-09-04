package main

import (
    "io"
    "fmt"
    "log"
    "flag"
    "time"
    "strconv"
    "net/http"
    "math/rand"
    "encoding/json"

    "github.com/gorilla/mux"
    "github.com/davecgh/go-spew/spew"
)

///// GLOBAL VARIABLES

var listenPort *int

type Product struct {
    Name string `json:"Name"`
    SerialNo int `json:"SerialNo"`
    CodeNo int `json:"CodeNo"`
    Location int `json:"Location"`
}

var ProductData []Product

///// LIST OF FUNCTIONS

func init() {
    fmt.Println("Welcome to NTU Blockchain Server!")
    listenPort = flag.Int("port", 8080, "MUX server listen port")
    flag.Parse()
    log.SetFlags(log.Lshortfile)
}

func main() {
    log.Fatal(launchMUXServer())
}

func launchMUXServer() error { // launch MUX server
    mux := makeMUXRouter()
    log.Println("HTTP Server Listening on port :", *listenPort) // listenPort is a global flag
    s := &http.Server{
        Addr:           ":" + strconv.Itoa(*listenPort),
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
    muxRouter.HandleFunc("/info", handleInfo).Methods("GET")
    muxRouter.HandleFunc("/next", handleNext).Methods("GET")
    return muxRouter
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
    bytes, err := json.MarshalIndent(ProductData, "", "  ")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    io.WriteString(w, string(bytes))
}

func handleNext(w http.ResponseWriter, r *http.Request) {
    newProduct := Product{}
    newProduct.Name = genRandString(4)
    newProduct.SerialNo = len(ProductData)
    newProduct.CodeNo = genRandInt(50)
    newProduct.Location = 0

    ProductData = append(ProductData, newProduct)
    spew.Dump(ProductData)
    respondWithJSON(w, r, http.StatusCreated, newProduct)
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

func genRandInt(n int) int { // generate Random Integer less than 'n'
    myRandSource := rand.NewSource(time.Now().UnixNano())
    myRand := rand.New(myRandSource)
    val := myRand.Intn(n) + 10
    return val
}
///////////////////////////////////////////////////////
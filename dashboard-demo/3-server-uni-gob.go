/* README

Written by Sumanta Bose, 4 Sept 2018

MUX server methods available are:
    http://localhost:port/next
    http://localhost:port/info
    http://localhost:port/info/{loc}

*/

package main

import (
    "io"
    "os"
    "fmt"
    "log"
    "flag"
    "time"
    "strconv"
    "runtime"
    "net/http"
    "math/rand"
    "io/ioutil"
    "encoding/gob"
    "encoding/json"

    "github.com/gorilla/mux"
    "github.com/davecgh/go-spew/spew"
)

///// GLOBAL FLAGS & VARIABLES

var listenPort, totalLocs *int // listen port & total locations in the supply chain
var dataDir *string // data directory where the gob files are stored

type Product struct {
    Name string `json:"Name"`
    SerialNo int `json:"SerialNo"`
    CodeNo int `json:"CodeNo"`
    Location int `json:"Location"`
}

var ProductData []Product // `Product` array, to be saved as gob file

///// LIST OF FUNCTIONS

func init() {
    log.SetFlags(log.Lshortfile)

    log.Printf("Welcome to NTU Blockchain Server!")
    listenPort = flag.Int("port", 8080, "mux server listen port")
    totalLocs = flag.Int("locs", 5, "total locations in supply chain")
    dataDir = flag.String("data", "data", "pathname of data directory")
    flag.Parse()

    LoadProductData() // load from existing files, if any
}

func main() {
    log.Fatal(launchMUXServer())
}

func launchMUXServer() error { // launch MUX server
    mux := makeMUXRouter()
    log.Println("HTTP Server Listening on port:", *listenPort) // listenPort is a global flag
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
    muxRouter.HandleFunc("/info", handleInfoAll).Methods("GET")
    muxRouter.HandleFunc("/info/{loc}", handleInfoLoc).Methods("GET")
    muxRouter.HandleFunc("/next", handleNext).Methods("GET")
    return muxRouter
}

func handleInfoAll(w http.ResponseWriter, r *http.Request) {
    bytes, err := json.MarshalIndent(ProductData, "", "  ")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    io.WriteString(w, string(bytes))
}

func handleInfoLoc(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    loc, err := strconv.Atoi(params["loc"])

    var productStart, productEnd int

    if err == nil {
        for i, _ := range ProductData {
            if ProductData[i].Location == loc {
                productStart = i
                break
            }
        }
        for i, _ := range ProductData {
            if ProductData[i].Location == loc {
                productEnd = i
            }
        }

        var ProductSlice []Product = ProductData[productStart : productEnd+1]
        bytes, err := json.MarshalIndent(ProductSlice, "", "  ")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        io.WriteString(w, string(bytes))
    }
}

func handleNext(w http.ResponseWriter, r *http.Request) {
    newProduct := Product{}
    newProduct.Name = genRandString(4)
    newProduct.SerialNo = len(ProductData) + 1
    newProduct.CodeNo = genRandInt(50)
    newProduct.Location = 0

    ProductData = append(ProductData, newProduct)
    for i, _ := range ProductData {
        if ProductData[i].Location < *totalLocs {
            ProductData[i].Location ++
        }
    }
    spew.Dump(ProductData)
    fmt.Println("-----------")
    gobCheck(writeGob(ProductData, len(ProductData)))
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

func writeGob(object interface{}, fileNoCount int) error {
    filePath := *dataDir + "/product-data-" + strconv.Itoa(fileNoCount) + ".gob"
    file, err := os.Create(filePath)
    if err == nil {
        encoder := gob.NewEncoder(file)
        encoder.Encode(object)
    }
    file.Close()
    return err
}

func LoadProductData() { // load from existing files, if any
    if _, err := os.Stat(*dataDir); os.IsNotExist(err) { // if *dataDir does not exist
        log.Println("`", *dataDir, "` does not exist. Creating directory.")
        os.Mkdir(*dataDir, 0755) // https://stackoverflow.com/questions/14249467/os-mkdir-and-os-mkdirall-permission-value
    }

    files, err := ioutil.ReadDir(*dataDir) // dataDir from flag
    if err != nil {
        log.Fatal(err)
    }

    mostRecectFileNo := 0

    for _, file := range files {
        fileNo, _ := strconv.Atoi(file.Name()[len("product-data")+1:len(file.Name())-4])
        if fileNo > mostRecectFileNo { mostRecectFileNo = fileNo }
    }

    if mostRecectFileNo == 0 {
        log.Println("No existing ProductData")
    } else {
        mostRecentFile := *dataDir + "/product-data-" + strconv.Itoa(mostRecectFileNo) + ".gob"
        log.Println("Loading existing ProductData from", mostRecentFile)
        gobCheck(readGob(&ProductData, mostRecentFile))
    }  
}

func readGob(object interface{}, filePath string) error {
       file, err := os.Open(filePath)
       if err == nil {
              decoder := gob.NewDecoder(file)
              err = decoder.Decode(object)
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
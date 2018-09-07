/* README

Written by Sumanta Bose, 5 Sept 2018

MUX server methods available are:
    http://localhost:port/next
    http://localhost:port/next/{loop}
    http://localhost:port/info
    http://localhost:port/info/{loc}

FLAGS are:
  -data string
        pathname of data directory (default "data")
  -locs int
        total locations in supply chain (default 6)
  -port int
        mux server listen port (default 8080)
*/

package main

import (
    "io"
    "os"
    "fmt"
    "log"
    "flag"
    "time"
    "math"
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

    log.Printf("Welcome to Sumanta's Supply Chain Dashboard Server!")
    listenPort = flag.Int("port", 8080, "mux server listen port")
    totalLocs = flag.Int("locs", 6, "total locations in supply chain")
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
    muxRouter.HandleFunc("/next/{loop}", handleNextLoop).Methods("GET")
    muxRouter.HandleFunc("/post", handlePOST).Methods("POST")
    return muxRouter
}

func handlePOST(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var pdt Product

    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&pdt); err != nil {
        respondWithJSON(w, r, http.StatusBadRequest, r.Body)
        return
    }
    defer r.Body.Close()

    fmt.Println("Received in POST:", pdt)
    respondWithJSON(w, r, http.StatusCreated, pdt)
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
    var ProductElements []Product

    params := mux.Vars(r)
    loc, err := strconv.Atoi(params["loc"])
    if err == nil {
        for i, _ := range ProductData {
            if ProductData[i].Location == loc {
                ProductElements = append(ProductElements, ProductData[i])
            }
        }
        bytes, err := json.MarshalIndent(ProductElements, "", "  ")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        io.WriteString(w, string(bytes))
    }
}

func handleNext(w http.ResponseWriter, r *http.Request) {
    newProduct := updateProductData()
    gobCheck(writeGob(ProductData, len(ProductData)))
    respondWithJSON(w, r, http.StatusCreated, newProduct)
}

func handleNextLoop(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    loop, err := strconv.Atoi(params["loop"])
    var newProducts []Product

    if err == nil {
        for i := 0 ; i < loop ; i++ {
            // handleNext(w,r) // RISKY (don't attempt): server.go:2923: http: multiple response.WriteHeader calls
            newProduct := updateProductData()
            newProducts = append(newProducts, newProduct)
            gobCheck(writeGob(ProductData, len(ProductData)))
        }
        respondWithJSON(w, r, http.StatusCreated, newProducts)
    }
}

func updateProductData() Product {
    newProduct := Product{}
    newProduct.Name = genRandString(4)
    newProduct.SerialNo = len(ProductData) + 1
    newProduct.CodeNo = genRandInt(900,6000)
    newProduct.Location = 1

    ProductData = append(ProductData, newProduct)
    for i, _ := range ProductData {
        if i < len(ProductData) - 1 { // ignore the newest element i.e. newProduct
            if ProductData[i].Location < *totalLocs{
                ProductData[i].Location = ProductData[i].Location + int(math.Floor(float64(genRandInt(100,0))/80)) // 20% random increment
            }
        }
    }
    fmt.Println("----ProductData---- len:", len(ProductData))
    spew.Dump(ProductData)
    return newProduct
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

func genRandInt(n int, offset int) int { // generate Random Integer less than 'n' with an offset
    myRandSource := rand.NewSource(time.Now().UnixNano())
    myRand := rand.New(myRandSource)
    val := myRand.Intn(n) + offset
    return val
}

///////////////////////////////////////////////////////
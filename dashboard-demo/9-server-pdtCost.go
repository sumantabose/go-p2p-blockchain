/* README

Written by Sumanta Bose, 8 Oct 2018

MUX server methods available are:
    http://localhost:port/
    http://localhost:port/next
    http://localhost:port/next/{loop}
    http://localhost:port/info
    http://localhost:port/info/{loc}
    http://localhost:port/move/{serial}
    http://localhost:port/post/{CodeNo}/{Name}/{Cost}
    http://localhost:port/query/{field}/{value}
    http://localhost:port/history/{SerialNo}
    http://localhost:port/blockchain

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
    "encoding/hex"
    "encoding/gob"
    "encoding/json"
    "crypto/sha256"

    "github.com/gorilla/mux"
    "github.com/davecgh/go-spew/spew"
)

///// GLOBAL FLAGS & VARIABLES

var StartTime time.Time

var listenPort, totalLocs *int // listen port & total locations in the supply chain
var dataDir *string // data directory where the gob files are stored

type Product struct {
    Name string `json:"Name"`
    Cost int `json:"Cost"`
    SerialNo int `json:"SerialNo"`
    CodeNo int `json:"CodeNo"`
    Location int `json:"Location"`
}

type ProductVals struct { // Used while Posting from Location 1
    Name string `json:"Name"`
    CodeNo int `json:"CodeNo"`
}

var ProductData []Product // `Product` array, to be saved as gob file

/////

type Event struct {
    Snapshot int
    Timestamp string
    Location int
}

type Block struct {
    Index int
    Timestamp string
    BlockProductData []Product
    PrevHash string
    ThisHash string
}

var Blockchain []Block

///// LIST OF FUNCTIONS

func init() {
    log.SetFlags(log.Lshortfile)

    log.Printf("Welcome to Sumanta's Supply Chain Dashboard Server!")
    listenPort = flag.Int("port", 8080, "mux server listen port")
    totalLocs = flag.Int("locs", 6, "total locations in supply chain")
    dataDir = flag.String("data", "data", "pathname of data directory")
    flag.Parse()

    LoadProductData() // load from existing files, if any

    StartTime = time.Now()
    StartTime = StartTime.AddDate(0, -6, 10) // random negative offset
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
    muxRouter.HandleFunc("/", handleHome).Methods("GET")
    muxRouter.HandleFunc("/info", handleInfoAll).Methods("GET")
    muxRouter.HandleFunc("/info/{loc}", handleInfoLoc).Methods("GET")
    muxRouter.HandleFunc("/next", handleNext).Methods("GET")
    muxRouter.HandleFunc("/next/{loop}", handleNextLoop).Methods("GET")
    muxRouter.HandleFunc("/post/{CodeNo}/{Name}/{Cost}", handlePost).Methods("GET")
    muxRouter.HandleFunc("/query/{field}/{value}", handleQuery).Methods("GET")
    muxRouter.HandleFunc("/history/{SerialNo}", handleHistory).Methods("GET")
    muxRouter.HandleFunc("/move/{serial}", handleMove).Methods("GET")
    muxRouter.HandleFunc("/blockchain", handleBlockchain).Methods("GET")
    return muxRouter
}

func handleHome(w http.ResponseWriter, r *http.Request) {
    log.Println("handleHome() API called")
    io.WriteString(w, "You have entered the restricted zone. Trespassing is strictly prohibited. Defaulters will be reported.")
}

func handleMove(w http.ResponseWriter, r *http.Request) {
    var pdt Product

    params := mux.Vars(r)
    serial, err := strconv.Atoi(params["serial"])
    if err == nil {
        if serial <= len(ProductData) && serial > 0 {
            if ProductData[serial-1].Location < *totalLocs {
                ProductData[serial-1].Location = ProductData[serial-1].Location + 1
                pdt = ProductData[serial-1]
                gobCheck(writeGob(ProductData, len(ProductData)))
                respondWithJSON(w, r, http.StatusCreated, pdt)
            }
        }
    } else {
        respondWithJSON(w, r, http.StatusBadRequest, "SerialNo invalid.")
    }
}

func handlePost(w http.ResponseWriter, r *http.Request) {
 
    params := mux.Vars(r)
 
    newProduct := Product{}
    newProduct.SerialNo = len(ProductData) + 1
    newProduct.Location = 1
    newProduct.Name = params["Name"]
    codeNo, err1 := strconv.Atoi(params["CodeNo"]) ; newProduct.CodeNo = codeNo
    cost, err2 := strconv.Atoi(params["Cost"]) ; newProduct.Cost = cost

    if err1 == nil && err2 == nil {
        fmt.Println("Adding to ProductData:", newProduct)
        ProductData = append(ProductData, newProduct)
        gobCheck(writeGob(ProductData, len(ProductData)))
        respondWithJSON(w, r, http.StatusCreated, newProduct)
    }    
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
    var ProductElements []Product
 
    params := mux.Vars(r)
    field := params["field"]
    if field == "Name" {
        value := params["value"]
        for i, _ := range ProductData {
            if ProductData[i].Name == value {
                ProductElements = append(ProductElements, ProductData[i])
            }
        }
    } else if field == "SerialNo" {
        value, err := strconv.Atoi(params["value"])
        if err == nil {
            for i, _ := range ProductData {
                if ProductData[i].SerialNo == value {
                    ProductElements = append(ProductElements, ProductData[i])
                }
            }
        }
    } else if field == "CodeNo" {
        value, err := strconv.Atoi(params["value"])
        if err == nil {
            for i, _ := range ProductData {
                if ProductData[i].CodeNo == value {
                    ProductElements = append(ProductElements, ProductData[i])
                }
            }
        }  
    }

    bytes, err := json.MarshalIndent(ProductElements, "", "  ")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    io.WriteString(w, string(bytes))
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

func handleHistory(w http.ResponseWriter, r *http.Request) {
    var History []Event
    var TempProductData []Product

    params := mux.Vars(r)
    SerialNo, err := strconv.Atoi(params["SerialNo"])
    if err == nil {
        files, err := ioutil.ReadDir(*dataDir) // dataDir from flag
        if err != nil { log.Fatal(err) }

        mostRecectFileNo := 0

        for _, file := range files {
            fileNo, _ := strconv.Atoi(file.Name()[len("product-data")+1:len(file.Name())-4])
            if fileNo > mostRecectFileNo { mostRecectFileNo = fileNo }
        }

        for i := 1; i <= mostRecectFileNo; i++ {
            readFile := *dataDir + "/product-data-" + strconv.Itoa(i) + ".gob"
            gobCheck(readGob(&TempProductData, readFile))
            for j, _ := range TempProductData {
                if TempProductData[j].SerialNo == SerialNo {
                    e := Event {
                        Snapshot: i,
                        Timestamp: StartTime.AddDate(0, 0, genRandInt(3,1)+(3*i)).String(), // random date increment
                        Location: TempProductData[j].Location,
                    }
                    History = append(History, e)
                }
            }
        }
        bytes, err := json.MarshalIndent(History, "", "  ")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        io.WriteString(w, string(bytes))
    } // ending if err == nil 
}

func handleBlockchain(w http.ResponseWriter, r *http.Request) {
    genesisBlock := Block {
        Index: 0,
        Timestamp: StartTime.String(),
        PrevHash: "GENESIS-BLOCK",
    }
    genesisBlock.ThisHash = calculateHash(genesisBlock)
    Blockchain = append(Blockchain, genesisBlock)

    files, err := ioutil.ReadDir(*dataDir) // dataDir from flag
    if err != nil { log.Fatal(err) }

    mostRecectFileNo := 0
    for _, file := range files {
        fileNo, _ := strconv.Atoi(file.Name()[len("product-data")+1:len(file.Name())-4])
        if fileNo > mostRecectFileNo { mostRecectFileNo = fileNo }
    }

    var TempProductData []Product

    for i := 1; i <= mostRecectFileNo; i++ {
        readFile := *dataDir + "/product-data-" + strconv.Itoa(i) + ".gob"
        gobCheck(readGob(&TempProductData, readFile))
        b := Block {
            Index: i,
            Timestamp: StartTime.AddDate(0, 0, genRandInt(3,1)+(3*i)).String(), // random date increment
            BlockProductData: TempProductData,
            PrevHash: Blockchain[len(Blockchain)-1].ThisHash,
        }
        b.ThisHash = calculateHash(b)
        Blockchain = append(Blockchain, b)
    }

    respondWithJSON(w, r, http.StatusCreated, Blockchain)
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
    newProduct.Name, newProduct.Cost = genRandPdt() //genRandString(4)
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
    if err != nil { log.Fatal(err) }

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

func genRandPdt() (string, int) {

    PdtMap := map[string]int{
        "Camera": 300, "Laptop": 1200, "Chair": 30, "Table": 40, "Pen": 5, "Pencil": 2, "Car": 150000, "Motorbike": 40000,
        "Beer": 15, "Wine": 25, "Shampoo": 5,"Soap": 3, "Radio": 60, "Fridge": 250, "Shirt": 40, "Pant": 50, "Blazer": 90,
        "Watch": 150, "Shoe": 180, "Tie": 70, "Cap": 35, "Scarf": 15, "Book": 20, "Chalk": 2, "Desk": 50, "Wardrobe": 120,
        "Chocolate": 10, "Cake": 20, "Pie": 5, "Chips": 2, "Ice-cream": 2, "Headphone": 130, "Bulb": 30, "Biycle": 90,
        "Guitar": 60, "Soft-drink": 10, "Pillow": 30, "Mattress": 90, "Goggles": 40, "Bag": 20,
    }

    randChoice := genRandInt(len(PdtMap),0) ; i := 0
    for key, value := range PdtMap {
        if i == randChoice {
            return key, value
        }
        i ++
    }
    return "Camera", 300 // default - although it will never reach here
}

// SHA256 hashing
func calculateHash(b Block) string {
    record := strconv.Itoa(b.Index) + b.Timestamp + spew.Sdump(b.BlockProductData) + b.PrevHash
    h := sha256.New()
    h.Write([]byte(record))
    hashed := h.Sum(nil)
    return hex.EncodeToString(hashed)
}

///////////////////////////////////////////////////////
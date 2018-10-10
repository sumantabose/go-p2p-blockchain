/* README

Written by Sumanta Bose, 9 Oct 2018

MUX server methods available are:
    http://localhost:port/

    http://localhost:port/info
    http://localhost:port/info/{location}

    http://localhost:port/next
    http://localhost:port/next/{loop}
    
    http://localhost:port/post/{CodeNo}/{Name}/{Cost}
    http://localhost:port/query/{field}/{value}
    http://localhost:port/move/{serial}

    http://localhost:port/history/{SerialNo}
    http://localhost:port/account/{direction}/{location}

    http://localhost:port/blockchain

FLAGS are:
  -pdts string
        pathname of product data directory (default "pdts")
  -accs string
        pathname of accounts data directory (default "accs")
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
var pdtsDir *string // product data directory where the gob files are stored
var accsDir *string // accounts data directory where the gob files are stored

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

type Event struct { // A Snapshot of History
    Snapshot int `json:"Snapshot"`
    Timestamp string `json:"Timestamp"`
    Location int `json:"Location"`
}

type Block struct { // An element of Blockchain
    Index int
    Timestamp string
    BlockProductData []Product
    PrevHash string
    ThisHash string
}

/////

type Entry struct { // A Snapshot of Account
    Snapshot int `json:"Snapshot"`
    Timestamp string `json:"Timestamp"`
    Name string `json:"Name"`
    Cost int `json:"Cost"`
    SerialNo int `json:"SerialNo"`
    CodeNo int `json:"CodeNo"`
    Party int `json:"Party"` // Transacting Party
    Direction string `json:"Direction"` // In or Out
}

var Blockchain []Block

///// LIST OF FUNCTIONS

func init() {

    gob.Register(Product{}) ; gob.Register(Event{}) ; gob.Register(Block{}) ; gob.Register(Entry{})
    gob.Register(map[string]interface{}{})

    log.SetFlags(log.Lshortfile)

    log.Printf("Welcome to Sumanta's Supply Chain Dashboard Server!")
    listenPort = flag.Int("port", 8080, "mux server listen port")
    totalLocs = flag.Int("locs", 6, "total locations in supply chain")
    pdtsDir = flag.String("pdts", "pdts", "pathname of product data directory")
    accsDir = flag.String("accs", "accs", "pathname of accounts data directory")
    flag.Parse()

    LoadProductData() // load from existing files, if any
    CreateAccounts() // create accounts for all totalLocs

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
    muxRouter.HandleFunc("/info/{location}", handleInfoLoc).Methods("GET")

    muxRouter.HandleFunc("/next", handleNext).Methods("GET")
    muxRouter.HandleFunc("/next/{loop}", handleNextLoop).Methods("GET")

    muxRouter.HandleFunc("/post/{CodeNo}/{Name}/{Cost}", handlePost).Methods("GET")
    muxRouter.HandleFunc("/query/{field}/{value}", handleQuery).Methods("GET")
    muxRouter.HandleFunc("/move/{serial}", handleMove).Methods("GET")

    muxRouter.HandleFunc("/history/{SerialNo}", handleHistory).Methods("GET")    
    muxRouter.HandleFunc("/account/{direction}/{location}", handleAccount).Methods("GET")

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
                updateAccount(ProductData[serial-1], len(ProductData))
                pdt = ProductData[serial-1]
                gobCheck(writePdtGob(ProductData, len(ProductData)))
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
        gobCheck(writePdtGob(ProductData, len(ProductData)))
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
    location, err := strconv.Atoi(params["location"])
    if err == nil {
        for i, _ := range ProductData {
            if ProductData[i].Location == location {
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
        files, err := ioutil.ReadDir(*pdtsDir) // pdtsDir from flag
        if err != nil { log.Fatal(err) }

        mostRecectFileNo := 0

        for _, file := range files {
            fileNo, _ := strconv.Atoi(file.Name()[len("product-data")+1:len(file.Name())-4])
            if fileNo > mostRecectFileNo { mostRecectFileNo = fileNo }
        }

        for i := 1; i <= mostRecectFileNo; i++ {
            readFile := *pdtsDir + "/product-data-" + strconv.Itoa(i) + ".gob"
            gobCheck(readGob(&TempProductData, readFile))
            for j, _ := range TempProductData {
                if TempProductData[j].SerialNo == SerialNo {
                    e := Event {
                        Snapshot: i,
                        Timestamp: StartTime.AddDate(0, 0, genRandInt(3,1)+(3*i)).Add(time.Duration(genRandInt(30000,0))*time.Second).Format("02-01-2006 15:04:05 Mon"), // random date increment
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
        Timestamp: StartTime.Add(time.Duration(genRandInt(30000,0))*time.Second).Format("02-01-2006 15:04:05 Mon"),
        PrevHash: "GENESIS-BLOCK",
    }
    genesisBlock.ThisHash = calculateHash(genesisBlock)
    Blockchain = append(Blockchain, genesisBlock)

    files, err := ioutil.ReadDir(*pdtsDir) // pdtsDir from flag
    if err != nil { log.Fatal(err) }

    mostRecectFileNo := 0
    for _, file := range files {
        fileNo, _ := strconv.Atoi(file.Name()[len("product-data")+1:len(file.Name())-4])
        if fileNo > mostRecectFileNo { mostRecectFileNo = fileNo }
    }

    var TempProductData []Product

    for i := 1; i <= mostRecectFileNo; i++ {
        readFile := *pdtsDir + "/product-data-" + strconv.Itoa(i) + ".gob"
        gobCheck(readGob(&TempProductData, readFile))
        b := Block {
            Index: i,
            Timestamp: StartTime.AddDate(0, 0, genRandInt(3,1)+(3*i)).Add(time.Duration(genRandInt(30000,0))*time.Second).Format("02-01-2006 15:04:05 Mon"), // random date increment
            BlockProductData: TempProductData,
            PrevHash: Blockchain[len(Blockchain)-1].ThisHash,
        }
        b.ThisHash = calculateHash(b)
        Blockchain = append(Blockchain, b)
    }

    respondWithJSON(w, r, http.StatusCreated, Blockchain)
}

func handleAccount(w http.ResponseWriter, r *http.Request) {
    var returnAcc []Entry

    params := mux.Vars(r)
    direction := params["direction"]
    location, err := strconv.Atoi(params["location"])
    
    if err == nil && (direction == "in" || direction == "out") {
        filePath := *accsDir + "/" + direction + "-" + strconv.Itoa(location) + ".gob"
        gobCheck(readGob(&returnAcc, filePath))
    }
    // respondWithJSON(w, r, http.StatusCreated, returnAcc)
    bytes, err := json.MarshalIndent(returnAcc, "", "  ")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    io.WriteString(w, string(bytes))
}

func handleNext(w http.ResponseWriter, r *http.Request) {
    newProduct := updateProductData()
    gobCheck(writePdtGob(ProductData, len(ProductData)))
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
            gobCheck(writePdtGob(ProductData, len(ProductData)))
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
                randIncrement := int(math.Floor(float64(genRandInt(100,0))/80)) // 20% random increment
                ProductData[i].Location = ProductData[i].Location + randIncrement
                if randIncrement == 1 {
                    updateAccount(ProductData[i], len(ProductData))
                }
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

func writePdtGob(object interface{}, fileNoCount int) error {
    filePath := *pdtsDir + "/product-data-" + strconv.Itoa(fileNoCount) + ".gob"
    file, err := os.Create(filePath)
    if err == nil {
        encoder := gob.NewEncoder(file)
        encoder.Encode(object)
    }
    file.Close()
    return err
}

func writeAccGob(object interface{}, direction string, location int) error {
    filePath := *accsDir + "/" + direction + "-" + strconv.Itoa(location) + ".gob"
    file, err := os.Create(filePath)
    if err == nil {
        encoder := gob.NewEncoder(file)
        encoder.Encode(object)
    }
    file.Close()
    return err
}

func LoadProductData() { // load from existing files, if any
    if _, err := os.Stat(*pdtsDir); os.IsNotExist(err) { // if *pdtsDir does not exist
        log.Println("`", *pdtsDir, "` does not exist. Creating directory.")
        os.Mkdir(*pdtsDir, 0755) // https://stackoverflow.com/questions/14249467/os-mkdir-and-os-mkdirall-permission-value
    }

    files, err := ioutil.ReadDir(*pdtsDir) // pdtsDir from flag
    if err != nil { log.Fatal(err) }

    mostRecectFileNo := 0

    for _, file := range files {
        fileNo, _ := strconv.Atoi(file.Name()[len("product-data")+1:len(file.Name())-4])
        if fileNo > mostRecectFileNo { mostRecectFileNo = fileNo }
    }

    if mostRecectFileNo == 0 {
        log.Println("No existing ProductData")
    } else {
        mostRecentFile := *pdtsDir + "/product-data-" + strconv.Itoa(mostRecectFileNo) + ".gob"
        log.Println("Loading existing ProductData from", mostRecentFile)
        gobCheck(readGob(&ProductData, mostRecentFile))
    }  
}

func CreateAccounts() {
    var blankAcc []Entry
    if _, err := os.Stat(*accsDir); os.IsNotExist(err) { // if *accsDir does not exist
        log.Println("`", *accsDir, "` does not exist. Creating directory.")
        os.Mkdir(*accsDir, 0755) // https://stackoverflow.com/questions/14249467/os-mkdir-and-os-mkdirall-permission-value
        for i := 1 ; i <= *totalLocs ; i ++ {
            gobCheck(writeAccGob(blankAcc, "out", i))
            gobCheck(writeAccGob(blankAcc, "in", i))
        } 
    }  
}

func updateAccount(p Product, snapshot int) { // AFTER location update. VERY VERY IMP "AFTER"
    newLocation := p.Location        ; outFile := *accsDir + "/out-" + strconv.Itoa(newLocation) + ".gob" 
    oldLocation := newLocation - 1   ; inFile := *accsDir + "/in-" + strconv.Itoa(oldLocation) + ".gob"

    var outAcc []Entry  ; gobCheck(readGob(&outAcc, outFile))  // ; fmt.Printf("\x1b[32m%s\x1b[0m> ", spew.Sdump(outAcc))
    var inAcc []Entry   ; gobCheck(readGob(&inAcc, inFile))    // ; fmt.Printf("\x1b[32m%s\x1b[0m> ", spew.Sdump(inAcc))

    e := Entry {
        Snapshot: snapshot,
        Timestamp: StartTime.AddDate(0, 0, genRandInt(3,1)+(3*snapshot)).Add(time.Duration(genRandInt(30000,0))*time.Second).Format("02-01-2006 15:04:05 Mon"), // random date increment
        Name: p.Name,
        Cost: p.Cost,
        SerialNo: p.SerialNo,
        CodeNo: p.CodeNo,
    }

    outE := e  ; outE.Party = oldLocation   ; outE.Direction = "out"    ; outAcc = append(outAcc, outE)
    inE := e   ; inE.Party = newLocation    ; inE.Direction = "in"      ; inAcc = append(inAcc, inE)

    gobCheck(writeAccGob(outAcc, "out", newLocation))
    gobCheck(writeAccGob(inAcc, "in", oldLocation))
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
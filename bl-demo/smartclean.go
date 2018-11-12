/* README

Written by Sumanta Bose, 13 Nov 2018

MUX server methods available are:
    http://localhost:port/

    http://localhost:port/info
    http://localhost:port/info/{member}

    http://localhost:port/add
    http://localhost:port/add/{loop}

    http://localhost:port/post

FLAGS are:
  -dataDir string
        pathname of BL data storage directory (default "bldata")
  -toilets int
        total number of members in the BL supply chain (default 5) // 4 + 1
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
    _"math"
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

var port, toilets *int  // listen port & total toilet member in the network
var dataDir *string     // BL data directory where the gob files are stored


type Entry struct {
    SerialNum int `json:"SerialNum"`
    Timestamp string `json:"Timestamp"`

    Temperature int `json:"Temperature"`
    Humidity int `json:"Humidity"`
    Luminescence int `json:"Luminescence"`
    Decibel int `json:"Decibel"`
    GasComposition string `json:"GasComposition"`
    Olfactometer int `json:"Olfactometer"`

    Toilet int
}

var Report []Entry // Array, to be saved as gob file

///// LIST OF FUNCTIONS

func init() {

    gob.Register(Entry{})
    gob.Register(map[string]interface{}{})

    log.SetFlags(log.Lshortfile)

    log.Printf("Welcome to Sumanta's SmartClean Dashboard Server!")
    port = flag.Int("port", 8080, "mux server listen port")
    toilets = flag.Int("toilets", 8, "total number of toilets in the network")
    dataDir = flag.String("dataDir", "dataDir", "pathname of data storage directory")
    flag.Parse()

    LoadData() // load from existing files, if any
}

func main() {
    log.Fatal(launchMUXServer())
}

func launchMUXServer() error { // launch MUX server
    mux := makeMUXRouter()
    log.Println("HTTP Server Listening on port:", *port) // port is a global flag
    s := &http.Server{
        Addr:           ":" + strconv.Itoa(*port),
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
    muxRouter.HandleFunc("/info/{toilet}", handleInfoMember).Methods("GET")

    muxRouter.HandleFunc("/add", handleAdd).Methods("GET")
    muxRouter.HandleFunc("/add/{loop}", handleAddLoop).Methods("GET")

    muxRouter.HandleFunc("/post", handlePOST).Methods("POST")

    return muxRouter
}

func handleHome(w http.ResponseWriter, r *http.Request) {
    log.Println("handleHome() API called")
    io.WriteString(w, "You have entered the restricted zone. Trespassing is strictly prohibited. Defaulters will be reported.\n")
}

func handleInfoAll(w http.ResponseWriter, r *http.Request) {
    bytes, err := json.MarshalIndent(Report, "", "  ")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    io.WriteString(w, string(bytes))
}

func handleInfoMember(w http.ResponseWriter, r *http.Request) {
    var EntryElements []Entry

    params := mux.Vars(r)
    toilet, err := strconv.Atoi(params["toilet"])
    if err == nil {
        for i, _ := range Report {
            if Report[i].Toilet == toilet {
                EntryElements = append(EntryElements, Report[i])
            }
        }
        bytes, err := json.MarshalIndent(EntryElements, "", "  ")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        io.WriteString(w, string(bytes))
    } else {
        respondWithJSON(w, r, http.StatusBadRequest, "Invalid Member.")
    }
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
    newEntry := updateReport()
    gobCheck(writeBLGob(Report, len(Report)))
    respondWithJSON(w, r, http.StatusCreated, newEntry)
}

func handleAddLoop(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    loop, err := strconv.Atoi(params["loop"])
    var newEntries []Entry

    if err == nil {
        for i := 0 ; i < loop ; i++ {
            newEntry := updateReport()
            newEntries = append(newEntries, newEntry)
            gobCheck(writeBLGob(Report, len(Report)))
        }
        respondWithJSON(w, r, http.StatusCreated, newEntries)
    }
}

func handlePOST(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var entry Entry

    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&entry); err != nil {
        respondWithJSON(w, r, http.StatusBadRequest, r.Body)
        return
    }
    defer r.Body.Close()

    entry.SerialNum = len(Report) + 1
    entry.Timestamp = time.Now().Format("02-01-2006 15:04:05 Mon")
    log.Println("Adding to Report in handlePOST():", entry)

    Report = append(Report, entry)
    gobCheck(writeBLGob(Report, len(Report)))
    respondWithJSON(w, r, http.StatusCreated, entry)
}

///// LIST OF HELPER FUNCTIONS

func LoadData() { // load from existing files, if any
    if _, err := os.Stat(*dataDir); os.IsNotExist(err) { // if *dataDir does not exist
        log.Println("`", *dataDir, "` does not exist. Creating directory.")
        os.Mkdir(*dataDir, 0755) // https://stackoverflow.com/questions/14249467/os-mkdir-and-os-mkdirall-permission-value
    }

    files, err := ioutil.ReadDir(*dataDir) // dataDir from flag
    if err != nil { log.Fatal(err) }

    mostRecectFileNo := 0

    for _, file := range files {
        fileNo, _ := strconv.Atoi(file.Name()[len("toilet-data")+1:len(file.Name())-4])
        if fileNo > mostRecectFileNo { mostRecectFileNo = fileNo }
    }

    if mostRecectFileNo == 0 {
        log.Println("No existing toilet data")
    } else {
        mostRecentFile := *dataDir + "/toilet-data-" + strconv.Itoa(mostRecectFileNo) + ".gob"
        log.Println("Loading existing toilet data from", mostRecentFile)
        gobCheck(readGob(&Report, mostRecentFile))
    }  
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

func readGob(object interface{}, filePath string) error {
       file, err := os.Open(filePath)
       if err == nil {
              decoder := gob.NewDecoder(file)
              err = decoder.Decode(object)
       }
       file.Close()
       return err
}

func writeBLGob(object interface{}, fileNoCount int) error {
    filePath := *dataDir + "/toilet-data-" + strconv.Itoa(fileNoCount) + ".gob"
    file, err := os.Create(filePath)
    if err == nil {
        encoder := gob.NewEncoder(file)
        encoder.Encode(object)
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

func updateReport() Entry {
    newEntry := Entry{}

    newEntry.SerialNum = len(Report) + 1
    newEntry.Timestamp = time.Now().Format("02-01-2006 15:04:05 Mon")

    newEntry.Temperature = genRandInt(30,10)
    newEntry.Humidity = genRandInt(30,10)
    newEntry.Luminescence = genRandInt(30,10)
    newEntry.Decibel = genRandInt(30,10)
    newEntry.GasComposition = "N [78], O [21], X [1]"
    newEntry.Olfactometer = genRandInt(30,10)

    newEntry.Toilet = genRandInt(8,1)

    Report = append(Report, newEntry)

    fmt.Println("----Report---- len:", len(Report))
    spew.Dump(newEntry)
    return newEntry
}

///// LIST OF TRIVIAL FUNCTIONS

func genRandString(n int) string { // generate Random String of length 'n'
    var letterRunes  = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
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
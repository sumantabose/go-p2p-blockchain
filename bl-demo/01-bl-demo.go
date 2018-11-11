/* README

Written by Sumanta Bose, 12 Nov 2018

MUX server methods available are:
    http://localhost:port/

    http://localhost:port/info
    http://localhost:port/info/{status}/{member}

    http://localhost:port/add
    http://localhost:port/add/{loop}

    http://localhost:port/move/{serial}
    http://localhost:port/post

FLAGS are:
  -bldata string
        pathname of BL data storage directory (default "bldata")
  -members int
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
    "math"
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

var port, members *int  // listen port & total members in the B/L supply chain
var bldataDir *string      // BL data directory where the gob files are stored

type BL struct {
    SerialNum int `json:"SerialNum"`
    BLuid int `json:"BLuid"`
    Timestamp string `json:"Timestamp"`

    Issuer string `json:"Issuer"`
    Shipper string `json:"Shipper"`
    Consignee string `json:"Consignee"`
    Releaser string `json:"Releaser"`

    PortOfLoading string `json:"PortOfLoading"`
    PortOfDischarge string `json:"PortOfDischarge"`
    VesselNum int `json:"VesselNum"`

    PkgQty int `json:"PkgQty"`
    PkgType string `json:"PkgType"`
    PkgDesc string `json:"PkgDesc"`
    PkgWeight int `json:"PkgWeight"`

    BLurl string `json:"BLurl"`

    Owner int `json:"Owner"`
}

var BLArray []BL // `BL` array, to be saved as gob file

///// LIST OF FUNCTIONS

func init() {

    gob.Register(BL{})
    gob.Register(map[string]interface{}{})

    log.SetFlags(log.Lshortfile)

    log.Printf("Welcome to Sumanta's Bill of Lading Dashboard Server!")
    port = flag.Int("port", 8080, "mux server listen port")
    members = flag.Int("members", 5, "total number of members in the BL supply chain")
    bldataDir = flag.String("bldata", "bldata", "pathname of BL data storage directory")
    flag.Parse()

    LoadBLData() // load from existing files, if any
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
    muxRouter.HandleFunc("/info/{status}/{member}", handleInfoStatusMember).Methods("GET")

    muxRouter.HandleFunc("/add", handleAdd).Methods("GET")
    muxRouter.HandleFunc("/add/{loop}", handleAddLoop).Methods("GET")

    muxRouter.HandleFunc("/move/{serial}", handleMove).Methods("GET")
    muxRouter.HandleFunc("/post", handlePOST).Methods("POST")

    return muxRouter
}

func handleHome(w http.ResponseWriter, r *http.Request) {
    log.Println("handleHome() API called")
    io.WriteString(w, "You have entered the restricted zone. Trespassing is strictly prohibited. Defaulters will be reported.\n")
}

func handleInfoAll(w http.ResponseWriter, r *http.Request) {
    bytes, err := json.MarshalIndent(BLArray, "", "  ")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    io.WriteString(w, string(bytes))
}

func handleInfoStatusMember(w http.ResponseWriter, r *http.Request) {
    var BLElements []BL

    params := mux.Vars(r)
    status := params["status"]
    member, err := strconv.Atoi(params["member"])
    if err == nil {
        for i, _ := range BLArray {
            switch status {
            case "incoming":
                if BLArray[i].Owner < member { BLElements = append(BLElements, BLArray[i]) }
            case "owned":
                if BLArray[i].Owner == member { BLElements = append(BLElements, BLArray[i]) }
            case "sent":
                if BLArray[i].Owner > member { BLElements = append(BLElements, BLArray[i]) }
            case "archived":
                if BLArray[i].Owner == 5 { BLElements = append(BLElements, BLArray[i]) }
            default:
                // do nothing
            }
        }
        bytes, err := json.MarshalIndent(BLElements, "", "  ")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        io.WriteString(w, string(bytes))
    } else {
        respondWithJSON(w, r, http.StatusBadRequest, "Invalid Status or Member.")
    }
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
    newBL := updateBLArray()
    gobCheck(writeBLGob(BLArray, len(BLArray)))
    respondWithJSON(w, r, http.StatusCreated, newBL)
}

func handleAddLoop(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    loop, err := strconv.Atoi(params["loop"])
    var newBLs []BL

    if err == nil {
        for i := 0 ; i < loop ; i++ {
            newBL := updateBLArray()
            newBLs = append(newBLs, newBL)
            gobCheck(writeBLGob(BLArray, len(BLArray)))
        }
        respondWithJSON(w, r, http.StatusCreated, newBLs)
    }
}

func handleMove(w http.ResponseWriter, r *http.Request) {
    var bl BL

    params := mux.Vars(r)
    serial, err := strconv.Atoi(params["serial"])
    if err == nil {
        if serial <= len(BLArray) && serial > 0 {
            if BLArray[serial-1].Owner < *members {
                BLArray[serial-1].Owner = BLArray[serial-1].Owner + 1
                bl = BLArray[serial-1]
                gobCheck(writeBLGob(BLArray, len(BLArray)))
                respondWithJSON(w, r, http.StatusCreated, bl)
            }
        }
    } else {
        respondWithJSON(w, r, http.StatusBadRequest, "SerialNo invalid.")
    }
}

func handlePOST(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var bl BL

    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&bl); err != nil {
        respondWithJSON(w, r, http.StatusBadRequest, r.Body)
        return
    }
    defer r.Body.Close()

    bl.SerialNum = len(BLArray) + 1
    bl.Timestamp = time.Now().Format("02-01-2006 15:04:05 Mon")
    bl.Owner = 1
    log.Println("Adding to BLArray in handlePOST():", bl)

    BLArray = append(BLArray, bl)
    gobCheck(writeBLGob(BLArray, len(BLArray)))
    respondWithJSON(w, r, http.StatusCreated, bl)
}

///// LIST OF HELPER FUNCTIONS

func LoadBLData() { // load from existing files, if any
    if _, err := os.Stat(*bldataDir); os.IsNotExist(err) { // if *bldataDir does not exist
        log.Println("`", *bldataDir, "` does not exist. Creating directory.")
        os.Mkdir(*bldataDir, 0755) // https://stackoverflow.com/questions/14249467/os-mkdir-and-os-mkdirall-permission-value
    }

    files, err := ioutil.ReadDir(*bldataDir) // bldataDir from flag
    if err != nil { log.Fatal(err) }

    mostRecectFileNo := 0

    for _, file := range files {
        fileNo, _ := strconv.Atoi(file.Name()[len("bl-data")+1:len(file.Name())-4])
        if fileNo > mostRecectFileNo { mostRecectFileNo = fileNo }
    }

    if mostRecectFileNo == 0 {
        log.Println("No existing BLData")
    } else {
        mostRecentFile := *bldataDir + "/bl-data-" + strconv.Itoa(mostRecectFileNo) + ".gob"
        log.Println("Loading existing BLData from", mostRecentFile)
        gobCheck(readGob(&BLArray, mostRecentFile))
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
    filePath := *bldataDir + "/bl-data-" + strconv.Itoa(fileNoCount) + ".gob"
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

func updateBLArray() BL {
    newBL := BL{}

    newBL.SerialNum = len(BLArray) + 1
    newBL.BLuid = genRandInt(9000,60000)
    newBL.Timestamp = time.Now().Format("02-01-2006 15:04:05 Mon")

    newBL.Issuer = genRandIssuer()
    newBL.Releaser = newBL.Issuer
    newBL.Shipper = genRandShipper()
    newBL.Consignee = genRandConsignee()

    newBL.PortOfLoading = genRandPortOfLoading()
    newBL.PortOfDischarge = genRandPortOfDischarge()
    newBL.VesselNum = genRandInt(8000,1000)

    newBL.PkgQty = genRandInt(8,1)
    newBL.PkgType = genRandString(2) + strconv.Itoa(genRandInt(60,10))
    newBL.PkgDesc = genRandPkgDesc()
    newBL.PkgWeight = genRandInt(80,10)

    newBL.BLurl = "http://dropbox.com/" + genRandString(8) + ".pdf"

    newBL.Owner = 1

    BLArray = append(BLArray, newBL)
    for i, _ := range BLArray {
        if i < len(BLArray) - 1 { // ignore the newest element ie. newBL
            if BLArray[i].Owner < *members{
                randIncrement := int(math.Floor(float64(genRandInt(100,0))/70)) // 30% random increment
                BLArray[i].Owner = BLArray[i].Owner + randIncrement
            }
        }
    }
    fmt.Println("----BLArray---- len:", len(BLArray))
    spew.Dump(BLArray)
    return newBL
}

///// LIST OF TRIVIAL FUNCTIONS

func genRandIssuer() string {
    IssuerArray := [...]string{"Maersk", "Cosco", "Hapag-Lloyd", "Hamburg Süd", "DHL", "Panalpina", "Damco", "Geodis"}
    return  IssuerArray[genRandInt(len(IssuerArray),0)]
}

func genRandShipper() string {
    ShipperArray := [...]string{"Lee Latex", "HalcyonAgri", "Armstrong", "Kimberly-Clark", "Georgia-Pacific",
     "Weyerhaeuser", "Stora-Enso", "Corrie-MacColl"}
    return  ShipperArray[genRandInt(len(ShipperArray),0)]
}

func genRandConsignee() string {
    ConsigneeArray := [...]string{"MRF", "Bridgestone", "Apollo", "Goodyear", "Ceat", "Deenstone", "Viking",
    "Pulpo", "NewPaper", "Astrix", "Newman-Watson"}
    return  ConsigneeArray[genRandInt(len(ConsigneeArray),0)]
}

func genRandPortOfLoading() string {
    PortOfLoadingArray := [...]string{"Shanhai", "Hong Kong", "Shenzhen", "Guangzhou", "Rotterdam", "Colombo"}
    return  PortOfLoadingArray[genRandInt(len(PortOfLoadingArray),0)]
}

func genRandPortOfDischarge() string {
    PortOfDischargeArray := [...]string{"Singapore", "Antwerp", "Hamburg", "Los Angeles", "Long Beach", "Manila"}
    return  PortOfDischargeArray[genRandInt(len(PortOfDischargeArray),0)]
}

func genRandPkgDesc() string {
    PkgDescArray := [...]string{"Grade A Natural Rubber", "Grade B Natural Rubber", "Grade A Synthetic Rubber", 
        "Grade B Synthetic Rubber", "Grade A Natural Pulp", "Grade B Natural Pulp", "Grade A Natural Cotton"}
    return  PkgDescArray[genRandInt(len(PkgDescArray),0)]
}

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
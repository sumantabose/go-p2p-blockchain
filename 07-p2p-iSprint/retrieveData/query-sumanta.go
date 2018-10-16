package main

import (
	"log"
	"os"
	_ "fmt"
	"encoding/gob"
	"encoding/json"
	"runtime"
	"github.com/tidwall/gjson"
	"github.com/davecgh/go-spew/spew"

)

///// GLOBAL VARIABLES

// Raw Material Transaction (Type 1)
type RawMaterialTransaction struct {
	SerialNo int
	ProductCode string
	ProductName string
	ProductBatchNo string
	Quantity int
	RawMaterial
}
type RawMaterial struct {
	RawMaterialBatchNo string
	RawMaterialsID string
	RawMaterialName string
	RawMaterialQuantity float32
	RawMaterialMeasurementUnit string
}

// Delivery Transaction (Type 2)
type DeliveryTransaction struct {
	SerialNo int
	RecGenerator string
	ShipmentID string
	Timestamp string	
	Longitude string
	Latitude string
	ShippedFromCompID string
	ShippedToCompID string
	LocationID string
	DeliveryStatus string 	
	DeliveryType string
	Product
	Document
}
type Product struct {
	ProductCode string
	ProductName string
	ProductBatch
}
type ProductBatch struct {
	ProductBatchNo string
	ProductBatchQuantity int
}
type Document struct {
	DocumentURL string
	DocumentType string
	DocumentHash string
	DocumentSign string
}

// Block represents each 'item' in the blockchain
type Block struct {
	Index		int
	Timestamp	string
	TxnType		int
	TxnPayload	interface{}
	Comment		string
	Proposer	string
	PrevHash 	string
	ThisHash	string
}

// Standard Input takes incoming string from terminal for `comment` in block (Type 0)
type StdInput struct {
	Comment string
}

// Blockchain is a series of validated Blocks
var Blockchain []Block

//////// FUNCTIONS

func init() {
	log.SetFlags(log.Lshortfile) // Idea from https://appliedgo.net/networking/
	gob.Register(RawMaterialTransaction{})
	gob.Register(DeliveryTransaction{})
	gob.Register(map[string]interface{}{})
}

func main() {
	dataFile := "../data5000/blockchain-14.gob"
	log.Println("Loading Blockchain from", dataFile)
	gobCheck(readGob(&Blockchain, dataFile))
	//log.Println(Blockchain)
	//spew.Dump(Blockchain)
	//queryRaw("TxnPayload.SerialNo", "88")
	spew.Dump(query("raw", "TxnPayload.SerialNo", "88"))
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

//////////////////



func queryRaw(field string, value string) {

	var rawTxns []interface{}

	for _, block := range Blockchain {
		if block.TxnType == 1 {
			blockBytes, _ := json.Marshal(block)
			blockString := string(blockBytes)

			val := gjson.Get(blockString, "TxnPayload.SerialNo")
			log.Printf("val %T, %s\n", val, val)

			log.Println(field, value)

			gg := gjson.Get(blockString, "TxnPayload.SerialNo").String()
			log.Printf("gg %T, %s\n", gg, gg)


			if gjson.Get(blockString, field).String() == value {
				log.Println("here")
				rawTxns = append(rawTxns, block.TxnPayload)
				// log.Printf("%T\n", block)
				// log.Printf("%T\n", block.TxnPayload)
				// rawTxn = block.TxnPayload
				// log.Printf("%T\n", rawTxn)
			}
		}
	}

	spew.Dump(rawTxns)
	//return txnArray
}


func query(txnType string, field string, value string) []interface {} {

	txnTypeMap := map[string]int {
        "raw": 1,
        "del": 2,
    }

	var txnArray []interface{}
	for _, block := range Blockchain {
		if block.TxnType == txnTypeMap[txnType] {
			blockBytes, _ := json.Marshal(block)
			if gjson.Get(string(blockBytes), field).String() == value {
				txnArray = append(txnArray, block.TxnPayload)
			}
		}
	}
	return txnArray
}








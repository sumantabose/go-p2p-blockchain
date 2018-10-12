package main

import (
	"log"
	"os"
	"encoding/gob"
	"runtime"
	"github.com/fatih/structs"
	"github.com/davecgh/go-spew/spew"
	"encoding/json"
//	"github.com/tidwall/gjson"
//	"strconv"
	

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

func init() { // Idea from https://appliedgo.net/networking/
	log.SetFlags(log.Lshortfile)
	gob.Register(RawMaterialTransaction{})
	gob.Register(DeliveryTransaction{})
    gob.Register(map[string]interface{}{})
}


func Query(field string, value interface{}) {
	var block Block
	//var rawtx RawMaterialTransaction
	

	for _,block = range Blockchain {
		s := structs.New(block)
		if s.Field(field).Value() == value {
			
			//log.Println(block)
			spew.Dump(block)
			// log.Println("\n")

			// log.Println((s.Field("TxnPayload")).Value())

			// raw := s.Field("TxnPayload").Value()

			// byte, _ := json.Marshal(raw)
			// _ = json.Unmarshal(byte, &rawtx)

			// //log.Println(rawtx.ProductCode)
			// log.Println(rawtx)

			
		}
	}
} 

func QueryTxInfo(txnType int, field string, value interface{}){
	var block Block
	var raw RawMaterialTransaction
	var del DeliveryTransaction
	
	for _,block = range Blockchain {
		s := structs.New(block)
		if s.Field("TxnType").Value() == txnType {

			if txnType == 1 {
				byte, _ := json.Marshal(s.Field("TxnPayload").Value())
				_ = json.Unmarshal(byte, &raw)
				// raw := s.Field("TxnPayload").Value()
				// //raw1,_ := raw.(RawMaterialTransaction)
				rawtx := structs.New(raw)
				// log.Println(rawtx)

				// log.Println("rawtx field value",rawtx.Field(field).Value())
				// log.Println("\n")
				// log.Println("Given field value",value)

				if rawtx.Field(field).Value() == value {
					//log.Println("Hiii")
					log.Println(block)
				}

			} else if txnType == 2 {
				byte, _ := json.Marshal(s.Field("TxnPayload").Value())
				_ = json.Unmarshal(byte, &del)
				// raw := s.Field("TxnPayload").Value()
				// //raw1,_ := raw.(RawMaterialTransaction)
				deltx := structs.New(del)
				// log.Println(deltx)

				// log.Println("rawtx field value",deltx.Field(field).Value())
				// log.Println("\n")
				// log.Println("Given field value",value)

				if deltx.Field(field).Value() == value {
					//log.Println("Hiii")
					log.Println(block)
				}
				
			}

		}
	}
}

func main() {
	dataFile := "../data5000/blockchain-14.gob"
	log.Println("Loading Blockchain from", dataFile)

	gobCheck(readGob(&Blockchain, dataFile))

	 //Query("TxnType", 0)
	 Query("TxnType", 1)
	 // QueryTxInfo(1,"Quantity",15)
	 // QueryTxInfo(1,"ProductCode","12")
	 // QueryTxInfo(2,"SerialNo",70)
	//QueryTxInfo(2,"ShipmentID","msc")
	
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
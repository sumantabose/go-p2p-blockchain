package main

import (
	"log"
	"io/ioutil"
	"time"
	"encoding/json"
	"bytes"
	"net/http"
	"os"

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

func init() { // Idea from https://appliedgo.net/networking/
	log.SetFlags(log.Lshortfile)
}

func main() {

	for {
		time.Sleep(3 * time.Second)
		rmTxn := RawMaterialTransaction{}
		Post(rmTxn, "localhost", "5000", "raw")
		time.Sleep(3 * time.Second)
		delTxn := DeliveryTransaction{}
		Post(delTxn, "localhost", "5000", "delivery")
	}
}

func Post(object interface{}, IP string, Port string, Tag string) {
	log.Println("Posting")

	jsonObject, err := json.Marshal(object)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	url := "http://" + IP + ":" + Port + "/" + Tag
	log.Println("Posting to " + url)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonObject))
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer response.Body.Close()

    log.Println("response Status:", response.Status)
    log.Println("response Headers:", response.Header)
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Println(err)
        os.Exit(1)
    }
    log.Println("response Body:", string(body))
}
package main

import (
    "crypto/sha256"

	"encoding/hex"
	"encoding/gob"
	"strconv"
	"time"
	"log"
	"strings"
	"encoding/json"
	"github.com/tidwall/gjson"

	"github.com/davecgh/go-spew/spew"
)

var LastSentBlockchainLen = 0
var LastRcvdBlockchainLen = 0

// make sure block is valid by checking index, and comparing the hash of the previous block
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		log.Println("BLOCKCHAIN ERROR: Index Mismatch")
		return false
	}

	if oldBlock.ThisHash != newBlock.PrevHash {
		log.Println("BLOCKCHAIN ERROR: Hash Inconsistent")
		return false
	}

	if calculateHash(newBlock) != newBlock.ThisHash {
		log.Println("BLOCKCHAIN ERROR: Hash Mismatch")
		return false
	}

	return true
}

// SHA256 hashing
func calculateHash(b Block) string {
	// record := strconv.Itoa(block.Index) + block.Timestamp + block.Comment + block.PrevHash // old
	record := strconv.Itoa(b.Index) + b.Timestamp + strconv.Itoa(b.TxnType) + spew.Sdump(b.TxnPayload) + b.Comment + b.Proposer + b.PrevHash
	if *verbose { log.Println("Hashing:", record) }
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// create a new block using previous block's hash
func generateBlock(oldBlock Block, comment string, txnPayload interface{}, txnType int) Block {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Comment = comment
	newBlock.Proposer = thisPeerFullAddr

	newBlock.TxnType = txnType
	newBlock.TxnPayload = txnPayload

	newBlock.PrevHash = oldBlock.ThisHash
	newBlock.ThisHash = calculateHash(newBlock)

	return newBlock
}

func generateGenesisBlock() Block {
	genesisBlock := Block{0, time.Now().String(), 0, "", "Genesis Block", thisPeerFullAddr, "BIG-BANG!", ""}
	genesisBlock.ThisHash = calculateHash(genesisBlock)
	return genesisBlock
}

func registerGOB() {
	gob.Register(RawMaterialTransaction{})
	gob.Register(DeliveryTransaction{})
	gob.Register(RawMaterial{})
	gob.Register(Product{})
	gob.Register(ProductBatch{})
	// gob.Register(interface{})
	gob.Register(map[string]interface{}{})
	gob.Register(map[interface{}]interface{}{})
	gob.Register(map[string]map[string]interface{}{})
	gob.Register(map[interface{}]map[string]interface{}{})
}

func query(txnType string, field string, value string) []interface {} {
	var txnArray []interface{}
	txnTypeMap := map[string]int {
        "raw": 1,
        "del": 2,
    }

    if *verbose { log.Println("field = ", field) }

    newField := ""
    for i, subField := range strings.SplitAfter(field, ".") {
        if i > 0 {
        	newField = newField + "#."
        }
        newField = newField + subField
    }
    if *verbose { log.Println("newField = ", newField) }

	mutex.Lock()
	tempBlockchain := Blockchain
	mutex.Unlock()
	for _, block := range tempBlockchain {
		if block.TxnType == txnTypeMap[txnType] {
			blockBytes, _ := json.Marshal(block)

			result := gjson.Get(string(blockBytes), "TxnPayload."+newField)
			for _, item := range result.Array() {
				if (item.String() == value) {
					txnArray = append(txnArray, block.TxnPayload)
					break
				}
			}
		}
	}
	return txnArray
}

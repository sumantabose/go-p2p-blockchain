package main

import (
    "crypto/sha256"

	"encoding/hex"
	"strconv"
	"time"
	"log"

	"github.com/davecgh/go-spew/spew"
)

// make sure block is valid by checking index, and comparing the hash of the previous block
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		log.Println("BLOCKCHAIN ERROR: Index Mismatch")
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		log.Println("BLOCKCHAIN ERROR: Hash Inconsistent")
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
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

	newBlock.TxnType = txnType
	newBlock.TxnPayload = txnPayload

	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock
}
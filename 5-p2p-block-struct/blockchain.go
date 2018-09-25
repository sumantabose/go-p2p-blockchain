package main

import (
	"github.com/cnf/structhash"
    // Creates a hash of arbitrary Go datastructures.
    // Dump takes a data structure and returns its byte representation.
    // Godoc available at https://godoc.org/github.com/cnf/structhash

    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"

	"encoding/hex"
	"strconv"
	"time"
	"log"
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

	if calculateHashOld(newBlock) != newBlock.Hash {
		log.Println("BLOCKCHAIN ERROR: Hash Mismatch")
		return false
	}

	return true
}

// SHA256 hashing
func calculateHashOld(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + block.Comment + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Multi Hash Module (from zero.zero)
func calculateHash(object interface{}, hashType string) string {
	switch hashType {
		case "md5", "MD5":
			h := md5.New()
	        byteObject := structhash.Dump(object, 1)
	        h.Write(byteObject)
	        hashed := h.Sum(nil)
			return hex.EncodeToString(hashed)
	    case "sha1", "SHA1":
	  //       hashByte := sha1.Sum(structhash.Dump(object, 1))
	  //       hashString := hex.EncodeToString((hashByte[:]))
			// return hashString
			h := sha1.New()
	        byteObject := structhash.Dump(object, 1)
	        h.Write(byteObject)
	        hashed := h.Sum(nil)
			return hex.EncodeToString(hashed)
	    case "sha256", "SHA256":
			h := sha256.New()
	        byteObject := structhash.Dump(object, 1)
	        h.Write(byteObject)
	        hashed := h.Sum(nil)
			return hex.EncodeToString(hashed)
	    default:
	    	return "000" // return error in later iteration of code
	}
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
	newBlock.Hash = calculateHashOld(newBlock)

	return newBlock
}
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

func GenerateHash(block *Blockmember) string {
	record := block.Datestamp + string(block.Index) + block.Previoushash + block.Ordernumber
	hash := sha256.New()
	hash.Write([]byte(record))
	hashed := hash.Sum(nil)
	block.Hash = hex.EncodeToString(hashed)
	return hex.EncodeToString(hashed)
}
func GenerateBlock(index int, previousHash string, Ordernumber string) Blockmember {
	var block Blockmember
	block.Index = index
	block.Previoushash = previousHash
	block.Ordernumber = Ordernumber
	block.Datestamp = time.Now().String()
	block.Hash = GenerateHash(&block)
	return block
}
func VerifyBlock() bool {
	for i := 0; i < len(blocks); i++ {
		if i != len(blocks)-1 {
			if blocks[i+1].Previoushash != blocks[i].Hash {
				return false
			}
		}
		if blocks[i].Index != i {
			return false
		}
		if blocks[i].Hash != GenerateHash(&blocks[i]) {
			return false
		}
	}
	return true
}

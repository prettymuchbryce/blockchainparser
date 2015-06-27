package parser

import (
	"blockchainparser/src/utils"
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type Block struct {
	length            uint32
	version           uint32
	hash              [32]byte
	previousBlockHash [32]byte
	merkleRoot        [32]byte
	timestamp         uint32
	difficulty        uint32
	nonce             uint32
	transactionCount  uint64
}

//A block's hash is not included in the binary
//data on disk. It needs to be computed from
//some of the fields by doing a sha256(sha256(fields))
func (block *Block) ComputeHash() {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, block.version)
	binary.Write(buffer, binary.LittleEndian, block.previousBlockHash)
	binary.Write(buffer, binary.LittleEndian, block.merkleRoot)
	binary.Write(buffer, binary.LittleEndian, block.timestamp)
	binary.Write(buffer, binary.LittleEndian, block.difficulty)
	binary.Write(buffer, binary.LittleEndian, block.nonce)

	shaHash := sha256.New()
	shaHash.Write(buffer.Bytes())
	var hash []byte = shaHash.Sum(nil)

	shaHash2 := sha256.New()
	shaHash2.Write(hash)
	var hash2 []byte = shaHash2.Sum(nil)

	copy(block.hash[:], hash2)
}

//Returns a string of the big endian hex hash of this block
func (block *Block) getBigEndianString(value [32]byte) string {
	b := new(bytes.Buffer)

	for i := 0; i < 32; i++ {
		b.WriteByte(value[31-i])
	}

	return hex.EncodeToString(b.Bytes())
}

func (block *Block) Save(db *sql.DB) error {
	fmt.Println(utils.GetBigEndianString(block.hash[:]))
	_, err := db.Exec(`INSERT INTO blocks(length, version, hash, previousBlockHash, merkleRoot, timestamp, difficulty, nonce, transactionCount)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`, block.length, block.version, utils.GetBigEndianString(block.hash[:]), utils.GetBigEndianString(block.previousBlockHash[:]), utils.GetBigEndianString(block.merkleRoot[:]), block.timestamp, block.difficulty, block.nonce, block.transactionCount)

	return err
}

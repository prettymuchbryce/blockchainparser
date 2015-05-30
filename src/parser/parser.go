package parser

import (
	"blockchainparser/src/utils"
	"database/sql"
	"encoding/binary"
	"os"

	_ "github.com/lib/pq"
)

var magicBytes = []byte{249, 190, 180, 217}

//Transactions table
//Block table
//Address Table (could be generated with the above)

//Endianness ? Little?

//magicID (4 bytes)
//Block length (4 bytes) (# of bytes ?)
//Version number (4 bytes)
//Previous block Hash (32 bytes)
//Merkle root (32 bytes)
//Timestamp (4 bytes)
//Difficulty (4 bytes)
//Nonce (4 bytes)
//Transaction count (variable)

type Transaction struct {
	version            int
	transactionVersion int
	inputCount         int
	transactionIndex   int
	inputScript        []byte
	sequence           int
	outputCount        int
	value              int
	outputScript       []byte
	lockTime           int
}

type Input struct {
}

type Output struct {
}

//---
//Version (4 bytes)
//Transaction version # (4 bytes)
//# of inputs (variable)
//Transaction Index (4 bytes) uint32
//Script length (variable)
//Script data (length bytes)
//Sequence # (4 bytes) uint32 always 0xFFFFFFFF
//Output count (variable)
//Value (8 bytes) uint64
//Script length (variable)
//Output script (length bytes)
//lock time (4 bytes) uint32 always 0
//---

/*
Once we have consumed the final transaction, this
brings us to the end of the logical block. However,
and this is important to note, we will not
necessarily bring us to the end of the physical
block!  The 'block length' specified at the
beginning of this block may actually go beyond the
end of the last transaction which was consumed.
That is why it is important that you read the entire
block into memory rather than just reading each
transaction and expecting the file pointer to be
in the correct location for the next block.
*/

var db *sql.DB

func Parse() {
	err := connect()

	defer db.Close()

	if err != nil {
		panic(err)
	}

	file, err := os.Open("../blk00000.dat")
	if err != nil {
		panic(err)
	}

	var blocks int = 0
	for {
		success, err := scrollToNextBlock(file)
		if err != nil {
			panic(err)
		}
		if success {
			blocks++
			parseNextBlock(file)
		}
	}
}

func connect() (err error) {
	db, err = sql.Open("postgres",
		"user=bryceneal dbname=blockchainparser connect_timeout=5 sslmode=disable")

	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func parseNextBlock(file *os.File) error {
	block := new(Block)
	binary.Read(file, binary.LittleEndian, &block.length)
	binary.Read(file, binary.LittleEndian, &block.version)
	binary.Read(file, binary.LittleEndian, &block.previousBlockHash)
	binary.Read(file, binary.LittleEndian, &block.merkleRoot)
	binary.Read(file, binary.LittleEndian, &block.timestamp)
	binary.Read(file, binary.LittleEndian, &block.difficulty)
	binary.Read(file, binary.LittleEndian, &block.nonce)

	block.ComputeHash()

	transactionCount, err := utils.GetVariableInteger(file)
	if err != nil {
		return err
	}

	block.transactionCount = transactionCount

	err = block.Save(db)

	if err != nil {
		return err
	}

	return nil
}

func scrollToNextBlock(file *os.File) (bool, error) {
	i := 0
	for i < len(magicBytes) {
		success, err := doesNextByteEqual(file, magicBytes[i])
		i++
		if err != nil {
			return false, err
		}
		if !success {
			return false, nil
		}
	}

	return true, nil
}

func doesNextByteEqual(file *os.File, value byte) (bool, error) {
	nextByte := make([]byte, 1)
	_, err := file.Read(nextByte)

	if err != nil {
		return false, err
	}

	if nextByte[0] != value {
		return false, err
	}

	return true, nil
}

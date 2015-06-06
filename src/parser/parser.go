package parser

import (
	"blockchainparser/src/utils"
	"database/sql"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"

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

//Things to do when limit is reached, or done.
//1. count wallets
//2. delete orphan blocks

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

	datFileNum := 0
	for {
		fmt.Println("===Next block dat file num")
		file, err := os.Open("./" + getBlockDatFileName(datFileNum))
		// file, err := os.Open("./blk00000.dat")
		if err != nil {
			panic(err)
		}

		defer file.Close()

		var blocks int = 0
		for {
			success, err := scrollToNextBlock(file)
			if err != nil {
				if err.Error() == "EOF" {
					datFileNum++
					break
				} else {
					panic(err)
				}
			}
			if success {
				fmt.Println("block #: " + strconv.Itoa(blocks))
				blocks++
				parseNextBlock(file)
			}
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

func getBlockDatFileName(count int) (name string) {
	if count < 10 {
		name = "0000" + strconv.Itoa(count)
	} else if count < 100 {
		name = "000" + strconv.Itoa(count)
	} else if count < 1000 {
		name = "00" + strconv.Itoa(count)
	} else if count < 10000 {
		name = "0" + strconv.Itoa(count)
	} else {
		name = strconv.Itoa(count)
	}
	return "blk" + name + ".dat"
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

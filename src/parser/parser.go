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
	fmt.Println("----")
	block := new(Block)
	binary.Read(file, binary.LittleEndian, &block.length)
	binary.Read(file, binary.LittleEndian, &block.version)
	binary.Read(file, binary.LittleEndian, &block.previousBlockHash)
	binary.Read(file, binary.LittleEndian, &block.merkleRoot)
	binary.Read(file, binary.LittleEndian, &block.timestamp)
	binary.Read(file, binary.LittleEndian, &block.difficulty)
	binary.Read(file, binary.LittleEndian, &block.nonce)

	block.ComputeHash()

	transactionCount, _, err := utils.GetVariableInteger(file)
	if err != nil {
		return err
	}

	block.transactionCount = transactionCount

	err = block.Save(db)

	fmt.Println("What the heck transaction counts", transactionCount)

	for i := uint64(0); i < transactionCount; i++ {
		_, err := getTransactionHash(file)
		if err != nil {
			return err
		}
		//fmt.Println(value)
	}

	if err != nil {
		return err
	}

	return nil
}

// func parseNextTransaction(file *os.File) error {
// 	transaction := New(Transaction)
// 	binary.Read(file, binary.LittleEndian, &transaction.version)

// 	inputCount, _, err := utils.GetVariableInteger(file)
// 	if err != nil {
// 		return err
// 	}

// 	for i := 0; i < inputCount; i++ {
// 		input := New(Input)

// 	}
// }

func getTransactionHash(file *os.File) (value []byte, err error) {
	//Transaction Version
	value, err = readByte(file, value, 4)
	if err != nil {
		return nil, err
	}

	//Number of inputs
	numInputs, varIntBytes, err := utils.GetVariableInteger(file)
	if err != nil {
		return nil, err
	}

	value = append(value, varIntBytes...)

	fmt.Println("What the heck inputs", numInputs)

	for i := uint64(0); i < numInputs; i++ {
		//hash
		// value, err = readByte(file, value, 32)
		// if err != nil {
		// 	return nil, err
		// }

		thash := make([]byte, 32)
		binary.Read(file, binary.LittleEndian, &thash)

		fmt.Println("Hash in input", utils.GetBigEndianString(thash))

		//index
		value, err = readByte(file, value, 4)
		if err != nil {
			return nil, err
		}

		//Script length
		varInt, varIntBytes, err := utils.GetVariableInteger(file)
		if err != nil {
			return nil, err
		}

		value = append(value, varIntBytes...)

		//Script
		value, err = readByte(file, value, varInt)
		if err != nil {
			return nil, err
		}

		//Sequence #
		value, err = readByte(file, value, 4)
		if err != nil {
			return nil, err
		}
	}

	//Number of outputs
	numOutputs, varIntBytes, err := utils.GetVariableInteger(file)
	if err != nil {
		return nil, err
	}

	value = append(value, varIntBytes...)
	fmt.Println("What the heck output?", numOutputs)

	for i := uint64(0); i < numOutputs; i++ {
		//Value (# of satoshis)
		// value, err = readByte(file, value, 8)
		// if err != nil {
		// 	return nil, err
		// }

		var satoshis uint64
		binary.Read(file, binary.LittleEndian, &satoshis)
		fmt.Println("Satoshis in output", satoshis)

		//output script length
		varInt, varIntBytes, err := utils.GetVariableInteger(file)
		if err != nil {
			return nil, err
		}

		value = append(value, varIntBytes...)

		//Output script
		value, err = readByte(file, value, varInt)
		if err != nil {
			return nil, err
		}
	}

	//Transaction lock time
	value, err = readByte(file, value, 4)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func readByte(file *os.File, buffer []byte, numBytes uint64) ([]byte, error) {
	value := make([]byte, numBytes)
	err := binary.Read(file, binary.LittleEndian, &value)

	if err != nil {
		return nil, err
	}

	return append(buffer, value...), nil
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

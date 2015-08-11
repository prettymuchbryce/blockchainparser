package parser

import (
	"blockchainparser/utils"
	"bytes"
	"database/sql"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

// Magic Bytes which denote the start of each block in the bitcoin protocol
var magicBytes = []byte{249, 190, 180, 217}

// A map of every transaction seen. Used to prevent parsing duplicate
// transactions in a block.
var transactionsSeen map[[32]byte]bool

// Last block seen tracks the very last hash of the very last block
// so that which know where to start iterating backwards for deleting orphan
// blocks.
var lastBlockSeen [32]byte

// A map of every block hashe, to it's parent's block hash.
// Used for deleting orphan blocks.
var blockHashToPreviousHash map[[32]byte][32]byte

var blockFinishChannel chan error

func Parse(dbuser string, path string) {
	blockHashToPreviousHash = make(map[[32]byte][32]byte)
	blockFinishChannel = make(chan error)

	// Counter for dat files
	datFileNum := 0
	datFileMax := 0
	for {
		fmt.Println("===Next block dat file num" + strconv.Itoa(datFileNum))
		if strings.LastIndex(path, "/") != len(path)-1 {
			path = strings.Join([]string{path, "/"}, "")
		}
		file, err := os.Open(path + getBlockDatFileName(datFileNum))
		if err != nil {
			datFileMax = datFileNum
			break
			// deleteOrphanBlocks()
			// return
		}

		datFileNum++

		go func(file *os.File) {
			db, err := connect(dbuser)
			if err != nil {
				panic(err)
			}

			defer db.Close()
			defer file.Close()

			var blocks int = 0
			transactionsSeen = make(map[[32]byte]bool)
			for {
				success, err := scrollToNextBlock(file)
				if err != nil {
					blockFinishChannel <- err
					break
					// if err.Error() == "EOF" {
					// 	break
					// } else {
					// 	panic(err)
					// }
				}
				if success {
					fmt.Println("block #: " + strconv.Itoa(blocks))
					blocks++
					err := parseNextBlock(file, db)
					if err != nil {
						blockFinishChannel <- err
						break
						// panic(err)
					}
				}
			}
		}(file)
	}

	datFileFinished := 0
	for {
		err := <-blockFinishChannel
		if err.Error() == "EOF" {
			datFileFinished++
		} else {
			panic(err)
		}

		if datFileFinished == datFileMax {
			// TODO delete orphan blocks
			fmt.Println("Finished parsing!")
			break
		}
	}

}

func connect(dbuser string) (db *sql.DB, err error) {
	db, err = sql.Open("postgres",
		"user="+dbuser+" dbname=blockchainparser connect_timeout=5 sslmode=disable")

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Parse the next block in the chain
func parseNextBlock(file *os.File, db *sql.DB) error {
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

	lastBlockSeen = block.hash
	blockHashToPreviousHash[block.hash] = block.previousBlockHash

	transactionCount, _, err := utils.GetVariableInteger(file)
	if err != nil {
		return err
	}

	block.transactionCount = transactionCount

	err = block.Save(db)
	if err != nil {
		return err
	}

	fmt.Println("transactionCount", transactionCount)

	for i := uint64(0); i < transactionCount; i++ {
		transaction, err := parseNextTransaction(file, db)
		if err != nil {
			return err
		}

		// There are a few cases early in the blockchain
		// where the same transaction is in multiple blocks.
		// The bitcoin network only cares about the first, and
		// the subsequent transactions are ignored.
		if transactionsSeen[transaction.hash] {
			continue
		}
		transactionsSeen[transaction.hash] = true
		err = transaction.Save(db)
		if err != nil {
			return err
		}
		fmt.Println("Transaction Inputs", transaction.inputs[0].hash)
	}

	return nil
}

// Parse the next transaction
func parseNextTransaction(file *os.File, db *sql.DB) (transaction *Transaction, err error) {
	transaction = new(Transaction)
	var transactionBytes []byte

	// Store TransactionBytes as we go along in addition to populating
	// the transaction struct. We need to do this in order to compute
	// the transaction hash at the end of this function.

	// Transaction Version
	transactionBytes, versionBytes, err := readByte(file, transactionBytes, 4)
	reader := bytes.NewReader(versionBytes)
	binary.Read(reader, binary.LittleEndian, &transaction.version)
	if err != nil {
		return transaction, err
	}

	// Number of inputs
	numInputs, numInputsBytes, err := utils.GetVariableInteger(file)
	if err != nil {
		return transaction, err
	}

	transactionBytes = append(transactionBytes, numInputsBytes...)

	fmt.Println("numInputs", numInputs)
	for i := uint64(0); i < numInputs; i++ {
		input := new(Input)
		transaction.inputs = append(transaction.inputs, input)

		// hash
		var hashBytes []byte
		transactionBytes, hashBytes, err = readByte(file, transactionBytes, 32)
		reader := bytes.NewReader(hashBytes)
		binary.Read(reader, binary.LittleEndian, &input.hash)
		if err != nil {
			return transaction, err
		}

		// index
		var indexBytes []byte
		transactionBytes, indexBytes, err = readByte(file, transactionBytes, 4)
		reader = bytes.NewReader(indexBytes)
		binary.Read(reader, binary.LittleEndian, &input.index)
		if err != nil {
			return transaction, err
		}

		// Script length
		inputScriptLength, inputScriptLengthBytes, err := utils.GetVariableInteger(file)
		if err != nil {
			return transaction, err
		}

		transactionBytes = append(transactionBytes, inputScriptLengthBytes...)

		// Script
		var scriptBytes []byte
		transactionBytes, scriptBytes, err = readByte(file, transactionBytes, inputScriptLength)
		reader = bytes.NewReader(scriptBytes)
		binary.Read(reader, binary.LittleEndian, &input.script)
		if err != nil {
			return transaction, err
		}

		// Sequence #
		var sequenceBytes []byte
		transactionBytes, sequenceBytes, err = readByte(file, transactionBytes, 4)
		reader = bytes.NewReader(sequenceBytes)
		binary.Read(reader, binary.LittleEndian, &input.sequence)
		if err != nil {
			return transaction, err
		}
	}

	// Number of outputs
	numOutputs, numOutputsBytes, err := utils.GetVariableInteger(file)
	if err != nil {
		return transaction, err
	}

	fmt.Println("numOutputs", numOutputs)

	transactionBytes = append(transactionBytes, numOutputsBytes...)

	for i := uint64(0); i < numOutputs; i++ {
		output := new(Output)
		transaction.outputs = append(transaction.outputs, output)

		// Value (# of satoshis)
		var valueBytes []byte
		transactionBytes, valueBytes, err = readByte(file, transactionBytes, 8)
		reader = bytes.NewReader(valueBytes)
		binary.Read(reader, binary.LittleEndian, &output.value)
		if err != nil {
			return transaction, err
		}

		// output script length
		outputScriptLength, outputScriptLengthBytes, err := utils.GetVariableInteger(file)
		if err != nil {
			return transaction, err
		}

		transactionBytes = append(transactionBytes, outputScriptLengthBytes...)

		// Output script
		var scriptBytes []byte
		transactionBytes, scriptBytes, err = readByte(file, transactionBytes, outputScriptLength)
		reader = bytes.NewReader(scriptBytes)
		binary.Read(reader, binary.LittleEndian, &output.script)
		if err != nil {
			return transaction, err
		}
		publicKey, err := utils.ExtractPublicKeyFromOutputScript(scriptBytes)
		var publicKeyBytes [20]byte
		copy(publicKeyBytes[:], publicKey[:])
		if err != nil {
			fmt.Println("Error! Can't find public key in output script")
			fmt.Println(err)
			//return transaction, err
		}
		output.publicKey = publicKeyBytes
	}

	// Transaction lock time
	transactionBytes, lockTimeBytes, err := readByte(file, transactionBytes, 4)
	reader = bytes.NewReader(lockTimeBytes)
	binary.Read(reader, binary.LittleEndian, &transaction.lock)
	if err != nil {
		return transaction, err
	}

	// Calculate the transaction hash
	dsha := utils.DoubleSha(transactionBytes)
	var transactionHashBytes [32]byte
	copy(transactionHashBytes[:], dsha[:])
	transaction.hash = transactionHashBytes
	fmt.Println(utils.GetBigEndianString(transaction.hash[:]))

	return transaction, nil
}

func deleteOrphanBlocks() {
	nextBlock := lastBlockSeen

	for {
		lastBlock := nextBlock
		nextBlock = blockHashToPreviousHash[nextBlock]
		delete(blockHashToPreviousHash, lastBlock)
		if isByteArrayZeroed(nextBlock[:]) {
			break
		}
	}

	fmt.Println("Orphan hashes")
	for k := range blockHashToPreviousHash {
		fmt.Println(k)
	}

}

func isByteArrayZeroed(a []byte) bool {
	for i := 0; i < len(a); i++ {
		if a[i] != byte(0) {
			return false
		}
	}

	return true
}

func readByte(file *os.File, buffer []byte, numBytes uint64) ([]byte, []byte, error) {
	value := make([]byte, numBytes)
	err := binary.Read(file, binary.LittleEndian, &value)

	if err != nil {
		return nil, nil, err
	}

	result := append(buffer, value...)

	return result, value, nil
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

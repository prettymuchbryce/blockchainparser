package parser

import (
	"blockchainparser/utils"
	"database/sql"
	"fmt"
)

//---
//Transaction version # (4 bytes)
//# of inputs (variable)
//  Hash
//  Transaction Index (4 bytes) uint32
//  Script length (variable)
//  Script data (length bytes)
//  Sequence # (4 bytes) uint32 always 0xFFFFFFFF
//Output count (variable)
//  Value (8 bytes) uint64
//  Script length (variable)
//  Output script (length bytes)
//lock time (4 bytes) uint32 always 0
//---

type Transaction struct {
	version uint32
	inputs  []*Input
	outputs []*Output
	lock    uint32
	hash    [32]byte
}

type Input struct {
	hash     [32]byte
	index    uint32
	script   []byte
	sequence uint32
}

type Output struct {
	value     uint64
	publicKey [20]byte
	script    []byte
}

func (transaction *Transaction) Save(db *sql.DB) error {
	_, err := db.Exec(`INSERT INTO transactions(hash, version, lock)
	VALUES($1, $2, $3)`, utils.GetBigEndianString(transaction.hash[:]), transaction.version, transaction.lock)
	if err != nil {
		return err
	}
	for i := 0; i < len(transaction.inputs); i++ {
		input := transaction.inputs[i]
		_, err = db.Exec(`INSERT INTO inputs(transaction, hash, index, script, sequence)
		VALUES($1, $2, $3, $4, $5)`, utils.GetBigEndianString(transaction.hash[:]), utils.GetBigEndianString(input.hash[:]), input.index, input.script, input.sequence)
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(transaction.outputs); i++ {
		output := transaction.outputs[i]
		fmt.Println(utils.ConvertPublicKeyToAscii(output.publicKey))
		_, err = db.Exec(`INSERT INTO outputs(transaction, publicKey, value, script)
		VALUES($1, $2, $3, $4)`, utils.GetBigEndianString(transaction.hash[:]), utils.ConvertPublicKeyToAscii(output.publicKey), output.value, output.script)
		if err != nil {
			return err
		}
	}

	return nil
}

package utils

import (
	"bytes"
	"encoding/hex"
)

//Returns a string of the big endian hex hash
func GetBigEndianString(value []byte) string {
	b := new(bytes.Buffer)

	for i := 0; i < 32; i++ {
		b.WriteByte(value[31-i])
	}

	return hex.EncodeToString(b.Bytes())
}

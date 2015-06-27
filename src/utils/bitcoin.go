package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"scriptcodes"

	"github.com/prettymuchbryce/hellobitcoin/base58check"

	"code.google.com/p/go.crypto/ripemd160"
)

//Returns a string of the big endian hex hash
func GetBigEndianString(value []byte) string {
	b := new(bytes.Buffer)

	for i := 0; i < len(value); i++ {
		b.WriteByte(value[len(value)-1-i])
	}

	return hex.EncodeToString(b.Bytes())
}

func Convert20BytePublicKeyToAscii(key [20]byte) string {
	mainNetPrefix := "00"
	asciiKey := base58check.Encode(mainNetPrefix, key[:])
	return asciiKey
}

func ConvertLongPublicKeyToShortPublicKey(key []byte) (newKey []byte) {
	shaHash := sha256.New()
	shaHash.Write(key)
	shadPublicKeyBytes := shaHash.Sum(nil)

	ripeHash := ripemd160.New()
	ripeHash.Write(shadPublicKeyBytes)
	ripeHashedBytes := ripeHash.Sum(nil)

	return ripeHashedBytes
}

func ExtractPublicKeyFromOutputScript(script []byte) (key []byte, err error) {
	if len(script) == 67 {
		//67 byte long output script containing a full ECDSA 65 byte public key address.
		for {
			if script[0] != byte(65) || script[66] != scriptcodes.OP_CHECKSIG {
				break
			}
			return ConvertLongPublicKeyToShortPublicKey(script[1:65]), nil
		}
	} else if len(script) == 66 {
		// 66 byte long output script.  Contains a 65 byte public key address.
		for {
			if script[65] != scriptcodes.OP_CHECKSIG {
				break
			}
			return ConvertLongPublicKeyToShortPublicKey(script[0:64]), nil
		}
	}

	if len(script) > 25 {
		// Script is 25 bytes long or more, contains a 20 byte public key hash address.
		for {
			if script[0] != scriptcodes.OP_DUP || script[1] != scriptcodes.OP_HASH160 {
				break
			}

			if script[2] != byte(20) {
				break
			}

			return script[3:24], nil
		}
	}

	//Search by pattern
	for i := 0; i < len(script); i++ {
		if i+24 > len(script) {
			break
		}

		if script[i] != scriptcodes.OP_DUP {
			continue
		}

		if script[i+1] != scriptcodes.OP_HASH160 {
			continue
		}

		if script[i+2] != byte(20) {
			continue
		}

		if script[i+23] != scriptcodes.OP_EQUALVERIFY {
			continue
		}

		if script[i+24] != scriptcodes.OP_CHECKSIG {
			continue
		}

		return script[i+3 : i+22], nil
	}
	return nil, errors.New("key not found in script")
}

func DoubleSha(value []byte) (finalHash []byte) {
	shaHash := sha256.New()
	shaHash.Write(value)
	var hash []byte = shaHash.Sum(nil)

	shaHash2 := sha256.New()
	shaHash2.Write(hash)
	var hash2 []byte = shaHash2.Sum(nil)

	return hash2
}

package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

func GetVariableInteger(reader io.Reader) (value uint64, err error) {
	firstByte := make([]byte, 1)
	_, err = reader.Read(firstByte)

	if err != nil {
		return 0, err
	}

	if firstByte[0] < 253 {
		return uint64(firstByte[0]), nil
	}

	if firstByte[0] == 253 {
		twoBytes := make([]byte, 2)
		_, err = reader.Read(twoBytes)

		if err != nil {
			return 0, err
		}

		return binary.ReadUvarint(bytes.NewReader(twoBytes))
	}

	if firstByte[0] == 254 {
		fourBytes := make([]byte, 4)
		_, err = reader.Read(fourBytes)

		if err != nil {
			return 0, err
		}

		return binary.ReadUvarint(bytes.NewReader(fourBytes))
	}

	if firstByte[0] == 255 {
		eightBytes := make([]byte, 8)
		_, err = reader.Read(eightBytes)

		if err != nil {
			return 0, err
		}

		return binary.ReadUvarint(bytes.NewReader(eightBytes))
	}

	return 0, errors.New("Unexpected value for variable integer")
}

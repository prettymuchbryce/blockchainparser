package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

func GetVariableInteger(reader io.Reader) (value uint64, b []byte, err error) {
	firstByte := make([]byte, 1)
	_, err = reader.Read(firstByte)

	if err != nil {
		return 0, nil, err
	}

	if firstByte[0] < 253 {
		return uint64(firstByte[0]), firstByte, nil
	}

	if firstByte[0] == 253 {
		twoBytes := make([]byte, 2)
		_, err = reader.Read(twoBytes)

		if err != nil {
			return 0, nil, err
		}

		value, err = binary.ReadUvarint(bytes.NewReader(twoBytes))

		return value, append(firstByte, twoBytes...), err
	}

	if firstByte[0] == 254 {
		fourBytes := make([]byte, 4)
		_, err = reader.Read(fourBytes)

		if err != nil {
			return 0, nil, err
		}

		value, err = binary.ReadUvarint(bytes.NewReader(fourBytes))

		return value, append(firstByte, fourBytes...), err
	}

	if firstByte[0] == 255 {
		eightBytes := make([]byte, 8)
		_, err = reader.Read(eightBytes)

		if err != nil {
			return 0, nil, err
		}

		value, err = binary.ReadUvarint(bytes.NewReader(eightBytes))

		return value, append(firstByte, eightBytes...), err
	}

	return 0, nil, errors.New("Unexpected value for variable integer")
}

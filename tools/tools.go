package tools

import (
	"bytes"
	"encoding/binary"
)

func ByteTouint16(bytes []byte) uint16 {
	return (uint16(bytes[0]) << 8) | uint16(bytes[1])
}

func Uint16ToByteArr(val uint16) []byte {
	return []byte{byte(val >> 8), byte(val)}
}

func ByteArraySetArr(bytes *[16]byte, ind int, arr []byte) {
	for i, item := range arr {
		(*bytes)[ind+i] = item
	}
}

func NameToByteArr(name []string) ([]byte, error) {
	retBytes := new(bytes.Buffer)

	var err error
	for _, item := range name {
		err = binary.Write(retBytes, binary.BigEndian, uint8(len(item)))
		if err != nil {
			return nil, err
		}

		err = binary.Write(retBytes, binary.BigEndian, []byte(item))
		if err != nil {
			return nil, err
		}
	}

	err = binary.Write(retBytes, binary.BigEndian, uint8(0))
	if err != nil {
		return nil, err
	}

	return retBytes.Bytes(), nil
}

func TXTToByteArr(value string) ([]byte, error) {
	retBytes := new(bytes.Buffer)
	var err error

	err = binary.Write(retBytes, binary.BigEndian, uint8(len(value)))
	if err != nil {
		return nil, err
	}

	err = binary.Write(retBytes, binary.BigEndian, []byte(value))
	if err != nil {
		return nil, err
	}

	return retBytes.Bytes(), nil
}

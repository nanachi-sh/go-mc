package util

import (
	"bytes"
	"encoding/binary"
)

func PutVarint(i int) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, uint64(i))
	return buf[:n]
}

func ReadVarint(b []byte) (uint64, error) {
	return binary.ReadUvarint(bytes.NewBuffer(b))
}

func Int(i int) *int {
	return &i
}

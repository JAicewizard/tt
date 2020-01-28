package v3

import (
	"encoding/binary"
	"math"
)

func Float32ToBytes(f *float32) []byte {
	var buf [4]byte
	bits := math.Float32bits(*f)
	binary.LittleEndian.PutUint32(buf[:], bits)
	return buf[:]
}

func Float32FromBytes(buf []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(buf))
}

func Float64ToBytes(f float64, buf *[8]byte) {
	bits := math.Float64bits(f)
	binary.LittleEndian.PutUint64(buf[:], bits)
}

func Float64FromBytes(buf []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(buf))
}

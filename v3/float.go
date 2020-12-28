package v3

import (
	"encoding/binary"
	"math"
)

//Float32ToBytes converts a float32 into bytes
func Float32ToBytes(f float32, buf []byte) {
	bits := math.Float32bits(f)
	binary.LittleEndian.PutUint32(buf, bits)
}

//Float32FromBytes bytes into a float32
func Float32FromBytes(buf []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(buf))
}

//Float64ToBytes converts a float64 into bytes
func Float64ToBytes(f float64, buf []byte) {
	bits := math.Float64bits(f)
	binary.LittleEndian.PutUint64(buf, bits)
}

//Float64FromBytes bytes into a float64
func Float64FromBytes(buf []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(buf))
}

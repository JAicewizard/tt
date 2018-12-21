package tt

import (
	"encoding/binary"
	"math"
	"reflect"
	"unsafe"
)

func stringToBytes(s string) []byte {
	x := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: x.Data, Len: x.Len, Cap: x.Len}))
}

func stringFromBytes(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func float32ToBytes(f float32) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint32(buf[:], math.Float32bits(f))
	return buf[:]
}

func float32FromBytes(buf []byte) float32 {
	return math.Float32frombits((binary.LittleEndian.Uint32(buf[:])))
}

func float64ToBytes(f float64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}

func float64FromBytes(buf []byte) float64 {
	return math.Float64frombits((binary.LittleEndian.Uint64(buf[:])))
}

func int64ToBytes(f int64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(f))
	return buf[:]
}

func int64FromBytes(buf []byte) int64 {
	return int64(binary.LittleEndian.Uint64(buf[:]))
}

func int32ToBytes(f int32) []byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(f))
	return buf[:]
}

func int32FromBytes(buf []byte) int32 {
	return int32(binary.LittleEndian.Uint32(buf[:]))
}

func int16ToBytes(f int16) []byte {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], uint16(f))
	return buf[:]
}

func int16FromBytes(buf []byte) int16 {
	return int16(binary.LittleEndian.Uint16(buf[:]))
}

func int8ToBytes(f int8) []byte {
	buf := [1]byte{byte(f)}
	return buf[:]
}

func int8FromBytes(buf byte) int8 {
	return int8(buf)
}

func uint64ToBytes(f uint64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], f)
	return buf[:]
}

func uint64FromBytes(buf []byte) uint64 {
	return binary.LittleEndian.Uint64(buf[:])
}

func uint32ToBytes(f uint32) []byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], f)
	return buf[:]
}

func uint32FromBytes(buf []byte) uint32 {
	return binary.LittleEndian.Uint32(buf[:])
}

func uint16ToBytes(f uint16) []byte {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], f)
	return buf[:]
}

func uint16FromBytes(buf []byte) uint16 {
	return binary.LittleEndian.Uint16(buf[:])
}

func uint8ToBytes(f uint8) []byte {
	buf := [1]byte{byte(f)}
	return buf[:]
}

func uint8FromBytes(buf byte) uint8 {
	return uint8(buf)
}

func ikeytobytes(key ikeytype) (buf [ikeylen]byte) {
	binary.LittleEndian.PutUint32(buf[:], uint32(key))
	return
}

func ikeyfrombytes(buf []byte) ikeytype {
	return ikeytype(binary.LittleEndian.Uint32(buf[:]))
}

func valuelentobyte(key valuelen) (buf [valuelenbytes]byte) {
	binary.LittleEndian.PutUint32(buf[:], uint32(key))
	return
}

func getvaluelen(buf []byte) valuelen {
	return valuelen(binary.LittleEndian.Uint32(buf[:]))
}

func keylentobyte(key keylen) (buf [keylenbytes]byte) {
	binary.LittleEndian.PutUint32(buf[:], uint32(key))
	return
}

func getkeylen(buf []byte) keylen {
	return keylen(binary.LittleEndian.Uint32(buf[:]))
}

package v3

import (
	"encoding/binary"
)

func Int8ToBytes(u int8) byte {
	return byte(u)
}

func Int8FromBytes(buf byte) int8 {
	return int8(buf)
}

func Int16ToBytes(f int16) []byte {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], uint16(f))
	return buf[:]
}

func Int16FromBytes(buf []byte) int16 {
	return int16(binary.LittleEndian.Uint16(buf))
}

func Int32ToBytes(f int32) []byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(f))
	return buf[:]
}

func Int32FromBytes(buf []byte) int32 {
	return int32(binary.LittleEndian.Uint32(buf))
}

//Passing in aa byte buffer only works for 8byte long values
func Int64ToBytes(f int64, buf *[8]byte) {
	binary.LittleEndian.PutUint64(buf[:], uint64(f))
}

func Int64FromBytes(buf []byte) int64 {
	return int64(binary.LittleEndian.Uint64(buf))
}

func Uint8ToBytes(u uint8) byte {
	return byte(u)
}

func Uint8FromBytes(buf byte) uint8 {
	return uint8(buf)
}

func Uint16ToBytes(f uint16) []byte {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], f)
	return buf[:]
}

func Uint16FromBytes(buf []byte) uint16 {
	return binary.LittleEndian.Uint16(buf)
}

func Uint32ToBytes(f uint32) []byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], f)
	return buf[:]
}

func Uint32FromBytes(buf []byte) uint32 {
	return binary.LittleEndian.Uint32(buf)
}

func Uint64ToBytes(f uint64, buf *[8]byte) {
	binary.LittleEndian.PutUint64(buf[:], f)
}

func Uint64FromBytes(buf []byte) uint64 {
	return binary.LittleEndian.Uint64(buf)
}

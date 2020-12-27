package v3

import (
	"encoding/binary"
)

//Int8ToBytes converts a int8 into bytes
func Int8ToBytes(u int8) byte {
	return byte(u)
}

//Int8FromBytes converts bytes into an int8
func Int8FromBytes(buf byte) int8 {
	return int8(buf)
}

//Int16ToBytes converts a int16 into bytes
func Int16ToBytes(f int16) []byte {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], uint16(f))
	return buf[:]
}

//Int16FromBytes converts bytes into an int16
func Int16FromBytes(buf []byte) int16 {
	return int16(binary.LittleEndian.Uint16(buf))
}

//Int32ToBytes converts a int32 into bytes
func Int32ToBytes(f int32) []byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(f))
	return buf[:]
}

//Int32FromBytes converts bytes into an int32
func Int32FromBytes(buf []byte) int32 {
	return int32(binary.LittleEndian.Uint32(buf))
}

//Int64ToBytes converts a int64 into bytes
func Int64ToBytes(f int64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(f))
	return buf[:]
}

//Int64FromBytes converts bytes into an int64
func Int64FromBytes(buf []byte) int64 {
	return int64(binary.LittleEndian.Uint64(buf))
}

//Uint8ToBytes converts a uint8 into bytes
func Uint8ToBytes(u uint8) byte {
	return byte(u)
}

//Uint8FromBytes converts bytes into an uint8
func Uint8FromBytes(buf byte) uint8 {
	return uint8(buf)
}

//Uint16ToBytes converts a uint16 into bytes
func Uint16ToBytes(f uint16) []byte {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], f)
	return buf[:]
}

//Uint16FromBytes converts bytes into an uint16
func Uint16FromBytes(buf []byte) uint16 {
	return binary.LittleEndian.Uint16(buf)
}

//Uint32ToBytes converts a uint32 into bytes
func Uint32ToBytes(f uint32) []byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], f)
	return buf[:]
}

//Uint32FromBytes converts bytes into an uint32
func Uint32FromBytes(buf []byte) uint32 {
	return binary.LittleEndian.Uint32(buf)
}

//Uint64ToBytes converts a uint64 into bytes
func Uint64ToBytes(f uint64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(f))
	return buf[:]
}

//Uint64FromBytes converts bytes into an uint64
func Uint64FromBytes(buf []byte) uint64 {
	return binary.LittleEndian.Uint64(buf)
}

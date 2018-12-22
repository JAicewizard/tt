package tt

import (
	"encoding/binary"
	"reflect"
	"unsafe"
)

func int64ToBytes(f int64, buf *[8]byte) {
	binary.LittleEndian.PutUint64(buf[:], uint64(f))
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

func int8ToBytes(f *int8, buf *[]byte) {
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(f)), Len: 1, Cap: 1} // we dont care about byte order since its only one byte!!!
	*buf = *(*[]byte)(unsafe.Pointer(&hdr))
}

func int8FromBytes(buf []byte) int8 {
	hrd := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	int := *(*int8)(unsafe.Pointer(hrd.Data))
	return int
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

func uint8ToBytes(f *uint8, buf *[]byte) {
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(f)), Len: 1, Cap: 1} // we dont care about byte order since its only one byte!!!
	*buf = *(*[]byte)(unsafe.Pointer(&hdr))
}

func uint8FromBytes(buf byte) uint8 {
	hrd := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	int := *(*uint8)(unsafe.Pointer(hrd.Data))
	return int
}

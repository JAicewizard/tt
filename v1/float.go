package v1

import (
	"reflect"
	"unsafe"
)

func Float32ToBytes(f *float32, ret *[]byte) {
	*ret = (*(*[4]byte)(unsafe.Pointer(f)))[:]
}

func Float32ToBytes3(f *float32, ret *[]byte) {
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(f)), Len: 4, Cap: 4}
	*ret = *(*[]byte)(unsafe.Pointer(&hdr))
}

func Float32FromBytes(buf []byte) float32 {
	hrd := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	f := *(*float32)(unsafe.Pointer(hrd.Data))
	return f
}

func Float64ToBytes(f *float64, ret *[]byte) {
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(f)), Len: 8, Cap: 8}
	*ret = *(*[]byte)(unsafe.Pointer(&hdr))
}

func Float64FromBytes(buf []byte) float64 {
	hrd := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	f := *(*float64)(unsafe.Pointer(hrd.Data))
	return f
}

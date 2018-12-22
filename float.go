package tt

import (
	"reflect"
	"unsafe"
)

func float32ToBytes(f *float32, ret *[]byte) {
	*ret = (*(*[4]byte)(unsafe.Pointer(f)))[:]
}
func float32ToBytes3(f *float32, ret *[]byte) {
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(f)), Len: 4, Cap: 4}
	*ret = *(*[]byte)(unsafe.Pointer(&hdr))
}
func float32FromBytes(buf []byte) float32 {
	hrd := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	f := *(*float32)(unsafe.Pointer(hrd.Data))
	return f
}

func float64ToBytes(f *float64, ret *[]byte) {
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(f)), Len: 8, Cap: 8}
	*ret = *(*[]byte)(unsafe.Pointer(&hdr))
}

func float64FromBytes(buf []byte) float64 {
	hrd := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	f := *(*float64)(unsafe.Pointer(hrd.Data))
	return f
}

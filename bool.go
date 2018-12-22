package tt

import (
	"reflect"
	"unsafe"
)

func boolToBytes(b *bool, buf *[]byte) {
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(b)), Len: 1, Cap: 1} // we dont care about byte order since its only one byte!!!
	*buf = *(*[]byte)(unsafe.Pointer(&hdr))
}

func boolFromBytes(buf []byte) bool {
	hrd := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	b := *(*bool)(unsafe.Pointer(hrd.Data))
	return b
}

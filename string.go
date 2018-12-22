package tt

import (
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

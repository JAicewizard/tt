package v3

import (
	"reflect"
	"unsafe"
)

//StringToBytes converts a string to bytes
func StringToBytes(s string) []byte {
	x := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: x.Data, Len: x.Len, Cap: x.Len}))
}

//StringFromBytes converts bytes to a string
func StringFromBytes(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

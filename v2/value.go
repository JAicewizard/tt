package v2

import (
	"bytes"
	"encoding/binary"
)

type (
	//Value is the type that contains the value data
	Value struct {
		Key      Key        //the key of this object in its final form
		Children []Ikeytype //the key to the values in the data send across
		Value    []byte     //the key to the data of this object in its final form
		Vtype    byte       //the type of the data of this object in its final form
	}
)

func (v *Value) Len() int {
	if v.Key.Value == nil {
		return len(v.Value) + len(v.Children)*ikeylen + 2 + valuelenbytes + keylenbytes
	}
	return len(v.Value) + len(v.Children)*ikeylen + len(v.Key.Value) + 3 + valuelenbytes + keylenbytes
}

func AddValue(slice *bytes.Buffer, v *Value) {
	slice.Grow(v.Len())
	v.Tobytes(slice)
}

func (v *Value) Tobytes(buf *bytes.Buffer) {
	var klen keylen
	if v.Key.Value == nil {
		klen = 0
	} else {
		klen = keylen(len(v.Key.Value) + 1)
	}
	var vlen = len(v.Value)
	buf.WriteByte(byte(len(v.Children)))
	vlenb := valuelentobyte(valuelen(vlen))
	buf.Write(vlenb[:])
	klenb := keylentobyte(keylen(klen))
	buf.Write(klenb[:])

	buf.Write(v.Value)

	buf.WriteByte(v.Vtype)

	if klen != 0 {
		v.Key.tobytes(buf)
	}

	for i := range v.Children {
		b := ikeytobytes(v.Children[i])
		buf.Write(b[:])
	}
}

func ikeytobytes(key Ikeytype) (buf [ikeylen]byte) {
	binary.LittleEndian.PutUint32(buf[:], uint32(key))
	return
}

func keylentobyte(key keylen) (buf [keylenbytes]byte) {
	binary.LittleEndian.PutUint32(buf[:], uint32(key))
	return
}

func Getkeylen(buf []byte) keylen {
	return keylen(binary.LittleEndian.Uint32(buf[:]))
}

func valuelentobyte(key valuelen) (buf [valuelenbytes]byte) {
	binary.LittleEndian.PutUint32(buf[:], uint32(key))
	return
}

func Getvaluelen(buf []byte) valuelen {
	return valuelen(binary.LittleEndian.Uint32(buf[:]))
}

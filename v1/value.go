package v1

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

func (v *Value) FromBytes(data []byte) {
	dlen := len(data)
	if dlen <= 0 {
		panic(corruptinputdata)
	}
	clen := int(data[0])
	vlen := int(Getvaluelen(data[1 : 1+valuelenbytes]))
	klen := int(Getkeylen(data[1+valuelenbytes : 1+valuelenbytes+keylenbytes]))

	start := 1 + valuelenbytes + keylenbytes
	if dlen < int(klen+vlen+clen*ikeylen)+2+valuelenbytes+keylenbytes {
		panic(corruptinputdata)
	}

	v.Value = data[start : vlen+start]
	v.Vtype = data[vlen+start]

	if klen != 0 {
		v.Key.fromBytes(data[vlen+1+start : klen+vlen+1+start])
	}
	if clen != 0 {
		v.Children = make([]Ikeytype, clen)
		for i := 0; i < clen*ikeylen; i = i + ikeylen {
			v.Children[i/ikeylen] = ikeyfrombytes(data[klen+vlen+i+1+start : klen+vlen+i+1+start+ikeylen])
		}
	} else {
		v.Children = nil
	}
}

func ikeytobytes(key Ikeytype) (buf [ikeylen]byte) {
	binary.LittleEndian.PutUint32(buf[:], uint32(key))
	return
}

func ikeyfrombytes(buf []byte) Ikeytype {
	return Ikeytype(binary.LittleEndian.Uint32(buf[:]))
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

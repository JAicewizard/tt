package v3

import (
	"bytes"
	"encoding/binary"
)

type (
	//Value is the type that contains the value data
	Value struct {
		Key       Key    //the key of this object in its final form
		Value     []byte //the key to the data of this object in its final form
		Vtype     ttType //the type of the data of this object in its final form
		Childrenn uint64
	}
)

func AddValue(slice *bytes.Buffer, v *Value) {
	v.Tobytes(slice)
}

func (v *Value) Tobytes(buf *bytes.Buffer) {
	var klen = len(v.Key.Value)
	var vlen = len(v.Value)

	varintBuf := [binary.MaxVarintLen64]byte{}

	varintBytes := binary.PutUvarint(varintBuf[:], uint64(vlen))
	buf.Write(varintBuf[:varintBytes])

	varintBytes = binary.PutUvarint(varintBuf[:], uint64(klen))
	buf.Write(varintBuf[:varintBytes])

	buf.WriteByte(byte(v.Vtype))
	buf.Write(v.Value)

	buf.WriteByte(byte(v.Key.Vtype))
	buf.Write(v.Key.Value)

	varintBytes = binary.PutUvarint(varintBuf[:], v.Childrenn)
	buf.Write(varintBuf[:varintBytes])
}

func (v *Value) FromBytes(data *bytes.Buffer) {
	vlen, err := binary.ReadUvarint(data)
	if err != nil {
		panic(corruptinputdata)
	}
	klen, err := binary.ReadUvarint(data)
	if err != nil {
		panic(corruptinputdata)
	}
	if uint64(data.Len()) < 1+vlen+1+klen {
		panic(corruptinputdata)
	}
	typ, _ := data.ReadByte()
	v.Vtype = ttType(typ)
	v.Value = data.Next(int(vlen))
	typ, _ = data.ReadByte()
	v.Key.Vtype = ttType(typ)
	v.Key.Value = data.Next(int(klen))
	children, err := binary.ReadUvarint(data)
	if err != nil {
		panic(corruptinputdata)
	}
	v.Childrenn = children
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

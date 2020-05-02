package v3

import (
	"bytes"
	"encoding/binary"
	"io"
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

//Reader merges io.reader and io.ByteReader
type Reader interface {
	io.Reader
	io.ByteReader
}

type readFirstByte struct{
	b byte
	r io.ByteReader
}

func (r readFirstByte)ReadByte() (byte, error){
	if r.b == 0{
		return r.ReadByte()
	}else{
		b := r.b
		r.b = 0
		return b, nil
	}
}

func AddValue(out io.Writer, v *Value, buf *bytes.Buffer) {
	v.Tobytes(out, buf)
}

func (v *Value) Tobytes(out io.Writer, buf *bytes.Buffer) {
	var klen = len(v.Key.Value)
	var vlen = len(v.Value)
	buf.Reset()
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
	buf.WriteTo(out)
}

func (v *Value) FromBytes(in Reader) {
	vlen, err := binary.ReadUvarint(in)
	if err != nil {
		panic(corruptinputdata)
	}
	klen, err := binary.ReadUvarint(in)
	if err != nil {
		panic(corruptinputdata)
	}
	data := make([]byte, 1+vlen+1+klen+1)
	_, err = io.ReadFull(in, data)
	if err != nil {
		panic(corruptinputdata)
	}

	v.Vtype = ttType(data[0])
	v.Value = data[1 : 1+vlen]
	v.Key.Vtype = ttType(data[1+vlen])
	v.Key.Value = data[1+vlen+1 : 1+vlen+1+klen]

	if data[1+vlen+1+klen] < 0x80 {
		v.Childrenn = uint64(data[1+vlen+1+klen])
	}else{
		children, err := binary.ReadUvarint(in)
		if err != nil {
			panic(corruptinputdata)
		}
		v.Childrenn = children
	}
}
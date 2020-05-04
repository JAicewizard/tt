package v3

import (
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

type readFirstByte struct {
	b byte
	r io.ByteReader
}

func (r readFirstByte) ReadByte() (byte, error) {
	if r.b == 0 {
		return r.ReadByte()
	}
	b := r.b
	r.b = 0
	return b, nil
}

func AddValue(out io.Writer, v *Value, varintbuf *[binary.MaxVarintLen64 + 1]byte) {
	v.Tobytes(out, varintbuf)
}

func (v *Value) Tobytes(out io.Writer, varintbuf *[binary.MaxVarintLen64 + 1]byte) {
	var klen = len(v.Key.Value)
	var vlen = len(v.Value)

	varintBytes := binary.PutUvarint(varintbuf[:], uint64(vlen))
	out.Write(varintbuf[:varintBytes])

	varintBytes = binary.PutUvarint(varintbuf[:], uint64(klen))
	varintbuf[varintBytes] = byte(v.Vtype)
	out.Write(varintbuf[:varintBytes+1])

	out.Write(v.Value)
	varintbuf[0] = byte(v.Key.Vtype)
	out.Write(varintbuf[:1])
	out.Write(v.Key.Value)

	varintBytes = binary.PutUvarint(varintbuf[:], v.Childrenn)
	out.Write(varintbuf[:varintBytes])
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
	_, err = io.ReadAtLeast(in, data, int(1+vlen+1+klen+1))
	if err != nil {
		panic(corruptinputdata)
	}

	v.Vtype = ttType(data[0])
	v.Value = data[1 : 1+vlen]
	v.Key.Vtype = ttType(data[1+vlen])
	v.Key.Value = data[1+vlen+1 : 1+vlen+1+klen]

	if data[1+vlen+1+klen] < 0x80 {
		v.Childrenn = uint64(data[1+vlen+1+klen])
	} else {
		children, err := binary.ReadUvarint(readFirstByte{
			b: data[1+vlen+1+klen],
			r: in,
		})
		if err != nil {
			panic(corruptinputdata)
		}
		v.Childrenn = children
	}
}

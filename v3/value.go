package v3

import (
	"encoding/binary"
	"errors"
	"io"
)

type (
	//Value is the type that contains the value data
	Value struct {
		Key   KeyValue //the key of this object in its final form
		Value KeyValue //the value of this object in its final form
		//Value  []byte   //the key to the data of this object in its final form
		//Vtype     ttType   //the type of the data of this object in its final form
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

//AddValue is just a wrapper around v.Tobytes
func AddValue(out io.Writer, v *Value, varintbuf *[2*binary.MaxVarintLen64 + 1]byte) {
	v.Tobytes(out, varintbuf)
}

//Tobytes writes te Value into a io.Writer
func (v *Value) Tobytes(out io.Writer, varintbuf *[2*binary.MaxVarintLen64 + 1]byte) {
	var klen = len(v.Key.Value)
	var vlen = len(v.Value.Value)

	varintBytes1 := PutUvarint(varintbuf[:], uint64(vlen))

	varintBytes := PutUvarint(varintbuf[varintBytes1:], uint64(klen))
	varintbuf[varintBytes1+varintBytes] = byte(v.Value.Vtype)
	out.Write(varintbuf[:varintBytes1+varintBytes+1])

	out.Write(v.Value.Value)

	varintbuf[0] = byte(v.Key.Vtype)
	out.Write(varintbuf[:1])
	out.Write(v.Key.Value)

	varintBytes = PutUvarint(varintbuf[:], v.Childrenn)
	out.Write(varintbuf[:varintBytes])
}

//FromBytes reads bytes from a v3.Reader into Value
func (v *Value) FromBytes(in Reader, limit uint64) error {
	vlen, err := readerReadUvarint(in)
	if err != nil {
		return errors.New(corruptinputdata)
	}
	klen, err := readerReadUvarint(in)
	if err != nil {
		return errors.New(corruptinputdata)
	}
	if 1+vlen+1+klen+1 > limit {
		return errors.New(oversizedInputData)
	}
	data := make([]byte, 1+vlen+1+klen+1)
	_, err = io.ReadAtLeast(in, data, int(1+vlen+1+klen+1))
	if err != nil {
		return errors.New(corruptinputdata)
	}

	v.Value.Vtype = ttType(data[0])
	v.Value.Value = data[1 : 1+vlen]
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
			return errors.New(corruptinputdata)
		}
		v.Childrenn = children
	}
	return nil
}

// copy of binary.ReadUvarint with diferent interface input

var errEverflow = errors.New("binary: varint overflows a 64-bit integer")

func readerReadUvarint(r Reader) (uint64, error) {
	var x uint64
	var s uint
	for i := 0; i < binary.MaxVarintLen64; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return x, err
		}
		if b < 0x80 {
			if i == 9 && b > 1 {
				return x, errEverflow
			}
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return x, errEverflow
}

// copy of binary.PutUvarint with diferent return

// PutUvarint encodes a uint64 into buf and returns the number of bytes written.
// If the buffer is too small, PutUvarint will panic.
func PutUvarint(buf []byte, x uint64) uint8 {
	i := uint8(0)
	for x >= 0x80 {
		buf[i] = byte(x) | 0x80
		x >>= 7
		i++
	}
	buf[i] = byte(x)
	return i + 1
}

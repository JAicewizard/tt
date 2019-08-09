package v1

import (
	"bytes"
)

func (k *Key) tobytes(buf *bytes.Buffer) {
	buf.Write(k.Value)
	buf.WriteByte(k.Vtype)
}

func (k *Key) Export() interface{} {
	switch k.Vtype {
	case StringT:
		return StringFromBytes(k.Value)
	case BytesT:
		return k.Value
	case Float64T:
		return Float64FromBytes(k.Value)
	case Float32T:
		return Float32FromBytes(k.Value)
	case Int64T:
		return Int64FromBytes(k.Value)
	case Int32T:
		return Int32FromBytes(k.Value)
	case Int16T:
		return Int16FromBytes(k.Value)
	case Int8T:
		return Int8FromBytes(k.Value)
	case Uint64T:
		return Uint64FromBytes(k.Value)
	case Uint32T:
		return Uint32FromBytes(k.Value)
	case Uint16T:
		return Uint16FromBytes(k.Value)
	case Uint8T:
		return Uint8FromBytes(k.Value[0])
	case BoolT:
		return BoolFromBytes(k.Value)
	}

	return nil
}

func (k *Key) fromBytes(data []byte) {
	dlen := len(data)
	if dlen <= 1 { //the key has to be at least of length one
		panic(corruptinputdata)
	}
	typeloc := dlen - 1
	k.Value = data[0:typeloc]
	k.Vtype = data[typeloc]
}

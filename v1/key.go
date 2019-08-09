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
	case stringT:
		return StringFromBytes(k.Value)
	case bytesT:
		return k.Value
	case float64T:
		return Float64FromBytes(k.Value)
	case float32T:
		return Float32FromBytes(k.Value)
	case int64T:
		return Int64FromBytes(k.Value)
	case int32T:
		return Int32FromBytes(k.Value)
	case int16T:
		return Int16FromBytes(k.Value)
	case int8T:
		return Int8FromBytes(k.Value)
	case uint64T:
		return Uint64FromBytes(k.Value)
	case uint32T:
		return Uint32FromBytes(k.Value)
	case uint16T:
		return Uint16FromBytes(k.Value)
	case uint8T:
		return Uint8FromBytes(k.Value[0])
	case boolT:
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

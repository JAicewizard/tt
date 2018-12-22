package tt

import (
	"bytes"
)

func (k *Key) tobytes(buf *bytes.Buffer) {
	buf.Write(k.Value)
	buf.WriteByte(k.Vtype)
}

func (k *Key) export() interface{} {
	switch k.Vtype {
	case v1String:
		return stringFromBytes(k.Value)
	case v1Bytes:
		return k.Value
	case v1Float64:
		return float64FromBytes(k.Value)
	case v1Float32:
		return float32FromBytes(k.Value)
	case v1Int64:
		return int64FromBytes(k.Value)
	case v1Int32:
		return int32FromBytes(k.Value)
	case v1Int16:
		return int16FromBytes(k.Value)
	case v1Int8:
		return int8FromBytes(k.Value)
	case v1Uint64:
		return uint64FromBytes(k.Value)
	case v1Uint32:
		return uint32FromBytes(k.Value)
	case v1Uint16:
		return uint16FromBytes(k.Value)
	case v1Uint8:
		return uint8FromBytes(k.Value[0])
	case v1Bool:
		return boolFromBytes(k.Value)

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

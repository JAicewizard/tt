package v3

import (
	"strconv"
)

type (
	//Key is the key used for storing the key of a Value
	Key struct {
		Value []byte //the key to the data of this object in its final form
		Vtype ttType //the type of the data of this object in its final form
	}
)

func (k *Key) ExportStructID() string {
	switch k.Vtype {
	case StringT:
		return StringFromBytes(k.Value)
	case BytesT:
		return string(k.Value)
	case Float32T:
		return strconv.FormatFloat(float64(Float32FromBytes(k.Value)), 'X', -1, 32)
	case Float64T:
		return strconv.FormatFloat(Float64FromBytes(k.Value), 'X', -1, 32)
	case Int8T:
		return strconv.FormatInt(int64(Int8FromBytes(k.Value[0])), 10)
	case Int16T:
		return strconv.FormatInt(int64(Int16FromBytes(k.Value)), 10)
	case Int32T:
		return strconv.FormatInt(int64(Int32FromBytes(k.Value)), 10)
	case Int64T:
		return strconv.FormatInt(int64(Int64FromBytes(k.Value)), 10)
	case Uint8T:
		return strconv.FormatUint(uint64(Uint8FromBytes(k.Value[0])), 10)
	case Uint16T:
		return strconv.FormatUint(uint64(Uint16FromBytes(k.Value)), 10)
	case Uint32T:
		return strconv.FormatUint(uint64(Uint32FromBytes(k.Value)), 10)
	case Uint64T:
		return strconv.FormatUint(uint64(Uint64FromBytes(k.Value)), 10)
	case BoolT:
		return strconv.FormatBool(BoolFromBytes(k.Value))
	}
	return ""
}

package v3

import (
	"reflect"
	"strconv"
)

type (
	//Key is the key used for storing the key of a Value
	Key struct {
		Value []byte //the key to the data of this object in its final form
		Vtype ttType //the type of the data of this object in its final form
	}
)

func (k *Key) Export() interface{} {
	switch k.Vtype {
	case StringT:
		return StringFromBytes(k.Value)
	case BytesT:
		return k.Value
	case Float32T:
		return Float32FromBytes(k.Value)
	case Float64T:
		return Float64FromBytes(k.Value)
	case Int8T:
		return Int8FromBytes(k.Value)
	case Int16T:
		return Int16FromBytes(k.Value)
	case Int32T:
		return Int32FromBytes(k.Value)
	case Int64T:
		return Int64FromBytes(k.Value)
	case Uint8T:
		return Uint8FromBytes(k.Value[0])
	case Uint16T:
		return Uint16FromBytes(k.Value)
	case Uint32T:
		return Uint32FromBytes(k.Value)
	case Uint64T:
		return Uint64FromBytes(k.Value)
	case BoolT:
		return BoolFromBytes(k.Value)
	}

	return nil
}
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
		return strconv.FormatInt(int64(Int8FromBytes(k.Value)), 10)
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

func (k *Key) Export_2() reflect.Value {
	if k.Vtype == StringT {
		return reflect.ValueOf(StringFromBytes(k.Value))
	} else if k.Vtype == BytesT {
		return reflect.ValueOf(k.Value)
	} else if k.Vtype == Float32T {
		return reflect.ValueOf(Float32FromBytes(k.Value))
	} else if k.Vtype == Float64T {
		return reflect.ValueOf(Float64FromBytes(k.Value))
	} else if k.Vtype == Int8T {
		return reflect.ValueOf(Int8FromBytes(k.Value))
	} else if k.Vtype == Int16T {
		return reflect.ValueOf(Int16FromBytes(k.Value))
	} else if k.Vtype == Int32T {
		return reflect.ValueOf(Int32FromBytes(k.Value))
	} else if k.Vtype == Int64T {
		return reflect.ValueOf(Int64FromBytes(k.Value))
	} else if k.Vtype == Uint8T {
		return reflect.ValueOf(Uint8FromBytes(k.Value[0]))
	} else if k.Vtype == Uint16T {
		return reflect.ValueOf(Uint16FromBytes(k.Value))
	} else if k.Vtype == Uint32T {
		return reflect.ValueOf(Uint32FromBytes(k.Value))
	} else if k.Vtype == Uint64T {
		return reflect.ValueOf(Uint64FromBytes(k.Value))
	} else if k.Vtype == BoolT {
		return reflect.ValueOf(BoolFromBytes(k.Value))
	} else {
		return reflect.Value{}
	}
}

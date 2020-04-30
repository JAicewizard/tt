package v3

import (
	"errors"
	"reflect"
)

type (
	Ikeytype uint32
	valuelen uint32
	keylen   uint32

	ttType byte
)

var (
	//ErrInvalidInput is used for when the input it invalid
	ErrInvalidInput = errors.New("invalid input")
)

const (
	corruptinputdata = "byte length not long enough, contact the authors for a solution"

	ikeylen       = 4
	valuelenbytes = 4
	keylenbytes   = 4
)

const (
	StringT = ttType(iota + 1)
	BytesT
	Int8T
	Int16T
	Int32T
	Int64T
	Uint8T
	Uint16T
	Uint32T
	Uint64T
	BoolT
	Float32T
	Float64T
	_
	_
	_
	_
	MapT
	ArrT
)

var (
//strType =
)

func (t ttType) GetReflectType() reflect.Type {
	switch t {
	case StringT:
		return reflect.TypeOf(string(""))
	case BytesT:
		return reflect.TypeOf([]byte{})
	case Int8T:
		return reflect.TypeOf(int8(0))
	case Int16T:
		return reflect.TypeOf(int16(0))
	case Int32T:
		return reflect.TypeOf(int32(0))
	case Int64T:
		return reflect.TypeOf(int64(0))
	case Uint8T:
		return reflect.TypeOf(uint8(0))
	case Uint16T:
		return reflect.TypeOf(uint16(0))
	case Uint32T:
		return reflect.TypeOf(uint32(0))
	case Uint64T:
		return reflect.TypeOf(uint64(0))
	case BoolT:
		return reflect.TypeOf(true)
	case Float32T:
		return reflect.TypeOf(float32(0.0))
	case Float64T:
		return reflect.TypeOf(float64(0.0))
	case MapT:
		return reflect.TypeOf(map[interface{}]interface{}{})
	case ArrT:
		return reflect.TypeOf([]interface{}{})
	}
	return nil
}

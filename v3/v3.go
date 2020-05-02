package v3

import (
	"errors"
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


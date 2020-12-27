package v3

import (
	"errors"
)

type ttType byte

var (
	//ErrInvalidInput is used for when the input it invalid
	ErrInvalidInput = errors.New("invalid input")
)

const (
	corruptinputdata   = "Not enough data in the datastream, imput might be corrupt."
	oversizedInputData = "imput too big."
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

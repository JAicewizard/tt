package v1

import (
	"errors"
)

type (
	//Key is the key used for storing the key of a Value
	Key struct {
		Value []byte //the key to the data of this object in its final form
		Vtype byte   //the type of the data of this object in its final form
	}

	Ikeytype uint32
	valuelen uint32
	keylen   uint32
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
	FinalValueT = iota + 1
	StringT
	BytesT
	Float64T
	Float32T
	MapT
	ArrT
	Int64T
	Int32T
	Int16T
	Int8T
	Uint64T
	Uint32T
	Uint16T
	Uint8T
	BoolT
)

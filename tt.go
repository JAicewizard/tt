package tt

import (
	"bytes"
	"errors"

	v2 "github.com/jaicewizard/tt/v2"
)

/*
	TODO: supported types: string, int, []byte, floats bool
	TODO: pre-allocate all data when encoding.
*/
type (
	//Data is the type
	Data map[interface{}]interface{}
	//Key is the key used for storing the key of a Value
	Key struct {
		Value []byte //the key to the data of this object in its final form
		Vtype byte   //the type of the data of this object in its final form
	}

	ikeytype uint32
	valuelen uint32
	keylen   uint32

	//Transmitter transmits the data
	Transmitter interface {
		Encode() ([]byte, error)
		Decode([]byte) (interface{}, error)
		GetCode() byte
	}
)

var (
	//ErrCodeUsed is for when the code for the transmitter is already used
	ErrCodeUsed = errors.New("code already used")

	//ErrInvalidInput is used for when the input it invalid
	ErrInvalidInput = errors.New("invalid input")

	transmitters = make(map[byte]Transmitter)
)

//RegisterTransmitter registers a new transmitter
func RegisterTransmitter(tr Transmitter) error {
	code := tr.GetCode()
	if code == v2.StringT || code == v2.FinalValueT || code == v2.MapT || code == v2.ArrT {
		return ErrCodeUsed
	}

	transmitters[code] = tr
	return nil
}

//GobEncode encodes the data of a map[interface{}]interface{}
func (d Data) GobEncode() ([]byte, error) {
	var byt []byte
	buf := bytes.NewBuffer(byt)
	Encodev3(d.(map[interface{}]interface{}), buf)
	return buf.Bytes(), nil
}

//GobDecode decodes the data back into a map[interface{}]interface{}
func (d *Data) GobDecode(data []byte) error {
	if len(data) == 0 {
		return ErrInvalidInput
	}
	switch data[0] {
	case version1:
		return Decodev1(data[1:], d)
	case version2:
		return Decodev2(data[1:], d)
	case version3:
		return Decodev3(bytes.NewBuffer(data), d)
	}
	return nil
}

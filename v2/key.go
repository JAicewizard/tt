package v2

import (
	"bytes"
)

type(
	//Key is the key used for storing the key of a Value
	Key struct {
	Value []byte //the key to the data of this object in its final form
		Vtype byte   //the type of the data of this object in its final form
	}	
)

func (k *Key) tobytes(buf *bytes.Buffer) {
	buf.Write(k.Value)
	buf.WriteByte(k.Vtype)
}
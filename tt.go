package tt

import (
	"bytes"
	"encoding/binary"

	"errors"
)

/*
	TODO:
	supported types: string, int
*/
type (
	Data  map[interface{}]interface{}
	Value struct {
		Key      Key        //the key of this object in its final form
		Children []ikeytype //the key to the values in the data send across
		Value    []byte     //the key to the data of this object in its final form
		Vtype    byte       //the type of the data of this object in its final form
		used     bool
	}
	Key struct {
		Value []byte //the key to the data of this object in its final form
		Vtype byte   //the type of the data of this object in its final form
	}

	ikeytype uint16

	Transmitter interface {
		Encode() ([]byte, error)
		Decode([]byte) (interface{}, error)
		GetCode() byte
	}
)

const ikeylen = 2

const version1 = 1

const (
	v1String     = 's'
	v1FinalValue = 'v'
	v1Map        = 'm'
)

const corruptinputdata = "byte length not long enough, contact the authors for a solution"

var transmitters = make([]Transmitter, 255)

var (
	//ErrCodeUsed is for when the code for the transmitter is already used
	ErrCodeUsed = errors.New("code already used")

	//ErrInvalidInput is used for when the input it invalid
	ErrInvalidInput = errors.New("invalid input")
)

//RegisterTransmitter registers a new transmitter
func RegisterTransmitter(tr Transmitter) error {
	code := tr.GetCode()
	if code == v1String || code == v1FinalValue || code == v1Map {
		return ErrCodeUsed
	}

	transmitters[int(code)] = tr
	return nil
}

//GobEncode encodes the data of a map[interface{}]interface{}
func (d Data) GobEncode() ([]byte, error) {
	var byt []byte
	buf := bytes.NewBuffer(byt)
	buf.WriteByte(byte(version1))
	encodev1(d, buf)
	return buf.Bytes(), nil
}

//GobDecode decodes the data back into a map[interface{}]interface{}
func (d *Data) GobDecode(data []byte) error {
	if len(data) == 0 {
		return ErrInvalidInput
	}
	switch data[0] {
	case version1:
		return decodev1(data[1:], d)
	}
	return nil
}

func encodev1(d Data, values *bytes.Buffer) {
	tv := ikeytype(0)
	firstChilds, err := encodemapv1(values, d, &tv)
	if err != nil {
		panic(err)
	}

	addValue(values, &Value{
		Children: firstChilds,
		Vtype:    v1FinalValue,
	})
	values.WriteByte(byte(tv + 1))
}

func decodev1(b []byte, d *Data) (err error) {
	vlen := int(b[len(b)-1])

	locs := make([]uint16, vlen)
	locs[0] = 0

	for i := 1; i < vlen; i++ {
		locs[i] = locs[i-1] + uint16(b[locs[i-1]]) + 1
	}

	//decoding the actual values
	var v Value

	v.fromBytes(b[locs[vlen-1]+1:])

	if *d == nil {
		*d = make(Data, len(v.Children)*3)
	}

	data := d
	childs := v.Children
	for ck := range childs {
		var err error
		v.fromBytes(b[locs[childs[ck]]+1:])

		err = valueToMapv1(&v, *data, locs, b)

		if err != nil {
			return err
		}
	}

	return nil

}

func encodemapv1(values *bytes.Buffer, d Data, nextValue *ikeytype) ([]ikeytype, error) {
	createdObjects := make([]ikeytype, len(d))
	i := 0

	for k := range d {
		var value Value

		switch v := k.(type) {
		case string:
			value.Key.Value = stringToBytes(v)
			value.Key.Vtype = v1String
		}

		switch v := d[k].(type) {
		case string:
			value.Value = stringToBytes(v)
			value.Vtype = v1String

		default:
			if v, ok := d[k].(Data); ok {
				childs, err := encodemapv1(values, v, nextValue)
				if err != nil {
					return nil, err
				}

				value.Children = childs
				value.Vtype = v1Map

			} else if v, ok := d[k].(Transmitter); ok {
				var err error
				value.Value, err = v.Encode()
				if err != nil {
					return nil, err
				}
				value.Vtype = v.GetCode()
			} else {
				return nil, ErrInvalidInput
			}
		}

		createdObjects[i] = *nextValue
		i++

		addValue(values, &value)
		*nextValue++
	}
	return createdObjects, nil
}

func stringToBytes(s string) []byte {
	return []byte(s)
}

func stringFromBytes(b []byte) interface{} {
	return interface{}(string(b))
}

func addValue(slice *bytes.Buffer, v *Value) {
	var ln int
	if v.Key.Value == nil {
		ln = len(v.Value) + len(v.Children)*2 + 4
	} else {
		ln = len(v.Value) + len(v.Children)*2 + len(v.Key.Value) + 6
	}
	slice.Grow(ln)
	slice.WriteByte(byte(ln))
	v.tobytes(slice)
}

func valueToMapv1(v *Value, data Data, locs []uint16, buf []byte) error {
	key := v.Key.export()
	if key == nil {
		return ErrInvalidInput
	}

	switch v.Vtype {
	case v1String:
		data[key] = stringFromBytes(v.Value)
	case v1Map:
		data[key] = make(Data, len(v.Children))
		childs := v.Children
		for ck := range childs {
			var err error
			v.fromBytes(buf[locs[childs[ck]]+1:])
			err = valueToMapv1(v, data[key].(Data), locs, buf)
			if err != nil {
				return err
			}
		}
	default:
		t := transmitters[int(v.Vtype)]
		if t == nil {
			return errors.New("no transmitter for type:" + string(v.Vtype))
		}
		var err error
		data[key], err = t.Decode(v.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Key) tobytes(buf *bytes.Buffer) {
	buf.WriteByte(byte(len(k.Value)))
	buf.Write(k.Value)
	buf.WriteByte(k.Vtype)
}

func (k *Key) export() interface{} {
	switch k.Vtype {
	case v1String:
		return stringFromBytes(k.Value)
	}

	return nil
}

func (k *Key) fromBytes(data []byte) {
	dlen := len(data)
	if dlen <= 0 {
		panic(corruptinputdata)
	}
	typeloc := data[0] + 1
	if dlen < int(typeloc+1) {
		panic(corruptinputdata)
	}
	k.Value = data[1:typeloc]
	k.Vtype = data[typeloc]
}

func (v *Value) tobytes(buf *bytes.Buffer) {

	buf.WriteByte(byte(len(v.Value)))
	buf.Write(v.Value)

	buf.WriteByte(v.Vtype)

	if v.Key.Value == nil {
		buf.WriteByte(byte(0))
	} else {
		klen := len(v.Key.Value)

		buf.WriteByte(byte(klen + 2))
		buf.Grow(klen + 2)

		v.Key.tobytes(buf)
	}

	buf.WriteByte(byte(len(v.Children)))

	for i := range v.Children {
		b := ikeytobytes(v.Children[i])
		buf.Write(b[:])
	}
}

func (v *Value) fromBytes(data []byte) {
	dlen := len(data)
	if dlen <= 0 {
		panic(corruptinputdata)
	}
	vlen := data[0]
	if dlen < int(vlen+1) {
		panic(corruptinputdata)
	}

	klen := data[vlen+2]
	if dlen < int(klen+vlen+2) {
		panic(corruptinputdata)
	}

	clen := int(data[klen+vlen+3])

	v.Value = data[1 : vlen+1]
	v.Vtype = data[vlen+1]

	if klen != 0 {
		v.Key.fromBytes(data[vlen+3 : klen+vlen+3])
	}
	if clen != 0 {
		v.Children = make([]ikeytype, clen)
		for i := 0; i < clen*2; i = i + 2 {
			v.Children[i/2] = ikeyfrombytes(data[int(klen+vlen+4)+i : int(klen+vlen+4)+i+2])
		}
	}
}

func ikeytobytes(key ikeytype) (buf [ikeylen]byte) {
	binary.LittleEndian.PutUint16(buf[:], uint16(key))
	return
}

func ikeyfrombytes(buf []byte) ikeytype {
	return ikeytype(binary.LittleEndian.Uint16(buf[:]))
}

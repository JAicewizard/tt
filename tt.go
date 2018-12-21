package tt

import (
	"bytes"
	"reflect"

	"errors"
)

/*
	TODO:
	supported types: string, int, []byte, floats
	TODO: pre-allocate all data when encoding.
*/
type (
	//Data is the type
	Data map[interface{}]interface{}
	//Value is the type that contains the value data
	Value struct {
		Key      Key        //the key of this object in its final form
		Children []ikeytype //the key to the values in the data send across
		Value    []byte     //the key to the data of this object in its final form
		Vtype    byte       //the type of the data of this object in its final form
	}
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

const ikeylen = 4

const valuelenbytes = 4
const keylenbytes = 4

const version1 = 1

const (
	v1String = iota
	v1Bytes
	v1Float64
	v1Float32
	v1Map
	v1Arr
	v1Int64
	v1Int32
	v1Int16
	v1Int8
	v1Uint64
	v1Uint32
	v1Uint16
	v1Uint8
	v1FinalValue
)

const corruptinputdata = "byte length not long enough, contact the authors for a solution"

var transmitters = make(map[byte]Transmitter)

var (
	//ErrCodeUsed is for when the code for the transmitter is already used
	ErrCodeUsed = errors.New("code already used")

	//ErrInvalidInput is used for when the input it invalid
	ErrInvalidInput = errors.New("invalid input")
)

//RegisterTransmitter registers a new transmitter
func RegisterTransmitter(tr Transmitter) error {
	code := tr.GetCode()
	if code == v1String || code == v1FinalValue || code == v1Map || code == v1Arr {
		return ErrCodeUsed
	}

	transmitters[code] = tr
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

func encodemapv1(values *bytes.Buffer, d Data, nextValue *ikeytype) ([]ikeytype, error) {
	createdObjects := make([]ikeytype, len(d))
	i := 0

	for k := range d {
		err := encodevaluev1(values, d[k], k, nextValue)
		if err != nil {
			return nil, err
		}
		createdObjects[i] = *nextValue
		i++
		*nextValue++
	}

	return createdObjects, nil
}

func encodevaluev1(values *bytes.Buffer, d interface{}, k interface{}, nextValue *ikeytype) error {
	var value Value

	if k != nil {
		switch v := k.(type) {
		case string:
			value.Key.Value = stringToBytes(v)
			value.Key.Vtype = v1String
		case []byte:
			value.Key.Value = v
			value.Key.Vtype = v1Bytes
		case float64:
			value.Key.Value = float64ToBytes(v)
			value.Key.Vtype = v1Float64
		case float32:
			value.Key.Value = float32ToBytes(v)
			value.Key.Vtype = v1Float32
		case int64:
			value.Key.Value = int64ToBytes(v)
			value.Key.Vtype = v1Int64
		case int32:
			value.Key.Value = int32ToBytes(v)
			value.Key.Vtype = v1Int32
		case int16:
			value.Key.Value = int16ToBytes(v)
			value.Key.Vtype = v1Int16
		case int8:
			value.Key.Value = int8ToBytes(v)
			value.Key.Vtype = v1Int8
		case uint64:
			value.Key.Value = uint64ToBytes(v)
			value.Key.Vtype = v1Uint64
		case uint32:
			value.Key.Value = uint32ToBytes(v)
			value.Key.Vtype = v1Uint32
		case uint16:
			value.Key.Value = uint16ToBytes(v)
			value.Key.Vtype = v1Uint16
		case uint8:
			value.Key.Value = uint8ToBytes(v)
			value.Key.Vtype = v1Uint8
		case int:
			value.Key.Value = int64ToBytes(int64(v))
			value.Key.Vtype = v1Int64
		}
	}

	switch v := d.(type) {
	case string:
		value.Value = stringToBytes(v)
		value.Vtype = v1String
	case []byte:
		value.Value = v
		value.Vtype = v1Bytes
	case float64:
		value.Value = float64ToBytes(v)
		value.Vtype = v1Float64
	case float32:
		value.Value = float32ToBytes(v)
		value.Vtype = v1Float32
	case int64:
		value.Value = int64ToBytes(v)
		value.Vtype = v1Int64
	case int32:
		value.Value = int32ToBytes(v)
		value.Vtype = v1Int32
	case int16:
		value.Value = int16ToBytes(v)
		value.Vtype = v1Int16
	case int8:
		value.Value = int8ToBytes(v)
		value.Vtype = v1Int8
	case uint64:
		value.Value = uint64ToBytes(v)
		value.Vtype = v1Uint64
	case uint32:
		value.Value = uint32ToBytes(v)
		value.Vtype = v1Uint32
	case uint16:
		value.Value = uint16ToBytes(v)
		value.Vtype = v1Uint16
	case uint8:
		value.Value = uint8ToBytes(v)
		value.Vtype = v1Uint8
	case int:
		value.Value = int64ToBytes(int64(v))
		value.Vtype = v1Int64
	case uint:
		value.Value = uint64ToBytes(uint64(v))
		value.Vtype = v1Uint64

	case []interface{}:
		value.Children = make([]ikeytype, len(v))

		for i := 0; i < len(v); i++ {
			err := encodevaluev1(values, v[i], nil, nextValue)
			if err != nil {
				return err
			}
			value.Children[i] = *nextValue
			*nextValue++
		}

		value.Vtype = v1Arr
	case map[string]interface{}:
		value.Children = make([]ikeytype, len(v))
		i := 0
		for k := range v {

			encodevaluev1(values, v[k], k, nextValue)

			value.Children[i] = *nextValue
			i++
			*nextValue++
		}
		value.Vtype = v1Map
	case Data:
		childs, err := encodemapv1(values, v, nextValue)
		if err != nil {
			return err
		}
		value.Children = childs
		value.Vtype = v1Map
	default:
		if val := reflect.ValueOf(d); val.Kind() == reflect.Array {
			value.Children = make([]ikeytype, val.Len())

			for i := 0; i < val.Len(); i++ {
				e := val.Index(i).Interface()
				err := encodevaluev1(values, e, nil, nextValue)
				if err != nil {
					return err
				}
				value.Children[i] = *nextValue
				*nextValue++
			}

			value.Vtype = v1Arr
		} else if v, ok := d.(Transmitter); ok {
			var err error
			value.Value, err = v.Encode()
			if err != nil {
				return err
			}
			value.Vtype = v.GetCode()
		} else {
			return ErrInvalidInput
		}
	}
	addValue(values, &value)
	return nil
}

func addValue(slice *bytes.Buffer, v *Value) {
	var ln int
	if v.Key.Value == nil {
		ln = len(v.Value) + len(v.Children)*ikeylen + 2 + valuelenbytes + keylenbytes
	} else {
		ln = len(v.Value) + len(v.Children)*ikeylen + len(v.Key.Value) + 3 + valuelenbytes + keylenbytes
	}
	slice.Grow(ln)
	v.tobytes(slice)
}

func decodev1(b []byte, d *Data) (err error) {
	vlen := int(b[len(b)-1])

	locs := make([]uint64, vlen)
	locs[0] = 0

	for i := 1; i < vlen; i++ {
		locs[i] = locs[i-1] + uint64(
			uint32(b[locs[i-1]]*ikeylen)+ //this is the length for the children
				uint32(getvaluelen(b[locs[i-1]+1:locs[i-1]+1+valuelenbytes]))+ //this is the value length
				uint32(getkeylen(b[locs[i-1]+1+valuelenbytes:locs[i-1]+1+valuelenbytes+keylenbytes]))+ //this is the key length
				2+valuelenbytes+keylenbytes) //add 4 so that we cound the length values as wel, +1 is for going to the next value
	}

	//decoding the actual values
	var v Value
	v.fromBytes(b[locs[vlen-1]:])

	if *d == nil {
		*d = make(Data, len(v.Children)*1)
	}

	data := d
	childs := v.Children

	for ck := range childs {
		var err error
		v.fromBytes(b[locs[childs[ck]]:])

		err = valueToMapv1(&v, *data, locs, b)

		if err != nil {
			return err
		}
	}
	return nil
}

func valueToMapv1(v *Value, data Data, locs []uint64, buf []byte) error {
	key := v.Key.export()
	if key == nil {
		return ErrInvalidInput
	}

	switch v.Vtype { //!make sure to update types in valueToArray as well!
	case v1String:
		data[key] = stringFromBytes(v.Value)
	case v1Bytes:
		data[key] = v.Value
	case v1Float64:
		data[key] = float64FromBytes(v.Value)
	case v1Float32:
		data[key] = float32FromBytes(v.Value)
	case v1Int64:
		data[key] = int64FromBytes(v.Value)
	case v1Int32:
		data[key] = int32FromBytes(v.Value)
	case v1Int16:
		data[key] = int16FromBytes(v.Value)
	case v1Int8:
		data[key] = int8FromBytes(v.Value[0])
	case v1Uint64:
		data[key] = uint64FromBytes(v.Value)
	case v1Uint32:
		data[key] = uint32FromBytes(v.Value)
	case v1Uint16:
		data[key] = uint16FromBytes(v.Value)
	case v1Uint8:
		data[key] = uint8FromBytes(v.Value[0])
	case v1Map:
		data[key] = make(Data, len(v.Children))
		childs := v.Children
		for ck := range childs {
			var err error
			v.fromBytes(buf[locs[childs[ck]]:])
			err = valueToMapv1(v, data[key].(Data), locs, buf)
			if err != nil {
				return err
			}
		}
	case v1Arr:
		val := make([]interface{}, len(v.Children))
		childs := v.Children
		for i := range childs {
			var err error
			v.fromBytes(buf[locs[childs[i]]:])
			err = valueToArrayv1(v, val, i, locs, buf)
			if err != nil {
				return err
			}
		}
		data[key] = val

	default:
		t := transmitters[v.Vtype]
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

func valueToArrayv1(v *Value, arr []interface{}, i int, locs []uint64, buf []byte) error {
	switch v.Vtype {
	case v1String:
		arr[i] = stringFromBytes(v.Value)
	case v1Bytes:
		arr[i] = v.Value
	case v1Float64:
		arr[i] = float64FromBytes(v.Value)
	case v1Float32:
		arr[i] = float32FromBytes(v.Value)
	case v1Int64:
		arr[i] = int64FromBytes(v.Value)
	case v1Int32:
		arr[i] = int32FromBytes(v.Value)
	case v1Int16:
		arr[i] = int16FromBytes(v.Value)
	case v1Int8:
		arr[i] = int8FromBytes(v.Value[0])
	case v1Uint64:
		arr[i] = uint64FromBytes(v.Value)
	case v1Uint32:
		arr[i] = uint32FromBytes(v.Value)
	case v1Uint16:
		arr[i] = uint16FromBytes(v.Value)
	case v1Uint8:
		arr[i] = uint8FromBytes(v.Value[0])
	case v1Map:
		arr[i] = make(Data, len(v.Children))
		childs := v.Children
		for ck := range childs {
			var err error
			v.fromBytes(buf[locs[childs[ck]]:])
			err = valueToMapv1(v, arr[i].(Data), locs, buf)
			if err != nil {
				return err
			}
		}
	case v1Arr:
		val := make([]interface{}, len(v.Children))
		childs := v.Children
		for i := range childs {
			var err error
			v.fromBytes(buf[locs[childs[i]]:])
			err = valueToArrayv1(v, interface{}((val)).([]interface{}), i, locs, buf)
			if err != nil {
				return err
			}
		}
		arr[i] = val
	default:
		t := transmitters[v.Vtype]
		if t == nil {
			return errors.New("no transmitter for type:" + string(v.Vtype))
		}
		var err error
		arr[i], err = t.Decode(v.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

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
		return int8FromBytes(k.Value[0])
	case v1Uint64:
		return uint64FromBytes(k.Value)
	case v1Uint32:
		return uint32FromBytes(k.Value)
	case v1Uint16:
		return uint16FromBytes(k.Value)
	case v1Uint8:
		return uint8FromBytes(k.Value[0])
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
func (v *Value) tobytes(buf *bytes.Buffer) {
	var klen keylen
	if v.Key.Value == nil {
		klen = 0
	} else {
		klen = keylen(len(v.Key.Value) + 1)
	}
	var vlen = len(v.Value)
	buf.WriteByte(byte(len(v.Children)))
	vlenb := valuelentobyte(valuelen(vlen))
	buf.Write(vlenb[:])
	klenb := keylentobyte(keylen(klen))
	buf.Write(klenb[:])

	buf.Write(v.Value)

	buf.WriteByte(v.Vtype)

	if klen != 0 {
		buf.Grow(int(klen + 1))
		v.Key.tobytes(buf)
	}

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
	clen := int(data[0])
	vlen := int(getvaluelen(data[1 : 1+valuelenbytes]))
	klen := int(getkeylen(data[1+valuelenbytes : 1+valuelenbytes+keylenbytes]))

	start := 1 + valuelenbytes + keylenbytes
	if dlen < int(klen+vlen+(clen-1)*ikeylen)+2+valuelenbytes+keylenbytes {
		panic(corruptinputdata)
	}

	v.Value = data[start : vlen+start]
	v.Vtype = data[vlen+start]

	if klen != 0 {
		v.Key.fromBytes(data[vlen+1+start : klen+vlen+1+start])
	}
	if clen != 0 {
		v.Children = make([]ikeytype, clen)
		for i := 0; i < clen*ikeylen; i = i + ikeylen {
			v.Children[i/ikeylen] = ikeyfrombytes(data[klen+vlen+i+1+start : klen+vlen+i+1+start+ikeylen])
		}
	}
}

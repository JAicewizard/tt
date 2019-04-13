package tt

import (
	"bytes"
	"reflect"
	"runtime"

	"errors"
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

const ikeylen = 4

const valuelenbytes = 4
const keylenbytes = 4

const version1 = 1

const (
	v1FinalValue = iota + 1
	v1String
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
	v1Bool
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

	values.WriteByte(byte(version1))

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
	var KeepAlive interface{}
	var KeepAlive1 interface{}
	if k != nil {
		switch v := k.(type) { //making this s seperate function will decrese performance, it won't be able to inline and make more allocations
		case string:
			value.Key.Value = stringToBytes(v)
			value.Key.Vtype = v1String
		case []byte:
			value.Key.Value = v
			value.Key.Vtype = v1Bytes
		case float64:
			float64ToBytes(&v, &value.Key.Value)
			value.Key.Vtype = v1Float64
			KeepAlive = &v
		case float32:
			float32ToBytes3(&v, &value.Key.Value)
			value.Key.Vtype = v1Float32
			KeepAlive = &v
		case int64:
			var buf [8]byte
			int64ToBytes(v, &buf)
			value.Key.Value = buf[:]
			value.Key.Vtype = v1Int64
		case int32:
			value.Key.Value = int32ToBytes(v)
			value.Key.Vtype = v1Int32
		case int16:
			value.Key.Value = int16ToBytes(v)
			value.Key.Vtype = v1Int16
		case int8:
			int8ToBytes(&v, &value.Key.Value)
			value.Key.Vtype = v1Int8
			KeepAlive = &v
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
			uint8ToBytes(&v, &value.Key.Value)
			value.Key.Vtype = v1Uint8
			KeepAlive = &v
		case int:
			var buf [8]byte
			int64ToBytes(int64(v), &buf)
			value.Key.Value = buf[:]
			value.Key.Vtype = v1Int64
		case uint:
			value.Key.Value = uint64ToBytes(uint64(v))
			value.Key.Vtype = v1Bool
			KeepAlive = &v
		case bool:
			boolToBytes(&v, &value.Key.Value)
			value.Key.Vtype = v1Bool
			KeepAlive = &v
		}

	}

	//this sets value.Value, it does al the basic types and some more
	switch v := d.(type) {
	case string:
		value.Value = stringToBytes(v)
		value.Vtype = v1String
	case []byte:
		value.Value = v
		value.Vtype = v1Bytes
	case float64:
		float64ToBytes(&v, &value.Value)
		value.Vtype = v1Float64
		KeepAlive1 = &v
	case float32:
		float32ToBytes3(&v, &value.Value)
		value.Vtype = v1Float32
		KeepAlive1 = &v
	case int64:
		var buf [8]byte
		int64ToBytes(v, &buf)
		value.Value = buf[:]
		value.Vtype = v1Int64
	case int32:
		value.Value = int32ToBytes(v)
		value.Vtype = v1Int32
	case int16:
		value.Value = int16ToBytes(v)
		value.Vtype = v1Int16
	case int8:
		int8ToBytes(&v, &value.Value)
		value.Vtype = v1Int8
		KeepAlive1 = &v
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
		uint8ToBytes(&v, &value.Value)
		value.Vtype = v1Uint8
		KeepAlive1 = &v
	case int:
		var buf [8]byte
		int64ToBytes(int64(v), &buf)
		value.Value = buf[:]
		value.Vtype = v1Int64
	case uint:
		value.Value = uint64ToBytes(uint64(v))
		value.Vtype = v1Bool
		KeepAlive1 = &v
	case bool:
		boolToBytes(&v, &value.Value)
		value.Vtype = v1Bool
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
	case map[interface{}]interface{}:
		childs, err := encodemapv1(values, Data(v), nextValue)
		if err != nil {
			return err
		}
		value.Children = childs
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
	runtime.KeepAlive(KeepAlive)
	runtime.KeepAlive(KeepAlive1)
	return nil
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

func valueToMapv1(v *Value, data Data, locs []uint64, buf []byte) (err error) {
	key := v.Key.export()
	if key == nil {
		return ErrInvalidInput
	}
	switch v.Vtype { //!make sure to update types in interfaceFromValue as well!
	//standard types
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
		data[key] = int8FromBytes(v.Value)
	case v1Uint64:
		data[key] = uint64FromBytes(v.Value)
	case v1Uint32:
		data[key] = uint32FromBytes(v.Value)
	case v1Uint16:
		data[key] = uint16FromBytes(v.Value)
	case v1Uint8:
		data[key] = uint8FromBytes(v.Value[0])
	case v1Bool:
		data[key] = boolFromBytes(v.Value)

	// special types
	case v1Map:
		val := make(Data, len(v.Children))
		childs := v.Children
		for ck := range childs {
			var err error
			v.fromBytes(buf[locs[childs[ck]]:])
			err = valueToMapv1(v, val, locs, buf)
			if err != nil {
				return err
			}
		}
		data[key] = val
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
	return err
}

func valueToArrayv1(v *Value, arr []interface{}, i int, locs []uint64, buf []byte) (err error) {
	return interfaceFromValue(v, &arr[i], locs, buf)
}

//interfaceFromValue converts a value into an interface{}, you should pass &interface{}
func interfaceFromValue(v *Value, ret *interface{}, locs []uint64, buf []byte) error {
	switch v.Vtype {
	//standard types
	case v1String:
		*ret = stringFromBytes(v.Value)
	case v1Bytes:
		*ret = v.Value
	case v1Float64:
		*ret = float64FromBytes(v.Value)
	case v1Float32:
		*ret = float32FromBytes(v.Value)
	case v1Int64:
		*ret = int64FromBytes(v.Value)
	case v1Int32:
		*ret = int32FromBytes(v.Value)
	case v1Int16:
		*ret = int16FromBytes(v.Value)
	case v1Int8:
		*ret = int8FromBytes(v.Value)
	case v1Uint64:
		*ret = uint64FromBytes(v.Value)
	case v1Uint32:
		*ret = uint32FromBytes(v.Value)
	case v1Uint16:
		*ret = uint16FromBytes(v.Value)
	case v1Uint8:
		*ret = uint8FromBytes(v.Value[0])
	case v1Bool:
		*ret = boolFromBytes(v.Value)

	// special types
	case v1Map:
		val := make(Data, len(v.Children))
		childs := v.Children
		for ck := range childs {
			var err error
			v.fromBytes(buf[locs[childs[ck]]:])
			err = valueToMapv1(v, val, locs, buf)
			if err != nil {
				return err
			}
		}
		*ret = val

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
		*ret = val
	default:
		t := transmitters[v.Vtype]
		if t == nil {
			return errors.New("no transmitter for type:" + string(v.Vtype))
		}
		var err error
		*ret, err = t.Decode(v.Value)
		return err
	}
	return nil
}

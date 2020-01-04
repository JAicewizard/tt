package tt

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"runtime"

	v2 "github.com/JAicewizard/tt/v2"
)

const (
	version1 = 1
	version2 = 2
)

func Encodev2(d Data, values *bytes.Buffer) {
	values.WriteByte(version2)

	tv := v2.Ikeytype(0)
	firstChilds, err := encodemapv1(values, d, &tv)
	if err != nil {
		panic(err)
	}

	v2.AddValue(values, &v2.Value{
		Children: firstChilds,
		Vtype:    v2.FinalValueT,
	})
	//this is bigEndian because then it is compatible with version one of tt
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(tv+1))
	values.Write(buf[:])
}

func encodemapv1(values *bytes.Buffer, d Data, nextValue *v2.Ikeytype) ([]v2.Ikeytype, error) {
	createdObjects := make([]v2.Ikeytype, len(d))
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

func encodevaluev1(values *bytes.Buffer, d interface{}, k interface{}, nextValue *v2.Ikeytype) error {
	var value v2.Value
	var KeepAlive interface{}
	var KeepAlive1 interface{}
	if k != nil {
		switch v := k.(type) { //making this s seperate function will decrese performance, it won't be able to inline and make more allocations
		case string:
			value.Key.Value = v2.StringToBytes(v)
			value.Key.Vtype = v2.StringT
		case []byte:
			value.Key.Value = v
			value.Key.Vtype = v2.BytesT
		case float64:
			v2.Float64ToBytes(&v, &value.Key.Value)
			value.Key.Vtype = v2.Float64T
			KeepAlive = &v
		case float32:
			v2.Float32ToBytes3(&v, &value.Key.Value)
			value.Key.Vtype = v2.Float32T
			KeepAlive = &v
		case int64:
			var buf [8]byte
			v2.Int64ToBytes(v, &buf)
			value.Key.Value = buf[:]
			value.Key.Vtype = v2.Int64T
		case int32:
			value.Key.Value = v2.Int32ToBytes(v)
			value.Key.Vtype = v2.Int32T
		case int16:
			value.Key.Value = v2.Int16ToBytes(v)
			value.Key.Vtype = v2.Int16T
		case int8:
			v2.Int8ToBytes(&v, &value.Key.Value)
			value.Key.Vtype = v2.Int8T
			KeepAlive = &v
		case uint64:
			value.Key.Value = v2.Uint64ToBytes(v)
			value.Key.Vtype = v2.Uint64T
		case uint32:
			value.Key.Value = v2.Uint32ToBytes(v)
			value.Key.Vtype = v2.Uint32T
		case uint16:
			value.Key.Value = v2.Uint16ToBytes(v)
			value.Key.Vtype = v2.Uint16T
		case uint8:
			v2.Uint8ToBytes(&v, &value.Key.Value)
			value.Key.Vtype = v2.Uint8T
			KeepAlive = &v
		case int:
			var buf [8]byte
			v2.Int64ToBytes(int64(v), &buf)
			value.Key.Value = buf[:]
			value.Key.Vtype = v2.Int64T
		case uint:
			value.Key.Value = v2.Uint64ToBytes(uint64(v))
			value.Key.Vtype = v2.BoolT
			KeepAlive = &v
		case bool:
			v2.BoolToBytes(&v, &value.Key.Value)
			value.Key.Vtype = v2.BoolT
			KeepAlive = &v
		}
	}

	//this sets value.Value, it does al the basic types and some more
	switch v := d.(type) {
	case string:
		value.Value = v2.StringToBytes(v)
		value.Vtype = v2.StringT
	case []byte:
		value.Value = v
		value.Vtype = v2.BytesT
	case float64:
		v2.Float64ToBytes(&v, &value.Value)
		value.Vtype = v2.Float64T
		KeepAlive1 = &v
	case float32:
		v2.Float32ToBytes3(&v, &value.Value)
		value.Vtype = v2.Float32T
		KeepAlive1 = &v
	case int64:
		var buf [8]byte
		v2.Int64ToBytes(v, &buf)
		value.Value = buf[:]
		value.Vtype = v2.Int64T
	case int32:
		value.Value = v2.Int32ToBytes(v)
		value.Vtype = v2.Int32T
	case int16:
		value.Value = v2.Int16ToBytes(v)
		value.Vtype = v2.Int16T
	case int8:
		v2.Int8ToBytes(&v, &value.Value)
		value.Vtype = v2.Int8T
		KeepAlive1 = &v
	case uint64:
		value.Value = v2.Uint64ToBytes(v)
		value.Vtype = v2.Uint64T
	case uint32:
		value.Value = v2.Uint32ToBytes(v)
		value.Vtype = v2.Uint32T
	case uint16:
		value.Value = v2.Uint16ToBytes(v)
		value.Vtype = v2.Uint16T
	case uint8:
		v2.Uint8ToBytes(&v, &value.Value)
		value.Vtype = v2.Uint8T
		KeepAlive1 = &v
	case int:
		var buf [8]byte
		v2.Int64ToBytes(int64(v), &buf)
		value.Value = buf[:]
		value.Vtype = v2.Int64T
	case uint:
		value.Value = v2.Uint64ToBytes(uint64(v))
		value.Vtype = v2.Uint64T
		KeepAlive1 = &v
	case bool:
		v2.BoolToBytes(&v, &value.Value)
		value.Vtype = v2.BoolT
	case []interface{}:
		value.Children = make([]v2.Ikeytype, len(v))

		for i := 0; i < len(v); i++ {
			err := encodevaluev1(values, v[i], nil, nextValue)
			if err != nil {
				return err
			}
			value.Children[i] = *nextValue
			*nextValue++
		}

		value.Vtype = v2.ArrT
	case map[string]interface{}:
		value.Children = make([]v2.Ikeytype, len(v))
		i := 0
		for k := range v {
			encodevaluev1(values, v[k], k, nextValue)

			value.Children[i] = *nextValue
			i++
			*nextValue++
		}
		value.Vtype = v2.MapT
	case map[interface{}]interface{}:
		childs, err := encodemapv1(values, Data(v), nextValue)
		if err != nil {
			return err
		}
		value.Children = childs
		value.Vtype = v2.MapT

	case Data:
		childs, err := encodemapv1(values, v, nextValue)
		if err != nil {
			return err
		}
		value.Children = childs
		value.Vtype = v2.MapT
	default:
		if val := reflect.ValueOf(d); val.Kind() == reflect.Array {
			value.Children = make([]v2.Ikeytype, val.Len())

			for i := 0; i < val.Len(); i++ {
				e := val.Index(i).Interface()
				err := encodevaluev1(values, e, nil, nextValue)
				if err != nil {
					return err
				}
				value.Children[i] = *nextValue
				*nextValue++
			}

			value.Vtype = v2.ArrT
		} else if v, ok := d.(Transmitter); ok {
			var err error
			value.Value, err = v.Encode()
			if err != nil {
				return err
			}
			value.Vtype = v.GetCode()
		} else {
			return v2.ErrInvalidInput
		}
	}
	v2.AddValue(values, &value)
	runtime.KeepAlive(KeepAlive)
	runtime.KeepAlive(KeepAlive1)
	return nil
}

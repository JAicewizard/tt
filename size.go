package tt

import (
	"reflect"
	"runtime"

	v2 "github.com/JAicewizard/tt/v2"
)

//TODO: maybe dont copy all the data over etc, it might be faster do just return the length and use that

//Size gives an accurate size value for the final size of buffer needed by encoding
func (d Data) Size() (int, error) {
	var size = 1 //the version ID is included
	var tv = v2.Ikeytype(0)
	firstChilds, err := sizeMapv2(&size, d, &tv)
	if err != nil {
		return 0, err
	}
	//see encodev1 for reference of why this is needed
	size += (&v2.Value{
		Children: firstChilds,
		Vtype:    v2.FinalValueT,
	}).Len()
	size += 4
	return size, nil
}

func sizeMapv2(size *int, d Data, nextValue *v2.Ikeytype) ([]v2.Ikeytype, error) {
	createdObjects := make([]v2.Ikeytype, len(d))
	i := 0

	for k := range d {
		err := sizeValuev2(size, d[k], k, nextValue)
		if err != nil {
			return nil, err
		}
		createdObjects[i] = *nextValue
		i++
		*nextValue++
	}

	return createdObjects, nil
}

func sizeValuev2(size *int, d interface{}, k interface{}, nextValue *v2.Ikeytype) error {
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
			value.Key.Vtype = v2.Uint64T
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
			err := sizeValuev2(size, v[i], nil, nextValue)
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
			sizeValuev2(size, v[k], k, nextValue)

			value.Children[i] = *nextValue
			i++
			*nextValue++
		}
		value.Vtype = v2.MapT
	case map[interface{}]interface{}:
		childs, err := sizeMapv2(size, Data(v), nextValue)
		if err != nil {
			return err
		}
		value.Children = childs
		value.Vtype = v2.MapT

	case Data:
		childs, err := sizeMapv2(size, v, nextValue)
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
				err := sizeValuev2(size, e, nil, nextValue)
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
	*size += value.Len()
	runtime.KeepAlive(KeepAlive)
	runtime.KeepAlive(KeepAlive1)
	return nil
}

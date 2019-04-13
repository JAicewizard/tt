package tt

import (
	"reflect"
	"runtime"
)

//TODO: maybe dont copy all the data over etc, it might be faster do just return the length and use that

//Sizev1 gives an accurate size value for the final size of buffer needed by encodeV1
func (d Data) Sizev1() (int, error) {
	var size = 1 //the version ID is included
	var tv = ikeytype(0)
	firstChilds, err := sizeMapv1(&size, d, &tv)
	if err != nil {
		return 0, err
	}
	//see encodev1 for reference of why this is needed
	size += (&Value{
		Children: firstChilds,
		Vtype:    v1FinalValue,
	}).len()
	size++
	return size, nil
}

func sizeMapv1(size *int, d Data, nextValue *ikeytype) ([]ikeytype, error) {
	createdObjects := make([]ikeytype, len(d))
	i := 0

	for k := range d {
		err := sizeValuev1(size, d[k], k, nextValue)
		if err != nil {
			return nil, err
		}
		createdObjects[i] = *nextValue
		i++
		*nextValue++
	}

	return createdObjects, nil
}

func sizeValuev1(size *int, d interface{}, k interface{}, nextValue *ikeytype) error {
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
			err := sizeValuev1(size, v[i], nil, nextValue)
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
			sizeValuev1(size, v[k], k, nextValue)

			value.Children[i] = *nextValue
			i++
			*nextValue++
		}
		value.Vtype = v1Map
	case map[interface{}]interface{}:
		childs, err := sizeMapv1(size, Data(v), nextValue)
		if err != nil {
			return err
		}
		value.Children = childs
		value.Vtype = v1Map

	case Data:
		childs, err := sizeMapv1(size, v, nextValue)
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
				err := sizeValuev1(size, e, nil, nextValue)
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
	*size += value.len()
	runtime.KeepAlive(KeepAlive)
	runtime.KeepAlive(KeepAlive1)
	return nil
}

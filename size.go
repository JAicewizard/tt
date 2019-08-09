package tt

import(
	"reflect"
	"runtime"
	"github.com/JAicewizard/tt/v1"
)

//TODO: maybe dont copy all the data over etc, it might be faster do just return the length and use that

//Size gives an accurate size value for the final size of buffer needed by encoding
func (d Data) Size() (int, error) {
	var size = 1 //the version ID is included
	var tv = v1.Ikeytype(0)
	firstChilds, err := sizeMapv1(&size, d, &tv)
	if err != nil {
		return 0, err
	}
	//see encodev1 for reference of why this is needed
	size += (&v1.Value{
		Children: firstChilds,
		Vtype:    finalValueT,
	}).Len()
	size++
	return size, nil
}

func sizeMapv1(size *int, d Data, nextValue *v1.Ikeytype) ([]v1.Ikeytype, error) {
	createdObjects := make([]v1.Ikeytype, len(d))
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

func sizeValuev1(size *int, d interface{}, k interface{}, nextValue *v1.Ikeytype) error {
	var value v1.Value
	var KeepAlive interface{}
	var KeepAlive1 interface{}
	if k != nil {
		switch v := k.(type) { //making this s seperate function will decrese performance, it won't be able to inline and make more allocations
		case string:
			value.Key.Value = v1.StringToBytes(v)
			value.Key.Vtype = stringT
		case []byte:
			value.Key.Value = v
			value.Key.Vtype = bytesT
		case float64:
			v1.Float64ToBytes(&v, &value.Key.Value)
			value.Key.Vtype = float64T
			KeepAlive = &v
		case float32:
			v1.Float32ToBytes3(&v, &value.Key.Value)
			value.Key.Vtype = float32T
			KeepAlive = &v
		case int64:
			var buf [8]byte
			v1.Int64ToBytes(v, &buf)
			value.Key.Value = buf[:]
			value.Key.Vtype = int64T
		case int32:
			value.Key.Value = v1.Int32ToBytes(v)
			value.Key.Vtype = int32T
		case int16:
			value.Key.Value = v1.Int16ToBytes(v)
			value.Key.Vtype = int16T
		case int8:
			v1.Int8ToBytes(&v, &value.Key.Value)
			value.Key.Vtype = int8T
			KeepAlive = &v
		case uint64:
			value.Key.Value = v1.Uint64ToBytes(v)
			value.Key.Vtype = uint64T
		case uint32:
			value.Key.Value = v1.Uint32ToBytes(v)
			value.Key.Vtype = uint32T
		case uint16:
			value.Key.Value = v1.Uint16ToBytes(v)
			value.Key.Vtype = uint16T
		case uint8:
			v1.Uint8ToBytes(&v, &value.Key.Value)
			value.Key.Vtype = uint8T
			KeepAlive = &v
		case int:
			var buf [8]byte
			v1.Int64ToBytes(int64(v), &buf)
			value.Key.Value = buf[:]
			value.Key.Vtype = int64T
		case uint:
			value.Key.Value = v1.Uint64ToBytes(uint64(v))
			value.Key.Vtype = boolT
			KeepAlive = &v
		case bool:
			v1.BoolToBytes(&v, &value.Key.Value)
			value.Key.Vtype = boolT
			KeepAlive = &v
		}

	}

	//this sets value.Value, it does al the basic types and some more
	switch v := d.(type) {
	case string:
		value.Value = v1.StringToBytes(v)
		value.Vtype = stringT
	case []byte:
		value.Value = v
		value.Vtype = bytesT
	case float64:
		v1.Float64ToBytes(&v, &value.Value)
		value.Vtype = float64T
		KeepAlive1 = &v
	case float32:
		v1.Float32ToBytes3(&v, &value.Value)
		value.Vtype = float32T
		KeepAlive1 = &v
	case int64:
		var buf [8]byte
		v1.Int64ToBytes(v, &buf)
		value.Value = buf[:]
		value.Vtype = int64T
	case int32:
		value.Value = v1.Int32ToBytes(v)
		value.Vtype = int32T
	case int16:
		value.Value = v1.Int16ToBytes(v)
		value.Vtype = int16T
	case int8:
		v1.Int8ToBytes(&v, &value.Value)
		value.Vtype = int8T
		KeepAlive1 = &v
	case uint64:
		value.Value = v1.Uint64ToBytes(v)
		value.Vtype = uint64T
	case uint32:
		value.Value = v1.Uint32ToBytes(v)
		value.Vtype = uint32T
	case uint16:
		value.Value = v1.Uint16ToBytes(v)
		value.Vtype = uint16T
	case uint8:
		v1.Uint8ToBytes(&v, &value.Value)
		value.Vtype = uint8T
		KeepAlive1 = &v
	case int:
		var buf [8]byte
		v1.Int64ToBytes(int64(v), &buf)
		value.Value = buf[:]
		value.Vtype = int64T
	case uint:
		value.Value = v1.Uint64ToBytes(uint64(v))
		value.Vtype = boolT
		KeepAlive1 = &v
	case bool:
		v1.BoolToBytes(&v, &value.Value)
		value.Vtype = boolT
	case []interface{}:
		value.Children = make([]v1.Ikeytype, len(v))

		for i := 0; i < len(v); i++ {
			err := sizeValuev1(size, v[i], nil, nextValue)
			if err != nil {
				return err
			}
			value.Children[i] = *nextValue
			*nextValue++
		}

		value.Vtype = arrT
	case map[string]interface{}:
		value.Children = make([]v1.Ikeytype, len(v))
		i := 0
		for k := range v {
			sizeValuev1(size, v[k], k, nextValue)

			value.Children[i] = *nextValue
			i++
			*nextValue++
		}
		value.Vtype = mapT
	case map[interface{}]interface{}:
		childs, err := sizeMapv1(size, Data(v), nextValue)
		if err != nil {
			return err
		}
		value.Children = childs
		value.Vtype = mapT

	case Data:
		childs, err := sizeMapv1(size, v, nextValue)
		if err != nil {
			return err
		}
		value.Children = childs
		value.Vtype = mapT
	default:
		if val := reflect.ValueOf(d); val.Kind() == reflect.Array {
			value.Children = make([]v1.Ikeytype, val.Len())

			for i := 0; i < val.Len(); i++ {
				e := val.Index(i).Interface()
				err := sizeValuev1(size, e, nil, nextValue)
				if err != nil {
					return err
				}
				value.Children[i] = *nextValue
				*nextValue++
			}

			value.Vtype = arrT
		} else if v, ok := d.(Transmitter); ok {
			var err error
			value.Value, err = v.Encode()
			if err != nil {
				return err
			}
			value.Vtype = v.GetCode()
		} else {
			return v1.ErrInvalidInput
		}
	}
	*size += value.Len()
	runtime.KeepAlive(KeepAlive)
	runtime.KeepAlive(KeepAlive1)
	return nil
}

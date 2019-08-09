package tt

import (
	"errors"

	v1 "github.com/JAicewizard/tt/v1"
)

var (
	transmitters = make(map[byte]Transmitter)
)

func Decodev1(b []byte, d *Data) (err error) {
	vlen := int(b[len(b)-1])

	locs := make([]uint64, vlen)
	v1.GetLocs(b, locs, vlen)

	//decoding the actual values
	var v v1.Value
	v.FromBytes(b[locs[vlen-1]:])

	if *d == nil {
		*d = make(Data, len(v.Children)*1)
	}

	data := d
	childs := v.Children

	for ck := range childs {
		var err error
		v.FromBytes(b[locs[childs[ck]]:])

		err = valueToMapv1(&v, *data, locs, b)

		if err != nil {
			return err
		}
	}
	return nil
}

func valueToMapv1(v *v1.Value, dat Data, locs []uint64, buf []byte) (err error) {
	key := v.Key.Export()
	if key == nil {
		return v1.ErrInvalidInput
	}
	switch v.Vtype { //!make sure to update types in interfaceFromValue as well!
	//standard types
	case stringT:
		dat[key] = v1.StringFromBytes(v.Value)
	case bytesT:
		dat[key] = v.Value
	case float64T:
		dat[key] = v1.Float64FromBytes(v.Value)
	case float32T:
		dat[key] = v1.Float32FromBytes(v.Value)
	case int64T:
		dat[key] = v1.Int64FromBytes(v.Value)
	case int32T:
		dat[key] = v1.Int32FromBytes(v.Value)
	case int16T:
		dat[key] = v1.Int16FromBytes(v.Value)
	case int8T:
		dat[key] = v1.Int8FromBytes(v.Value)
	case uint64T:
		dat[key] = v1.Uint64FromBytes(v.Value)
	case uint32T:
		dat[key] = v1.Uint32FromBytes(v.Value)
	case uint16T:
		dat[key] = v1.Uint16FromBytes(v.Value)
	case uint8T:
		dat[key] = v1.Uint8FromBytes(v.Value[0])
	case boolT:
		dat[key] = v1.BoolFromBytes(v.Value)

	// special types
	case mapT:
		val := make(Data, len(v.Children))
		childs := v.Children
		for ck := range childs {
			var err error
			v.FromBytes(buf[locs[childs[ck]]:])
			err = valueToMapv1(v, val, locs, buf)
			if err != nil {
				return err
			}
		}
		dat[key] = val
	case arrT:
		val := make([]interface{}, len(v.Children))
		childs := v.Children
		for i := range childs {
			var err error
			v.FromBytes(buf[locs[childs[i]]:])
			err = interfaceFromValue(v, &val[i], locs, buf)
			if err != nil {
				return err
			}
		}
		dat[key] = val

	default:
		t := transmitters[v.Vtype]
		if t == nil {
			return errors.New("no transmitter for type:" + string(v.Vtype))
		}
		var err error
		dat[key], err = t.Decode(v.Value)
		if err != nil {
			return err
		}
	}
	return err
}

//interfaceFromValue converts a value into an interface{}, you should pass &interface{}
func interfaceFromValue(v *v1.Value, ret *interface{}, locs []uint64, buf []byte) error {
	switch v.Vtype {
	//standard types
	case stringT:
		*ret = v1.StringFromBytes(v.Value)
	case bytesT:
		*ret = v.Value
	case float64T:
		*ret = v1.Float64FromBytes(v.Value)
	case float32T:
		*ret = v1.Float32FromBytes(v.Value)
	case int64T:
		*ret = v1.Int64FromBytes(v.Value)
	case int32T:
		*ret = v1.Int32FromBytes(v.Value)
	case int16T:
		*ret = v1.Int16FromBytes(v.Value)
	case int8T:
		*ret = v1.Int8FromBytes(v.Value)
	case uint64T:
		*ret = v1.Uint64FromBytes(v.Value)
	case uint32T:
		*ret = v1.Uint32FromBytes(v.Value)
	case uint16T:
		*ret = v1.Uint16FromBytes(v.Value)
	case uint8T:
		*ret = v1.Uint8FromBytes(v.Value[0])
	case boolT:
		*ret = v1.BoolFromBytes(v.Value)

	// special types
	case mapT:
		val := make(Data, len(v.Children))
		childs := v.Children
		for ck := range childs {
			var err error
			v.FromBytes(buf[locs[childs[ck]]:])
			err = valueToMapv1(v, val, locs, buf)
			if err != nil {
				return err
			}
		}
		*ret = val

	case arrT:
		val := make([]interface{}, len(v.Children))
		childs := v.Children
		for i := range childs {
			var err error
			v.FromBytes(buf[locs[childs[i]]:])
			err = interfaceFromValue(v, &val[i], locs, buf)
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

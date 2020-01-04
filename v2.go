package tt

import (
	"encoding/binary"
	"errors"

	v2 "github.com/JAicewizard/tt/v2"
)

func Decodev2(b []byte, d *Data) (err error) {
	vlen := binary.BigEndian.Uint32(b[len(b)-4 : len(b)])
	locs := make([]uint64, vlen)
	v2.GetLocs(b, locs, vlen)

	//decoding the actual values
	var v v2.Value
	v.FromBytes(b[locs[vlen-1]:])

	if *d == nil {
		*d = make(Data, len(v.Children)*1)
	}

	data := d
	childs := v.Children

	for ck := range childs {
		var err error
		v.FromBytes(b[locs[childs[ck]]:])

		err = valueToMapv2(&v, *data, locs, b)

		if err != nil {
			return err
		}
	}
	return nil
}

func valueToMapv2(v *v2.Value, dat Data, locs []uint64, buf []byte) (err error) {
	key := v.Key.Export()
	if key == nil {
		return v2.ErrInvalidInput
	}
	switch v.Vtype { //!make sure to update types in interfaceFromValuev1 as well!
	//standard types
	case v2.StringT:
		dat[key] = v2.StringFromBytes(v.Value)
	case v2.BytesT:
		dat[key] = v.Value
	case v2.Float64T:
		dat[key] = v2.Float64FromBytes(v.Value)
	case v2.Float32T:
		dat[key] = v2.Float32FromBytes(v.Value)
	case v2.Int64T:
		dat[key] = v2.Int64FromBytes(v.Value)
	case v2.Int32T:
		dat[key] = v2.Int32FromBytes(v.Value)
	case v2.Int16T:
		dat[key] = v2.Int16FromBytes(v.Value)
	case v2.Int8T:
		dat[key] = v2.Int8FromBytes(v.Value)
	case v2.Uint64T:
		dat[key] = v2.Uint64FromBytes(v.Value)
	case v2.Uint32T:
		dat[key] = v2.Uint32FromBytes(v.Value)
	case v2.Uint16T:
		dat[key] = v2.Uint16FromBytes(v.Value)
	case v2.Uint8T:
		dat[key] = v2.Uint8FromBytes(v.Value[0])
	case v2.BoolT:
		dat[key] = v2.BoolFromBytes(v.Value)

	// special types
	case v2.MapT:
		val := make(Data, len(v.Children))
		childs := v.Children
		for ck := range childs {
			var err error
			v.FromBytes(buf[locs[childs[ck]]:])
			err = valueToMapv2(v, val, locs, buf)
			if err != nil {
				return err
			}
		}
		dat[key] = val
	case v2.ArrT:
		val := make([]interface{}, len(v.Children))
		childs := v.Children
		for i := range childs {
			var err error
			v.FromBytes(buf[locs[childs[i]]:])
			err = interfaceFromValuev2(v, &val[i], locs, buf)
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

//interfaceFromValuev1 converts a value into an interface{}, you should pass &interface{}
func interfaceFromValuev2(v *v2.Value, ret *interface{}, locs []uint64, buf []byte) error {
	switch v.Vtype {
	//standard types
	case v2.StringT:
		*ret = v2.StringFromBytes(v.Value)
	case v2.BytesT:
		*ret = v.Value
	case v2.Float64T:
		*ret = v2.Float64FromBytes(v.Value)
	case v2.Float32T:
		*ret = v2.Float32FromBytes(v.Value)
	case v2.Int64T:
		*ret = v2.Int64FromBytes(v.Value)
	case v2.Int32T:
		*ret = v2.Int32FromBytes(v.Value)
	case v2.Int16T:
		*ret = v2.Int16FromBytes(v.Value)
	case v2.Int8T:
		*ret = v2.Int8FromBytes(v.Value)
	case v2.Uint64T:
		*ret = v2.Uint64FromBytes(v.Value)
	case v2.Uint32T:
		*ret = v2.Uint32FromBytes(v.Value)
	case v2.Uint16T:
		*ret = v2.Uint16FromBytes(v.Value)
	case v2.Uint8T:
		*ret = v2.Uint8FromBytes(v.Value[0])
	case v2.BoolT:
		*ret = v2.BoolFromBytes(v.Value)

	// special types
	case v2.MapT:
		val := make(Data, len(v.Children))
		childs := v.Children
		for ck := range childs {
			var err error
			v.FromBytes(buf[locs[childs[ck]]:])
			err = valueToMapv2(v, val, locs, buf)
			if err != nil {
				return err
			}
		}
		*ret = val

	case v2.ArrT:
		val := make([]interface{}, len(v.Children))
		childs := v.Children
		for i := range childs {
			var err error
			v.FromBytes(buf[locs[childs[i]]:])
			err = interfaceFromValuev2(v, &val[i], locs, buf)
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

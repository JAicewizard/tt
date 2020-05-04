package tt

import (
	"encoding/binary"
	"encoding/gob"
	"errors"
	"io"
	"reflect"
	"runtime"
	"sync"

	v3 "github.com/jaicewizard/tt/v3"
)

//V3Encoder is the encoder used to encode a ttv3 data stream
type V3Encoder struct {
	out       io.Writer
	varintbuf *[binary.MaxVarintLen64 + 1]byte
	sync.Mutex
}

var v3StreamHeader = []byte{version3, 1 << 7}
var v3NoStreamHeader = []byte{version3, 0}

//NewV3Encoder creates a new encoder to encode a ttv3 data stream
func NewV3Encoder(out io.Writer, isStream bool) *V3Encoder {
	if isStream {
		out.Write(v3StreamHeader)
	} else {
		out.Write(v3NoStreamHeader)
	}
	return &V3Encoder{
		out:       out,
		varintbuf: &[binary.MaxVarintLen64 + 1]byte{},
	}
}

//Encodev3 encodes an `interface{}`` into a bytebuffer using ttv3
func Encodev3(d interface{}, out io.Writer) error {
	out.Write(v3NoStreamHeader)

	enc := &V3Encoder{
		out:       out,
		varintbuf: &[binary.MaxVarintLen64 + 1]byte{},
	}
	//We dont have to lock/unlock since we know we are the only one witha acces
	return enc.encodeValuev3(d, v3.Key{})
}

//Encode encodes an `interface{}`` into a bytebuffer using ttv3
func (enc *V3Encoder) Encode(d interface{}) error {
	enc.Lock()
	ret := enc.encodeValuev3(d, v3.Key{})
	enc.Unlock()
	return ret
}

func encodeKeyv3(k interface{}) v3.Key {
	var key v3.Key
	if k != nil {
		switch v := k.(type) { //making this s seperate function will decrese performance, it won't be able to inline and make more allocations
		case string:
			key.Value = v3.StringToBytes(v)
			key.Vtype = v3.StringT
		case []byte:
			key.Value = v
			key.Vtype = v3.BytesT
		case int8:
			key.Value = []byte{v3.Int8ToBytes(v)}
			key.Vtype = v3.Int8T
		case int16:
			key.Value = v3.Int16ToBytes(v)
			key.Vtype = v3.Int16T
		case int32:
			key.Value = v3.Int32ToBytes(v)
			key.Vtype = v3.Int32T
		case int64:
			var buf [8]byte
			v3.Int64ToBytes(v, &buf)
			key.Value = buf[:]
			key.Vtype = v3.Int64T
		case int:
			var buf [8]byte
			v3.Int64ToBytes(int64(v), &buf)
			key.Value = buf[:]
			key.Vtype = v3.Int64T
		case uint8:
			key.Value = []byte{v3.Uint8ToBytes(v)}
			key.Vtype = v3.Uint8T
		case uint16:
			key.Value = v3.Uint16ToBytes(v)
			key.Vtype = v3.Uint16T
		case uint32:
			key.Value = v3.Uint32ToBytes(v)
			key.Vtype = v3.Uint32T
		case uint64:
			var buf [8]byte
			v3.Uint64ToBytes(v, &buf)
			key.Value = buf[:]
			key.Vtype = v3.Uint64T
		case uint:
			var buf [8]byte
			v3.Uint64ToBytes(uint64(v), &buf)
			key.Value = buf[:]
			key.Vtype = v3.Uint64T
		case float32:
			key.Value = v3.Float32ToBytes(&v)
			key.Vtype = v3.Float32T
		case float64:
			var buf [8]byte
			v3.Float64ToBytes(v, &buf)
			key.Value = buf[:]
			key.Vtype = v3.Float64T
		case bool:
			key.Value = v3.BoolToBytes(v)
			key.Vtype = v3.BoolT
		}
	}
	return key
}

func encodeString(s string) v3.Key {
	return v3.Key{
		Value: v3.StringToBytes(s),
		Vtype: v3.StringT,
	}
}
func encodeBytes(b []byte) v3.Key {
	return v3.Key{
		Value: b,
		Vtype: v3.BytesT,
	}
}

func (enc *V3Encoder) encodeValuev3(d interface{}, k v3.Key) error {
	value := v3.Value{
		Key: k,
	}
	var KeepAlive interface{}
	alreadyEncoded := false
	//this sets value.Value, it does al the basic types and some more
	switch v := d.(type) {
	case string:
		value.Value = v3.StringToBytes(v)
		value.Vtype = v3.StringT
	case []byte:
		value.Value = v
		value.Vtype = v3.BytesT
	case int8:
		value.Value = []byte{v3.Int8ToBytes(v)}
		value.Vtype = v3.Int8T
		KeepAlive = &v
	case int16:
		value.Value = v3.Int16ToBytes(v)
		value.Vtype = v3.Int16T
	case int32:
		value.Value = v3.Int32ToBytes(v)
		value.Vtype = v3.Int32T
	case int64:
		var buf [8]byte
		v3.Int64ToBytes(v, &buf)
		value.Value = buf[:]
		value.Vtype = v3.Int64T
	case int:
		var buf [8]byte
		v3.Int64ToBytes(int64(v), &buf)
		value.Value = buf[:]
		value.Vtype = v3.Int64T
	case uint8:
		value.Value = []byte{v3.Uint8ToBytes(v)}
		value.Vtype = v3.Uint8T
		KeepAlive = &v
	case uint16:
		value.Value = v3.Uint16ToBytes(v)
		value.Vtype = v3.Uint16T
	case uint32:
		value.Value = v3.Uint32ToBytes(v)
		value.Vtype = v3.Uint32T
	case uint64:
		var buf [8]byte
		v3.Uint64ToBytes(uint64(v), &buf)
		value.Value = buf[:]
		value.Vtype = v3.Uint64T
	case uint:
		var buf [8]byte
		v3.Uint64ToBytes(uint64(v), &buf)
		value.Value = buf[:]
		value.Vtype = v3.Uint64T
	case float32:
		value.Value = v3.Float32ToBytes(&v)
		value.Vtype = v3.Float32T
		KeepAlive = &v
	case float64:
		var buf [8]byte
		v3.Float64ToBytes(v, &buf)
		value.Value = buf[:]
		value.Vtype = v3.Float64T
	case bool:
		value.Value = v3.BoolToBytes(v)
		value.Vtype = v3.BoolT
		KeepAlive = &v

	default:
		val := reflect.ValueOf(d)
		kind := val.Kind()
		if kind == reflect.Map {
			value.Childrenn = uint64(val.Len())
			value.Vtype = v3.MapT
			alreadyEncoded = true
			v3.AddValue(enc.out, &value, enc.varintbuf)

			switch v := d.(type) {
			case map[string]string:
				for k, v := range v {
					enc.encodeValuev3(v, encodeString(k))
				}
			case map[string]interface{}:
				for k, v := range v {
					enc.encodeValuev3(v, encodeString(k))
				}
			case map[interface{}]interface{}:
				for k, v := range v {
					enc.encodeValuev3(v, encodeKeyv3(k))
				}
			default:
				//if its not a specific map type
				mapRange := val.MapRange()
				for mapRange.Next() {
					enc.encodeValuev3(mapRange.Value().Interface(), encodeKeyv3(mapRange.Key().Interface()))
				}
			}
		} else if kind == reflect.Array || kind == reflect.Slice {
			value.Childrenn = uint64(val.Len())
			value.Vtype = v3.ArrT
			alreadyEncoded = true
			v3.AddValue(enc.out, &value, enc.varintbuf)
			switch s := d.(type) {
			case []string:
				for _, v := range s {
					enc.encodeValuev3(v, v3.Key{})
				}
			default:
				//if its not a specific slice type
				for i := 0; i < val.Len(); i++ {
					err := enc.encodeValuev3(val.Index(i).Interface(), v3.Key{})
					if err != nil {
						return err
					}
				}
			}
		} else if v, ok := d.(gob.GobEncoder); ok {
			var err error
			value.Value, err = v.GobEncode()
			if err != nil {
				return err
			}
			value.Vtype = v3.BytesT
		} else if kind == reflect.Struct {
			usableFields := getStructFields(val)
			value.Childrenn = uint64(len(usableFields))
			value.Vtype = v3.MapT
			alreadyEncoded = true
			v3.AddValue(enc.out, &value, enc.varintbuf)
			for fieldName, fieldID := range usableFields {
				field := val.Field(fieldID)
				e := field.Interface()
				err := enc.encodeValuev3(e, encodeBytes([]byte(fieldName)))
				if err != nil {
					return err
				}
			}
		} else {
			return v3.ErrInvalidInput
		}
	}

	if !alreadyEncoded {
		v3.AddValue(enc.out, &value, enc.varintbuf)
	}
	runtime.KeepAlive(KeepAlive)

	return nil
}

func getStructFields(val reflect.Value) map[string]int {
	usableFields := make(map[string]int, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if field.PkgPath != "" || !val.Field(i).CanInterface() {
			continue
		}
		usableFields[getFieldName(field)] = i
	}
	return usableFields
}

//V3Decoder is the decoder used to decode a ttv3 data stream
type V3Decoder struct {
	didInit   bool
	isStream  bool
	didDecode bool
	in        v3.Reader
	sync.Mutex
}

//NewV3Decoder creates aa new V3Decoder to decode a ttv3 data stream.
//The init flag specifies wether it should initialize the decoder.
//Initializing the decoder blocks until at least the first 2 bytes are read.
func NewV3Decoder(in v3.Reader, init bool) *V3Decoder {
	dec := V3Decoder{
		didInit: !init,
		in:      in,
	}
	if init {
		dec.Init()
	}
	return &dec
}

//Init initizes the decoder, initizlizing blocks until at least the first
//2 bytes are read.
func (dec *V3Decoder) Init() error {
	version, err := dec.in.ReadByte()
	if version != 3 || err != nil {
		return v3.ErrInvalidInput
	}
	b, err := dec.in.ReadByte()
	if err != nil {
		return v3.ErrInvalidInput
	}
	dec.isStream = b&(1<<7) != 0
	dec.didDecode = false
	return nil
}

//Decode decodes a one ttv3 encoded value from a stream.
//Note that a stream of one value is the same as one value just with
//the stream bit set
func (dec *V3Decoder) Decode(e interface{}) error {
	dec.Lock()
	ret := dec.decode(e)
	dec.Unlock()
	return ret
}

//Decodev3 decodes a ttv3 encoded byte-slice into tt.Data
func Decodev3(buf v3.Reader, e interface{}) error {
	dec := NewV3Decoder(buf, true)
	//We dont have to lock/unlock since we know we are the only one witha access
	return dec.decode(e)
}

func (dec *V3Decoder) decode(e interface{}) error {
	if !dec.isStream {
		if dec.didDecode {
			return errors.New("TT: attempt to decode twice from non-stream")
		}
		dec.didDecode = true
	}
	if e == nil {
		var v v3.Value
		v.FromBytes(dec.in)
		clearNextValues(dec.in, v.Childrenn)
		return nil
	}
	value := reflect.ValueOf(e)

	// If e represents a value as opposed to a pointer, the answer won't
	// get back to the caller. Make sure it's a pointer.
	if value.Type().Kind() != reflect.Ptr {
		return errors.New("TT: attempt to decode into a non-pointer")
	}
	if value.IsValid() {
		if value.Kind() == reflect.Ptr && !value.IsNil() {
			// That's okay, we'll store through the pointer.
		} else if !value.CanSet() {
			return errors.New("TT: DecodeValue of unassignable value")
		}
	} else {
		var v v3.Value
		v.FromBytes(dec.in)
		clearNextValues(dec.in, v.Childrenn)
		return nil
	}

	var v v3.Value
	v.FromBytes(dec.in)
	yetToRead := v.Childrenn

	err := dec.decodeValuev3(v, value, &yetToRead)
	if yetToRead != 0 {
		clearNextValues(dec.in, yetToRead)
	}

	return err
}

func decodeKeyv3(k v3.Key, buf v3.Reader, e reflect.Value) error {
	if k.Vtype == v3.StringT {
		val := v3.StringFromBytes(k.Value)
		if e.Kind() != reflect.String {
			if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
				return errors.New("TT: cannot unmarshal string into " + e.Type().String() + " Go type")
			}
			e.Set(reflect.ValueOf(val))
		} else {
			e.SetString(val)
		}
	} else if k.Vtype == v3.BytesT {
		val := k.Value
		if e.Kind() != reflect.Slice || e.Type().Elem().Kind() != reflect.Uint8 {
			if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
				return errors.New("TT: cannot unmarshal bytes into " + e.Type().String() + " Go type")
			}
			e.Set(reflect.ValueOf(val))
		} else {
			e.SetBytes(val)
		}
	} else if k.Vtype == v3.Float32T {
		e.SetFloat(float64(v3.Float32FromBytes(k.Value)))
	} else if k.Vtype == v3.Float64T {
		e.SetFloat(v3.Float64FromBytes(k.Value))
	} else if k.Vtype == v3.Int8T {
		e.SetInt(int64(v3.Int8FromBytes(k.Value[0])))
	} else if k.Vtype == v3.Int16T {
		e.SetInt(int64(v3.Int16FromBytes(k.Value)))
	} else if k.Vtype == v3.Int32T {
		e.SetInt(int64(v3.Int32FromBytes(k.Value)))
	} else if k.Vtype == v3.Int64T {
		e.SetInt(int64(v3.Int64FromBytes(k.Value)))
	} else if k.Vtype == v3.Uint8T {
		e.SetUint(uint64(v3.Uint8FromBytes(k.Value[0])))
	} else if k.Vtype == v3.Uint16T {
		e.SetUint(uint64(v3.Uint16FromBytes(k.Value)))
	} else if k.Vtype == v3.Uint32T {
		e.SetUint(uint64(v3.Uint32FromBytes(k.Value)))
	} else if k.Vtype == v3.Uint64T {
		e.SetUint(uint64(v3.Uint64FromBytes(k.Value)))
	} else if k.Vtype == v3.BoolT {
		e.SetBool(v3.BoolFromBytes(k.Value))
	}
	return nil
}
func (dec *V3Decoder) decodeValuev3(v v3.Value, e reflect.Value, yetToRead *uint64) error {
	//copy from json/decode.go:indirect
	haveAddr := false
	e0 := e

	if e.Kind() != reflect.Ptr && e.Type().Name() != "" && e.CanAddr() {
		haveAddr = true
		e = e.Addr()
	}
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if e.Kind() == reflect.Interface && !e.IsNil() {
			te := e.Elem()
			if te.Kind() == reflect.Ptr && !te.IsNil() {
				haveAddr = false
				e = te
				continue
			}
		}

		if e.Kind() != reflect.Ptr {
			break
		}

		// Prevent infinite loop if v is an interface pointing to its own address:
		//     var v interface{}
		//     v = &v
		if e.Elem().Kind() == reflect.Interface && e.Elem().Elem() == e {
			break
		}
		if e.IsNil() {
			e.Set(reflect.New(e.Type().Elem()))
		}

		if e.Type().NumMethod() > 0 && e.CanInterface() {
			if u, ok := e.Interface().(gob.GobDecoder); ok {
				return u.GobDecode(v.Value)
			}
		}

		if haveAddr {
			e = e0 // restore original value after round-trip Value.Addr().Elem()
			haveAddr = false
		} else {
			e = e.Elem()
		}
	}

	switch v.Vtype {
	case v3.StringT:
		val := v3.StringFromBytes(v.Value)
		if e.Kind() != reflect.String {
			if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
				return errors.New("TT: cannot unmarshal string into " + e.Type().String() + " Go type")
			}
			e.Set(reflect.ValueOf(val))
		} else {
			e.SetString(val)
		}
	case v3.BytesT:
		val := v.Value
		if e.Kind() != reflect.Slice || e.Type().Elem().Kind() != reflect.Uint8 {
			if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
				return errors.New("TT: cannot unmarshal bytes into " + e.Type().String() + " Go type")
			}
			e.Set(reflect.ValueOf(val))
		} else {
			e.SetBytes(val)
		}
	case v3.Int8T:
		val := v3.Int8FromBytes(v.Value[0])
		if e.Kind() != reflect.Int8 {
			if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
				return errors.New("TT: cannot unmarshal int8 into " + e.Kind().String() + " Go type")
			}
			e.Set(reflect.ValueOf(val))
		} else {
			e.SetInt(int64(val))
		}
	case v3.Int16T:
		val := v3.Int16FromBytes(v.Value)
		if e.Kind() != reflect.Int16 {
			if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
				return errors.New("TT: cannot unmarshal int16 into " + e.Kind().String() + " Go type")
			}
			e.Set(reflect.ValueOf(val))
		} else {
			e.SetInt(int64(val))
		}
	case v3.Int32T:
		val := v3.Int32FromBytes(v.Value)
		if e.Kind() != reflect.Int32 {
			if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
				return errors.New("TT: cannot unmarshal int32 into " + e.Kind().String() + " Go type")
			}
			e.Set(reflect.ValueOf(val))
		} else {
			e.SetInt(int64(val))
		}
	case v3.Int64T:
		val := v3.Int64FromBytes(v.Value)
		if e.Kind() != reflect.Int64 || e.Kind() != reflect.Int {
			if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
				return errors.New("TT: cannot unmarshal int64 into " + e.Kind().String() + " Go type")
			}
			e.Set(reflect.ValueOf(val))
		} else {
			e.SetInt(val)
		}

	case v3.Uint8T:
		val := v3.Uint8FromBytes(v.Value[0])
		if e.Kind() != reflect.Uint8 {
			if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
				return errors.New("TT: cannot unmarshal uint8 into " + e.Kind().String() + " Go type")
			}
			e.Set(reflect.ValueOf(val))
		} else {
			e.SetUint(uint64(val))
		}
	case v3.Uint16T:
		val := v3.Uint16FromBytes(v.Value)
		if e.Kind() != reflect.Uint16 {
			if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
				return errors.New("TT: cannot unmarshal uint16 into " + e.Kind().String() + " Go type")
			}
			e.Set(reflect.ValueOf(val))
		} else {
			e.SetUint(uint64(val))
		}
	case v3.Uint32T:
		val := v3.Uint32FromBytes(v.Value)
		if e.Kind() != reflect.Uint32 {
			if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
				return errors.New("TT: cannot unmarshal uint32 into " + e.Kind().String() + " Go type")
			}
			e.Set(reflect.ValueOf(val))
		} else {
			e.SetUint(uint64(val))
		}
	case v3.Uint64T:
		val := v3.Uint64FromBytes(v.Value)
		if e.Kind() != reflect.Uint64 {
			if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
				return errors.New("TT: cannot unmarshal uint64 into " + e.Kind().String() + " Go type")
			}
			e.Set(reflect.ValueOf(val))
		} else {
			e.SetUint(val)
		}

	case v3.Float32T:
		val := v3.Float32FromBytes(v.Value)
		if e.Kind() != reflect.Float32 {
			if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
				return errors.New("TT: cannot unmarshal float32 into " + e.Kind().String() + " Go type")
			}
			e.Set(reflect.ValueOf(val))
		} else {
			e.SetFloat(float64(val))
		}
	case v3.Float64T:
		val := v3.Float64FromBytes(v.Value)
		if e.Kind() != reflect.Float64 {
			if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
				return errors.New("TT: cannot unmarshal float64 into " + e.Kind().String() + " Go type")
			}
			e.Set(reflect.ValueOf(val))
		} else {
			e.SetFloat(val)
		}

	case v3.BoolT:
		val := v3.BoolFromBytes(v.Value)
		if e.Kind() != reflect.Bool {
			if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
				return errors.New("TT: cannot unmarshal bytes into " + e.Kind().String() + " Go type")
			}
			e.Set(reflect.ValueOf(val))
		} else {
			e.SetBool(val)
		}

	// special types
	case v3.MapT:
		if e.Kind() == reflect.Interface && e.NumMethod() == 0 {
			children := v.Childrenn
			m := make(map[interface{}]interface{}, children)
			var err error
			key := reflect.New(reflect.TypeOf(m).Key()).Elem()
			for i := uint64(0); i < children; i++ {
				v.FromBytes(dec.in)
				*yetToRead += v.Childrenn - 1
				err = decodeKeyv3(v.Key, dec.in, key)
				if err != nil {
					return err
				}
				k := key.Interface()
				err = dec.decodeValuev3(v, key, yetToRead)
				if err != nil {
					return err
				}
				if v, ok := k.([]byte); ok {
					m[string(v)] = key.Interface()
				} else {
					m[k] = key.Interface()
				}
			}
			e.Set(reflect.ValueOf(m))
		}

		if e.Kind() == reflect.Map {
			children := v.Childrenn
			if e.IsNil() {
				e.Set(reflect.MakeMap(e.Type()))
			}

			var err error
			value := reflect.New(e.Type().Elem()).Elem()
			key := reflect.New(e.Type().Key()).Elem()

			for i := uint64(0); i < children; i++ {
				v.FromBytes(dec.in)
				*yetToRead += v.Childrenn - 1

				err = decodeKeyv3(v.Key, dec.in, key)
				if err != nil {
					return err
				}
				err = dec.decodeValuev3(v, value, yetToRead)
				if err != nil {
					return err
				}

				e.SetMapIndex(key, value)
			}
		} else if e.Kind() == reflect.Struct {
			children := v.Childrenn
			usableFields := getStructFields(e)

			for i := uint64(0); i < children; i++ {
				v.FromBytes(dec.in)
				*yetToRead += v.Childrenn - 1

				key := v.Key.ExportStructID()
				if key == "" {
					clearNextValues(dec.in, v.Childrenn)
					continue
				}
				fieldIndex, ok := usableFields[key]
				if !ok {
					clearNextValues(dec.in, v.Childrenn)
					continue
				}

				field := e.Field(fieldIndex)

				err := dec.decodeValuev3(v, field, yetToRead)
				if err != nil {
					return err
				}
			}
		}
	case v3.ArrT:
		children := v.Childrenn

		if e.Kind() == reflect.Array {
			len := e.Len()
			if len < int(children) {
				break
			}
			for i := 0; i < int(children); i++ {
				v.FromBytes(dec.in)
				*yetToRead += v.Childrenn - 1

				err := dec.decodeValuev3(v, e.Index(i), yetToRead)
				if err != nil {
					return err
				}
			}
			break
		} else if e.Kind() == reflect.Slice {
			len := e.Len()
			if len < int(children) {
				e.Set(reflect.MakeSlice(e.Type(), int(children), int(children)))
			}
			for i := 0; i < int(children); i++ {
				v.FromBytes(dec.in)
				*yetToRead += v.Childrenn - 1
				err := dec.decodeValuev3(v, e.Index(i), yetToRead)
				if err != nil {
					return err
				}
			}
			break
		} else if e.Kind() == reflect.Map {
			//TODO: if kind is also a numeric value it is usable
			e.Type().Key().Kind()
		}

		//if all special cases fail we fall back to []interface{}
		arr := make([]interface{}, children)
		var err error
		value := reflect.New(reflect.TypeOf(arr).Elem()).Elem()

		for i := 0; i < int(children); i++ {
			v.FromBytes(dec.in)
			*yetToRead += v.Childrenn - 1

			err = dec.decodeValuev3(v, value, yetToRead)
			if err != nil {
				return err
			}
			arr[i] = value.Interface()
		}
		e.Set(reflect.ValueOf(arr))
	}
	return nil
}

func getFieldName(field reflect.StructField) string {
	name := field.Tag.Get("TT")

	if name == "" {
		name = field.Name
	}
	return name
}

func clearNextValues(buf v3.Reader, values uint64) {
	var value v3.Value
	for ; values > 0; values-- {
		value.FromBytes(buf)
		values += value.Childrenn
	}
}

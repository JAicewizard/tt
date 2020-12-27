package tt

import (
	"encoding/binary"
	"encoding/gob"
	"errors"
	"io"
	"reflect"
	"strconv"
	"sync"

	v3 "github.com/jaicewizard/tt/v3"
)

var (
	decodersSlice = [...]func(v3.KeyValue, reflect.Value) error{
		decodeUndefined,
		decodeString,
		decodeBytes,
		decodeInt8,
		decodeInt16,
		decodeInt32,
		decodeInt64,
		decodeUint8,
		decodeUint16,
		decodeUint32,
		decodeUint64,
		decodeBool,
		decodeFloat32,
		decodeFloat64,
		decodeUndefined,
		decodeUndefined,
		decodeUndefined,
		decodeUndefined,
	}
	valueDecodersSlice [2]func(*V3Decoder, v3.Value, reflect.Value) error
)

func init() {
	valueDecodersSlice = [2]func(*V3Decoder, v3.Value, reflect.Value) error{
		decodeMap,
		decodeArr,
	}
}

/*
no-map
BenchmarkV3               254824              4413 ns/op            1186 B/op         34 allocs/op
BenchmarkV3Decode         347130              3244 ns/op             810 B/op         29 allocs/op
BenchmarkV3Encode        1000000              1032 ns/op             374 B/op          5 allocs/op
BenchmarkV3int64         1239445               932 ns/op             510 B/op         18 allocs/op
slice
BenchmarkV3              2838296              4335 ns/op            1300 B/op         34 allocs/op
BenchmarkV3Decode        3853204              3147 ns/op             810 B/op         29 allocs/op
BenchmarkV3Encode       12269514               970 ns/op             445 B/op          5 allocs/op
BenchmarkV3int64        13352065               898 ns/op             626 B/op         18 allocs/op
*/

//V3Encoder is the encoder used to encode a ttv3 data stream
type V3Encoder struct {
	out       io.Writer
	varintbuf *[binary.MaxVarintLen64 + 1]byte
	sync.Mutex
	typeCache map[string]map[string]int
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
		typeCache: map[string]map[string]int{},
	}
}

//Encodev3 encodes an `interface{}`` into a bytebuffer using ttv3
func Encodev3(d interface{}, out io.Writer) error {
	out.Write(v3NoStreamHeader)

	enc := &V3Encoder{
		out:       out,
		varintbuf: &[binary.MaxVarintLen64 + 1]byte{},
		typeCache: map[string]map[string]int{},
	}
	//We dont have to lock/unlock since we know we are the only one witha acces
	return enc.encodeValuev3(d, v3.KeyValue{})
}

//Encode encodes an `interface{}`` into a bytebuffer using ttv3
func (enc *V3Encoder) Encode(d interface{}) error {
	enc.Lock()
	ret := enc.encodeValuev3(d, v3.KeyValue{})
	enc.Unlock()
	return ret
}

func encodeKeyv3(k interface{}) v3.KeyValue {
	var key v3.KeyValue
	if k != nil {
		switch v := k.(type) { //making this a seperate function will decrese performance, it won't be able to inline and make more allocations
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
			key.Value = v3.Int64ToBytes(v)
			key.Vtype = v3.Int64T
		case int:
			key.Value = v3.Int64ToBytes(int64(v))
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
			key.Value = v3.Uint64ToBytes(v)
			key.Vtype = v3.Uint64T
		case uint:
			key.Value = v3.Uint64ToBytes(uint64(v))
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

func encodeKeyv3_reflect(d reflect.Value) v3.KeyValue {
	var key v3.KeyValue
	switch d.Type().Kind() {
	case reflect.Interface, reflect.Ptr:
		encodeKeyv3_reflect(d.Elem())
	case reflect.String:
		key.Value = v3.StringToBytes(d.String())
		key.Vtype = v3.StringT
	case reflect.Slice:
		if d.Type().Elem().Kind() == reflect.Int8 {
			key.Value = d.Bytes()
			key.Vtype = v3.StringT
		}
	case reflect.Int8:
		key.Value = []byte{v3.Int8ToBytes(int8(d.Int()))}
		key.Vtype = v3.Int8T
	case reflect.Int16:
		key.Value = v3.Int16ToBytes(int16(d.Int()))
		key.Vtype = v3.Int16T
	case reflect.Int32:
		key.Value = v3.Int32ToBytes(int32(d.Int()))
		key.Vtype = v3.Int32T
	case reflect.Int64:
		key.Value = v3.Int64ToBytes(d.Int())
		key.Vtype = v3.Int64T
	case reflect.Int:
		key.Value = v3.Int64ToBytes(d.Int())
		key.Vtype = v3.Int64T
	case reflect.Uint8:
		key.Value = []byte{v3.Uint8ToBytes(uint8(d.Uint()))}
		key.Vtype = v3.Uint8T
	case reflect.Uint16:
		key.Value = v3.Uint16ToBytes(uint16(d.Uint()))
		key.Vtype = v3.Uint16T
	case reflect.Uint32:
		key.Value = v3.Uint32ToBytes(uint32(d.Uint()))
		key.Vtype = v3.Uint32T
	case reflect.Uint64:
		key.Value = v3.Uint64ToBytes(d.Uint())
		key.Vtype = v3.Uint64T
	case reflect.Uint:
		key.Value = v3.Uint64ToBytes(d.Uint())
		key.Vtype = v3.Uint64T
	case reflect.Bool:
		key.Value = v3.BoolToBytes(d.Bool())
		key.Vtype = v3.BoolT
	case reflect.Float32:
		v := float32(d.Float())
		key.Value = v3.Float32ToBytes(&v)
		key.Vtype = v3.Float32T
	case reflect.Float64:
		v := d.Float()
		var buf [8]byte
		v3.Float64ToBytes(v, &buf)
		key.Value = buf[:]
		key.Vtype = v3.Float64T
	}
	return key
}

func encodeString(s string) v3.KeyValue {
	return v3.KeyValue{
		Value: v3.StringToBytes(s),
		Vtype: v3.StringT,
	}
}
func (enc *V3Encoder) encodeString(s string, k v3.KeyValue) error {
	value := v3.Value{
		Key: k,
	}
	value.Value.Value = v3.StringToBytes(s)
	value.Value.Vtype = v3.StringT
	v3.AddValue(enc.out, &value, enc.varintbuf)
	return nil
}

func encodeBytes(b []byte) v3.KeyValue {
	return v3.KeyValue{
		Value: b,
		Vtype: v3.BytesT,
	}
}

func (enc *V3Encoder) encodeValuev3(d interface{}, k v3.KeyValue) error {
	value := v3.Value{
		Key: k,
	}
	alreadyEncoded := false
	//this sets value.Value, it does al the basic types and some more
	switch v := d.(type) {
	case string:
		value.Value.Value = v3.StringToBytes(v)
		value.Value.Vtype = v3.StringT
	case []byte:
		value.Value.Value = v
		value.Value.Vtype = v3.BytesT
	case int8:
		value.Value.Value = []byte{v3.Int8ToBytes(v)}
		value.Value.Vtype = v3.Int8T
	case int16:
		value.Value.Value = v3.Int16ToBytes(v)
		value.Value.Vtype = v3.Int16T
	case int32:
		value.Value.Value = v3.Int32ToBytes(v)
		value.Value.Vtype = v3.Int32T
	case int64:
		value.Value.Value = v3.Int64ToBytes(v)
		value.Value.Vtype = v3.Int64T
	case int:
		value.Value.Value = v3.Int64ToBytes(int64(v))
		value.Value.Vtype = v3.Int64T
	case uint8:
		value.Value.Value = []byte{v3.Uint8ToBytes(v)}
		value.Value.Vtype = v3.Uint8T
	case uint16:
		value.Value.Value = v3.Uint16ToBytes(v)
		value.Value.Vtype = v3.Uint16T
	case uint32:
		value.Value.Value = v3.Uint32ToBytes(v)
		value.Value.Vtype = v3.Uint32T
	case uint64:
		value.Value.Value = v3.Uint64ToBytes(v)
		value.Value.Vtype = v3.Uint64T
	case uint:
		value.Value.Value = v3.Uint64ToBytes(uint64(v))
		value.Value.Vtype = v3.Uint64T
	case float32:
		value.Value.Value = v3.Float32ToBytes(&v)
		value.Value.Vtype = v3.Float32T
	case float64:
		var buf [8]byte
		v3.Float64ToBytes(v, &buf)
		value.Value.Value = buf[:]
		value.Value.Vtype = v3.Float64T
	case bool:
		value.Value.Value = v3.BoolToBytes(v)
		value.Value.Vtype = v3.BoolT
	default:
		val := reflect.ValueOf(d)
		kind := val.Kind()
		if kind == reflect.Map {
			value.Childrenn = uint64(val.Len())
			value.Value.Vtype = v3.MapT
			alreadyEncoded = true
			v3.AddValue(enc.out, &value, enc.varintbuf)

			switch v := d.(type) {
			case map[string]string:
				for k, v := range v {
					enc.encodeString(v, encodeString(k))
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
					enc.encodeValuev3_reflect(mapRange.Value(), encodeKeyv3_reflect(mapRange.Key()))
				}
			}
		} else if kind == reflect.Array || kind == reflect.Slice {
			value.Childrenn = uint64(val.Len())
			value.Value.Vtype = v3.ArrT
			alreadyEncoded = true
			v3.AddValue(enc.out, &value, enc.varintbuf)
			switch s := d.(type) {
			case []string:
				for _, v := range s {
					enc.encodeString(v, v3.KeyValue{})
				}
			default:
				//if its not a specific slice type
				for i := 0; i < int(value.Childrenn); i++ {
					err := enc.encodeValuev3_reflect(val.Index(i), v3.KeyValue{})
					if err != nil {
						return err
					}
				}
			}
		} else if v, ok := d.(gob.GobEncoder); ok {
			var err error
			value.Value.Value, err = v.GobEncode()
			if err != nil {
				return err
			}
			value.Value.Vtype = v3.BytesT
		} else if kind == reflect.Struct {
			name := val.Type().String()
			var usableFields map[string]int
			if v, ok := enc.typeCache[name]; ok {
				usableFields = v
			} else {
				usableFields = getStructFields(val)
				enc.typeCache[name] = usableFields
			}

			value.Childrenn = uint64(len(usableFields))
			value.Value.Vtype = v3.MapT
			alreadyEncoded = true
			v3.AddValue(enc.out, &value, enc.varintbuf)
			for fieldName, fieldID := range usableFields {
				field := val.Field(fieldID)
				err := enc.encodeValuev3_reflect(field, encodeBytes([]byte(fieldName)))
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

	return nil
}

func (enc *V3Encoder) encodeValuev3_reflect(d reflect.Value, k v3.KeyValue) error {
	value := v3.Value{
		Key: k,
	}
	alreadyEncoded := false
	//this sets value.Value, it does al the basic types and some more
	switch d.Type().Kind() {
	case reflect.Interface, reflect.Ptr:
		enc.encodeValuev3_reflect(d.Elem(), k)
		alreadyEncoded = true

	case reflect.String:
		value.Value.Value = v3.StringToBytes(d.String())
		value.Value.Vtype = v3.StringT
	case reflect.Slice:
		if d.Type().Elem().Kind() == reflect.Int8 {
			value.Value.Value = d.Bytes()
			value.Value.Vtype = v3.StringT
		} else if d.Type().Elem().Kind() == reflect.Int8 {
			value.Childrenn = uint64(d.Len())
			value.Value.Vtype = v3.ArrT
			i := d.Interface()
			for _, v := range i.([]string) {
				enc.encodeString(v, v3.KeyValue{})
			}
			//If its not a specific type
		} else {
			value.Childrenn = uint64(d.Len())
			value.Value.Vtype = v3.ArrT
			alreadyEncoded = true
			v3.AddValue(enc.out, &value, enc.varintbuf)
			for i := 0; i < int(value.Childrenn); i++ {
				err := enc.encodeValuev3_reflect(d.Index(i), v3.KeyValue{})
				if err != nil {
					return err
				}
			}
		}
	case reflect.Int8:
		value.Value.Value = []byte{v3.Int8ToBytes(int8(d.Int()))}
		value.Value.Vtype = v3.Int8T
	case reflect.Int16:
		value.Value.Value = v3.Int16ToBytes(int16(d.Int()))
		value.Value.Vtype = v3.Int16T
	case reflect.Int32:
		value.Value.Value = v3.Int32ToBytes(int32(d.Int()))
		value.Value.Vtype = v3.Int32T
	case reflect.Int64:
		value.Value.Value = v3.Int64ToBytes(d.Int())
		value.Value.Vtype = v3.Int64T
	case reflect.Int:
		value.Value.Value = v3.Int64ToBytes(d.Int())
		value.Value.Vtype = v3.Int64T
	case reflect.Uint8:
		value.Value.Value = []byte{v3.Uint8ToBytes(uint8(d.Uint()))}
		value.Value.Vtype = v3.Uint8T
	case reflect.Uint16:
		value.Value.Value = v3.Uint16ToBytes(uint16(d.Uint()))
		value.Value.Vtype = v3.Uint16T
	case reflect.Uint32:
		value.Value.Value = v3.Uint32ToBytes(uint32(d.Uint()))
		value.Value.Vtype = v3.Uint32T
	case reflect.Uint64:
		value.Value.Value = v3.Uint64ToBytes(d.Uint())
		value.Value.Vtype = v3.Uint64T
	case reflect.Uint:
		value.Value.Value = v3.Uint64ToBytes(d.Uint())
		value.Value.Vtype = v3.Uint64T
	case reflect.Bool:
		value.Value.Value = v3.BoolToBytes(d.Bool())
		value.Value.Vtype = v3.BoolT
	case reflect.Float32:
		v := float32(d.Float())
		value.Value.Value = v3.Float32ToBytes(&v)
		value.Value.Vtype = v3.Float32T
	case reflect.Float64:
		v := d.Float()
		var buf [8]byte
		v3.Float64ToBytes(v, &buf)
		value.Value.Value = buf[:]
		value.Value.Vtype = v3.Float64T
	case reflect.Map:
		//TODO Only check the types first and only then convert??
		i := d.Interface()
		value.Childrenn = uint64(d.Len())
		value.Value.Vtype = v3.MapT
		alreadyEncoded = true
		v3.AddValue(enc.out, &value, enc.varintbuf)

		switch v := i.(type) {
		case map[string]string:
			for k, v := range v {
				enc.encodeString(v, encodeString(k))
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
			mapRange := d.MapRange()
			for mapRange.Next() {
				enc.encodeValuev3_reflect(mapRange.Value(), encodeKeyv3_reflect(mapRange.Key()))
			}
		}
	case reflect.Array:
		//TODO Only check the types first and only then convert??
		i := d.Interface()

		value.Childrenn = uint64(d.Len())
		value.Value.Vtype = v3.ArrT
		alreadyEncoded = true
		v3.AddValue(enc.out, &value, enc.varintbuf)
		switch s := i.(type) {
		case []string:
			for _, v := range s {
				enc.encodeString(v, v3.KeyValue{})
			}
		default:
			//if its not a specific slice type
			for i := 0; i < int(value.Childrenn); i++ {
				err := enc.encodeValuev3_reflect(d.Index(i), v3.KeyValue{})
				if err != nil {
					return err
				}
			}
		}
	case reflect.Struct:
		modelType := reflect.TypeOf((*gob.GobEncoder)(nil)).Elem()
		if d.Type().Implements(modelType) {
			var err error
			value.Value.Value, err = d.Interface().(gob.GobEncoder).GobEncode()
			if err != nil {
				return err
			}
			value.Value.Vtype = v3.BytesT
		} else {
			name := d.Type().String()
			var usableFields map[string]int
			if v, ok := enc.typeCache[name]; ok {
				usableFields = v
			} else {
				usableFields = getStructFields(d)
				enc.typeCache[name] = usableFields
			}

			value.Childrenn = uint64(len(usableFields))
			value.Value.Vtype = v3.MapT
			alreadyEncoded = true
			v3.AddValue(enc.out, &value, enc.varintbuf)
			for fieldName, fieldID := range usableFields {
				field := d.Field(fieldID)
				err := enc.encodeValuev3_reflect(field, encodeBytes([]byte(fieldName)))
				if err != nil {
					return err
				}
			}
		}
	}

	if !alreadyEncoded {
		v3.AddValue(enc.out, &value, enc.varintbuf)
	}

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
	typeCache map[string]map[string]int
	yetToRead uint64
	sync.Mutex
}

//NewV3Decoder creates aa new V3Decoder to decode a ttv3 data stream.
//The init flag specifies wether it should initialize the decoder.
//Initializing the decoder blocks until at least the first 2 bytes are read.
func NewV3Decoder(in v3.Reader, init bool) *V3Decoder {
	dec := V3Decoder{
		didInit:   !init,
		in:        in,
		typeCache: map[string]map[string]int{},
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
	var err error
	if !dec.isStream {
		if dec.didDecode {
			return errors.New("TT: attempt to decode twice from non-stream")
		}
		dec.didDecode = true
	}
	if e == nil {
		var v v3.Value
		err = v.FromBytes(dec.in)
		if err != nil {
			return err
		}
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
		err = v.FromBytes(dec.in)
		if err != nil {
			return err
		}

		clearNextValues(dec.in, v.Childrenn)
		return nil
	}

	var v v3.Value
	err = v.FromBytes(dec.in)
	if err != nil {
		return err
	}

	dec.yetToRead = v.Childrenn
	err = dec.decodeValuev3(v, value)
	if dec.yetToRead != 0 {
		clearNextValues(dec.in, dec.yetToRead)
	}

	return err
}

func decodeUndefined(data v3.KeyValue, e reflect.Value) error {
	return errors.New("TT: cannot unmarshal invalid type:" + strconv.Itoa(int(data.Vtype)))
}

func decodeString(data v3.KeyValue, e reflect.Value) error {
	val := v3.StringFromBytes(data.Value)
	if e.Kind() != reflect.String {
		if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
			return errors.New("TT: cannot unmarshal string into " + e.Type().String() + " Go type")
		}
		e.Set(reflect.ValueOf(val))
	} else {
		e.SetString(val)
	}
	return nil
}

func decodeBytes(data v3.KeyValue, e reflect.Value) error {
	if e.Kind() != reflect.Slice || e.Type().Elem().Kind() != reflect.Uint8 {
		if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
			return errors.New("TT: cannot unmarshal bytes into " + e.Type().String() + " Go type")
		}
		e.Set(reflect.ValueOf(data.Value))
	} else {
		e.SetBytes(data.Value)
	}
	return nil
}

func decodeInt8(data v3.KeyValue, e reflect.Value) error {
	val := v3.Int8FromBytes(data.Value[0])
	if e.Kind() != reflect.Int8 {
		if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
			return errors.New("TT: cannot unmarshal int8 into " + e.Kind().String() + " Go type")
		}
		e.Set(reflect.ValueOf(val))
	} else {
		e.SetInt(int64(val))
	}
	return nil
}

func decodeInt16(data v3.KeyValue, e reflect.Value) error {
	val := v3.Int16FromBytes(data.Value)
	if e.Kind() != reflect.Int16 {
		if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
			return errors.New("TT: cannot unmarshal int16 into " + e.Kind().String() + " Go type")
		}
		e.Set(reflect.ValueOf(val))
	} else {
		e.SetInt(int64(val))
	}
	return nil
}

func decodeInt32(data v3.KeyValue, e reflect.Value) error {
	val := v3.Int32FromBytes(data.Value)
	if e.Kind() != reflect.Int32 {
		if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
			return errors.New("TT: cannot unmarshal int32 into " + e.Kind().String() + " Go type")
		}
		e.Set(reflect.ValueOf(val))
	} else {
		e.SetInt(int64(val))
	}
	return nil
}

func decodeInt64(data v3.KeyValue, e reflect.Value) error {
	val := v3.Int64FromBytes(data.Value)
	if e.Kind() != reflect.Int64 && e.Kind() != reflect.Int {
		if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
			return errors.New("TT: cannot unmarshal int64 into " + e.Kind().String() + " Go type")
		}
		e.Set(reflect.ValueOf(val))
	} else {
		e.SetInt(val)
	}
	return nil
}

func decodeUint8(data v3.KeyValue, e reflect.Value) error {
	val := v3.Uint8FromBytes(data.Value[0])
	if e.Kind() != reflect.Uint8 {
		if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
			return errors.New("TT: cannot unmarshal uint8 into " + e.Kind().String() + " Go type")
		}
		e.Set(reflect.ValueOf(val))
	} else {
		e.SetUint(uint64(val))
	}
	return nil
}

func decodeUint16(data v3.KeyValue, e reflect.Value) error {
	val := v3.Uint16FromBytes(data.Value)
	if e.Kind() != reflect.Uint16 {
		if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
			return errors.New("TT: cannot unmarshal uint16 into " + e.Kind().String() + " Go type")
		}
		e.Set(reflect.ValueOf(val))
	} else {
		e.SetUint(uint64(val))
	}
	return nil
}

func decodeUint32(data v3.KeyValue, e reflect.Value) error {
	val := v3.Uint32FromBytes(data.Value)
	if e.Kind() != reflect.Uint32 {
		if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
			return errors.New("TT: cannot unmarshal uint32 into " + e.Kind().String() + " Go type")
		}
		e.Set(reflect.ValueOf(val))
	} else {
		e.SetUint(uint64(val))
	}
	return nil
}

func decodeUint64(data v3.KeyValue, e reflect.Value) error {
	val := v3.Uint64FromBytes(data.Value)
	if e.Kind() != reflect.Uint64 && e.Kind() != reflect.Uint {
		if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
			return errors.New("TT: cannot unmarshal uint64 into " + e.Kind().String() + " Go type")
		}
		e.Set(reflect.ValueOf(val))
	} else {
		e.SetUint(val)
	}
	return nil
}

func decodeBool(data v3.KeyValue, e reflect.Value) error {
	val := v3.BoolFromBytes(data.Value)
	if e.Kind() != reflect.Bool {
		if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
			return errors.New("TT: cannot unmarshal bytes into " + e.Kind().String() + " Go type")
		}
		e.Set(reflect.ValueOf(val))
	} else {
		e.SetBool(val)
	}
	return nil
}

func decodeFloat32(data v3.KeyValue, e reflect.Value) error {
	val := v3.Float32FromBytes(data.Value)
	if e.Kind() != reflect.Float32 {
		if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
			return errors.New("TT: cannot unmarshal float32 into " + e.Kind().String() + " Go type")
		}
		e.Set(reflect.ValueOf(val))
	} else {
		e.SetFloat(float64(val))
	}
	return nil
}

func decodeFloat64(data v3.KeyValue, e reflect.Value) error {
	val := v3.Float64FromBytes(data.Value)
	if e.Kind() != reflect.Float64 {
		if e.Kind() != reflect.Interface || e.Type().NumMethod() != 0 {
			return errors.New("TT: cannot unmarshal float64 into " + e.Kind().String() + " Go type")
		}
		e.Set(reflect.ValueOf(val))
	} else {
		e.SetFloat(val)
	}
	return nil
}

func decodeMap(dec *V3Decoder, v v3.Value, e reflect.Value) error {
	if e.Kind() == reflect.Interface && e.NumMethod() == 0 {
		children := v.Childrenn
		m := make(map[interface{}]interface{}, children)
		var err error
		key := reflect.New(reflect.TypeOf(m).Key()).Elem()
		for i := uint64(0); i < children; i++ {
			err = v.FromBytes(dec.in)
			if err != nil {
				return err
			}
			dec.yetToRead += v.Childrenn - 1

			err = decodeKeyv3(v.Key, key)
			if err != nil {
				return err
			}
			k := key.Interface()
			err = dec.decodeValuev3(v, key)
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

	} else if e.Kind() == reflect.Map {
		children := v.Childrenn
		if e.IsNil() {
			e.Set(reflect.MakeMap(e.Type()))
		}

		var err error
		var value reflect.Value
		key := reflect.New(e.Type().Key()).Elem()

		elem := e.Type().Elem()
		ValueKind := elem.Kind()
		shouldReplace := ValueKind == reflect.Array || ValueKind == reflect.Slice || ValueKind == reflect.Map

		if !shouldReplace {
			value = reflect.New(elem).Elem()
		}
		for i := uint64(0); i < children; i++ {
			err = v.FromBytes(dec.in)
			if err != nil {
				return err
			}
			dec.yetToRead += v.Childrenn - 1

			err = decodeKeyv3(v.Key, key)
			if err != nil {
				return err
			}
			if shouldReplace {
				value = reflect.New(elem).Elem()
			}
			err = dec.decodeValuev3(v, value)
			if err != nil {
				return err
			}
			e.SetMapIndex(key, value)
		}
	} else if e.Kind() == reflect.Struct {
		children := v.Childrenn
		name := e.Type().String()
		var usableFields map[string]int
		if v, ok := dec.typeCache[name]; ok {
			usableFields = v
		} else {
			usableFields = getStructFields(e)
			dec.typeCache[name] = usableFields
		}

		for i := uint64(0); i < children; i++ {
			err := v.FromBytes(dec.in)
			if err != nil {
				return err
			}

			dec.yetToRead += v.Childrenn - 1

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

			err = dec.decodeValuev3(v, field)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func decodeArr(dec *V3Decoder, v v3.Value, e reflect.Value) error {
	children := v.Childrenn

	if e.Kind() == reflect.Array {
		len := e.Len()
		if len < int(children) {
			return nil
		}
		for i := 0; i < int(children); i++ {
			err := v.FromBytes(dec.in)
			if err != nil {
				return err
			}
			dec.yetToRead += v.Childrenn - 1

			err = dec.decodeValuev3(v, e.Index(i))
			if err != nil {
				return err
			}
		}
	} else if e.Kind() == reflect.Slice {
		len := e.Len()
		if len < int(children) {
			e.Set(reflect.MakeSlice(e.Type(), int(children), int(children)))
		} else if len > int(children) {
			e.SetLen(int(children))
		}
		for i := 0; i < int(children); i++ {
			err := v.FromBytes(dec.in)
			if err != nil {
				return err
			}
			dec.yetToRead += v.Childrenn - 1

			err = dec.decodeValuev3(v, e.Index(i))
			if err != nil {
				return err
			}
		}
		/*} else if e.Kind() == reflect.Map {
		//TODO: if kind is also a numeric value it is usable
		e.Type().Key().Kind()
		*/
	} else {
		//if all special cases fail we fall back to []interface{}
		arr := make([]interface{}, children)
		var err error
		var value reflect.Value
		valueElem := reflect.TypeOf(arr).Elem()
		ValueKind := valueElem.Kind()
		shouldReplace := ValueKind == reflect.Array || ValueKind == reflect.Slice || ValueKind == reflect.Map

		if !shouldReplace {
			value = reflect.New(valueElem).Elem()
		}
		for i := 0; i < int(children); i++ {
			err = v.FromBytes(dec.in)
			if err != nil {
				return err
			}
			dec.yetToRead += v.Childrenn - 1

			if shouldReplace {
				value = reflect.New(valueElem).Elem()
			}

			err = dec.decodeValuev3(v, value)
			if err != nil {
				return err
			}
			arr[i] = value.Interface()
		}
		e.Set(reflect.ValueOf(arr))
	}
	return nil
}

func decodeKeyv3(k v3.KeyValue, e reflect.Value) error {
	return decodersSlice[int(k.Vtype)](k, e)
}

func (dec *V3Decoder) decodeValuev3(v v3.Value, e reflect.Value) error {
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
				return u.GobDecode(v.Value.Value)
			}
		}

		if haveAddr {
			e = e0 // restore original value after round-trip Value.Addr().Elem()
			haveAddr = false
		} else {
			e = e.Elem()
		}
	}
	if int(v.Value.Vtype) < 18 {
		return decodersSlice[int(v.Value.Vtype)](v.Value, e)
	} else if v.Value.Vtype == v3.MapT {
		return decodeMap(dec, v, e)
	} else if v.Value.Vtype == v3.ArrT {
		return decodeArr(dec, v, e)
	}
	return decodeUndefined(v.Value, e)
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

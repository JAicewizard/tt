package tt

import (
	"encoding/binary"
	"encoding/gob"
	"errors"
	"io"
	"math"
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

type fieldsstruct struct {
	canGob bool
	usable uint
	fields []string
}

//V3Encoder is the encoder used to encode a ttv3 data stream
type V3Encoder struct {
	out        io.Writer
	varintbuf  *[2*binary.MaxVarintLen64 + 1]byte
	typeCache  map[reflect.Type]fieldsstruct
	isStream   bool
	fsKeyBuf   [8]byte //Fixed Size buffer for keys
	fsValueBuf [8]byte //Fixed Size buffer for values
	sync.Mutex
}

var v3StreamHeader = []byte{version3, 1 << 7}
var v3NoStreamHeader = []byte{version3, 0}

//NewV3Encoder creates a new encoder to encode a ttv3 data stream
func NewV3Encoder(out io.Writer, isStream bool) *V3Encoder {
	if isStream {
		out.Write(v3StreamHeader)
	}
	return &V3Encoder{
		isStream:  isStream,
		out:       out,
		varintbuf: &[2*binary.MaxVarintLen64 + 1]byte{},
		typeCache: map[reflect.Type]fieldsstruct{},
	}
}

//Encodev3 encodes an `interface{}`` into a bytebuffer using ttv3
func Encodev3(d interface{}, out io.Writer) error {
	out.Write(v3NoStreamHeader)

	enc := &V3Encoder{
		out:       out,
		varintbuf: &[2*binary.MaxVarintLen64 + 1]byte{},
		typeCache: map[reflect.Type]fieldsstruct{},
	}
	//We dont have to lock/unlock since we know we are the only one witha acces
	return enc.encodeValuev3(d, v3.KeyValue{})
}

//Encode encodes an `interface{}` into a bytebuffer using ttv3
func (enc *V3Encoder) Encode(d interface{}) error {
	enc.Lock()
	if !enc.isStream {
		enc.out.Write(v3NoStreamHeader)
	}
	ret := enc.encodeValuev3(d, v3.KeyValue{})
	enc.Unlock()
	return ret
}

func encodeKeyv3(k interface{}, fsKeyBuf *[8]byte) v3.KeyValue {
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
			key.Value = fsKeyBuf[:1]
			fsKeyBuf[0] = v3.Int8ToBytes(v)
			key.Vtype = v3.Int8T
		case int16:
			key.Value = fsKeyBuf[:2]
			v3.Int16ToBytes(v, key.Value)
			key.Vtype = v3.Int16T
		case int32:
			key.Value = fsKeyBuf[:4]
			v3.Int32ToBytes(v, key.Value)
			key.Vtype = v3.Int32T
		case int64:
			key.Value = fsKeyBuf[:8]
			v3.Int64ToBytes(v, key.Value)
			key.Vtype = v3.Int64T
		case int:
			key.Value = fsKeyBuf[:8]
			v3.Int64ToBytes(int64(v), key.Value)
			key.Vtype = v3.Int64T
		case uint8:
			key.Value = fsKeyBuf[:1]
			fsKeyBuf[0] = v3.Uint8ToBytes(v)
			key.Vtype = v3.Uint8T
		case uint16:
			key.Value = fsKeyBuf[:2]
			v3.Uint16ToBytes(v, key.Value)
			key.Vtype = v3.Uint16T
		case uint32:
			key.Value = fsKeyBuf[:4]
			v3.Uint32ToBytes(v, key.Value)
			key.Vtype = v3.Uint32T
		case uint64:
			key.Value = fsKeyBuf[:8]
			v3.Uint64ToBytes(v, key.Value)
			key.Vtype = v3.Uint64T
		case uint:
			key.Value = fsKeyBuf[:8]
			v3.Uint64ToBytes(uint64(v), key.Value)
			key.Vtype = v3.Uint64T
		case float32:
			key.Value = fsKeyBuf[:4]
			v3.Float32ToBytes(v, key.Value)
			key.Vtype = v3.Float32T
		case float64:
			key.Value = fsKeyBuf[:8]
			v3.Float64ToBytes(v, key.Value)
			key.Vtype = v3.Float64T
		case bool:
			key.Value = fsKeyBuf[:1]
			fsKeyBuf[0] = v3.BoolToBytes(v)
			key.Vtype = v3.BoolT
		}
	}
	return key
}

func encodeKeyv3_reflect(d reflect.Value, fsKeyBuf *[8]byte) v3.KeyValue {
	var key v3.KeyValue
	switch d.Type().Kind() {
	case reflect.Interface, reflect.Ptr:
		encodeKeyv3_reflect(d.Elem(), fsKeyBuf)
	case reflect.String:
		key.Value = v3.StringToBytes(d.String())
		key.Vtype = v3.StringT
	case reflect.Slice:
		if d.Type().Elem().Kind() == reflect.Int8 {
			key.Value = d.Bytes()
			key.Vtype = v3.StringT
		}
	case reflect.Int8:
		key.Value = fsKeyBuf[:1]
		fsKeyBuf[0] = v3.Int8ToBytes(int8(d.Int()))
		key.Vtype = v3.Int8T
	case reflect.Int16:
		key.Value = fsKeyBuf[:2]
		v3.Int16ToBytes(int16(d.Int()), key.Value)
		key.Vtype = v3.Int16T
	case reflect.Int32:
		key.Value = fsKeyBuf[:4]
		v3.Int32ToBytes(int32(d.Int()), key.Value)
		key.Vtype = v3.Int32T
	case reflect.Int64, reflect.Int:
		key.Value = fsKeyBuf[:8]
		v3.Int64ToBytes(d.Int(), key.Value)
		key.Vtype = v3.Int64T
	case reflect.Uint8:
		key.Value = fsKeyBuf[:1]
		fsKeyBuf[0] = v3.Uint8ToBytes(uint8(d.Uint()))
		key.Vtype = v3.Uint8T
	case reflect.Uint16:
		key.Value = fsKeyBuf[:2]
		v3.Uint16ToBytes(uint16(d.Uint()), key.Value)
		key.Vtype = v3.Uint16T
	case reflect.Uint32:
		key.Value = fsKeyBuf[:4]
		v3.Uint32ToBytes(uint32(d.Uint()), key.Value)
		key.Vtype = v3.Uint32T
	case reflect.Uint64, reflect.Uint:
		key.Value = fsKeyBuf[:8]
		v3.Uint64ToBytes(d.Uint(), key.Value)
		key.Vtype = v3.Uint64T
	case reflect.Bool:
		key.Value = fsKeyBuf[:1]
		fsKeyBuf[0] = v3.BoolToBytes(d.Bool())
		key.Vtype = v3.BoolT
	case reflect.Float32:
		key.Value = fsKeyBuf[:4]
		v3.Float32ToBytes(float32(d.Float()), key.Value)
		key.Vtype = v3.Float32T
	case reflect.Float64:
		key.Value = fsKeyBuf[:8]
		v3.Float64ToBytes(d.Float(), key.Value)
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
		Value: v3.KeyValue{
			Value: v3.StringToBytes(s),
			Vtype: v3.StringT,
		},
	}
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
		value.Value.Value = enc.fsValueBuf[:1]
		enc.fsValueBuf[0] = v3.Int8ToBytes(v)
		value.Value.Vtype = v3.Int8T
	case int16:
		value.Value.Value = enc.fsValueBuf[:2]
		v3.Int16ToBytes(v, value.Value.Value)
		value.Value.Vtype = v3.Int16T
	case int32:
		value.Value.Value = enc.fsValueBuf[:4]
		v3.Int32ToBytes(v, value.Value.Value)
		value.Value.Vtype = v3.Int32T
	case int64:
		value.Value.Value = enc.fsValueBuf[:8]
		v3.Int64ToBytes(v, value.Value.Value)
		value.Value.Vtype = v3.Int64T
	case int:
		value.Value.Value = enc.fsValueBuf[:8]
		v3.Int64ToBytes(int64(v), value.Value.Value)
		value.Value.Vtype = v3.Int64T
	case uint8:
		value.Value.Value = enc.fsValueBuf[:1]
		enc.fsValueBuf[0] = v3.Uint8ToBytes(v)
		value.Value.Vtype = v3.Uint8T
	case uint16:
		value.Value.Value = enc.fsValueBuf[:2]
		v3.Uint16ToBytes(v, value.Value.Value)
		value.Value.Vtype = v3.Uint16T
	case uint32:
		value.Value.Value = enc.fsValueBuf[:4]
		v3.Uint32ToBytes(v, value.Value.Value)
		value.Value.Vtype = v3.Uint32T
	case uint64:
		value.Value.Value = enc.fsValueBuf[:8]
		v3.Uint64ToBytes(uint64(v), value.Value.Value)
		value.Value.Vtype = v3.Uint64T
	case uint:
		value.Value.Value = enc.fsValueBuf[:8]
		v3.Uint64ToBytes(uint64(v), value.Value.Value)
		value.Value.Vtype = v3.Uint64T
	case float32:
		value.Value.Value = enc.fsValueBuf[:4]
		v3.Float32ToBytes(v, value.Value.Value)
		value.Value.Vtype = v3.Float32T
	case float64:
		value.Value.Value = enc.fsValueBuf[:8]
		v3.Float64ToBytes(v, value.Value.Value)
		value.Value.Vtype = v3.Float64T
	case bool:
		value.Value.Value = enc.fsValueBuf[:1]
		enc.fsValueBuf[0] = v3.BoolToBytes(v)
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
					enc.encodeValuev3(v, encodeKeyv3(k, &enc.fsKeyBuf))
				}
			default:
				//if its not a specific map type
				mapRange := val.MapRange()
				for mapRange.Next() {
					enc.encodeValuev3_reflect(mapRange.Value(), encodeKeyv3_reflect(mapRange.Key(), &enc.fsKeyBuf))
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
			typ := val.Type()
			var usableFields fieldsstruct
			if v, ok := enc.typeCache[typ]; ok {
				usableFields = v
			} else {
				usableFields = getStructFields2(val)
				usableFields.canGob = false
				enc.typeCache[typ] = usableFields
			}

			value.Childrenn = uint64(usableFields.usable)
			value.Value.Vtype = v3.MapT
			alreadyEncoded = true
			v3.AddValue(enc.out, &value, enc.varintbuf)
			for fieldID, fieldName := range usableFields.fields {
				if fieldName == "" {
					continue
				}
				field := val.Field(fieldID)
				err := enc.encodeValuev3_reflect(field, encodeString(fieldName))
				if err != nil {
					return err
				}
			}
		} else if kind == reflect.Interface || kind == reflect.Ptr {
			enc.encodeValuev3_reflect(val.Elem(), k)
			alreadyEncoded = true
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
	typ := d.Type()
	//this sets value.Value, it does al the basic types and some more
	switch typ.Kind() {
	case reflect.Interface, reflect.Ptr:
		enc.encodeValuev3_reflect(d.Elem(), k)
		alreadyEncoded = true
	case reflect.String:
		value.Value.Value = v3.StringToBytes(d.String())
		value.Value.Vtype = v3.StringT
	case reflect.Slice:
		if typ.Elem().Kind() == reflect.Int8 {
			value.Value.Value = d.Bytes()
			value.Value.Vtype = v3.StringT
		} else if typ.Elem().Kind() == reflect.String {
			value.Childrenn = uint64(d.Len())
			value.Value.Vtype = v3.ArrT
			for i := 0; i < int(value.Childrenn); i++ {
				enc.encodeString(d.Index(i).String(), v3.KeyValue{})
			}
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
		value.Value.Value = enc.fsValueBuf[:1]
		enc.fsValueBuf[0] = v3.Int8ToBytes(int8(d.Int()))
		value.Value.Vtype = v3.Int8T
	case reflect.Int16:
		value.Value.Value = enc.fsValueBuf[:2]
		v3.Int16ToBytes(int16(d.Int()), value.Value.Value)
		value.Value.Vtype = v3.Int16T
	case reflect.Int32:
		value.Value.Value = enc.fsValueBuf[:4]
		v3.Int32ToBytes(int32(d.Int()), value.Value.Value)
		value.Value.Vtype = v3.Int32T
	case reflect.Int64, reflect.Int:
		value.Value.Value = enc.fsValueBuf[:8]
		v3.Int64ToBytes(d.Int(), value.Value.Value)
		value.Value.Vtype = v3.Int64T
	case reflect.Uint8:
		value.Value.Value = enc.fsValueBuf[:1]
		enc.fsValueBuf[0] = v3.Uint8ToBytes(uint8(d.Uint()))
		value.Value.Vtype = v3.Uint8T
	case reflect.Uint16:
		value.Value.Value = enc.fsValueBuf[:2]
		v3.Uint16ToBytes(uint16(d.Uint()), value.Value.Value)
		value.Value.Vtype = v3.Uint16T
	case reflect.Uint32:
		value.Value.Value = enc.fsValueBuf[:4]
		v3.Uint32ToBytes(uint32(d.Uint()), value.Value.Value)
		value.Value.Vtype = v3.Uint32T
	case reflect.Uint64, reflect.Uint:
		value.Value.Value = enc.fsValueBuf[:8]
		v3.Uint64ToBytes(d.Uint(), value.Value.Value)
		value.Value.Vtype = v3.Uint64T
	case reflect.Bool:
		value.Value.Value = enc.fsValueBuf[:1]
		enc.fsValueBuf[0] = v3.BoolToBytes(d.Bool())
		value.Value.Vtype = v3.BoolT
	case reflect.Float32:
		value.Value.Value = enc.fsValueBuf[:4]
		v3.Float32ToBytes(float32(d.Float()), value.Value.Value)
		value.Value.Vtype = v3.Float32T
	case reflect.Float64:
		value.Value.Value = enc.fsValueBuf[:8]
		v3.Float64ToBytes(d.Float(), value.Value.Value)
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
				enc.encodeValuev3(v, encodeKeyv3(k, &enc.fsKeyBuf))
			}
		default:
			//if its not a specific map type
			mapRange := d.MapRange()
			for mapRange.Next() {
				enc.encodeValuev3_reflect(mapRange.Value(), encodeKeyv3_reflect(mapRange.Key(), &enc.fsKeyBuf))
			}
		}
	case reflect.Array:
		value.Childrenn = uint64(d.Len())
		value.Value.Vtype = v3.ArrT
		alreadyEncoded = true
		v3.AddValue(enc.out, &value, enc.varintbuf)
		//if its not a specific slice type
		for i := 0; i < int(value.Childrenn); i++ {
			err := enc.encodeValuev3_reflect(d.Index(i), v3.KeyValue{})
			if err != nil {
				return err
			}
		}
	case reflect.Struct:
		var usableFields fieldsstruct
		if v, ok := enc.typeCache[typ]; ok {
			usableFields = v
		} else {
			usableFields = getStructFields2(d)
			usableFields.canGob = canGobEncode(d)
			enc.typeCache[typ] = usableFields
		}

		if usableFields.canGob {
			if i, ok := d.Interface().(gob.GobEncoder); ok {
				var err error
				value.Value.Value, err = i.GobEncode()
				if err != nil {
					return err
				}
				value.Value.Vtype = v3.BytesT
			}
			break
		}

		value.Childrenn = uint64(usableFields.usable)
		value.Value.Vtype = v3.MapT
		alreadyEncoded = true
		v3.AddValue(enc.out, &value, enc.varintbuf)
		for fieldID, fieldName := range usableFields.fields {
			if fieldName == "" {
				continue
			}
			field := d.Field(fieldID)
			err := enc.encodeValuev3_reflect(field, encodeString(fieldName))
			if err != nil {
				return err
			}
		}
	}

	if !alreadyEncoded {
		v3.AddValue(enc.out, &value, enc.varintbuf)
	}

	return nil
}

func getStructFields2(val reflect.Value) fieldsstruct {
	fields := fieldsstruct{}
	fields.fields = make([]string, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if field.PkgPath != "" || !val.Field(i).CanInterface() {
			continue
		}
		fields.fields[i] = getFieldName(field)
		fields.usable++
	}
	return fields
}

func getStructFieldsDec(val reflect.Value) map[string]int {
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
	didInit     bool
	isStream    bool
	didDecode   bool
	in          v3.Reader
	typeCache   map[string]map[string]int
	yetToRead   uint64
	allocLimmit uint64
	sync.Mutex
}

//NewV3Decoder creates aa new V3Decoder to decode a ttv3 data stream.
//The init flag specifies wether it should initialize the decoder.
//Initializing the decoder blocks until at least the first 2 bytes are read.
func NewV3Decoder(in v3.Reader, init bool) *V3Decoder {
	dec := V3Decoder{
		didInit:     !init,
		in:          in,
		typeCache:   map[string]map[string]int{},
		allocLimmit: math.MaxUint64,
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

//SetAllocLimmit sets the limmit of allocations. This does not induce
//a global limmit in tt but only for individual allocations
func (dec *V3Decoder) SetAllocLimmit(limit uint64) {
	dec.allocLimmit = limit
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
		err = v.FromBytes(dec.in, dec.allocLimmit)
		if err != nil {
			return err
		}

		return clearNextValues(dec.in, v.Childrenn, dec.allocLimmit)
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
		err = v.FromBytes(dec.in, dec.allocLimmit)
		if err != nil {
			return err
		}

		return clearNextValues(dec.in, v.Childrenn, dec.allocLimmit)
	}

	var v v3.Value
	err = v.FromBytes(dec.in, dec.allocLimmit)
	if err != nil {
		return err
	}

	dec.yetToRead = v.Childrenn
	err = dec.decodeValuev3(v, value)
	if dec.yetToRead != 0 {
		err2 := clearNextValues(dec.in, dec.yetToRead, dec.allocLimmit)
		if err == nil {
			return err2
		}
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
	if len(data.Value) < 1 {
		var buf [1]byte
		copy(buf[:], data.Value)
		data.Value = buf[:]
	}
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
	if len(data.Value) < 2 {
		var buf [2]byte
		copy(buf[:], data.Value)
		data.Value = buf[:]
	}
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
	if len(data.Value) < 4 {
		var buf [4]byte
		copy(buf[:], data.Value)
		data.Value = buf[:]
	}
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
	if len(data.Value) < 8 {
		var buf [8]byte
		copy(buf[:], data.Value)
		data.Value = buf[:]
	}
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
	if len(data.Value) < 1 {
		var buf [1]byte
		copy(buf[:], data.Value)
		data.Value = buf[:]
	}
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
	if len(data.Value) < 2 {
		var buf [2]byte
		copy(buf[:], data.Value)
		data.Value = buf[:]
	}
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
	if len(data.Value) < 4 {
		var buf [4]byte
		copy(buf[:], data.Value)
		data.Value = buf[:]
	}
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
	if len(data.Value) < 8 {
		var buf [8]byte
		copy(buf[:], data.Value)
		data.Value = buf[:]
	}
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
	val := false
	if len(data.Value) < 1 {
		val = v3.BoolFromBytes([]byte{0})
	} else {
		val = v3.BoolFromBytes(data.Value)
	}

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
	if len(data.Value) < 4 {
		var buf [4]byte
		copy(buf[:], data.Value)
		data.Value = buf[:]
	}
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
	if len(data.Value) < 8 {
		var buf [8]byte
		copy(buf[:], data.Value)
		data.Value = buf[:]
	}
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
			err = v.FromBytes(dec.in, dec.allocLimmit)
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
			err = v.FromBytes(dec.in, dec.allocLimmit)
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
			usableFields = getStructFieldsDec(e)
			dec.typeCache[name] = usableFields
		}

		for i := uint64(0); i < children; i++ {
			err := v.FromBytes(dec.in, dec.allocLimmit)
			if err != nil {
				return err
			}

			dec.yetToRead += v.Childrenn - 1

			key := v.Key.ExportStructID()
			if key == "" {
				err = clearNextValues(dec.in, v.Childrenn, dec.allocLimmit)
				if err != nil {
					return err
				}

				continue
			}
			fieldIndex, ok := usableFields[key]
			if !ok {
				err = clearNextValues(dec.in, v.Childrenn, dec.allocLimmit)
				if err != nil {
					return err
				}

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
			err := v.FromBytes(dec.in, dec.allocLimmit)
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
			err := v.FromBytes(dec.in, dec.allocLimmit)
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
			err = v.FromBytes(dec.in, dec.allocLimmit)
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
	if int(k.Vtype) > 18 {
		return errors.New("TT: cannot unmarshal invalid key type:" + strconv.Itoa(int(k.Vtype)))
	}
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

func clearNextValues(buf v3.Reader, values uint64, limit uint64) error {
	var value v3.Value
	for ; values > 0; values-- {
		err := value.FromBytes(buf, limit)
		if err != nil {
			return err
		}
		values += value.Childrenn
	}
	return nil
}

func canGobEncode(d reflect.Value) bool {
	if _, ok := d.Type().MethodByName("GobEncode"); ok {
		if _, ok := d.Interface().(gob.GobEncoder); ok {
			return true
		}
	}
	return false
}

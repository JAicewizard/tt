package v3

type (
	//KeyValue is the key used for storing the key/value
	KeyValue struct {
		Value []byte //the key to the data of this object in its final form
		Vtype ttType //the type of the data of this object in its final form
	}
)

//ExportStructID returns a string identifying the value
func (k *KeyValue) ExportStructID() string {
	if k.Vtype == StringT || k.Vtype == BytesT {
		return StringFromBytes(k.Value)
	}
	return ""
	/*switch k.Vtype {
	case StringT, BytesT:
		return StringFromBytes(k.Value)
	case Float32T:
		return "", errors.New("cannot load float32 as struct field name")
	case Float64T:
		return "", errors.New("cannot load float64 as struct field name")
	case Int8T:
		return "", errors.New("cannot load int8 as struct field name")
	case Int16T:
		return "", errors.New("cannot load int16 as struct field name")
	case Int32T:
		return "", errors.New("cannot load int32 as struct field name")
	case Int64T:
		return "", errors.New("cannot load int64 as struct field name")
	case Uint8T:
		return "", errors.New("cannot load uint8 as struct field name")
	case Uint16T:
		return "", errors.New("cannot load uint16 as struct field name")
	case Uint32T:
		return "", errors.New("cannot load uint32 as struct field name")
	case Uint64T:
		return "", errors.New("cannot load uint64 as struct field name")
	case BoolT:
		return "", errors.New("cannot load bool as struct field name")
	}
	return "", errors.New("cannot load invalid type as struct field name")*/
}

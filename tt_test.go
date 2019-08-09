package tt

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"testing"

	v1 "github.com/JAicewizard/tt/v1"
)

var testData = Data{
	"Da5ta": "n0thing",
	"Data2": Data{
		"more": "d5ata89",
	},
	"Data4": []interface{}{
		"hey",
		"jude",
	},
	"1": []byte{
		'h',
		'i',
	},
	"2": float64(0.64),
	"3": int64(99),
	"4": true,
}

var testDataGobOnly = map[interface{}]interface{}{
	"Da5ta": "n0thing",
	"Data2": map[interface{}]interface{}{
		"more": "d5ata89",
	},
	"Data4": []interface{}{
		"hey",
		"jude",
	},
	"1": []byte{
		'h',
		'i',
	},
	"2": float64(0.64),
	"3": int64(99),
	"4": true,
}
var testDataMapii = Data{
	"1": Data{
		"hey": "jude",
	},
}

var testDataSlice = Data{
	"1": []interface{}{
		"hey",
		"jude",
	},
}
var testDataBytes = Data{
	"1": []byte{
		'h',
		'i',
	},
}
var testDataFloat64 = Data{
	"1": float64(0.64),
	"2": float32(0.64),
}

var testDataIntUint = Data{
	"1": int64(99),
	"2": int32(99),
	"3": int16(99),
	"4": int8(99),
	"5": int64(99),
	"6": int32(99),
	"7": int16(99),
	"8": int8(99),
}

var testDataBool = Data{
	"1": true,
	"2": false,
}

var testEmpty = Data{}

var testEmptyMap = Data{
	"none": map[interface{}]interface{}{},
}

func init() {
	gob.Register(map[interface{}]interface{}{})
	gob.Register([]interface{}{})
}

func TestGob(t *testing.T) {
	var data Data
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	enc.Encode(testData)

	dec := gob.NewDecoder(buf)
	dec.Decode(&data)

	runtime.GC()
	if !reflect.DeepEqual(data, testData) {
		fmt.Println(data)
		fmt.Println(testData)
		t.FailNow()
	}
}
func TestMapII(t *testing.T) {
	var data Data
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	enc.Encode(testDataMapii)

	dec := gob.NewDecoder(buf)
	dec.Decode(&data)

	runtime.GC()

	if !reflect.DeepEqual(data, testDataMapii) {
		fmt.Println(data)
		fmt.Println(testDataMapii)
		t.FailNow()
	}
}

func TestMapEmpty(t *testing.T) {
	var data Data
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	enc.Encode(testEmpty)

	dec := gob.NewDecoder(buf)
	dec.Decode(&data)

	runtime.GC()

	if !reflect.DeepEqual(data, testEmpty) {
		fmt.Println(data)
		fmt.Println(testEmpty)
		t.FailNow()
	}
}

func TestInterfaceSlice(t *testing.T) {
	var data Data
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	dec := gob.NewDecoder(buf)

	enc.Encode(testDataSlice)
	dec.Decode(&data)

	runtime.GC()
	if !reflect.DeepEqual(data, testDataSlice) {
		fmt.Println(data)
		fmt.Println(testDataSlice)
		t.FailNow()
	}
}

func TestMapIISize(t *testing.T) {
	bytes, _ := testDataMapii.GobEncode()
	size, _ := testDataMapii.Size()
	if len(bytes) != size {
		fmt.Println(size)
		fmt.Println(len(bytes))
	}
}
func TestMapEmptySize(t *testing.T) {
	bytes, _ := testEmpty.GobEncode()
	size, _ := testEmpty.Size()
	if len(bytes) != size {
		fmt.Println(size)
		fmt.Println(len(bytes))
	}
}
func TestMapNestedSize(t *testing.T) {
	bytes, _ := testEmptyMap.GobEncode()
	size, _ := testEmptyMap.Size()
	if len(bytes) != size {
		fmt.Println(size)
		fmt.Println(len(bytes))
	}
}
func TestInterfaceSliceSize(t *testing.T) {
	bytes, _ := testDataSlice.GobEncode()
	size, _ := testDataSlice.Size()
	if len(bytes) != size {
		fmt.Println(size)
		fmt.Println(len(bytes))
	}
}

func TestBytes(t *testing.T) {
	var data Data
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	enc.Encode(testDataBytes)

	dec := gob.NewDecoder(buf)
	dec.Decode(&data)

	runtime.GC()
	if !reflect.DeepEqual(data, testDataBytes) {
		fmt.Println(data)
		fmt.Println(testDataBytes)
		t.FailNow()
	}
}

func TestFloat(t *testing.T) {
	var data Data
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	enc.Encode(testDataFloat64)

	dec := gob.NewDecoder(buf)
	dec.Decode(&data)

	runtime.GC()
	if !reflect.DeepEqual(data, testDataFloat64) {
		fmt.Println(data)
		fmt.Println(testDataFloat64)
		t.FailNow()
	}
}
func TestIntUint(t *testing.T) {
	var data Data
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	enc.Encode(testDataIntUint)

	dec := gob.NewDecoder(buf)
	dec.Decode(&data)

	runtime.GC()
	if !reflect.DeepEqual(data, testDataIntUint) {
		fmt.Println(data)
		fmt.Println(testDataIntUint)
		t.FailNow()
	}
}

func TestBool(t *testing.T) {
	var data Data
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	enc.Encode(testDataBool)

	dec := gob.NewDecoder(buf)
	dec.Decode(&data)

	runtime.GC()
	if !reflect.DeepEqual(data, testDataBool) {
		fmt.Println(data)
		fmt.Println(testDataBool)
		t.FailNow()
	}
}

func TestValue(t *testing.T) {
	value := v1.Value{
		Key: v1.Key{
			Value: []byte("hey2"),
			Vtype: 's',
		},
		Children: []v1.Ikeytype{1, 4, 5},
		Value:    []byte("hey"),
		Vtype:    's',
	}
	v := v1.Value{}
	buf := new(bytes.Buffer)
	value.Tobytes(buf)
	v.FromBytes(buf.Bytes())

	runtime.GC()
	if !reflect.DeepEqual(value, v) {
		enc := json.NewEncoder(os.Stdout)
		enc.Encode(v)
		enc.Encode(value)
		t.Fail()
	}
}

func BenchmarkGobData(b *testing.B) {
	b.StopTimer()
	var data Data
	var byt []byte

	buf := bytes.NewBuffer(byt)
	enc := gob.NewEncoder(buf)
	for n := 0; n < b.N; n++ {
		enc.Encode(testData)
	}

	enc = gob.NewEncoder(ioutil.Discard)
	dec := gob.NewDecoder(buf)
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		err := enc.Encode(testData)
		if err != nil {
			panic(err)
		}

		err = dec.Decode(&data)
		if err != nil {
			panic(err)
		}
	}
}
func BenchmarkGobDataDecode(b *testing.B) {
	b.StopTimer()
	var data Data
	var byt []byte

	buf := bytes.NewBuffer(byt)
	enc := gob.NewEncoder(buf)
	for n := 0; n < b.N; n++ {
		enc.Encode(testData)
	}

	dec := gob.NewDecoder(buf)

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		err := dec.Decode(&data)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkGobDataEncode(b *testing.B) {
	b.StopTimer()
	enc := gob.NewEncoder(ioutil.Discard)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		err := enc.Encode(testData)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkGobMap(b *testing.B) {
	b.StopTimer()
	var data map[interface{}]interface{}
	var byt []byte

	buf := bytes.NewBuffer(byt)
	enc := gob.NewEncoder(buf)
	for n := 0; n < b.N; n++ {
		enc.Encode(testDataGobOnly)
	}

	enc = gob.NewEncoder(ioutil.Discard)
	dec := gob.NewDecoder(buf)
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		err := enc.Encode(testDataGobOnly)
		if err != nil {
			panic(err)
		}

		err = dec.Decode(&data)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkGobMapDecode(b *testing.B) {
	b.StopTimer()
	var data map[interface{}]interface{}
	var byt []byte

	buf := bytes.NewBuffer(byt)
	enc := gob.NewEncoder(buf)
	for n := 0; n < b.N; n++ {
		enc.Encode(testDataGobOnly)
	}

	dec := gob.NewDecoder(buf)
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		err := dec.Decode(&data)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkGobMapEncode(b *testing.B) {
	b.StopTimer()
	enc := gob.NewEncoder(ioutil.Discard)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		enc.Encode(testDataGobOnly)
	}
}

func TestText(t *testing.T) {
	d := Data{
		"hello": "world",
	}
	fmt.Println(d.GobEncode())
}

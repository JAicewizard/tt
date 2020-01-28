package tt

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"reflect"
	"runtime"
	"testing"
	"time"

	v1 "github.com/JAicewizard/tt/v1"
	"github.com/go-test/deep"
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
var testDataGobOnly = map[string]interface{}{
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

/* var testDataFLOAT64 = map[interface{}]interface{}{
	"1": float64(0.64),
	"2": float64(0.64),
	"3": float64(0.64),
	"4": float64(0.64),
	"5": float64(0.64),
	"6": float64(0.64),
	"7": float64(0.64),
	"8": float64(0.64),
	"9": float64(0.64),
}
*/
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

/* func TestMap(t *testing.T) {
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
} */
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
		t.FailNow()
	}
}
func TestMapEmptySize(t *testing.T) {
	bytes, _ := testEmpty.GobEncode()
	size, _ := testEmpty.Size()
	if len(bytes) != size {
		fmt.Println(size)
		fmt.Println(len(bytes))
		t.FailNow()
	}
}
func TestMapNestedSize(t *testing.T) {
	bytes, _ := testEmptyMap.GobEncode()
	size, _ := testEmptyMap.Size()
	if len(bytes) != size {
		fmt.Println(size)
		fmt.Println(len(bytes))
		t.FailNow()
	}
}
func TestInterfaceSliceSize(t *testing.T) {
	bytes, _ := testDataSlice.GobEncode()
	size, _ := testDataSlice.Size()
	if len(bytes) != size {
		fmt.Println(size)
		fmt.Println(len(bytes))
		t.FailNow()
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
func BenchmarkV3(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	var data map[string]interface{}
	var byt []byte

	buf := bytes.NewBuffer(byt)
	for n := 0; n < b.N; n++ {
		Encodev3(testDataGobOnly, buf)
	}

	buf2 := bytes.NewBuffer(nil)

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		Encodev3(testDataGobOnly, buf2)
		err := Decodev3(buf, &data)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkV3Decode(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	var data map[string]interface{}
	var byt []byte
	buf := bytes.NewBuffer(byt)
	for n := 0; n < b.N; n++ {
		Encodev3(testDataGobOnly, buf)
	}

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		Decodev3(buf, &data)
	}
}

func BenchmarkV3Encode(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()

	buf := bytes.NewBuffer(nil)

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		Encodev3(testDataGobOnly, buf)
	}
}

func BenchmarkGobData(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	var data Data
	var byt []byte
	buf := bytes.NewBuffer(byt)
	Encodev2(testData, buf)

	buf2 := bytes.NewBuffer(nil)

	b.StartTimer()

	for n := 0; n < b.N; n++ {
		Encodev2(testData, buf2)

		err := Decodev2(buf.Bytes()[1:], &data)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkGobDataDecode(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	var data Data
	var byt []byte

	buf := bytes.NewBuffer(byt)
	Encodev2(testData, buf)

	b.StartTimer()
	for n := 0; n < b.N; n++ {

		err := Decodev2(buf.Bytes()[1:], &data)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkGobDataEncode(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	buf := bytes.NewBuffer(nil)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		Encodev2(testData, buf)
	}
}

func BenchmarkGobMap(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	var data map[string]interface{}
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
	b.ReportAllocs()
	b.StopTimer()
	var data map[string]interface{}
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
	b.ReportAllocs()
	b.StopTimer()
	enc := gob.NewEncoder(ioutil.Discard)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		enc.Encode(testDataGobOnly)
	}
}

type testBasicStruct struct {
	Name string `TT:"hi"`
}

type testembeddedPrivateStruct struct {
	Pame string `TT:"oops"`
	testBasicStruct
}
type testembeddedStruct struct {
	Pame  string `TT:"oops"`
	Embed testBasicStruct
}

type testCase struct {
	name  string
	data  interface{}
	bytes [][]byte
}

var testCases = []testCase{
	testCase{
		name:  "testBasicStruct",
		data:  testBasicStruct{"hello"},
		bytes: [][]byte{[]byte{3, 0, 0, 18, 0, 1, 5, 2, 1, 104, 101, 108, 108, 111, 2, 104, 105, 0}},
	},
	testCase{
		name:  "testBasicSlice",
		data:  []string{"hello", "world"},
		bytes: [][]byte{[]byte{3, 0, 0, 19, 0, 2, 5, 0, 1, 104, 101, 108, 108, 111, 0, 0, 5, 0, 1, 119, 111, 114, 108, 100, 0, 0}},
	},
	testCase{
		name:  "testcomplexSlice",
		data:  []interface{}{"hello", int64(5)},
		bytes: [][]byte{[]byte{3, 0, 0, 19, 0, 2, 5, 0, 1, 104, 101, 108, 108, 111, 0, 0, 8, 0, 6, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
	},
	testCase{
		name:  "testBasicFloatSlice",
		data:  []float32{0.14, 0.20},
		bytes: [][]byte{[]byte{3, 0, 0, 19, 0, 2, 4, 0, 12, 41, 92, 15, 62, 0, 0, 4, 0, 12, 205, 204, 76, 62, 0, 0}},
	},
	testCase{
		name:  "testBasicGobEncoder",
		data:  time.Date(2020, time.January, 26, 21, 6, 30, 5, time.UTC),
		bytes: [][]byte{[]byte{3, 15, 0, 2, 1, 0, 0, 0, 14, 213, 191, 246, 86, 0, 0, 0, 5, 255, 255, 0, 0}},
	},
	testCase{
		name: "testBasicMap",
		data: map[string]string{
			"hey": "jude",
			"bye": "pi",
		},
		bytes: [][]byte{[]byte{3, 0, 0, 18, 0, 2, 4, 3, 1, 106, 117, 100, 101, 1, 104, 101, 121, 0, 2, 3, 1, 112, 105, 1, 98, 121, 101, 0},
			[]byte{3, 0, 0, 18, 0, 2, 2, 3, 1, 112, 105, 1, 98, 121, 101, 0, 4, 3, 1, 106, 117, 100, 101, 1, 104, 101, 121, 0}},
	},
	testCase{
		name: "testComplexMap",
		data: map[interface{}]interface{}{
			"hey": "jude",
			"bye": math.Pi,
		},
		bytes: [][]byte{[]byte{3, 0, 0, 18, 0, 2, 8, 3, 13, 24, 45, 68, 84, 251, 33, 9, 64, 1, 98, 121, 101, 0, 4, 3, 1, 106, 117, 100, 101, 1, 104, 101, 121, 0},
			[]byte{3, 0, 0, 18, 0, 2, 4, 3, 1, 106, 117, 100, 101, 1, 104, 101, 121, 0, 8, 3, 13, 24, 45, 68, 84, 251, 33, 9, 64, 1, 98, 121, 101, 0}},
	},
	testCase{
		name:  "testEmbeddedPrivateStruct",
		data:  testembeddedPrivateStruct{Pame: "hi", testBasicStruct: testBasicStruct{Name: "lol"}},
		bytes: [][]byte{[]byte{3, 0, 0, 18, 0, 1, 2, 4, 1, 104, 105, 2, 111, 111, 112, 115, 0}},
	},
	testCase{
		name: "testEmbeddedStruct",
		data: testembeddedStruct{Pame: "hi", Embed: testBasicStruct{Name: "lol"}},
		bytes: [][]byte{[]byte{3, 0, 0, 18, 0, 2, 2, 4, 1, 104, 105, 2, 111, 111, 112, 115, 0, 0, 5, 18, 2, 69, 109, 98, 101, 100, 1, 3, 2, 1, 108, 111, 108, 2, 104, 105, 0},
			[]byte{3, 0, 0, 18, 0, 2, 0, 5, 18, 2, 69, 109, 98, 101, 100, 1, 3, 2, 1, 108, 111, 108, 2, 104, 105, 0, 2, 4, 1, 104, 105, 2, 111, 111, 112, 115, 0}},
	},
	//3, 0, 0, 18, 0, 2, 0, 5, 18, 2, 69, 109, 98, 101, 100, 1, 3, 2, 1, 108, 111, 108, 2, 104, 105, 0, 2, 4, 1, 104, 105, 2, 111, 111, 112, 115, 0
	//3, 0, 0, 18, 0, 2, 2, 4, 1, 104, 105, 2, 111, 111, 112, 115, 0, 0, 5, 18, 2, 69, 109, 98, 101, 100, 1, 3, 2, 1, 108, 111, 108, 2, 104, 105, 0
}

func testStructDecode(t *testing.T, testcase testCase) {
	buf := &bytes.Buffer{}
	Encodev3(testcase.data, buf)
	isCorrect := false
	var diff []string
	for _, b := range testcase.bytes {
		if diff = deep.Equal(string(b), buf.String()); diff == nil {
			isCorrect = true
		}
	}
	if !isCorrect {
		t.Error(buf.Bytes())
		t.Error(diff)
	}

	after := reflect.New(reflect.TypeOf(testcase.data)).Elem()
	err := Decodev3(buf, after.Addr().Interface())
	if err != nil {
		t.Error(err)
	}
	fmt.Println(testcase.data)
	fmt.Println(after.Interface())
	//this only tests public fields
	if diff := deep.Equal(testcase.data, after.Interface()); diff != nil {
		fmt.Println(testcase.data)
		fmt.Println(after.Interface())
		t.Error(diff)
	}
}

func TestEncodeDecode(t *testing.T) {
	for _, c := range testCases {
		t.Run(c.name, func(te *testing.T) {
			testStructDecode(te, c)
		})
	}
}

type testStructEmbeded struct {
	Name string `TT:"hi"`
}

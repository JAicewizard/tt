package tt

import (
	"bufio"
	"bytes"
	"crypto/rand"
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

	"github.com/go-test/deep"
	fuzz "github.com/google/gofuzz"
	v1 "github.com/jaicewizard/tt/v1"
)

func init() {
	runtime.GOMAXPROCS(1)
}

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
	"1": map[interface{}]interface{}{
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

	if dif := deep.Equal(data, testDataMapii); len(dif) != 0 {
		fmt.Println(data)
		fmt.Println(testDataMapii)
		fmt.Println(dif)
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

/* func TestMapIISize(t *testing.T) {
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
} */

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
	if dif := deep.Equal(data, testDataIntUint); len(dif) != 0 {
		fmt.Println(data)
		fmt.Println(testDataIntUint)
		fmt.Println(dif)
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
	if dif := deep.Equal(data, testDataBool); len(dif) != 0 {
		fmt.Println(data)
		fmt.Println(testDataBool)
		fmt.Println(dif)
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
	bu := bufio.NewReader(buf)
	for n := 0; n < b.N; n++ {
		err := Decodev3(bu, &data)
		if err != nil {
			panic(err)
		}
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

func BenchmarkV3int64(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()

	buf := bytes.NewBuffer(nil)

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		Encodev3([]int64{6734, 5, 34598, 2354983045, 8732, 37, 23492}, buf)
	}
}

func BenchmarkV2(b *testing.B) {
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

func BenchmarkV2Decode(b *testing.B) {
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

func BenchmarkV2Encode(b *testing.B) {
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
	{
		name:  "testBasicStruct",
		data:  testBasicStruct{"hello"},
		bytes: [][]byte{{3, 0, 0, 0, 18, 0, 1, 5, 2, 1, 104, 101, 108, 108, 111, 2, 104, 105, 0}},
	},
	{
		name:  "testBasicSlice",
		data:  []string{"hello", "world"},
		bytes: [][]byte{{3, 0, 0, 0, 19, 0, 2, 5, 0, 1, 104, 101, 108, 108, 111, 0, 0, 5, 0, 1, 119, 111, 114, 108, 100, 0, 0}},
	},
	{
		name:  "testcomplexSlice",
		data:  []interface{}{"hello", int64(5)},
		bytes: [][]byte{{3, 0, 0, 0, 19, 0, 2, 5, 0, 1, 104, 101, 108, 108, 111, 0, 0, 8, 0, 6, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
	},
	{
		name:  "testBasicFloatSlice",
		data:  []float32{0.14, 0.20},
		bytes: [][]byte{{3, 0, 0, 0, 19, 0, 2, 4, 0, 12, 41, 92, 15, 62, 0, 0, 4, 0, 12, 205, 204, 76, 62, 0, 0}},
	},
	{
		name:  "testBasicGobEncoder",
		data:  time.Date(2020, time.January, 26, 21, 6, 30, 5, time.UTC),
		bytes: [][]byte{{3, 0, 15, 0, 2, 1, 0, 0, 0, 14, 213, 191, 246, 86, 0, 0, 0, 5, 255, 255, 0, 0}},
	},
	{
		name: "testBasicMap",
		data: map[string]string{
			"hey": "jude",
			"bye": "pi",
		},
		bytes: [][]byte{
			{3, 0, 0, 0, 18, 0, 2, 4, 3, 1, 106, 117, 100, 101, 1, 104, 101, 121, 0, 2, 3, 1, 112, 105, 1, 98, 121, 101, 0},
			{3, 0, 0, 0, 18, 0, 2, 2, 3, 1, 112, 105, 1, 98, 121, 101, 0, 4, 3, 1, 106, 117, 100, 101, 1, 104, 101, 121, 0}},
	},
	{
		name: "testComplexMap",
		data: map[interface{}]interface{}{
			"hey": "jude",
			"bye": math.Pi,
		},
		bytes: [][]byte{
			{3, 0, 0, 0, 18, 0, 2, 8, 3, 13, 24, 45, 68, 84, 251, 33, 9, 64, 1, 98, 121, 101, 0, 4, 3, 1, 106, 117, 100, 101, 1, 104, 101, 121, 0},
			{3, 0, 0, 0, 18, 0, 2, 4, 3, 1, 106, 117, 100, 101, 1, 104, 101, 121, 0, 8, 3, 13, 24, 45, 68, 84, 251, 33, 9, 64, 1, 98, 121, 101, 0}},
	},
	{
		name:  "testEmbeddedPrivateStruct",
		data:  testembeddedPrivateStruct{Pame: "hi", testBasicStruct: testBasicStruct{Name: "lol"}},
		bytes: [][]byte{{3, 0, 0, 0, 18, 0, 1, 2, 4, 1, 104, 105, 2, 111, 111, 112, 115, 0}},
	},
	{
		name: "testEmbeddedStruct",
		data: testembeddedStruct{Pame: "hi", Embed: testBasicStruct{Name: "lol"}},
		bytes: [][]byte{
			{3, 0, 0, 0, 18, 0, 2, 2, 4, 1, 104, 105, 2, 111, 111, 112, 115, 0, 0, 5, 18, 2, 69, 109, 98, 101, 100, 1, 3, 2, 1, 108, 111, 108, 2, 104, 105, 0},
			{3, 0, 0, 0, 18, 0, 2, 0, 5, 18, 2, 69, 109, 98, 101, 100, 1, 3, 2, 1, 108, 111, 108, 2, 104, 105, 0, 2, 4, 1, 104, 105, 2, 111, 111, 112, 115, 0}},
	},
	{
		name: "testSingleValueToInterface",
		data: "hello",
		bytes: [][]byte{
			{3, 0, 5, 0, 1, 104, 101, 108, 108, 111, 0, 0}},
	},
}

func testStructDecode(t *testing.T, testcase testCase) {
	buf := &bytes.Buffer{}
	err := Encodev3(testcase.data, buf)
	if err != nil {
		t.Error(err)
	}
	equal := false
	bytes := []byte{}
	var diff []string
	for _, bytes = range testcase.bytes {
		if diff = deep.Equal(string(bytes), buf.String()); diff == nil {
			equal = true
			break
		}
	}
	if !equal {
		t.Error(bytes)
		t.Error(buf.Bytes())
		t.Error(diff)
	}
	after := reflect.New(reflect.TypeOf(testcase.data)).Elem()
	err = Decodev3(buf, after.Addr().Interface())
	if err != nil {
		t.Error(err)
	}
	//this only tests public fields
	if diff := deep.Equal(testcase.data, after.Interface()); diff != nil {
		t.Error(testcase.data)
		t.Error(after.Interface())
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

type FuzzStruct struct {
	M       map[int32]Point
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Uint    uint
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	Int     int
	Float32 float32
	Float64 float64
	Map2    map[string]Point
}
type Point struct {
	X float32
	Y float64
}

func testFuz(t *testing.T, testcase FuzzStruct) {
	buf := &bytes.Buffer{}
	err := Encodev3(testcase, buf)
	if err != nil {
		t.Error(err)
	}
	after := reflect.New(reflect.TypeOf(testcase)).Elem()
	err = Decodev3(buf, after.Addr().Interface())
	if err != nil {
		t.Error(err)
	}
	//this only tests public fields
	deep.NilMapsAreEmpty = true
	if diff := deep.Equal(testcase, after.Interface()); diff != nil {
		t.Error(testcase)
		t.Error(after.Interface())
		t.Error(diff)
	}
}

func TestEncodeDecodeFuzz(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long-running test.")
	}
	data := make([]byte, 1000000)
	s := FuzzStruct{}
	for i := 0; i < 10000; i++ {
		t.Run(fmt.Sprintf("%d", i), func(te *testing.T) {
			rand.Read(data)
			fuzz.NewFromGoFuzz(data).Fuzz(&s)
			if i%100 == 0 {
				fmt.Printf("%d\n", i)
			}
			testFuz(t, s)
		})
	}
}

func TestDecodeFuzz(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long-running test.")
	}
	data := make([]byte, 10000)
	for i := 0; i < 50000; i++ {
		t.Run(fmt.Sprintf("%d", i), func(te *testing.T) {
			rand.Read(data)
			if i%500 == 0 {
				fmt.Printf("%d\n", i)
			}
			buf := bytes.NewBuffer(data[:10000])
			after := interface{}(nil)
			dec := NewV3Decoder(buf, true)
			dec.SetAllocLimmit(2 << 30) //1 GiB
			err := dec.Decode(&after)
			if err != nil {
				t.Log(err)
			}
		})
	}
}

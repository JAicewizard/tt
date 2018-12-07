package tt

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var testData = Data{
	"Da5ta": "n0thing",
	"Data2": Data{
		"more": "d5ata",
	},
}
var testDataSlice = Data{
	"Da5ta": "n0thing",
	"Data2": []interface{}{
		"hey",
		"jude",
	},
}

type loop struct {
	Data    []byte
	pointer int
	length  int
}

func (l *loop) Read(p []byte) (n int, err error) {
	for i := 0; i < len(p); i++ {
		p[i] = l.Data[l.pointer]
		l.pointer++
		if l.pointer >= l.length {
			l.pointer = 0
		}
	}
	return len(p), nil
}

func (l *loop) Reset() {
	l.pointer = 0
}

func TestGob(t *testing.T) {
	var data Data
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	dec := gob.NewDecoder(buf)
	enc.Encode(testData)
	dec.Decode(&data)

	if !reflect.DeepEqual(data, testData) {
		fmt.Println(data)
		fmt.Println(testData)
		t.FailNow()
	}

	data = Data{}
	buf = new(bytes.Buffer)
	enc = gob.NewEncoder(buf)
	dec = gob.NewDecoder(buf)

	enc.Encode(testDataSlice)
	dec.Decode(&data)

	if !reflect.DeepEqual(data, testDataSlice) {
		fmt.Println(data)
		fmt.Println(testDataSlice)
		t.FailNow()
	}
}

func TestValue(t *testing.T) {
	value := Value{
		Key: Key{
			Value: []byte("hey2"),
			Vtype: 's',
		},
		Children: []ikeytype{1, 4, 5},
		Value:    []byte("hey"),
		Vtype:    's',
	}
	v := Value{}
	buf := new(bytes.Buffer)
	value.tobytes(buf)
	v.fromBytes(buf.Bytes())
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
	enc.Encode(testData)
	dat := buf.Bytes()

	l := loop{
		Data:    dat,
		pointer: 0,
		length:  len(dat),
	}

	enc = gob.NewEncoder(ioutil.Discard)
	dec := gob.NewDecoder(&l)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		enc.Encode(testData)
		dec.Decode(&data)
	}
}

func BenchmarkGobDataEncode(b *testing.B) {
	b.StopTimer()
	enc := gob.NewEncoder(ioutil.Discard)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		enc.Encode(testData)
	}
}

func BenchmarkGobMap(b *testing.B) {
	b.StopTimer()

	var data map[interface{}]interface{}
	var testdata = map[interface{}]interface{}(testData)
	var byt []byte

	buf := bytes.NewBuffer(byt)
	enc := gob.NewEncoder(buf)
	enc.Encode(testdata)

	dat := buf.Bytes()
	l := loop{
		Data:    dat,
		pointer: 0,
		length:  len(dat),
	}

	enc = gob.NewEncoder(ioutil.Discard)
	dec := gob.NewDecoder(&l)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		enc.Encode(testdata)
		dec.Decode(&data)
	}
}

func BenchmarkGobMapEncode(b *testing.B) {
	b.StopTimer()
	var testdata = map[interface{}]interface{}(testData)
	enc := gob.NewEncoder(ioutil.Discard)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		enc.Encode(testdata)
	}
}

func BenchmarkGobDataS(b *testing.B) {
	b.StopTimer()
	var data Data
	var byt []byte

	buf := bytes.NewBuffer(byt)
	enc := gob.NewEncoder(buf)
	enc.Encode(testData)
	dat := buf.Bytes()

	l := loop{
		Data:    dat,
		pointer: 0,
		length:  len(dat),
	}

	enc = gob.NewEncoder(ioutil.Discard)
	dec := gob.NewDecoder(&l)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		enc.Encode(testDataSlice)
		dec.Decode(&data)
	}
}

func BenchmarkGobDataEncodeS(b *testing.B) {
	b.StopTimer()
	enc := gob.NewEncoder(ioutil.Discard)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		enc.Encode(testDataSlice)
	}
}

func BenchmarkGobMapS(b *testing.B) {
	b.StopTimer()

	var data map[interface{}]interface{}
	var testdata = map[interface{}]interface{}(testDataSlice)
	var byt []byte

	buf := bytes.NewBuffer(byt)
	enc := gob.NewEncoder(buf)
	enc.Encode(testdata)

	dat := buf.Bytes()
	l := loop{
		Data:    dat,
		pointer: 0,
		length:  len(dat),
	}

	enc = gob.NewEncoder(ioutil.Discard)
	dec := gob.NewDecoder(&l)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		enc.Encode(testdata)
		dec.Decode(&data)
	}
}

func BenchmarkGobMapEncodeS(b *testing.B) {
	b.StopTimer()
	var testdata = map[interface{}]interface{}(testDataSlice)
	enc := gob.NewEncoder(ioutil.Discard)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		enc.Encode(testdata)
	}
}

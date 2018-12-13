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
		"more": "d5ata89",
	},
	"Data4": []interface{}{
		"hey",
		"jude",
	},
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
}

var testDataSlice = Data{
	"Da5ta": "n0thing",
	"Data2": []interface{}{
		"hey",
		"jude",
	},
}
var testDataSliceGobOnly = map[interface{}]interface{}{
	"Da5ta": "n0thing",
	"Data2": []interface{}{
		"hey",
		"jude",
	},
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

func BenchmarkGobMapEncode(b *testing.B) {
	b.StopTimer()
	enc := gob.NewEncoder(ioutil.Discard)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		enc.Encode(testDataGobOnly)
	}
}

package main

import "reflect"

var m = map[string]string{
	"hey": "byew",
}

func main() {
	reflect.PtrTo(reflect.TypeOf(m["hey"]))
}

func changeValue(val *string) {
	*val = "byetbuewrgjhdfbgjkhdsfgkjhdsbfgjkhdsfkhdsfkg"
}

package otira

import (
	"errors"
	"testing"
)

func TestAllSupportedTypes(t *testing.T) {
	var i int32
	var i2 int64
	var f float32
	var f2 float64
	var b []byte
	var s string
	if supportedType(i) && supportedType(i2) && supportedType(f) && supportedType(f2) && supportedType(b) && supportedType(s) {
		return
	}
	t.Fatal(errors.New("Type should be supported"))

}

func TestSomeUnSupportedTypes(t *testing.T) {
	var i []int
	var f []float32
	var s []string
	if !supportedType(t) && !supportedType(i) && !supportedType(f) && !supportedType(s) {
		return
	}
	t.Fatal(errors.New("Type should NOT be supported"))

}

package main

import (
	"testing"
)

type serializeTestStruct struct {
	A string
	B int
	C map[string]int
}

func TestSerialize(t *testing.T) {
	test := &serializeTestStruct{"yes", 42, make(map[string]int)}
	test.C["potato"] = 68
	s := Serialize(test)
	res := &serializeTestStruct{}
	UnSerialize(res, s)
	if test.A != res.A || test.B != res.B || test.C["potato"] != res.C["potato"] {
		t.Errorf("Expected %v but got %v", test, res)
	}
}

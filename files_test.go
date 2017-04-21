package main

import (
	"os"
	"testing"
)

type testdatatype struct {
	A string
	B rune
	C int
	D struct {
		Boo bool
		In  int
	}
	E []string
}

func TestDataSaveLoadNonExistent(t *testing.T) {
	os.Remove("testdata.tmp")
	DataSaveLoadAnything(t)
	os.Remove("testdata.tmp")
}

func TestDataSaveLoadExistent(t *testing.T) {
	os.Remove("testdata.tmp")
	c, _ := os.Create("testdata.tmp")
	c.Close()
	DataSaveLoadAnything(t)
	os.Remove("testdata.tmp")
}

func DataSaveLoadAnything(t *testing.T) {
	os.Remove("testdata.tmp")
	testdata := testdatatype{
		"more than a fish!",
		'€',
		19283791,
		struct {
			Boo bool
			In  int
		}{Boo: false, In: 42},
		[]string{
			"lol",
			"rofl",
			"caz€€€",
		},
	}

	err := SaveAs(testdata, "testdata.tmp")
	if err != nil {
		t.Errorf(err.Error())
	}
	res := &testdatatype{}
	err = LoadFrom(res, "testdata.tmp")
	if err != nil {
		t.Errorf(err.Error())
	}
	switch {
	case res.A != testdata.A:
		t.Errorf("Wrong string")
	case res.B != testdata.B:
		t.Errorf("Wrong rune")
	case res.C != testdata.C:
		t.Errorf("Wrong int")
	case res.D.Boo != testdata.D.Boo || res.D.In != testdata.D.In:
		t.Errorf("Wrong struct")
	default:
		for i, _ := range testdata.E {
			if testdata.E[i] != res.E[i] {
				t.Errorf("Wrong string slice element")
			}
		}
		if len(testdata.E) != len(res.E) {
			t.Errorf("Wrong string slice length")
		}
	}
}
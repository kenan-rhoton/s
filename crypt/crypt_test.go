package crypt

import (
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	testdata := []string{
		"There was a boy... a very strange enchanted boy..",
		"Cantábamos alegres: sólo veíamos acentos",
		"¿Cuánto cuesta? ¿50€? ¡Flipa!",
	}
	for _, v := range testdata {
		key, err := GenerateKey(1024)
		t.Log(err)
		encoded, err := Encode(v, &key.PublicKey)
		t.Log(err)
		decoded, err := Decode(encoded, key)
		t.Log(err)
		if decoded != v {
			t.Errorf("Expected \"%s\" but got \"%s\"\n", v, decoded)
		}
	}
}

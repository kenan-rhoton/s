package client

import (
	"bytes"
	"encoding/gob"
)

func Serialize(obj interface{}) string {
	res := &bytes.Buffer{}
	gobber := gob.NewEncoder(res)
	gobber.Encode(obj)
	return res.String()
}
func UnSerialize(obj interface{}, serial string) {
	reader := bytes.NewBufferString(serial)
	gobber := gob.NewDecoder(reader)
	gobber.Decode(obj)
}

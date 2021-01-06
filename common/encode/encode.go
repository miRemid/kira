package encode

import (
	"bytes"
	"encoding/gob"
)

func StructToBytes(st interface{}) ([]byte, error) {
	var buf bytes.Buffer
	env := gob.NewEncoder(&buf)
	if err := env.Encode(st); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func BytesToStruct(data []byte, st interface{}) error {
	var buf bytes.Buffer
	buf.Write(data)
	den := gob.NewDecoder(&buf)
	return den.Decode(st)
}

package cache

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

type Converter interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, v interface{}) error
}

type JsonConverter struct{}

func (JsonConverter) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func (JsonConverter) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type GobConverter struct{}

func (GobConverter) Encode(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
func (GobConverter) Decode(data []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(v)
}

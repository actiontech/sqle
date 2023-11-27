package utils

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
)

func Marshal(i interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	encoder := jsoniter.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(i)
	return buffer.Bytes(), err
}

func Unmarshal(data []byte, i interface{}) error {
	reader := bytes.NewReader(data)
	decoder := jsoniter.NewDecoder(reader)
	decoder.UseNumber()
	return decoder.Decode(i)
}

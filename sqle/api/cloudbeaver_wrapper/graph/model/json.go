package model

import (
	"encoding/json"
	"io"
)

type JSON map[string]interface{}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (b *JSON) UnmarshalGQL(v interface{}) error {
	*b = make(map[string]interface{})
	byteData, err := json.Marshal(v)
	if err != nil {
		panic("FAIL WHILE MARSHAL SCHEME")
	}
	tmp := make(map[string]interface{})
	err = json.Unmarshal(byteData, &tmp)
	if err != nil {
		panic("FAIL WHILE UNMARSHAL SCHEME")
		//return fmt.Errorf("%v", err)
	}
	*b = tmp
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (b JSON) MarshalGQL(w io.Writer) {
	byteData, err := json.Marshal(b)
	if err != nil {
		panic("FAIL WHILE MARSHAL SCHEME")
	}
	_, _ = w.Write(byteData)
}

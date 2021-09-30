package json

import (
	"bytes"

	"encoding/json"
	"strings"
)

func ToJson(s string) string {
	var out bytes.Buffer
	json.Indent(&out, []byte(s), "", "    ")
	return out.String()
}
func Marshal(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return string(b)
}
func MarshalAndFormat(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return string(ToJson(string(b)))
}
func Unmarshal(s string) interface{} {
	var result interface{}
	d := json.NewDecoder(strings.NewReader(s))
	d.UseNumber()
	err := d.Decode(result)
	if err != nil {
		panic(err)
	}
	return result
}

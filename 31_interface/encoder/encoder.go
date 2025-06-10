// This program illustrates the use of interfaces and custom types to encode data to JSON,
// which can be use to extend or overwrite the built-in JSON-Marshaller.

// Note that: none of the basic types (like int, string, bool, float64, map, slice, etc.)
// implement the json.Marshaler interface directly.
// The encoding/json package has built-in support for marshalling basic types, but these types
// do not satisfy json.Marshaler on their own, so a type switch with a case json.Marshaler will
// not match them unless they are wrapped in a custom type that explicitly implements MarshalJSON().

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func main() {
	var result []byte
	result = encode("one string")
	fmt.Println(string(result))

	result = encode(42)
	fmt.Println(string(result))

	// uncomment to see that the anonymous struct has not a defined marshalling function
	// result = encode(struct{ A string }{A: "one string"})
	// fmt.Println(string(result))

	result = encode(map[string]string{"color": "yellow", "weight": "123", "family": "citrus"})
	fmt.Println(string(result))

	result = encode(Email{name: "rob", domain: "go.dev"})
	fmt.Println(string(result))
}

func encode(v interface{}) []byte {
	var result []byte

	switch d := v.(type) {
	case string:
		result = encodeString(d)

	case int:
		result = encodeInt(d)

	case map[string]string:
		result = encodeMap(d)

	case json.Marshaler:
		result = encodeMarshaler(d)

	default:
		_ = d
		panic(fmt.Sprintf("error encoding %T: %s\n", v, "unknown data type"))
	}

	return result
}

func encodeMarshaler(m json.Marshaler) []byte {
	b, err := m.MarshalJSON()
	if err != nil {
		panic(fmt.Sprintf("error encoding %T: %s\n", m, err))
	}
	return b
}

func encodeMap(m map[string]string) []byte {
	js := "{"
	for k, v := range m {
		js = js + `"` + k + `":"` + v + `" `
	}
	js = js + "}"
	return []byte(js)
}

func encodeString(s string) []byte {
	fmt := `{"value":"` + s + `"}"`
	return []byte(fmt)
}

func encodeInt(i int) []byte {
	fmt := `{"value":"` + strconv.Itoa(i) + `"}"`
	return []byte(fmt)
}

type Email struct {
	name   string
	domain string
}

func (e Email) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("{\"email\":\"%s@%s\"}", e.name, e.domain)), nil
}

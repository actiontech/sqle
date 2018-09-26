package dry

import (
	// "strings"
	// "fmt"
	"testing"
)

func Test_ReflectSort(t *testing.T) {
	ints := []int{3, 5, 0, 2, 1, 4}
	ReflectSort(ints, func(a, b int) bool {
		return a < b
	})
	for i := range ints {
		if i != ints[i] {
			t.Fail()
		}
	}

	strings := []string{"aaa", "bbb", "abb", "aab"}
	ReflectSort(strings, func(a, b string) bool {
		return a < b
	})
	if strings[0] != "aaa" {
		t.Fail()
	}
	if strings[1] != "aab" {
		t.Fail()
	}
	if strings[2] != "abb" {
		t.Fail()
	}
	if strings[3] != "bbb" {
		t.Fail()
	}
}

type TestStruct struct {
	String  string
	Int     int
	Uint8   uint8
	Float32 float32
	Bool    bool
}

func Test_ReflectSetStructFieldsFromStringMap(t *testing.T) {
	structPtr := new(TestStruct)
	m := map[string]string{
		"String":  "Hello World",
		"Int":     "666",
		"Uint8":   "234",
		"Float32": "0.01",
		"Bool":    "true",
	}
	err := ReflectSetStructFieldsFromStringMap(structPtr, m, true)
	if err != nil {
		t.Fatal(err)
	}
	if structPtr.String != "Hello World" ||
		structPtr.Int != 666 ||
		structPtr.Uint8 != 234 ||
		structPtr.Float32 != 0.01 ||
		structPtr.Bool != true {
		t.Fatalf("Invalid values: %#v", structPtr)
	}

	m["NotExisting"] = "xxx"

	structPtr = new(TestStruct)
	err = ReflectSetStructFieldsFromStringMap(structPtr, m, true)
	if err == nil {
		t.Fail()
	}

	structPtr = new(TestStruct)
	err = ReflectSetStructFieldsFromStringMap(structPtr, m, false)
	if err != nil {
		t.Fatal(err)
	}
	if structPtr.String != "Hello World" ||
		structPtr.Int != 666 ||
		structPtr.Uint8 != 234 ||
		structPtr.Float32 != 0.01 ||
		structPtr.Bool != true {
		t.Fatalf("Invalid values: %#v", structPtr)
	}
}

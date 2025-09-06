package zaptext_test

import (
	"bytes"
	"testing"

	. "github.com/kaiiak/zaptext"
)

func TestNewReflectEncoder(t *testing.T) {
	t.Run("Create new reflect encoder", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)

		if encoder == nil {
			t.Errorf("NewReflectEncoder returned nil")
		}
	})
}

func TestReflectEncoderEncode(t *testing.T) {
	t.Run("Encode nil value", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)

		err := encoder.Encode(nil)
		if err != nil {
			t.Errorf("Encode(nil) returned error: %v", err)
		}
	})

	t.Run("Encode string value", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)

		err := encoder.Encode("test string")
		if err != nil {
			t.Errorf("Encode('test string') returned error: %v", err)
		}
	})

	t.Run("Encode integer value", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)

		err := encoder.Encode(42)
		if err != nil {
			t.Errorf("Encode(42) returned error: %v", err)
		}
	})

	t.Run("Encode struct value", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)

		type TestStruct struct {
			Name  string
			Value int
		}

		testObj := TestStruct{Name: "test", Value: 123}
		err := encoder.Encode(testObj)
		if err != nil {
			t.Errorf("Encode(struct) returned error: %v", err)
		}
	})

	t.Run("Encode map value", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)

		testMap := map[string]any{
			"key1": "value1",
			"key2": 42,
			"key3": true,
		}

		err := encoder.Encode(testMap)
		if err != nil {
			t.Errorf("Encode(map) returned error: %v", err)
		}
	})

	t.Run("Encode slice value", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)

		testSlice := []string{"item1", "item2", "item3"}

		err := encoder.Encode(testSlice)
		if err != nil {
			t.Errorf("Encode(slice) returned error: %v", err)
		}
	})
}

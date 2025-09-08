package zaptext_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	. "github.com/kaiiak/zaptext"
)

func TestNewReflectEncoder(t *testing.T) {
	t.Run("Create new reflect encoder", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)

		if encoder == nil {
			t.Errorf("NewReflectEncoder returned nil")
		}
		
		// Test Release method
		encoder.Release()
	})
}

func TestReflectEncoderEncode(t *testing.T) {
	t.Run("Encode nil value", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		err := encoder.Encode(nil)
		if err != nil {
			t.Errorf("Encode(nil) returned error: %v", err)
		}
		
		output := w.String()
		if output != "null" {
			t.Errorf("Expected 'null', got '%s'", output)
		}
	})

	t.Run("Encode string value", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		err := encoder.Encode("test string")
		if err != nil {
			t.Errorf("Encode('test string') returned error: %v", err)
		}
		
		output := w.String()
		expected := `"test string"`
		if output != expected {
			t.Errorf("Expected '%s', got '%s'", expected, output)
		}
	})

	t.Run("Encode integer value", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		err := encoder.Encode(42)
		if err != nil {
			t.Errorf("Encode(42) returned error: %v", err)
		}
		
		output := w.String()
		if output != "42" {
			t.Errorf("Expected '42', got '%s'", output)
		}
	})

	t.Run("Encode boolean values", func(t *testing.T) {
		w1 := &bytes.Buffer{}
		encoder1 := NewReflectEncoder(w1)
		defer encoder1.Release()

		err := encoder1.Encode(true)
		if err != nil {
			t.Errorf("Encode(true) returned error: %v", err)
		}
		
		if w1.String() != "true" {
			t.Errorf("Expected 'true', got '%s'", w1.String())
		}

		w2 := &bytes.Buffer{}
		encoder2 := NewReflectEncoder(w2)
		defer encoder2.Release()

		err = encoder2.Encode(false)
		if err != nil {
			t.Errorf("Encode(false) returned error: %v", err)
		}
		
		if w2.String() != "false" {
			t.Errorf("Expected 'false', got '%s'", w2.String())
		}
	})

	t.Run("Encode struct value", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		type TestStruct struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}

		testObj := TestStruct{Name: "test", Value: 123}
		err := encoder.Encode(testObj)
		if err != nil {
			t.Errorf("Encode(struct) returned error: %v", err)
		}
		
		output := w.String()
		// Should contain both fields (now expects proper JSON format)
		if !strings.Contains(output, `"name":"test"`) || !strings.Contains(output, `"value":123`) {
			t.Errorf("Expected struct with name and value fields, got '%s'", output)
		}
	})

	t.Run("Encode map value", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		testMap := map[string]any{
			"key1": "value1",
			"key2": 42,
			"key3": true,
		}

		err := encoder.Encode(testMap)
		if err != nil {
			t.Errorf("Encode(map) returned error: %v", err)
		}
		
		output := w.String()
		// Should start and end with braces and contain the keys
		if !strings.HasPrefix(output, "{") || !strings.HasSuffix(output, "}") {
			t.Errorf("Expected map format with braces, got '%s'", output)
		}
	})

	t.Run("Encode slice value", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		testSlice := []string{"item1", "item2", "item3"}

		err := encoder.Encode(testSlice)
		if err != nil {
			t.Errorf("Encode(slice) returned error: %v", err)
		}
		
		output := w.String()
		// Should start and end with brackets
		if !strings.HasPrefix(output, "[") || !strings.HasSuffix(output, "]") {
			t.Errorf("Expected array format with brackets, got '%s'", output)
		}
	})

	t.Run("Encode float value", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		err := encoder.Encode(3.14159)
		if err != nil {
			t.Errorf("Encode(3.14159) returned error: %v", err)
		}
		
		output := w.String()
		if !strings.Contains(output, "3.14159") {
			t.Errorf("Expected float representation, got '%s'", output)
		}
	})

	t.Run("Encode time value", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		testTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
		err := encoder.Encode(testTime)
		if err != nil {
			t.Errorf("Encode(time) returned error: %v", err)
		}
		
		output := w.String()
		// Should be in RFC3339 format with quotes
		expected := `"2023-01-01T12:00:00Z"`
		if output != expected {
			t.Errorf("Expected '%s', got '%s'", expected, output)
		}
	})

	t.Run("Encode pointer to value", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		value := 42
		ptr := &value
		err := encoder.Encode(ptr)
		if err != nil {
			t.Errorf("Encode(pointer) returned error: %v", err)
		}
		
		output := w.String()
		if output != "42" {
			t.Errorf("Expected '42', got '%s'", output)
		}
	})

	t.Run("Encode nil pointer", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		var ptr *int
		err := encoder.Encode(ptr)
		if err != nil {
			t.Errorf("Encode(nil pointer) returned error: %v", err)
		}
		
		output := w.String()
		if output != "null" {
			t.Errorf("Expected 'null', got '%s'", output)
		}
	})
}

func TestReflectEncoderErrorHandling(t *testing.T) {
	t.Run("Multiple encodes on same encoder", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		// First encode
		err := encoder.Encode("first")
		if err != nil {
			t.Errorf("First Encode returned error: %v", err)
		}

		// Reset buffer for second encode
		w.Reset()
		err = encoder.Encode("second")
		if err != nil {
			t.Errorf("Second Encode returned error: %v", err)
		}
		
		output := w.String()
		if output != `"second"` {
			t.Errorf("Expected '\"second\"', got '%s'", output)
		}
	})

	t.Run("Deep recursion protection", func(t *testing.T) {
		type recursive struct {
			Next *recursive `json:"next"`
			Data string     `json:"data"`
		}

		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		// Create a deep structure that will exceed maxDepth
		var root *recursive
		current := &recursive{Data: "root"}
		root = current
		
		// Create a chain longer than maxDepth
		for i := 0; i < 35; i++ {
			current.Next = &recursive{Data: fmt.Sprintf("level%d", i)}
			current = current.Next
		}

		err := encoder.Encode(root)
		if err == nil {
			t.Errorf("Expected error for deep recursion, but got none")
		}
		if !strings.Contains(err.Error(), "maximum encoding depth exceeded") {
			t.Errorf("Expected depth error, got: %v", err)
		}
	})

	t.Run("Encoder error state preservation", func(t *testing.T) {
		w := &failingWriter{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		// First encode should fail due to writer error
		err := encoder.Encode("test")
		if err == nil {
			t.Errorf("Expected error from failing writer")
		}

		// Second encode should return the same error (preserved state)
		err2 := encoder.Encode("test2")
		if err2 == nil {
			t.Errorf("Expected error to be preserved")
		}
	})
}

// failingWriter always returns an error on Write
type failingWriter struct{}

func (fw *failingWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("write failed")
}

func TestReflectEncoderComplexTypes(t *testing.T) {
	t.Run("Encode complex nested structure", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		type NestedStruct struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		type ComplexStruct struct {
			Title    string            `json:"title"`
			Count    int               `json:"count"`
			Tags     []string          `json:"tags"`
			Metadata map[string]string `json:"metadata"`
			Nested   NestedStruct      `json:"nested"`
		}

		testObj := ComplexStruct{
			Title: "Test Object",
			Count: 5,
			Tags:  []string{"tag1", "tag2"},
			Metadata: map[string]string{
				"author": "test",
				"version": "1.0",
			},
			Nested: NestedStruct{
				ID:   1,
				Name: "nested",
			},
		}

		err := encoder.Encode(testObj)
		if err != nil {
			t.Errorf("Encode(complex struct) returned error: %v", err)
		}
		
		output := w.String()
		// Verify the structure is properly encoded
		if !strings.Contains(output, `"title":"Test Object"`) {
			t.Errorf("Expected title field in output: %s", output)
		}
		if !strings.Contains(output, `"count":5`) {
			t.Errorf("Expected count field in output: %s", output)
		}
	})
}

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

func TestReflectEncoderSetMaxDepth(t *testing.T) {
	t.Run("SetMaxDepth configuration", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		// Set a very low max depth
		encoder.SetMaxDepth(2)

		// Create a deeply nested structure that should exceed the limit
		type DeepStruct struct {
			Name string
			Next *DeepStruct
		}

		deep := &DeepStruct{
			Name: "level1",
			Next: &DeepStruct{
				Name: "level2",
				Next: &DeepStruct{
					Name: "level3", // This should exceed depth 2
				},
			},
		}

		err := encoder.Encode(deep)
		if err == nil {
			t.Errorf("Expected depth limit error, but got nil")
		}
		if !strings.Contains(err.Error(), "maximum encoding depth exceeded") {
			t.Errorf("Expected depth limit error, got: %v", err)
		}
	})

	t.Run("SetMaxDepth with invalid value", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		// Try to set invalid depth (should be ignored)
		encoder.SetMaxDepth(0)
		encoder.SetMaxDepth(-1)

		// Should still work with simple structure (using default depth)
		err := encoder.Encode("test")
		if err != nil {
			t.Errorf("Expected no error for simple encoding, got: %v", err)
		}
	})
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
				"author":  "test",
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

func TestReflectEncoderStringEscaping(t *testing.T) {
	t.Run("HTML escaping with default behavior", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		// Test HTML characters that should be escaped by default
		testString := `<html>& "quotes" </div>`
		err := encoder.Encode(testString)
		if err != nil {
			t.Errorf("Encode(HTML string) returned error: %v", err)
		}

		output := w.String()
		// Should escape HTML characters
		if !strings.Contains(output, `\u003chtml\u003e`) {
			t.Errorf("Expected escaped HTML tags, got: %s", output)
		}
		if !strings.Contains(output, `\u0026`) {
			t.Errorf("Expected escaped ampersand, got: %s", output)
		}
	})

	t.Run("HTML escaping disabled", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		// Disable HTML escaping
		encoder.SetEscapeHTML(false)

		// Test HTML characters that should NOT be escaped when disabled
		testString := `<html>& "quotes" </div>`
		err := encoder.Encode(testString)
		if err != nil {
			t.Errorf("Encode(HTML string) returned error: %v", err)
		}

		output := w.String()
		// Should NOT escape HTML characters
		if strings.Contains(output, `\u003c`) || strings.Contains(output, `\u003e`) || strings.Contains(output, `\u0026`) {
			t.Errorf("Expected unescaped HTML characters, got: %s", output)
		}
		// But should still escape quotes
		if !strings.Contains(output, `\"quotes\"`) {
			t.Errorf("Expected escaped quotes, got: %s", output)
		}
	})

	t.Run("Special character escaping", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		// Test various special characters
		testString := "line1\nline2\rline3\tline4\\backslash\"quote"
		err := encoder.Encode(testString)
		if err != nil {
			t.Errorf("Encode(special chars) returned error: %v", err)
		}

		output := w.String()
		// Check for proper escaping
		if !strings.Contains(output, `\n`) {
			t.Errorf("Expected escaped newline, got: %s", output)
		}
		if !strings.Contains(output, `\r`) {
			t.Errorf("Expected escaped carriage return, got: %s", output)
		}
		if !strings.Contains(output, `\t`) {
			t.Errorf("Expected escaped tab, got: %s", output)
		}
		if !strings.Contains(output, `\\`) {
			t.Errorf("Expected escaped backslash, got: %s", output)
		}
		if !strings.Contains(output, `\"`) {
			t.Errorf("Expected escaped quote, got: %s", output)
		}
	})

	t.Run("Control character escaping", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		// Test control characters (ASCII 0-31)
		testString := "\x00\x01\x1f" // null, start of heading, unit separator
		err := encoder.Encode(testString)
		if err != nil {
			t.Errorf("Encode(control chars) returned error: %v", err)
		}

		output := w.String()
		// Should escape control characters as unicode
		if !strings.Contains(output, `\u0000`) {
			t.Errorf("Expected escaped null character, got: %s", output)
		}
		if !strings.Contains(output, `\u0001`) {
			t.Errorf("Expected escaped SOH character, got: %s", output)
		}
		if !strings.Contains(output, `\u001f`) {
			t.Errorf("Expected escaped US character, got: %s", output)
		}
	})

	t.Run("Complex number format", func(t *testing.T) {
		w := &bytes.Buffer{}
		encoder := NewReflectEncoder(w)
		defer encoder.Release()

		// Test complex64
		c64 := complex(float32(1.5), float32(2.5))
		err := encoder.Encode(c64)
		if err != nil {
			t.Errorf("Encode(complex64) returned error: %v", err)
		}

		output := w.String()
		// Should use JSON-compatible format
		if !strings.Contains(output, `{"real":1.5,"imag":2.5}`) {
			t.Errorf("Expected JSON-compatible complex number format, got: %s", output)
		}

		// Test complex128
		w.Reset()
		c128 := complex(3.14, 2.71)
		err = encoder.Encode(c128)
		if err != nil {
			t.Errorf("Encode(complex128) returned error: %v", err)
		}

		output = w.String()
		if !strings.Contains(output, `{"real":3.14,"imag":2.71}`) {
			t.Errorf("Expected JSON-compatible complex number format, got: %s", output)
		}
	})
}

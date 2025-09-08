package zaptext_test

import (
	"bytes"
	"testing"
	"time"

	. "github.com/kaiiak/zaptext"
)

type benchStruct struct {
	Name      string         `json:"name"`
	Age       int            `json:"age"`
	IsActive  bool           `json:"is_active"`
	Score     float64        `json:"score"`
	CreatedAt time.Time      `json:"created_at"`
	Tags      []string       `json:"tags"`
	Settings  map[string]any `json:"settings"`
}

func BenchmarkReflectEncoderEncode(b *testing.B) {
	testObj := benchStruct{
		Name:      "Test User",
		Age:       30,
		IsActive:  true,
		Score:     95.5,
		CreatedAt: time.Now(),
		Tags:      []string{"admin", "power_user", "beta_tester"},
		Settings: map[string]any{
			"theme":         "dark",
			"notifications": true,
			"limit":         100,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := bytes.NewBuffer(make([]byte, 0, 1024))
		encoder := NewReflectEncoder(w)
		if err := encoder.Encode(testObj); err != nil {
			b.Fatal(err)
		}
		encoder.Release()
	}
}

func BenchmarkReflectEncoderEncodeMap(b *testing.B) {
	testMap := map[string]any{
		"string_field": "hello world",
		"int_field":    42,
		"bool_field":   true,
		"float_field":  3.14159,
		"array_field":  []int{1, 2, 3, 4, 5},
		"nested_map": map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := bytes.NewBuffer(make([]byte, 0, 512))
		encoder := NewReflectEncoder(w)
		if err := encoder.Encode(testMap); err != nil {
			b.Fatal(err)
		}
		encoder.Release()
	}
}

func BenchmarkReflectEncoderEncodeSlice(b *testing.B) {
	testSlice := []any{
		"string",
		42,
		true,
		3.14159,
		[]int{1, 2, 3},
		map[string]string{"key": "value"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := bytes.NewBuffer(make([]byte, 0, 256))
		encoder := NewReflectEncoder(w)
		if err := encoder.Encode(testSlice); err != nil {
			b.Fatal(err)
		}
		encoder.Release()
	}
}

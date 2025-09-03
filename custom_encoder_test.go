package zaptext_test

import (
	"testing"
	"time"

	. "github.com/kaiiak/zaptext"
)

func TestCustomTimeEncoderFactory(t *testing.T) {
	t.Run("Create custom time encoder with RFC3339 layout", func(t *testing.T) {
		encoder := CustomTimeEncoderFactory(time.RFC3339)
		if encoder == nil {
			t.Errorf("CustomTimeEncoderFactory returned nil")
		}

		// Test the encoder with a mock PrimitiveArrayEncoder
		mockEncoder := &mockPrimitiveArrayEncoder{}
		testTime := time.Date(2023, 9, 2, 10, 30, 15, 0, time.UTC)

		encoder(testTime, mockEncoder)

		if !mockEncoder.called {
			t.Errorf("Expected encoder to call AppendString")
		}
		if mockEncoder.value != "2023-09-02T10:30:15Z" {
			t.Errorf("Expected time format '2023-09-02T10:30:15Z', got '%s'", mockEncoder.value)
		}
	})

	t.Run("Create custom time encoder with custom layout", func(t *testing.T) {
		layout := "2006-01-02 15:04:05"
		encoder := CustomTimeEncoderFactory(layout)

		mockEncoder := &mockPrimitiveArrayEncoder{}
		testTime := time.Date(2023, 9, 2, 10, 30, 15, 0, time.UTC)

		encoder(testTime, mockEncoder)

		if !mockEncoder.called {
			t.Errorf("Expected encoder to call AppendString")
		}
		if mockEncoder.value != "2023-09-02 10:30:15" {
			t.Errorf("Expected time format '2023-09-02 10:30:15', got '%s'", mockEncoder.value)
		}
	})

	t.Run("Create custom time encoder with kitchen layout", func(t *testing.T) {
		encoder := CustomTimeEncoderFactory(time.Kitchen)

		mockEncoder := &mockPrimitiveArrayEncoder{}
		testTime := time.Date(2023, 9, 2, 10, 30, 15, 0, time.UTC)

		encoder(testTime, mockEncoder)

		if !mockEncoder.called {
			t.Errorf("Expected encoder to call AppendString")
		}
		if mockEncoder.value != "10:30AM" {
			t.Errorf("Expected time format '10:30AM', got '%s'", mockEncoder.value)
		}
	})
}

// Mock implementation of zapcore.PrimitiveArrayEncoder for testing
type mockPrimitiveArrayEncoder struct {
	called bool
	value  string
}

func (m *mockPrimitiveArrayEncoder) AppendBool(v bool) {
	m.called = true
}

func (m *mockPrimitiveArrayEncoder) AppendByteString(v []byte) {
	m.called = true
	m.value = string(v)
}

func (m *mockPrimitiveArrayEncoder) AppendComplex128(v complex128) {
	m.called = true
}

func (m *mockPrimitiveArrayEncoder) AppendComplex64(v complex64) {
	m.called = true
}

func (m *mockPrimitiveArrayEncoder) AppendFloat64(v float64) {
	m.called = true
}

func (m *mockPrimitiveArrayEncoder) AppendFloat32(v float32) {
	m.called = true
}

func (m *mockPrimitiveArrayEncoder) AppendInt(v int) {
	m.called = true
}

func (m *mockPrimitiveArrayEncoder) AppendInt64(v int64) {
	m.called = true
}

func (m *mockPrimitiveArrayEncoder) AppendInt32(v int32) {
	m.called = true
}

func (m *mockPrimitiveArrayEncoder) AppendInt16(v int16) {
	m.called = true
}

func (m *mockPrimitiveArrayEncoder) AppendInt8(v int8) {
	m.called = true
}

func (m *mockPrimitiveArrayEncoder) AppendString(v string) {
	m.called = true
	m.value = v
}

func (m *mockPrimitiveArrayEncoder) AppendUint(v uint) {
	m.called = true
}

func (m *mockPrimitiveArrayEncoder) AppendUint64(v uint64) {
	m.called = true
}

func (m *mockPrimitiveArrayEncoder) AppendUint32(v uint32) {
	m.called = true
}

func (m *mockPrimitiveArrayEncoder) AppendUint16(v uint16) {
	m.called = true
}

func (m *mockPrimitiveArrayEncoder) AppendUint8(v uint8) {
	m.called = true
}

func (m *mockPrimitiveArrayEncoder) AppendUintptr(v uintptr) {
	m.called = true
}

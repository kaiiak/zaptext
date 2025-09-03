package zaptext

import (
	"testing"
	"time"

	"go.uber.org/zap/zapcore"
)

// These tests directly access the encoder methods to ensure 100% coverage
func TestTextEncoderDirectMethods(t *testing.T) {
	cfg := zapcore.EncoderConfig{
		TimeKey:    "time",
		LevelKey:   "level",
		MessageKey: "msg",
	}

	t.Run("Test AddUint method directly", func(t *testing.T) {
		encoder := NewTextEncoder(cfg).(*TextEncoder)
		encoder.AddUint("direct_uint", uint(123))

		output := encoder.buf.String()
		if output != "direct_uint=123" {
			t.Errorf("Expected 'direct_uint=123', got: %s", output)
		}
	})

	t.Run("Test OpenNamespace method directly", func(t *testing.T) {
		encoder := NewTextEncoder(cfg).(*TextEncoder)
		encoder.OpenNamespace("test_namespace") // This is a no-op, just testing it doesn't crash
		encoder.AddString("key", "value")

		output := encoder.buf.String()
		if output != "key=value" {
			t.Errorf("Expected 'key=value', got: %s", output)
		}
	})

	t.Run("Test AppendArray method directly", func(t *testing.T) {
		encoder := NewTextEncoder(cfg).(*TextEncoder)
		encoder.inArray = true // Set to simulate being inside an array

		arr := zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
			enc.AppendString("item1")
			enc.AppendInt(42)
			return nil
		})

		err := encoder.AppendArray(arr)
		if err != nil {
			t.Errorf("AppendArray failed: %v", err)
		}

		output := encoder.buf.String()
		if output != "[item1,42]" {
			t.Errorf("Expected '[item1,42]', got: %s", output)
		}
	})

	t.Run("Test AppendObject method directly", func(t *testing.T) {
		encoder := NewTextEncoder(cfg).(*TextEncoder)

		obj := zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
			enc.AddString("key", "value")
			enc.AddInt("num", 42)
			return nil
		})

		err := encoder.AppendObject(obj)
		if err != nil {
			t.Errorf("AppendObject failed: %v", err)
		}

		output := encoder.buf.String()
		if output != "key=value num=42" {
			t.Errorf("Expected 'key=value num=42', got: '%s'", output)
		}
	})

	t.Run("Test AppendReflected method directly", func(t *testing.T) {
		encoder := NewTextEncoder(cfg).(*TextEncoder)

		// AppendReflected is currently a no-op, just test it doesn't crash
		err := encoder.AppendReflected("test value")
		if err != nil {
			t.Errorf("AppendReflected failed: %v", err)
		}
	})

	t.Run("Test AppendUint and AppendUintptr directly", func(t *testing.T) {
		encoder := NewTextEncoder(cfg).(*TextEncoder)
		encoder.inArray = true // Set to simulate being inside an array

		encoder.AppendUint(uint(123))
		encoder.AppendUintptr(uintptr(0x1000))

		output := encoder.buf.String()
		if output != "123,4096" {
			t.Errorf("Expected '123,4096', got: '%s'", output)
		}
	})
}

// Test edge cases for UTF-8 handling to improve coverage
func TestTextEncoderUTFEdgeCases(t *testing.T) {
	cfg := zapcore.EncoderConfig{}

	t.Run("Test tryAddRuneError edge case", func(t *testing.T) {
		encoder := NewTextEncoder(cfg).(*TextEncoder)

		// Create a byte string with invalid UTF-8 to trigger edge cases
		invalidUTF8 := []byte{0xFF, 0xFE, 0xFD, 'h', 'e', 'l', 'l', 'o'}
		encoder.AddByteString("invalid", invalidUTF8)

		// The actual output will depend on how the encoder handles invalid UTF-8
		output := encoder.buf.String()
		if output == "" {
			t.Errorf("Expected some output for invalid UTF-8, got empty string")
		}
	})
}

func TestAddTimeAndDurationEdgeCases(t *testing.T) {
	cfg := zapcore.EncoderConfig{
		EncodeTime:     nil, // No encoder - should use default
		EncodeDuration: nil, // No encoder - should use default
	}

	t.Run("Test AddTime with no encoder", func(t *testing.T) {
		encoder := NewTextEncoder(cfg).(*TextEncoder)
		testTime := time.Date(2023, 9, 2, 10, 30, 15, 0, time.UTC)

		encoder.AddTime("timestamp", testTime)

		output := encoder.buf.String()
		if output == "" {
			t.Errorf("Expected some time output, got empty string")
		}
	})

	t.Run("Test AddDuration with no encoder", func(t *testing.T) {
		encoder := NewTextEncoder(cfg).(*TextEncoder)
		testDuration := 5 * time.Minute

		encoder.AddDuration("duration", testDuration)

		output := encoder.buf.String()
		if output == "" {
			t.Errorf("Expected some duration output, got empty string")
		}
	})

	t.Run("Test AppendTime with no encoder", func(t *testing.T) {
		encoder := NewTextEncoder(cfg).(*TextEncoder)
		testTime := time.Date(2023, 9, 2, 10, 30, 15, 0, time.UTC)

		encoder.AppendTime(testTime)

		output := encoder.buf.String()
		if output == "" {
			t.Errorf("Expected some time output, got empty string")
		}
	})

	t.Run("Test AppendDuration with no encoder", func(t *testing.T) {
		encoder := NewTextEncoder(cfg).(*TextEncoder)
		testDuration := 5 * time.Minute

		encoder.AppendDuration(testDuration)

		output := encoder.buf.String()
		if output == "" {
			t.Errorf("Expected some duration output, got empty string")
		}
	})
}

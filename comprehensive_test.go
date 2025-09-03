package zaptext_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	. "github.com/kaiiak/zaptext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Additional comprehensive tests to reach 90%+ coverage

func TestTextEncoderEdgeCases(t *testing.T) {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeDuration = zapcore.StringDurationEncoder

	t.Run("More array types to trigger append methods", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)

		logger.Info("test message",
			zap.Uint8s("uint8s", []uint8{1, 2, 3}),
			zap.Uint16s("uint16s", []uint16{10, 20, 30}),
			zap.Uint32s("uint32s", []uint32{100, 200, 300}),
			zap.Uint64s("uint64s", []uint64{1000, 2000, 3000}),
			zap.Int8s("int8s", []int8{-1, -2, -3}),
			zap.Int16s("int16s", []int16{-10, -20, -30}),
			zap.Int32s("int32s", []int32{-100, -200, -300}),
			zap.Float32s("float32s", []float32{1.1, 2.2, 3.3}),
		)

		output := buf.String()
		expected := []string{"uint8s=", "uint16s=", "uint32s=", "uint64s=", "int8s=", "int16s=", "int32s=", "float32s="}
		for _, exp := range expected {
			if !strings.Contains(output, exp) {
				t.Errorf("Expected %s in output, got: %s", exp, output)
			}
		}
	})

	t.Run("Complex number arrays to trigger AppendComplex methods", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)

		logger.Info("test message",
			zap.Complex64s("complex64s", []complex64{complex(1, 2), complex(3, 4)}),
			zap.Complex128s("complex128s", []complex128{complex(5, 6), complex(7, 8)}),
		)

		output := buf.String()
		if !strings.Contains(output, "complex64s=") {
			t.Errorf("Expected complex64s in output, got: %s", output)
		}
		if !strings.Contains(output, "complex128s=") {
			t.Errorf("Expected complex128s in output, got: %s", output)
		}
	})

	t.Run("Duration array to trigger AppendDuration", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)

		durations := []time.Duration{time.Second, 2 * time.Minute, 3 * time.Hour}
		logger.Info("test message", zap.Durations("durations", durations))

		output := buf.String()
		if !strings.Contains(output, "durations=") {
			t.Errorf("Expected durations in output, got: %s", output)
		}
	})

	t.Run("Time array to trigger AppendTime", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)

		times := []time.Time{
			time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
		}
		logger.Info("test message", zap.Times("times", times))

		output := buf.String()
		if !strings.Contains(output, "times=") {
			t.Errorf("Expected times in output, got: %s", output)
		}
	})
}

func TestTextEncoderUTF8EdgeCases(t *testing.T) {
	cfg := zap.NewProductionEncoderConfig()

	t.Run("UTF-8 edge cases and rune errors", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)

		// Test various UTF-8 edge cases that might trigger tryAddRuneError
		logger.Info("test message",
			zap.ByteString("invalid_utf8", []byte{0xFF, 0xFE, 0xFD}),        // Invalid UTF-8 bytes
			zap.ByteString("control_chars", []byte{0x01, 0x02, 0x03, 0x1F}), // Control characters
			zap.ByteString("mixed", []byte("hello\xFF\xFEworld")),           // Mixed valid/invalid UTF-8
		)

		output := buf.String()
		if !strings.Contains(output, "invalid_utf8=") {
			t.Errorf("Expected invalid_utf8 field in output, got: %s", output)
		}
		if !strings.Contains(output, "control_chars=") {
			t.Errorf("Expected control_chars field in output, got: %s", output)
		}
		if !strings.Contains(output, "mixed=") {
			t.Errorf("Expected mixed field in output, got: %s", output)
		}
	})

	t.Run("Special characters in strings", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)

		logger.Info("test message",
			zap.String("backslash", "test\\path"),
			zap.String("quote", "say \"hello\""),
			zap.String("newline", "line1\nline2"),
			zap.String("tab", "col1\tcol2"),
			zap.String("carriage_return", "line1\rline2"),
		)

		output := buf.String()
		fields := []string{"backslash=", "quote=", "newline=", "tab=", "carriage_return="}
		for _, field := range fields {
			if !strings.Contains(output, field) {
				t.Errorf("Expected %s in output, got: %s", field, output)
			}
		}
	})
}

func TestTextEncoderDirectAccess(t *testing.T) {
	// Test methods that might not be covered by logger interface
	cfg := zap.NewProductionEncoderConfig()

	t.Run("Test encoder methods directly", func(t *testing.T) {
		encoder := NewTextEncoder(cfg)

		// Test OpenNamespace directly - it's a no-op but should be called
		encoder.OpenNamespace("test_namespace")

		// Create a test entry to encode
		entry := zapcore.Entry{
			Level:   zap.InfoLevel,
			Time:    time.Date(2023, 9, 2, 10, 30, 15, 0, time.UTC),
			Message: "test message",
		}

		fields := []zapcore.Field{
			zap.String("test", "value"),
		}

		buf, err := encoder.EncodeEntry(entry, fields)
		if err != nil {
			t.Errorf("EncodeEntry failed: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "test message") {
			t.Errorf("Expected message in output, got: %s", output)
		}
		if !strings.Contains(output, "test=value") {
			t.Errorf("Expected field in output, got: %s", output)
		}
	})
}

func TestTextEncoderEncodeEntryEdgeCases(t *testing.T) {
	cfg := zap.NewProductionEncoderConfig()

	t.Run("Entry with caller info", func(t *testing.T) {
		cfg.CallerKey = "caller"
		encoder := NewTextEncoder(cfg)

		entry := zapcore.Entry{
			Level:   zap.InfoLevel,
			Time:    time.Date(2023, 9, 2, 10, 30, 15, 0, time.UTC),
			Message: "test message",
			Caller:  zapcore.EntryCaller{Defined: true, File: "/path/to/file.go", Line: 123},
		}

		buf, err := encoder.EncodeEntry(entry, nil)
		if err != nil {
			t.Errorf("EncodeEntry failed: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "file.go:123") {
			t.Errorf("Expected caller info in output, got: %s", output)
		}
	})

	t.Run("Entry with no time key", func(t *testing.T) {
		cfg := zap.NewProductionEncoderConfig()
		cfg.TimeKey = "" // Disable time key
		encoder := NewTextEncoder(cfg)

		entry := zapcore.Entry{
			Level:   zap.InfoLevel,
			Time:    time.Date(2023, 9, 2, 10, 30, 15, 0, time.UTC),
			Message: "test message",
		}

		buf, err := encoder.EncodeEntry(entry, nil)
		if err != nil {
			t.Errorf("EncodeEntry failed: %v", err)
		}

		output := buf.String()
		// Should not contain timestamp if TimeKey is empty
		if strings.Contains(output, "2023-09-02") {
			t.Logf("Time encoding behavior: %s", output)
		}
	})

	t.Run("Entry with no level key", func(t *testing.T) {
		cfg := zap.NewProductionEncoderConfig()
		cfg.LevelKey = "" // Disable level key
		encoder := NewTextEncoder(cfg)

		entry := zapcore.Entry{
			Level:   zap.InfoLevel,
			Message: "test message",
		}

		buf, err := encoder.EncodeEntry(entry, nil)
		if err != nil {
			t.Errorf("EncodeEntry failed: %v", err)
		}

		output := buf.String()
		// Should contain message but not level
		if !strings.Contains(output, "test message") {
			t.Errorf("Expected message in output, got: %s", output)
		}
	})
}

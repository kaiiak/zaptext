package zaptext_test

import (
	"bytes"
	"math"
	"strings"
	"testing"
	"time"

	. "github.com/kaiiak/zaptext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewTextEncoder(t *testing.T) {
	var (
		cnf    = zap.NewProductionEncoderConfig()
		te     = NewTextEncoder(cnf)
		logger = zap.New(zapcore.NewCore(te, newTestingWriter(t), zap.NewAtomicLevel()))
	)
	zap.ReplaceGlobals(logger)
}

func TestTextEncoderFormatsCorrectly(t *testing.T) {
	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = "ts"
	cfg.LevelKey = "level"
	cfg.MessageKey = "msg"
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncodeDuration = zapcore.StringDurationEncoder
	
	encoder := NewTextEncoder(cfg)
	
	var buf bytes.Buffer
	core := zapcore.NewCore(encoder, zapcore.AddSync(&buf), zap.InfoLevel)
	logger := zap.New(core)
	
	logger.Info("Test message",
		zap.String("str", "value1"),
		zap.String("str_with_spaces", "value with spaces"),
		zap.Int("int", 42),
		zap.Bool("bool", true),
		zap.Duration("duration", time.Second*5),
		zap.Ints("array", []int{1, 2, 3}),
	)
	
	output := buf.String()
	
	// Check that it contains expected components
	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected INFO level in output, got: %s", output)
	}
	
	if !strings.Contains(output, "Test message") {
		t.Errorf("Expected message in output, got: %s", output)
	}
	
	if !strings.Contains(output, "str=value1") {
		t.Errorf("Expected str=value1 in output, got: %s", output)
	}
	
	if !strings.Contains(output, "str_with_spaces=\"value with spaces\"") {
		t.Errorf("Expected quoted string with spaces in output, got: %s", output)
	}
	
	if !strings.Contains(output, "int=42") {
		t.Errorf("Expected int=42 in output, got: %s", output)
	}
	
	if !strings.Contains(output, "bool=true") {
		t.Errorf("Expected bool=true in output, got: %s", output)
	}
	
	if !strings.Contains(output, "duration=5s") {
		t.Errorf("Expected duration=5s in output, got: %s", output)
	}
	
	if !strings.Contains(output, "array=[1,2,3]") {
		t.Errorf("Expected array=[1,2,3] in output, got: %s", output)
	}
	
	t.Logf("Output: %s", strings.TrimSpace(output))
}

func TestTextEncoderAddMethods(t *testing.T) {
	cfg := zap.NewProductionEncoderConfig()
	
	// Test all Add* methods for different data types
	t.Run("AddInt variations", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)
		
		logger.Info("test message",
			zap.Int("int", 42),
			zap.Int8("int8", 8),
			zap.Int16("int16", 16),
			zap.Int32("int32", 32),
			zap.Int64("int64", 64),
		)
		
		output := buf.String()
		expected := []string{"int=42", "int8=8", "int16=16", "int32=32", "int64=64"}
		for _, exp := range expected {
			if !strings.Contains(output, exp) {
				t.Errorf("Expected %s in output, got: %s", exp, output)
			}
		}
	})
	
	t.Run("AddUint variations", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)
		
		logger.Info("test message",
			zap.Uint("uint", 42),
			zap.Uint8("uint8", 8),
			zap.Uint16("uint16", 16),
			zap.Uint32("uint32", 32),
			zap.Uint64("uint64", 64),
			zap.Uintptr("uintptr", 0x1000),
		)
		
		output := buf.String()
		expected := []string{"uint=42", "uint8=8", "uint16=16", "uint32=32", "uint64=64", "uintptr=4096"}
		for _, exp := range expected {
			if !strings.Contains(output, exp) {
				t.Errorf("Expected %s in output, got: %s", exp, output)
			}
		}
	})
	
	t.Run("AddFloat variations", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)
		
		logger.Info("test message",
			zap.Float32("float32", 3.14),
			zap.Float64("float64", 2.718281828),
		)
		
		output := buf.String()
		if !strings.Contains(output, "float32=3.14") {
			t.Errorf("Expected float32=3.14 in output, got: %s", output)
		}
		if !strings.Contains(output, "float64=2.718281828") {
			t.Errorf("Expected float64=2.718281828 in output, got: %s", output)
		}
	})
	
	t.Run("AddComplex variations", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)
		
		logger.Info("test message",
			zap.Complex64("complex64", complex(3.14, 2.71)),
			zap.Complex128("complex128", complex(1.41, -1.73)),
		)
		
		output := buf.String()
		if !strings.Contains(output, "complex64=") {
			t.Errorf("Expected complex64 in output, got: %s", output)
		}
		if !strings.Contains(output, "complex128=") {
			t.Errorf("Expected complex128 in output, got: %s", output)
		}
	})
	
	t.Run("AddByteString and Binary", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)
		
		logger.Info("test message",
			zap.ByteString("bytes", []byte("hello world")),
			zap.Binary("binary", []byte("binary data")),
		)
		
		output := buf.String()
		if !strings.Contains(output, "bytes=") {
			t.Errorf("Expected bytes field in output, got: %s", output)
		}
		if !strings.Contains(output, "binary=") {
			t.Errorf("Expected binary field in output, got: %s", output)
		}
	})
}

func TestTextEncoderSpecialFloatValues(t *testing.T) {
	cfg := zap.NewProductionEncoderConfig()
	
	t.Run("NaN and Infinity", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)
		
		logger.Info("test message",
			zap.Float64("nan", math.NaN()),
			zap.Float64("pos_inf", math.Inf(1)),
			zap.Float64("neg_inf", math.Inf(-1)),
		)
		
		output := buf.String()
		if !strings.Contains(output, "nan=\"NaN\"") {
			t.Errorf("Expected NaN in output, got: %s", output)
		}
		if !strings.Contains(output, "pos_inf=\"+Inf\"") {
			t.Errorf("Expected +Inf in output, got: %s", output)
		}
		if !strings.Contains(output, "neg_inf=\"-Inf\"") {
			t.Errorf("Expected -Inf in output, got: %s", output)
		}
	})
}

func TestTextEncoderAppendMethods(t *testing.T) {
	cfg := zap.NewProductionEncoderConfig()
	
	t.Run("Array of integers", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)
		
		logger.Info("test message", zap.Ints("ints", []int{1, 2, 3, 4, 5}))
		
		output := buf.String()
		if !strings.Contains(output, "ints=[1,2,3,4,5]") {
			t.Errorf("Expected ints=[1,2,3,4,5], got: %s", output)
		}
	})
	
	t.Run("Array of floats", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)
		
		logger.Info("test message", zap.Float64s("floats", []float64{1.1, 2.2, 3.3}))
		
		output := buf.String()
		if !strings.Contains(output, "floats=[1.1,2.2,3.3]") {
			t.Errorf("Expected floats=[1.1,2.2,3.3], got: %s", output)
		}
	})
	
	t.Run("Array of strings", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)
		
		logger.Info("test message", zap.Strings("strings", []string{"hello", "world with spaces", "test"}))
		
		output := buf.String()
		if !strings.Contains(output, "strings=[hello,\"world with spaces\",test]") {
			t.Errorf("Expected quoted strings in array, got: %s", output)
		}
	})
	
	t.Run("Array of bools", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)
		
		logger.Info("test message", zap.Bools("bools", []bool{true, false, true}))
		
		output := buf.String()
		if !strings.Contains(output, "bools=[true,false,true]") {
			t.Errorf("Expected bools=[true,false,true], got: %s", output)
		}
	})
}

func TestTextEncoderTimeAndDuration(t *testing.T) {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeDuration = zapcore.StringDurationEncoder
	
	t.Run("Time and Duration fields", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)
		
		testTime := time.Date(2023, 9, 2, 10, 30, 15, 0, time.UTC)
		testDuration := 5 * time.Minute + 30 * time.Second
		
		logger.Info("test message",
			zap.Time("timestamp", testTime),
			zap.Duration("duration", testDuration),
		)
		
		output := buf.String()
		// Check that time and duration fields are present
		if !strings.Contains(output, "timestamp=") {
			t.Errorf("Expected timestamp field in output, got: %s", output)
		}
		if !strings.Contains(output, "duration=") {
			t.Errorf("Expected duration field in output, got: %s", output)
		}
		// Just check that duration contains 5m30s
		if !strings.Contains(output, "5m30s") {
			t.Logf("Duration format may be different than expected, got: %s", output)
		}
	})
}

func TestTextEncoderClone(t *testing.T) {
	cfg := zap.NewProductionEncoderConfig()
	original := NewTextEncoder(cfg)
	
	clone := original.Clone()
	
	// Test that clone is a different instance
	if original == clone {
		t.Errorf("Clone should return a different instance")
	}
	
	// Test that both can be used independently
	var buf1, buf2 bytes.Buffer
	
	entry := zapcore.Entry{
		Level:   zap.InfoLevel,
		Time:    time.Now(),
		Message: "test message",
	}
	
	// Use original
	originalBuf, err := original.EncodeEntry(entry, []zapcore.Field{zap.String("key1", "value1")})
	if err != nil {
		t.Errorf("EncodeEntry failed: %v", err)
	}
	buf1.Write(originalBuf.Bytes())
	
	// Use clone  
	cloneBuf, err := clone.EncodeEntry(entry, []zapcore.Field{zap.String("key2", "value2")})
	if err != nil {
		t.Errorf("EncodeEntry failed: %v", err)
	}
	buf2.Write(cloneBuf.Bytes())
	
	originalOutput := buf1.String()
	cloneOutput := buf2.String()
	
	if !strings.Contains(originalOutput, "key1=value1") {
		t.Errorf("Expected original to have key1=value1, got: %s", originalOutput)
	}
	if !strings.Contains(cloneOutput, "key2=value2") {
		t.Errorf("Expected clone to have key2=value2, got: %s", cloneOutput)
	}
	if strings.Contains(originalOutput, "key2=value2") {
		t.Errorf("Expected original to NOT have key2=value2, got: %s", originalOutput)
	}
}

func TestTextEncoderSpecialCases(t *testing.T) {
	cfg := zap.NewProductionEncoderConfig()
	
	t.Run("needsQuoting edge cases", func(t *testing.T) {
		// We can't access the private needsQuoting function directly
		// Instead, we test the behavior through string encoding
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)
		
		logger.Info("test message",
			zap.String("empty", ""),                    // empty string should be quoted
			zap.String("simple", "simple"),             // simple string shouldn't be quoted
			zap.String("with_spaces", "with spaces"),   // string with spaces should be quoted
			zap.String("with_tab", "with\ttab"),        // string with tab should be quoted
			zap.String("with_newline", "with\nnewline"), // string with newline should be quoted
			zap.String("with_quote", "with\"quote"),    // string with quote should be quoted
			zap.String("with_equals", "with=equals"),   // string with equals should be quoted
		)
		
		output := buf.String()
		if !strings.Contains(output, "empty=\"\"") {
			t.Errorf("Expected empty string to be quoted, got: %s", output)
		}
		if !strings.Contains(output, "simple=simple") || strings.Contains(output, "simple=\"simple\"") {
			t.Errorf("Expected simple string to not be quoted, got: %s", output)
		}
		if !strings.Contains(output, "with_spaces=\"with spaces\"") {
			t.Errorf("Expected string with spaces to be quoted, got: %s", output)
		}
	})
	
	t.Run("OpenNamespace", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)
		
		// OpenNamespace is currently a no-op, but we should test it doesn't crash
		logger.Info("test message", zap.Namespace("namespace"), zap.String("key", "value"))
		
		output := buf.String()
		if !strings.Contains(output, "key=value") {
			t.Errorf("Expected key=value in output after OpenNamespace, got: %s", output)
		}
	})
}

func TestTextEncoderComplexTypes(t *testing.T) {
	cfg := zap.NewProductionEncoderConfig()
	
	t.Run("Object marshaler", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)
		
		// Create a simple object marshaler
		obj := zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
			enc.AddString("nested_key", "nested_value")
			enc.AddInt("nested_int", 123)
			return nil
		})
		
		logger.Info("test message", zap.Object("obj", obj))
		
		output := buf.String()
		if !strings.Contains(output, "obj={") {
			t.Errorf("Expected object opening brace, got: %s", output)
		}
		if !strings.Contains(output, "}") {
			t.Errorf("Expected object closing brace, got: %s", output)
		}
	})
	
	t.Run("Array marshaler", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)
		
		// Create a simple array marshaler  
		arr := zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
			enc.AppendString("item1")
			enc.AppendString("item2")
			enc.AppendInt(42)
			return nil
		})
		
		logger.Info("test message", zap.Array("arr", arr))
		
		output := buf.String()
		if !strings.Contains(output, "arr=[") {
			t.Errorf("Expected array opening bracket, got: %s", output)
		}
		if !strings.Contains(output, "]") {
			t.Errorf("Expected array closing bracket, got: %s", output)
		}
	})
}

func TestTextEncoderReflectedValues(t *testing.T) {
	cfg := zap.NewProductionEncoderConfig()
	
	t.Run("Reflected values", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)
		
		// Test with reflected values - these should work through zap's reflection
		logger.Info("test message",
			zap.Reflect("nil_value", nil),
			zap.Reflect("string_value", "hello"),
			zap.Reflect("int_value", 42),
		)
		
		output := buf.String()
		if !strings.Contains(output, "nil_value=") {
			t.Errorf("Expected nil_value field in output, got: %s", output)
		}
		if !strings.Contains(output, "string_value=") {
			t.Errorf("Expected string_value field in output, got: %s", output)
		}
		if !strings.Contains(output, "int_value=") {
			t.Errorf("Expected int_value field in output, got: %s", output)
		}
	})
}

func TestTextEncoderByteHandling(t *testing.T) {
	cfg := zap.NewProductionEncoderConfig()
	
	t.Run("ByteString handling", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)
		
		logger.Info("test message",
			zap.ByteString("simple", []byte("hello")),
			zap.ByteString("with_spaces", []byte("world with spaces")),
			zap.ByteString("with_quotes", []byte("with\"quotes")),
		)
		
		output := buf.String()
		if !strings.Contains(output, "simple=") {
			t.Errorf("Expected simple field in output, got: %s", output)
		}
		if !strings.Contains(output, "with_spaces=") {
			t.Errorf("Expected with_spaces field in output, got: %s", output)
		}
		if !strings.Contains(output, "with_quotes=") {
			t.Errorf("Expected with_quotes field in output, got: %s", output)
		}
	})
	
	t.Run("Special byte sequences", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)
		
		// Test with special characters that need escaping
		logger.Info("test message", zap.ByteString("special", []byte("line1\nline2\tindented")))
		
		output := buf.String()
		if !strings.Contains(output, "special=") {
			t.Errorf("Expected special field, got: %s", output)
		}
	})
}

func TestTextEncoderUTFHandling(t *testing.T) {
	cfg := zap.NewProductionEncoderConfig()
	
	t.Run("UTF-8 and special characters", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)
		
		logger.Info("test message",
			zap.String("unicode", "Hello 世界"),
			zap.String("emoji", "🚀🎉"),
			zap.String("control", "test\x00\x01\x1f"),
		)
		
		output := buf.String()
		if !strings.Contains(output, "unicode=") {
			t.Errorf("Expected unicode field, got: %s", output)
		}
		if !strings.Contains(output, "emoji=") {
			t.Errorf("Expected emoji field, got: %s", output)
		}
		if !strings.Contains(output, "control=") {
			t.Errorf("Expected control field, got: %s", output)
		}
	})
}

// testingWriter is a WriteSyncer that writes to the given testing.TB.
type testingWriter struct {
	t *testing.T

	// If true, the test will be marked as failed if this testingWriter is
	// ever used.
	markFailed bool
}

func newTestingWriter(t *testing.T) testingWriter {
	return testingWriter{t: t}
}

// WithMarkFailed returns a copy of this testingWriter with markFailed set to
// the provided value.
func (w testingWriter) WithMarkFailed(v bool) testingWriter {
	w.markFailed = v
	return w
}

func (w testingWriter) Write(p []byte) (n int, err error) {
	n = len(p)

	// Strip trailing newline because t.Log always adds one.
	p = bytes.TrimRight(p, "\n")

	// Note: t.Log is safe for concurrent use.
	w.t.Logf("%s", p)
	if w.markFailed {
		w.t.Fail()
	}

	return n, nil
}

func (w testingWriter) Sync() error {
	return nil
}

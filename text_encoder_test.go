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

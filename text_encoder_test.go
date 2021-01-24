package zaptext_test

import (
	"bytes"
	"testing"

	. "github.com/kaiiak/zaptext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewTextEncoder(t *testing.T) {
	var (
		cnf    = zap.NewProductionEncoderConfig()
		te     = NewTextEncoder(cnf)
		logger = zap.New(zapcore.NewCore(te, nil, zap.NewAtomicLevel()))
	)
	zap.ReplaceGlobals(logger)
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

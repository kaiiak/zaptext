package zaptext_test

import (
	"bytes"
	"strings"
	"testing"

	. "github.com/kaiiak/zaptext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestTextEncoderAddUint(t *testing.T) {
	cfg := zap.NewProductionEncoderConfig()

	t.Run("AddUint specific test", func(t *testing.T) {
		var buf bytes.Buffer
		core := zapcore.NewCore(NewTextEncoder(cfg), zapcore.AddSync(&buf), zap.InfoLevel)
		logger := zap.New(core)

		// Test specifically with Uint to ensure it's covered
		logger.Info("test message", zap.Uint("test_uint", uint(42)))

		output := buf.String()
		if !strings.Contains(output, "test_uint=42") {
			t.Errorf("Expected test_uint=42 in output, got: %s", output)
		}
	})
}

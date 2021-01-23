package zaptext

import (
	"time"

	"go.uber.org/zap/zapcore"
)

// CustomTimeEncoderFactory return a zapcore.TimeEncoder format time with custom layout
func CustomTimeEncoderFactory(layout string) zapcore.TimeEncoder {
	return func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format(layout))
	}
}

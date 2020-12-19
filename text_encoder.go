package zaptext

import (
	"sync"
	"time"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type (
	textEncoder struct {
		*zapcore.EncoderConfig
		buf   *buffer.Buffer
		level zapcore.LevelEncoder
	}
)

var (
	textPool = sync.Pool{New: func() interface{} {
		return &textEncoder{}
	}}
)

var _ zapcore.Encoder = (*textEncoder)(nil)

func NewTextEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return &textEncoder{EncoderConfig: &cfg}
}

// Clone copies the encoder, ensuring that adding fields to the copy doesn't
// affect the original.
func (t *textEncoder) Clone() zapcore.Encoder { return t }

// EncodeEntry encodes an entry and fields, along with any accumulated
// context, into a byte buffer and returns it. Any fields that are empty,
// including fields on the `Entry` type, should be omitted.
func (t *textEncoder) EncodeEntry(zapcore.Entry, []zapcore.Field) (buf *buffer.Buffer, err error) {
	return
}

// Logging-specific marshalers.
func (t *textEncoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) (err error)   { return }
func (t *textEncoder) AddObject(key string, marshaler zapcore.ObjectMarshaler) (err error) { return }

// Built-in types.
// for arbitrary bytes
func (t *textEncoder) AddBinary(key string, value []byte) {}

// for UTF-8 encoded bytes
func (t *textEncoder) AddByteString(key string, value []byte)      {}
func (t *textEncoder) AddBool(key string, value bool)              {}
func (t *textEncoder) AddComplex128(key string, value complex128)  {}
func (t *textEncoder) AddComplex64(key string, value complex64)    {}
func (t *textEncoder) AddDuration(key string, value time.Duration) {}
func (t *textEncoder) AddFloat64(key string, value float64)        {}
func (t *textEncoder) AddFloat32(key string, value float32)        {}
func (t *textEncoder) AddInt(key string, value int)                {}
func (t *textEncoder) AddInt64(key string, value int64)            {}
func (t *textEncoder) AddInt32(key string, value int32)            {}
func (t *textEncoder) AddInt16(key string, value int16)            {}
func (t *textEncoder) AddInt8(key string, value int8)              {}
func (t *textEncoder) AddString(key, value string)                 {}
func (t *textEncoder) AddTime(key string, value time.Time)         {}
func (t *textEncoder) AddUint(key string, value uint)              {}
func (t *textEncoder) AddUint64(key string, value uint64)          {}
func (t *textEncoder) AddUint32(key string, value uint32)          {}
func (t *textEncoder) AddUint16(key string, value uint16)          {}
func (t *textEncoder) AddUint8(key string, value uint8)            {}
func (t *textEncoder) AddUintptr(key string, value uintptr)        {}

// AddReflected uses reflection to serialize arbitrary objects, so it can be
// slow and allocation-heavy.
func (t *textEncoder) AddReflected(key string, value interface{}) (err error) { return }

// OpenNamespace opens an isolated namespace where all subsequent fields will
// be added. Applications can use namespaces to prevent key collisions when
// injecting loggers into sub-components or third-party libraries.
func (t *textEncoder) OpenNamespace(key string) {}

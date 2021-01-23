package zaptext

import (
	"encoding/base64"
	"math"
	"sync"
	"time"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type (
	textEncoder struct {
		*zapcore.EncoderConfig
		buf            *buffer.Buffer
		openNamespaces int
		spaced         bool

		// for encoding generic values by reflection
		reflectBuf *buffer.Buffer
	}
)

var (
	textpool = sync.Pool{New: func() interface{} {
		return &textEncoder{}
	}}
	buffpoll = buffer.NewPool()
)

var _ zapcore.Encoder = (*textEncoder)(nil)
var _ zapcore.ArrayEncoder = (*textEncoder)(nil)

func NewTextEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return &textEncoder{EncoderConfig: &cfg, buf: buffpoll.Get()}
}

func (t *textEncoder) addKey(key string) {

}

func (t *textEncoder) addElementSeparator() {
	last := t.buf.Len() - 1
	if last < 0 {
		return
	}
	switch t.buf.Bytes()[last] {
	case '{', '[', ':', ',', ' ':
		return
	default:
		t.buf.AppendByte(',')
		if t.spaced {
			t.buf.AppendByte(' ')
		}
	}
}

func (t *textEncoder) appendFloat(val float64, bitSize int) {
	t.addElementSeparator()
	switch {
	case math.IsNaN(val):
		t.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		t.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		t.buf.AppendString(`"-Inf"`)
	default:
		t.buf.AppendFloat(val, bitSize)
	}
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

func (t *textEncoder) AddComplex64(key string, value complex64) {
	t.AddComplex128(key, complex128(value))
}
func (t *textEncoder) AddFloat32(key string, value float32) {
	t.AddFloat64(key, float64(value))
}
func (t *textEncoder) AddInt(key string, value int) {
	t.AddInt64(key, int64(value))
}
func (t *textEncoder) AddInt32(key string, value int32) {
	t.AddInt64(key, int64(value))
}
func (t *textEncoder) AddInt16(key string, value int16) {
	t.AddInt64(key, int64(value))
}
func (t *textEncoder) AddInt8(key string, value int8) {
	t.AddInt64(key, int64(value))
}
func (t *textEncoder) AddUint32(key string, value uint32) {
	t.AddUint64(key, uint64(value))
}
func (t *textEncoder) AddUint(key string, value uint) {
	t.AddUint64(key, uint64(value))
}
func (t *textEncoder) AddUint16(key string, value uint16) {
	t.AddUint64(key, uint64(value))
}
func (t *textEncoder) AddUint8(key string, value uint8) {
	t.AddUint64(key, uint64(value))
}
func (t *textEncoder) AddUintptr(key string, value uintptr) {
	t.AddUint64(key, uint64(value))
}

// AddReflected uses reflection to serialize arbitrary objects, so it can be
// slow and allocation-heavy.
func (t *textEncoder) AddReflected(key string, value interface{}) (err error) { return }

// OpenNamespace opens an isolated namespace where all subsequent fields will
// be added. Applications can use namespaces to prevent key collisions when
// injecting loggers into sub-components or third-party libraries.
func (t *textEncoder) OpenNamespace(key string) {}

// Built-in types.
// for arbitrary bytes
func (t *textEncoder) AddBinary(key string, value []byte) {
	t.AddString(key, base64.StdEncoding.EncodeToString(value))
}
func (t *textEncoder) AddDuration(key string, value time.Duration) {
	cur := t.buf.Len()
	if e := t.EncodeDuration; e != nil {
		e(value, t)
	}
	if cur == t.buf.Len() {
		t.AppendInt64(int64(value))
	}
}
func (t *textEncoder) AddComplex128(key string, value complex128) {
	t.addElementSeparator()
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(value)), float64(imag(value))
	t.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	t.buf.AppendFloat(r, 64)
	t.buf.AppendByte('+')
	t.buf.AppendFloat(i, 64)
	t.buf.AppendByte('i')
	t.buf.AppendByte('"')
}
func (t *textEncoder) AddByteString(key string, value []byte) {
	t.addKey(key)
	t.AppendByteString(value)
}
func (t *textEncoder) AddFloat64(key string, value float64) {
	t.addKey(key)
	t.appendFloat(value, 64)
}
func (t *textEncoder) AddTime(key string, value time.Time) {
	t.addKey(key)
	t.buf.AppendTime(value, time.RFC3339)
}
func (t *textEncoder) AddUint64(key string, value uint64) {
	t.addKey(key)
	t.buf.AppendUint(value)
}
func (t *textEncoder) AddInt64(key string, value int64) {
	t.addKey(key)
	t.buf.AppendInt(value)
}
func (t *textEncoder) AddBool(key string, value bool) {
	t.addKey(key)
	t.buf.AppendBool(value)
}
func (t *textEncoder) AddString(key, value string) {
	t.addKey(key)
	t.buf.AppendString(value)
}

// ArrayEncoder

// Time-related types.
func (t *textEncoder) AppendDuration(time.Duration) {}
func (t *textEncoder) AppendTime(time.Time)         {}

// Logging-specific marshalers.{}
func (t *textEncoder) AppendArray(zapcore.ArrayMarshaler) (err error)   { return }
func (t *textEncoder) AppendObject(zapcore.ObjectMarshaler) (err error) { return }

// AppendReflected uses reflection to serialize arbitrary objects, so it's{}
// slow and allocation-heavy.{}
func (t *textEncoder) AppendReflected(value interface{}) (err error) { return }

func (t *textEncoder) AppendBool(value bool) {
	t.addElementSeparator()
	t.buf.AppendBool(value)
}
func (t *textEncoder) AppendByteString(value []byte)     {} // for UTF-8 encoded bytes
func (t *textEncoder) AppendComplex128(value complex128) {}
func (t *textEncoder) AppendUint64(value uint64)         {}
func (t *textEncoder) AppendString(value string)         {}
func (t *textEncoder) AppendInt64(value int64)           {}
func (t *textEncoder) AppendFloat64(value float64)       {}

func (t *textEncoder) AppendComplex64(value complex64) {
	t.AppendComplex128(complex128(value))
}
func (t *textEncoder) AppendFloat32(value float32) {
	t.AppendFloat64(float64(value))
}
func (t *textEncoder) AppendInt32(value int32) {
	t.AppendInt64(int64(value))
}
func (t *textEncoder) AppendInt16(value int16) {
	t.AppendInt64(int64(value))
}
func (t *textEncoder) AppendInt(value int) {
	t.AppendInt64(int64(value))
}
func (t *textEncoder) AppendInt8(value int8) {
	t.AppendInt64(int64(value))
}
func (t *textEncoder) AppendUint(value uint) {
	t.AppendUint64(uint64(value))
}
func (t *textEncoder) AppendUint32(value uint32) {
	t.AppendUint64(uint64(value))
}
func (t *textEncoder) AppendUint16(value uint16) {
	t.AppendUint64(uint64(value))
}
func (t *textEncoder) AppendUint8(value uint8) {
	t.AppendUint64(uint64(value))
}
func (t *textEncoder) AppendUintptr(value uintptr) {
	t.AppendUint64(uint64(value))
}

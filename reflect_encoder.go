package zaptext

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// Object pool for ReflectEncoder to optimize memory allocation
	reflectEncoderPool = sync.Pool{
		New: func() any {
			return &ReflectEncoder{}
		},
	}
	
	// Buffer pool for temporary buffers
	bufferPool = sync.Pool{
		New: func() any {
			return &bytes.Buffer{}
		},
	}
)

type ReflectEncoder struct {
	w          io.Writer
	err        error
	escapeHTML bool
	depth      int

	buf *bytes.Buffer
}

func NewReflectEncoder(w io.Writer) *ReflectEncoder {
	enc := reflectEncoderPool.Get().(*ReflectEncoder)
	enc.w = w
	enc.escapeHTML = true
	enc.depth = 0
	enc.err = nil
	
	if enc.buf == nil {
		enc.buf = bufferPool.Get().(*bytes.Buffer)
	}
	enc.buf.Reset()
	
	return enc
}

// Release returns the encoder back to the pool for reuse
func (enc *ReflectEncoder) Release() {
	if enc.buf != nil {
		enc.buf.Reset()
		bufferPool.Put(enc.buf)
		enc.buf = nil
	}
	enc.w = nil
	enc.err = nil
	enc.depth = 0
	reflectEncoderPool.Put(enc)
}

func (enc *ReflectEncoder) Encode(obj any) error {
	if enc.err != nil {
		return enc.err
	}
	
	if enc.buf == nil {
		enc.buf = bufferPool.Get().(*bytes.Buffer)
	}
	
	// Reset buffer for fresh encoding
	enc.buf.Reset()
	
	// Encode the object using reflection
	if err := enc.encodeValue(reflect.ValueOf(obj)); err != nil {
		enc.err = err
		return err
	}
	
	// Write the buffer to the writer
	if _, err := enc.w.Write(enc.buf.Bytes()); err != nil {
		enc.err = err
		return err
	}
	
	return nil
}

// Maximum depth to prevent infinite recursion
const maxDepth = 32

func (enc *ReflectEncoder) encodeValue(v reflect.Value) error {
	// Prevent infinite recursion
	if enc.depth > maxDepth {
		return fmt.Errorf("maximum encoding depth exceeded: %d", maxDepth)
	}
	
	// Handle invalid values
	if !v.IsValid() {
		enc.buf.WriteString("null")
		return nil
	}
	
	// Handle pointers and interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			enc.buf.WriteString("null")
			return nil
		}
		v = v.Elem()
	}
	
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			enc.buf.WriteString("true")
		} else {
			enc.buf.WriteString("false")
		}
		
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		enc.buf.WriteString(strconv.FormatInt(v.Int(), 10))
		
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		enc.buf.WriteString(strconv.FormatUint(v.Uint(), 10))
		
	case reflect.Float32, reflect.Float64:
		f := v.Float()
		enc.buf.WriteString(strconv.FormatFloat(f, 'g', -1, v.Type().Bits()))
		
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		enc.buf.WriteString(fmt.Sprintf("(%g+%gi)", real(c), imag(c)))
		
	case reflect.String:
		enc.buf.WriteByte('"')
		enc.writeEscapedString(v.String())
		enc.buf.WriteByte('"')
		
	case reflect.Array, reflect.Slice:
		if v.IsNil() {
			enc.buf.WriteString("null")
			return nil
		}
		return enc.encodeArray(v)
		
	case reflect.Map:
		if v.IsNil() {
			enc.buf.WriteString("null")
			return nil
		}
		return enc.encodeMap(v)
		
	case reflect.Struct:
		return enc.encodeStruct(v)
		
	default:
		// For types we don't handle specifically, try to convert to string
		enc.encodeString(fmt.Sprintf("%v", v.Interface()))
	}
	
	return nil
}

func (enc *ReflectEncoder) encodeString(s string) {
	if needsQuoting(s) {
		enc.buf.WriteByte('"')
		enc.writeEscapedString(s)
		enc.buf.WriteByte('"')
	} else {
		enc.buf.WriteString(s)
	}
}

func (enc *ReflectEncoder) writeEscapedString(s string) {
	for _, r := range s {
		switch r {
		case '"':
			enc.buf.WriteString(`\"`)
		case '\\':
			enc.buf.WriteString(`\\`)
		case '\n':
			enc.buf.WriteString(`\n`)
		case '\r':
			enc.buf.WriteString(`\r`)
		case '\t':
			enc.buf.WriteString(`\t`)
		default:
			if r < 32 {
				enc.buf.WriteString(fmt.Sprintf(`\u%04x`, r))
			} else {
				enc.buf.WriteRune(r)
			}
		}
	}
}

func (enc *ReflectEncoder) encodeArray(v reflect.Value) error {
	enc.buf.WriteByte('[')
	
	enc.depth++
	defer func() { enc.depth-- }()
	
	length := v.Len()
	for i := 0; i < length; i++ {
		if i > 0 {
			enc.buf.WriteByte(',')
		}
		if err := enc.encodeValue(v.Index(i)); err != nil {
			return err
		}
	}
	
	enc.buf.WriteByte(']')
	return nil
}

func (enc *ReflectEncoder) encodeMap(v reflect.Value) error {
	enc.buf.WriteByte('{')
	
	enc.depth++
	defer func() { enc.depth-- }()
	
	keys := v.MapKeys()
	
	// Sort keys for deterministic output
	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprintf("%v", keys[i].Interface()) < fmt.Sprintf("%v", keys[j].Interface())
	})
	
	for i, key := range keys {
		if i > 0 {
			enc.buf.WriteByte(',')
		}
		
		// Encode key
		keyStr := fmt.Sprintf("%v", key.Interface())
		enc.buf.WriteByte('"')
		enc.buf.WriteString(keyStr)
		enc.buf.WriteByte('"')
		enc.buf.WriteByte(':')
		
		// Encode value
		if err := enc.encodeValue(v.MapIndex(key)); err != nil {
			return err
		}
	}
	
	enc.buf.WriteByte('}')
	return nil
}

func (enc *ReflectEncoder) encodeStruct(v reflect.Value) error {
	t := v.Type()
	
	// Handle time.Time specially
	if t == reflect.TypeOf(time.Time{}) {
		timeVal := v.Interface().(time.Time)
		enc.encodeString(timeVal.Format(time.RFC3339))
		return nil
	}
	
	enc.buf.WriteByte('{')
	
	enc.depth++
	defer func() { enc.depth-- }()
	
	fieldCount := 0
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		
		// Skip unexported fields
		if !field.IsExported() {
			continue
		}
		
		// Skip nil pointers
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}
		
		if fieldCount > 0 {
			enc.buf.WriteByte(',')
		}
		
		// Use json tag if available, otherwise use field name
		fieldName := field.Name
		if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
			if commaIdx := strings.Index(tag, ","); commaIdx != -1 {
				fieldName = tag[:commaIdx]
			} else {
				fieldName = tag
			}
		}
		
		// Always quote field names in JSON format
		enc.buf.WriteByte('"')
		enc.buf.WriteString(fieldName)
		enc.buf.WriteByte('"')
		enc.buf.WriteByte(':')
		
		if err := enc.encodeValue(fieldValue); err != nil {
			return err
		}
		
		fieldCount++
	}
	
	enc.buf.WriteByte('}')
	return nil
}

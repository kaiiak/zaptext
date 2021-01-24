package zaptext

import (
	"bytes"
	"io"
)

type ReflectEncoder struct {
	w   io.Writer
	err error
	buf *bytes.Buffer
}

func NewReflectEncoder(w io.Writer) *ReflectEncoder {
	return &ReflectEncoder{w: w}
}

func (enc *ReflectEncoder) Encode() error {
	if enc.err != nil {
		return enc.err
	}
	if enc.buf == nil {
		enc.buf = new(bytes.Buffer)
	}

	return nil
}

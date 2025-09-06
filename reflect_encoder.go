package zaptext

import (
	"bytes"
	"io"
)

type ReflectEncoder struct {
	w          io.Writer
	err        error
	escapeHTML bool

	indentBuf *bytes.Buffer
}

func NewReflectEncoder(w io.Writer) *ReflectEncoder {
	return &ReflectEncoder{w: w, escapeHTML: true}
}

func (enc *ReflectEncoder) Encode(obj any) error {
	if enc.err != nil {
		return enc.err
	}
	if enc.indentBuf == nil {
		enc.indentBuf = new(bytes.Buffer)
	}

	return nil
}

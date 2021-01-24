package zaptext_test

import (
	"bytes"
	"testing"

	. "github.com/kaiiak/zaptext"
)

func BenchmarkReflectEncoderEncode(b *testing.B) {
	var w = bytes.NewBuffer(nil)
	NewReflectEncoder(w)
}

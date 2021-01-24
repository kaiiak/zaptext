package zaptext_test

import (
	"bytes"
	"testing"

	. "github.com/kaiiak/zaptext"
)

func TestNewReflectEncoder(t *testing.T) {
	tests := []struct {
		name  string
		wantW string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			NewReflectEncoder(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("NewReflectEncoder() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

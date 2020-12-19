package zaptext_test

import (
	"reflect"
	"testing"

	"github.com/kaiiak/zaptext"
	"go.uber.org/zap/zapcore"
)

func TestNewTextEncoder(t *testing.T) {
	type args struct {
		cfg zapcore.EncoderConfig
	}
	tests := []struct {
		name string
		args args
		want zapcore.Encoder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := zaptext.NewTextEncoder(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTextEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

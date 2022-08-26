package main

import (
	"os"

	"github.com/kaiiak/zaptext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var (
		logger *zap.Logger
		cfg    = zap.NewProductionEncoderConfig()
		te     = zaptext.NewTextEncoder(cfg)
	)

	logger = zap.New(zapcore.NewCore(te, os.Stdout, zap.DebugLevel))
	zap.ReplaceGlobals(logger)

	logger.Info("msg", zap.Int64s("ids", []int64{1, 2, 3}))
}

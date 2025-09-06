package main

import (
	"os"
	"time"

	"github.com/kaiiak/zaptext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Configure encoder for human-readable text output
	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = "time"
	cfg.LevelKey = "level"
	cfg.MessageKey = "msg"
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncodeDuration = zapcore.StringDurationEncoder

	// Create text encoder
	encoder := zaptext.NewTextEncoder(cfg)

	// Create logger with text encoder
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.InfoLevel)
	logger := zap.New(core, zap.AddCaller())

	defer func() {
		// Only log sync errors if they are not "invalid argument" (common on stdout/stderr)
		err := logger.Sync()
		if err != nil && err.Error() != "invalid argument" {
			logger.Error("Failed to sync logger", zap.Error(err))
		}
	}()

	// Log various types of data
	logger.Info("Application started successfully",
		zap.String("version", "1.0.0"),
		zap.String("env", "production"),
		zap.Duration("startup_time", 150*time.Millisecond),
	)

	logger.Info("Processing user request",
		zap.String("user_id", "12345"),
		zap.String("action", "create user profile"),
		zap.Int("request_size", 1024),
		zap.Bool("authenticated", true),
		zap.Strings("permissions", []string{"read", "write", "admin"}),
		zap.Ints("scores", []int{95, 87, 92}),
	)

	logger.Warn("High memory usage detected",
		zap.Float64("memory_usage_percent", 85.7),
		zap.String("recommendation", "consider scaling up"),
	)

	logger.Error("Database connection failed",
		zap.String("error", "connection timeout"),
		zap.String("database", "users_db"),
		zap.Duration("timeout", 30*time.Second),
		zap.Int("retry_count", 3),
	)
}

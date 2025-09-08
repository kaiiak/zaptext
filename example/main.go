package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kaiiak/zaptext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type User struct {
	ID       int               `json:"id"`
	Name     string            `json:"name"`
	Email    string            `json:"email"`
	IsActive bool              `json:"is_active"`
	Profile  UserProfile       `json:"profile"`
	Tags     []string          `json:"tags"`
	Settings map[string]string `json:"settings"`
}

type UserProfile struct {
	Age      int       `json:"age"`
	Location string    `json:"location"`
	JoinDate time.Time `json:"join_date"`
}

func main() {
	fmt.Println("=== ZapText TextEncoder Example ===")
	textEncoderExample()
	
	fmt.Println("\n=== ZapText ReflectEncoder Example ===")
	reflectEncoderExample()
}

func textEncoderExample() {
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

func reflectEncoderExample() {
	// Create a sample user object
	user := User{
		ID:       123,
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		IsActive: true,
		Profile: UserProfile{
			Age:      30,
			Location: "San Francisco",
			JoinDate: time.Date(2023, 1, 15, 10, 30, 0, 0, time.UTC),
		},
		Tags:     []string{"premium", "verified", "early_adopter"},
		Settings: map[string]string{
			"theme":         "dark",
			"notifications": "enabled",
			"language":      "en",
		},
	}

	// Demonstrate ReflectEncoder with different data types
	fmt.Println("1. Encoding User struct:")
	encoder1 := zaptext.NewReflectEncoder(os.Stdout)
	encoder1.Encode(user)
	encoder1.Release()
	fmt.Println()

	fmt.Println("\n2. Encoding simple map:")
	simpleMap := map[string]any{
		"message": "Hello World",
		"count":   42,
		"active":  true,
		"pi":      3.14159,
	}
	encoder2 := zaptext.NewReflectEncoder(os.Stdout)
	encoder2.Encode(simpleMap)
	encoder2.Release()
	fmt.Println()

	fmt.Println("\n3. Encoding array of mixed types:")
	mixedArray := []any{
		"string",
		123,
		true,
		[]int{1, 2, 3},
		map[string]string{"key": "value"},
	}
	encoder3 := zaptext.NewReflectEncoder(os.Stdout)
	encoder3.Encode(mixedArray)
	encoder3.Release()
	fmt.Println()

	fmt.Println("\n4. Encoding primitive values:")
	primitives := []any{
		"Hello, World!",
		42,
		3.14159,
		true,
		time.Now(),
		nil,
	}
	
	for i, val := range primitives {
		fmt.Printf("  %d. ", i+1)
		encoder := zaptext.NewReflectEncoder(os.Stdout)
		encoder.Encode(val)
		encoder.Release()
		fmt.Println()
	}
}

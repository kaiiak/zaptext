# Zaptext
Zap text encoder for human-readable logging

[![Build status](https://github.com/kaiiak/zaptext/workflows/build/badge.svg)](https://github.com/kaiiak/zaptext/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/kaiiak/zaptext.svg)](https://pkg.go.dev/github.com/kaiiak/zaptext)
[![GoReport](https://goreportcard.com/badge/github.com/kaiiak/zaptext)](https://goreportcard.com/report/github.com/kaiiak/zaptext)

Zaptext provides a text encoder for [Zap](https://github.com/uber-go/zap) that outputs human-readable logs instead of JSON. Perfect for development, debugging, or when you need logs that are easy to read by humans.

## Features

- Human-readable text output instead of JSON
- Proper key=value formatting for fields  
- Automatic quoting of strings containing spaces
- Array and object formatting: `array=[1,2,3]`
- Support for all Zap data types: strings, numbers, booleans, durations, timestamps
- Configurable time and duration formatting
- Thread-safe and performant

## Quick Start

```golang
package main

import (
    "os"
    "github.com/kaiiak/zaptext"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func main() {
    // Configure the encoder
    cfg := zap.NewProductionEncoderConfig()
    cfg.EncodeTime = zapcore.ISO8601TimeEncoder
    cfg.EncodeDuration = zapcore.StringDurationEncoder
    
    // Create text encoder
    encoder := zaptext.NewTextEncoder(cfg)
    
    // Create logger
    core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.InfoLevel)
    logger := zap.New(core)
    
    logger.Info("User logged in",
        zap.String("user_id", "12345"),
        zap.Bool("admin", true),
        zap.Duration("session_duration", time.Minute*30),
    )
    // Output: time=2023-09-02T10:30:15Z INFO User logged in user_id=12345 admin=true session_duration=30m
}
```

## Output Format

The text encoder produces logs in this format:
```
time=2023-09-02T10:30:15Z INFO User action user_id=12345 action="create profile" success=true duration=150ms tags=[user,profile,create]
```

Where:
- Field values are formatted as `key=value`
- Strings with spaces are automatically quoted: `action="create profile"`
- Arrays use brackets: `tags=[user,profile,create]`
- All standard Zap field types are supported

## Comparison with JSON Encoder

**JSON Output (default zap):**
```json
{"level":"info","ts":1693648215,"msg":"User action","user_id":"12345","action":"create profile","success":true,"duration":150000000,"tags":["user","profile","create"]}
```

**Text Output (zaptext):**
```
time=2023-09-02T10:30:15Z INFO User action user_id=12345 action="create profile" success=true duration=150ms tags=[user,profile,create]
```

## Thanks

- [zap](https://github.com/uber-go/zap): The excellent structured logging library this encoder extends
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
- Support for all Zap data types: strings, numbers, booleans, durations, timestamps, complex numbers
- Configurable time and duration formatting
- Thread-safe and performant
- Reflection-based encoder for arbitrary Go data structures
  - Object pooling for memory efficiency
  - HTML character escaping for web-safe output
  - JSON-compatible complex number formatting
  - Configurable maximum recursion depth (default: 32)
  - Comprehensive error handling and infinite recursion protection

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

## ReflectEncoder Usage

The ReflectEncoder provides reflection-based encoding of arbitrary Go data structures into JSON-like format:

```golang
package main

import (
    "bytes"
    "fmt"
    "github.com/kaiiak/zaptext"
)

type User struct {
    ID       int      `json:"id"`
    Name     string   `json:"name"`
    Email    string   `json:"email"`
    Tags     []string `json:"tags"`
    IsActive bool     `json:"active"`
}

func main() {
    // Create a user object
    user := User{
        ID:       123,
        Name:     "John Doe",
        Email:    "john@example.com",
        Tags:     []string{"developer", "admin"},
        IsActive: true,
    }
    
    // Create ReflectEncoder
    var buf bytes.Buffer
    encoder := zaptext.NewReflectEncoder(&buf)
    defer encoder.Release() // Always release back to pool
    
    // Configure options
    encoder.SetEscapeHTML(true)      // Enable HTML escaping
    encoder.SetMaxDepth(50)          // Set maximum recursion depth
    
    // Encode the object
    err := encoder.Encode(user)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Println(buf.String())
    // Output: {"id":123,"name":"John Doe","email":"john@example.com","tags":["developer","admin"],"active":true}
}
```

### ReflectEncoder Features

- **Object Pooling**: Uses sync.Pool for efficient memory allocation
- **HTML Escaping**: Optional HTML character escaping for web-safe output
- **Complex Numbers**: JSON-compatible format `{"real": 1.0, "imag": 2.0}`
- **Maximum Depth Protection**: Configurable recursion depth limit (default: 32)
- **Comprehensive Type Support**: All Go primitive types, arrays, slices, maps, structs
- **Special Handling**: time.Time formatted as RFC3339, nil pointers skipped
- **Error Handling**: Maintains error state and provides detailed error messages

## Output Format

The text encoder produces logs in this format:
```
time=2023-09-02T10:30:15Z INFO User action user_id=12345 action="create profile" success=true duration=150ms tags=[user,profile,create] complex=1+2i
```

Where:
- Field values are formatted as `key=value`
- Strings with spaces are automatically quoted: `action="create profile"`
- Arrays use brackets: `tags=[user,profile,create]`
- Complex numbers: `complex=1+2i`
- All standard Zap field types are supported

### ReflectEncoder Output Examples

**Complex Object:**
```json
{"id":123,"name":"John Doe","email":"john@example.com","tags":["developer","admin"],"active":true,"balance":1234.56,"metadata":{"last_login":"2023-09-02T10:30:15Z","login_count":42}}
```

**Complex Numbers (JSON-compatible):**
```json
{"real":1.5,"imag":2.5}
```

**HTML Escaping:**
```json
{"description":"User \u003cscript\u003ealert('xss')\u003c/script\u003e input","safe":"\u0026 \u003c \u003e"}
```

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
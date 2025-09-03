# Zaptext
Zap 人类可读文本日志编码器

[![Build status](https://github.com/kaiiak/zaptext/workflows/build/badge.svg)](https://github.com/kaiiak/zaptext/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/kaiiak/zaptext.svg)](https://pkg.go.dev/github.com/kaiiak/zaptext)
[![GoReport](https://goreportcard.com/badge/github.com/kaiiak/zaptext)](https://goreportcard.com/report/github.com/kaiiak/zaptext)

Zaptext 为 [Zap](https://github.com/uber-go/zap) 提供了一个文本编码器，输出人类可读的日志而不是 JSON 格式。非常适合开发、调试或需要易于阅读的日志的场景。

## 特性

- 输出人类可读的文本格式而不是 JSON
- 字段使用 key=value 格式
- 包含空格的字符串自动加引号
- 数组和对象格式化：`array=[1,2,3]`
- 支持所有 Zap 数据类型：字符串、数字、布尔值、时长、时间戳
- 可配置的时间和时长格式
- 线程安全和高性能

## 快速开始

```golang
package main

import (
    "os"
    "github.com/kaiiak/zaptext"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func main() {
    // 配置编码器
    cfg := zap.NewProductionEncoderConfig()
    cfg.EncodeTime = zapcore.ISO8601TimeEncoder
    cfg.EncodeDuration = zapcore.StringDurationEncoder
    
    // 创建文本编码器
    encoder := zaptext.NewTextEncoder(cfg)
    
    // 创建日志器
    core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.InfoLevel)
    logger := zap.New(core)
    
    logger.Info("用户登录",
        zap.String("user_id", "12345"),
        zap.Bool("admin", true),
        zap.Duration("session_duration", time.Minute*30),
    )
    // 输出: time=2023-09-02T10:30:15Z INFO 用户登录 user_id=12345 admin=true session_duration=30m
}
```

## 输出格式

文本编码器产生这样格式的日志：
```
time=2023-09-02T10:30:15Z INFO User action user_id=12345 action="create profile" success=true duration=150ms tags=[user,profile,create]
```

其中：
- 字段值格式化为 `key=value`
- 包含空格的字符串自动加引号：`action="create profile"`
- 数组使用方括号：`tags=[user,profile,create]`
- 支持所有标准 Zap 字段类型

## 与 JSON 编码器的对比

**JSON 输出 (默认 zap):**
```json
{"level":"info","ts":1693648215,"msg":"User action","user_id":"12345","action":"create profile","success":true,"duration":150000000,"tags":["user","profile","create"]}
```

**文本输出 (zaptext):**
```
time=2023-09-02T10:30:15Z INFO User action user_id=12345 action="create profile" success=true duration=150ms tags=[user,profile,create]
```

## 致谢

- [zap](https://github.com/uber-go/zap): 这个编码器扩展的优秀结构化日志库
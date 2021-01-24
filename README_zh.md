# Zaptext
![[Build status](https://github.com/kaiiak/zaptext/actions)](https://github.com/kaiiak/zaptext/workflows/build/badge.svg)
Zap文本日志



## 快速开始

```golang
cfg := zapcore.NewProductionEncoderConfig()
zaptext.NewTextencoder(cfg)
```
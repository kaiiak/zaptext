# Zaptext
Zap text encoder instance


[![Build status](https://github.com/kaiiak/zaptext/workflows/build/badge.svg)](https://github.com/kaiiak/zaptext/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/kaiiak/zaptext.svg)](https://pkg.go.dev/github.com/kaiiak/zaptext)
[![GoReport](https://goreportcard.com/badge/github.com/kaiiak/zaptext)](https://goreportcard.com/report/github.com/kaiiak/zaptext)

## Using example

```golang
cfg := zapcore.NewProductionEncoderConfig()
zaptext.NewTextencoder(cfg)
```
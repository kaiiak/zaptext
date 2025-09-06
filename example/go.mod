module github.com/kaiiak/zaptext/example

go 1.21

replace github.com/kaiiak/zaptext => ../

require (
	github.com/kaiiak/zaptext v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.21.0
)

require (
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/text v0.3.8 // indirect
)

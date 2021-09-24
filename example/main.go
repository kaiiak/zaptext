package main

import "go.uber.org/zap"

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	logger.Info("msg", zap.Int64s("ids", []int64{1, 2, 3}))
}

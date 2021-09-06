package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func init() {
	l, err := zap.NewProduction()
	if err != nil {
		fmt.Println("error creating logger")
		os.Exit(1)
	}
	Log = l.Sugar()
}

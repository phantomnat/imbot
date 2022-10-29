//go:build release
// +build release

package pkg

import (
	"log"
	"strings"
	"time"

	"go.uber.org/zap"
)

func InitLogger() {
	fn := strings.ReplaceAll(time.Now().Format(time.RFC3339Nano), ":", "-") + ".log"
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	cfg.Development = false
	cfg.OutputPaths = []string{
		"logs/" + fn,
	}
	cfg.ErrorOutputPaths = []string{
		"logs/" + fn,
	}
	l, err := cfg.Build()
	if err != nil {
		log.Fatalf("cannot init logger: %+v", err)
		return
	}
	zap.ReplaceGlobals(l)
}

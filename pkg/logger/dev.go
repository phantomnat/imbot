//go:build !release
// +build !release

package pkg

import (
	"log"

	"go.uber.org/zap"
)

func InitLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{
		"stdout",
	}
	l, err := cfg.Build()
	if err != nil {
		log.Fatalf("cannot init logger: %+v", err)
		return
	}
	zap.ReplaceGlobals(l)
}

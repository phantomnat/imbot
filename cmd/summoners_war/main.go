package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	summoners_war_chronicles "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles"
	"github.com/phantomnat/imbot/pkg/im"
	pkg "github.com/phantomnat/imbot/pkg/logger"
	"github.com/phantomnat/imbot/pkg/ui"
)

func main() {
	pkg.InitLogger()
	log := zap.S().Named("main")

	signalCtx := SetupSignalHandler()
	globalCtx, globalCancel := context.WithCancel(signalCtx)

	imgProcManager := im.GetImageManager()

	game, err := summoners_war_chronicles.New(summoners_war_chronicles.Option{
		Ctx:          globalCtx,
		ImageManager: imgProcManager,
	})
	if err != nil {
		log.Fatalf("cannot init summoners war: chronicles: %+v", err)
		return
	}

	uiHandler := ui.New(game)

	go game.Run(globalCtx.Done())

	uiHandler.Run(globalCtx.Done(), globalCancel)

	<-globalCtx.Done()

	log.Infof("exiting")
}

var onlyOneSignalHandler = make(chan struct{})

func SetupSignalHandler() context.Context {
	close(onlyOneSignalHandler) // panics when called twice

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return ctx
}

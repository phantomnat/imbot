package main

import (
	"context"

	"go.uber.org/zap"

	summoners_war_chronicles "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles"
	"github.com/phantomnat/imbot/pkg/im"
	pkg "github.com/phantomnat/imbot/pkg/logger"
	"github.com/phantomnat/imbot/pkg/ui"
)

func main() {
	pkg.InitLogger()
	log := zap.S().Named("main")

	ctx, cancel := context.WithCancel(context.Background())

	imgProcManager := im.GetImageManager()

	game, err := summoners_war_chronicles.New(summoners_war_chronicles.Option{
		Ctx:          ctx,
		ImageManager: imgProcManager,
	})
	if err != nil {
		log.Fatalf("cannot init summoners war: chronicles: %+v", err)
		return
	}

	uiHandler := ui.New(game)

	go game.Run(ctx.Done())

	uiHandler.Run()

	cancel()
}

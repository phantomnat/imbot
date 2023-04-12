package main

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	hook "github.com/robotn/gohook"
	"go.uber.org/zap"

	summoners_war_chronicles "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles"
	"github.com/phantomnat/imbot/pkg/im"
	pkg "github.com/phantomnat/imbot/pkg/logger"
	"github.com/phantomnat/imbot/pkg/screen/mumu"
	"github.com/phantomnat/imbot/pkg/ui"
)

func main() {
	pkg.InitLogger()
	log := zap.S().Named("main")
	imgProcManager := im.GetImageManager()

	mumu.NewFromTitle("Chronicles - MuMu Player")
	game, err := summoners_war_chronicles.New(imgProcManager)
	if err != nil {
		log.Fatalf("cannot init summoners war: chronicles: %+v", err)
		return
	}

	uiHandler := ui.New(game)
	log.Infof("registering hooks")

	hook.Register(hook.KeyDown, []string{"p", "ctrl"}, func(event hook.Event) {
		uiHandler.OnBtnToggleRunClicked()
	})

	// hook.Register(hook.KeyDown, []string{"f", "ctrl"}, func(event hook.Event) {
	// 	game.GetScreen().GetCurrentCursorPos()
	// })

	hook.Register(hook.KeyDown, []string{"w", "ctrl"}, func(event hook.Event) {
		log.Infof("screen capturing...")
		today := strings.ReplaceAll(time.Now().Format(time.RFC3339Nano), ":", "-")
		filePath := filepath.Join("cap", today+".png")
		game.GetScreen().CaptureMatAndSave(filePath)
	})

	s := hook.Start()
	hook.Process(s)

	ctx, cancel := context.WithCancel(context.Background())

	go game.Run(ctx.Done())

	uiHandler.Run()

	// hook.StopEvent()
	// close(s)
	// <-done
	cancel()

}

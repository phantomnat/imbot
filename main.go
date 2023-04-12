package main

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	hook "github.com/robotn/gohook"
	"go.uber.org/zap"

	"github.com/phantomnat/imbot/pkg/game/bloon_td6"
	"github.com/phantomnat/imbot/pkg/im"
	pkg "github.com/phantomnat/imbot/pkg/logger"
	"github.com/phantomnat/imbot/pkg/ui"
)

func main() {
	pkg.InitLogger()
	log := zap.S().Named("main")
	imgProcManager := im.GetImageManager()

	game, err := bloon_td6.New(imgProcManager)
	if err != nil {
		log.Fatalf("cannot init bloons td 6: %+v", err)
		return
	}

	uiHandler := ui.New(game)
	log.Infof("registering hooks")

	hook.Register(hook.KeyDown, []string{"p", "ctrl"}, func(event hook.Event) {
		uiHandler.OnBtnToggleRunClicked()
	})

	hook.Register(hook.KeyDown, []string{"f", "ctrl"}, func(event hook.Event) {
		// game.GetScreen().GetCurrentCursorPos()
	})

	hook.Register(hook.KeyDown, []string{"w", "ctrl"}, func(event hook.Event) {
		log.Infof("screen capturing...")
		today := strings.ReplaceAll(time.Now().Format(time.RFC3339Nano), ":", "-")
		filePath := filepath.Join("cap", today+".png")
		game.GetScreen().CaptureMatAndSave(filePath)
	})

	// hook.Register(hook.KeyDown, []string{"t"}, func(event hook.Event) {
	// 	log.Infof("testing...")

	// 	mat, err := game.GetMat()
	// 	if err != nil {
	// 		log.Errorf("get mat: %+v", err)
	// 		return
	// 	}
	// 	defer mat.Close()

	// 	ok, pt := imgProcManager.MatchDefault(mat, "btn-play")
	// 	if ok {
	// 		log.Info("found on %v", pt)
	// 	}
	// 	img, err := mat.ToImage()
	// 	if err == nil {
	// 		uiHandler.OnImageUpdated(img)
	// 	}

	// 	if ok {
	// 		game.GetScreen().MouseMoveAndClick(pt.X+10, pt.Y+10)
	// 	}
	// 	// rect := game.WindowSize()
	// 	// log.Infof("screen: %s, click: x: %d, y: %d", rect, rect.X+pt.X+10, rect.Y+pt.Y+10)
	// 	// robotgo.MoveClick(rect.X+, rect.Y+)
	// })

	s := hook.Start()
	hook.Process(s)

	ctx, cancel := context.WithCancel(context.Background())

	go game.Run(ctx.Done())

	uiHandler.Run()

	hook.StopEvent()
	close(s)
	// <-done
	cancel()

	// TODO: graceful shutdown?
}

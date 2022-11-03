package main

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"go.uber.org/zap"
	"gocv.io/x/gocv"

	"github.com/phantomnat/imbot/pkg/game/bloon_td6"
	"github.com/phantomnat/imbot/pkg/im"
	pkg "github.com/phantomnat/imbot/pkg/logger"
	"github.com/phantomnat/imbot/pkg/ui"
)

func main() {
	pkg.InitLogger()
	log := zap.S().Named("main")

	game, err := bloon_td6.New()
	if err != nil {
		log.Fatalf("cannot init bloons td 6: %+v", err)
		return
	}
	imgProcManager := im.GetImageManager()

	uiHandler := ui.New()
	log.Infof("registering hooks")

	hook.Register(hook.KeyDown, []string{"q"}, func(event hook.Event) {
		log.Infof("exiting...")
		hook.End()
	})

	hook.Register(hook.KeyDown, []string{"c"}, func(event hook.Event) {
		log.Infof("screen capturing...")
		screen := game.WindowSize()
		if screen == nil {
			return
		}
		log.Debugf("screen: %v", screen)

		bit := robotgo.CaptureScreen(screen.X, screen.Y, screen.Width, screen.Height)
		defer robotgo.FreeBitmap(bit)

		img := robotgo.ToImage(bit)
		mat, err := gocv.ImageToMatRGBA(img)
		if err != nil {
			log.Errorf("cannot convert to gocv mat: %+v", err)
		}
		defer mat.Close()
		today := strings.ReplaceAll(time.Now().Format(time.RFC3339Nano), ":", "-")
		fp := filepath.Join("cap", today+".png")
		gocv.IMWrite(fp, mat)

		uiHandler.OnImageUpdated(img)
	})

	hook.Register(hook.KeyDown, []string{"t"}, func(event hook.Event) {
		log.Infof("testing...")

		mat, err := game.GetMat()
		if err != nil {
			log.Errorf("get mat: %+v", err)
			return
		}
		defer mat.Close()

		ok, pt := imgProcManager.MatchDefault(mat, "btn-play")
		if ok {
			log.Info("found on %v", pt)
		}

		rect := game.WindowSize()
		log.Infof("screen: %s, click: x: %d, y: %d", rect, rect.X+pt.X+10, rect.Y+pt.Y+10)
		robotgo.MoveClick(rect.X+pt.X+10, rect.Y+pt.Y+10)
	})

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

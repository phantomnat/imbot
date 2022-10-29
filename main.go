package main

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/phantomnat/imbot/pkg/game/bloon_td6"
	hook "github.com/robotn/gohook"
	"go.uber.org/zap"
	"gocv.io/x/gocv"

	pkg "github.com/phantomnat/imbot/pkg/logger"
)

func main() {
	pkg.InitLogger()
	log := zap.S().Named("main")

	game, err := bloon_td6.New()
	if err != nil {
		log.Fatalf("cannot init bloons td 6: %+v", err)
		return
	}

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
	})

	s := hook.Start()
	done := hook.Process(s)

	ctx, cancel := context.WithCancel(context.Background())

	go game.Run(ctx.Done())
	<-done
	cancel()

	// TODO: graceful shutdown?
}

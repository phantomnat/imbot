package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	dir, err := os.MkdirTemp("", "chromedp-example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.UserDataDir(dir),
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1280, 720),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// also set up a custom logger
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// ensure that the browser process is started
	if err := chromedp.Run(taskCtx); err != nil {
		log.Fatal(err)
	}

	//path := filepath.Join(dir, "DevToolsActivePort")
	//bs, err := os.ReadFile(path)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//lines := bytes.Split(bs, []byte("\n"))
	//fmt.Printf("DevToolsActivePort has %d lines\n", len(lines))

	// capture screenshot of an element
	var buf []byte
	//if err := chromedp.Run(taskCtx, elementScreenshot(`https://pkg.go.dev/`, `img.Homepage-logo`, &buf)); err != nil {
	//	log.Fatal(err)
	//}
	//if err := os.WriteFile("elementScreenshot.png", buf, 0o644); err != nil {
	//	log.Fatal(err)
	//}

	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(taskCtx, fullScreenshot(`https://www.cloudemulator.net/app/sign-in/email`, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("fullScreenshot.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	log.Printf("wrote elementScreenshot.png and fullScreenshot.png")
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Note: chromedp.FullScreenshot overrides the device's emulation settings. Use
// device.Reset to reset the emulation and viewport settings.
func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(`input[name="ion-input-0"]`, chromedp.ByQuery),
		chromedp.WaitVisible(`div div.login-btns`, chromedp.ByQuery),
		chromedp.Sleep(10 * time.Second),
		chromedp.FullScreenshot(res, quality),
	}
}

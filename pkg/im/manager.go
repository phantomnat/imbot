//go:generate go run gen.go

package im

import (
	"bytes"
	"image"
	"sync"

	"go.uber.org/zap"
	"gocv.io/x/gocv"
)

type ImageManager struct {
	log       *zap.SugaredLogger
	images    map[string]*gocv.Mat
	bgrImages map[string]*gocv.Mat
}

var (
	name   = "image.manager"
	initor sync.Once
	_im    *ImageManager
)

func Init() {
	im := GetImageManager()
	log := im.log
	for name, data := range images {
		png, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			log.Fatalf("loading image: %+v\n", err)
			continue
		}

		m, err := gocv.ImageToMatRGBA(png)
		if err != nil {
			log.Fatalf("init mat: %+v\n", err)
			continue
		}
		defer m.Close()

		log.Debugf("load name")
		gocv.CvtColor(m, &m, gocv.ColorRGBAToBGR)
		rgb := gocv.NewMatWithSize(m.Rows(), m.Cols(), gocv.MatTypeCV32FC3)
		m.ConvertTo(&rgb, gocv.MatTypeCV32FC3)
		im.bgrImages[name] = &rgb
		gray := gocv.NewMatWithSize(m.Rows(), m.Cols(), gocv.MatTypeCV32FC1)
		m.ConvertTo(&gray, gocv.MatTypeCV32FC1)
		im.images[name] = &gray
	}
}

func GetImageManager() *ImageManager {
	initor.Do(func() {
		_im = &ImageManager{
			images:    map[string]*gocv.Mat{},
			bgrImages: map[string]*gocv.Mat{},
			log:       zap.S().Named(name),
		}
	})
	return _im
}

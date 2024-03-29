package summonerswar

import (
	"image"
	"strings"

	"gocv.io/x/gocv"

	"github.com/phantomnat/imbot/pkg/domain"
)

func (b *SummonersWar) ToggleSendCaptureImage(isSend bool, cb ...func(image.Image)) {
	b.muSendCaptureImage.Lock()
	defer b.muSendCaptureImage.Unlock()

	b.isSendCaptureImage = true
	if len(cb) > 0 && cb[0] != nil {
		b.cbSendCaptureImage = cb[0]
	}
}

func (b *SummonersWar) sendCaptureImage(m gocv.Mat) {
	b.muSendCaptureImage.RLock()
	defer b.muSendCaptureImage.RUnlock()

	if !b.isSendCaptureImage {
		return
	}

	img, err := m.ToImage()
	if err != nil {
		// ignore error
		return
	}
	b.cbSendCaptureImage(img)
}

func Rect(x, y, w, h int) image.Rectangle {
	return image.Rect(x, y, x+w, y+h)
}

func (b *SummonersWar) imMatchDefault(m gocv.Mat, path ...string) (bool, image.Point) {
	paths := append([]string{srcImageDir}, path...)
	return b.im.MatchDefault(m, paths...)
}

func (b *SummonersWar) imMatchDefaultInROI(m gocv.Mat, roi image.Rectangle, path ...string) (bool, image.Point) {
	mROI := m.Region(roi)
	defer mROI.Close()
	ok, pt := b.imMatchDefault(mROI, path...)
	if ok {
		return ok, image.Point{X: pt.X + roi.Min.X, Y: pt.Y + roi.Min.Y}
	}
	return ok, pt
}

func (b *SummonersWar) MatchInROI(m gocv.Mat, roi image.Rectangle, o domain.MatchOption) (bool, image.Point) {
	mROI := m.Region(roi)
	defer mROI.Close()
	o.Path = b.GetImagePath(o.Path)
	ok, pt := b.im.MatchWithCenterROI(mROI, o)
	if ok {
		return ok, image.Point{X: pt.X + roi.Min.X, Y: pt.Y + roi.Min.Y}
	}
	return ok, pt
}

func (b *SummonersWar) GetImagePath(path ...string) string {
	return srcImageDir + "." + strings.Join(path, ".")
}

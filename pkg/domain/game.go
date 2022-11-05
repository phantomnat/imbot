package domain

import "image"

type Game interface {
	Start()
	Pause()
	Reset()
	ToggleSendCaptureImage(isSend bool, cb ...func(image.Image))
}

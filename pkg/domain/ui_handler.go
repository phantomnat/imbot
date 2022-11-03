package domain

import "image"

type UIHandler interface {
	UpdateImage(in image.Image)
	Run()
}

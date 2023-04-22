package domain

import (
	"image"

	"gocv.io/x/gocv"
)

type MatchOption struct {
	Path      string
	Mask      *gocv.Mat
	HasMask   bool
	Th        float32
	PrintVal  bool
}

type ImageManager interface {
	MatchWithOption(src gocv.Mat, opt MatchOption) (bool, image.Point)
	MatchMultiPoints(src gocv.Mat, opt MatchOption) []image.Point
}

type ImageBuffer struct {
	Width  int
	Height int
	Data   []byte
}

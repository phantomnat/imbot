package tasks

import (
	"image"
	"strings"
	"time"

	"gocv.io/x/gocv"

	"github.com/phantomnat/imbot/pkg/domain"
)

type SearchROIOption struct {
	needClick   bool
	path        string
	roi         image.Rectangle
	wait        *time.Duration
	nextState   *domain.TaskState
	clickOffset image.Point

	debugTemplateMatch bool
}

func (o *SearchROIOption) getOffsetPt(pt image.Point) image.Point {
	if o.clickOffset == (image.Point{}) {
		return pt
	}
	return image.Pt(pt.X+o.clickOffset.X, pt.Y+o.clickOffset.Y)
}

func (o *SearchROIOption) apply(opts ...SearchROIOptionFunc) {
	for _, fn := range opts {
		fn(o)
	}
	if o.wait == nil {
		v := time.Millisecond * 1000
		o.wait = &v
	}
}

func (o *SearchROIOption) isEmpty() bool {
	if o == nil {
		return true
	}
	emptyPath := o.path == ""
	return emptyPath
}

type SearchROIOptionFunc func(o *SearchROIOption)

func WithClick() SearchROIOptionFunc {
	return func(o *SearchROIOption) {
		o.needClick = true
	}
}

func WithPath(paths ...string) SearchROIOptionFunc {
	return func(o *SearchROIOption) {
		o.path = strings.Join(paths, ".")
	}
}

func WithROI(roi image.Rectangle) SearchROIOptionFunc {
	return func(o *SearchROIOption) {
		o.roi = roi
	}
}

func WithDebugMatch() SearchROIOptionFunc {
	return func(o *SearchROIOption) {
		o.debugTemplateMatch = true
	}
}

func WithWaitMs(ms int) SearchROIOptionFunc {
	return func(o *SearchROIOption) {
		v := time.Duration(ms) * time.Millisecond
		o.wait = &v
	}
}

func WithNoWait() SearchROIOptionFunc {
	return WithWaitMs(0)
}

func WithNextState(s domain.TaskState) SearchROIOptionFunc {
	return func(o *SearchROIOption) {
		o.nextState = &s
	}
}

func WithClickOffset(roi image.Rectangle) SearchROIOptionFunc {
	return func(o *SearchROIOption) {
		o.clickOffset = roi.Min
	}
}

func (t *BaseTask) WaitMs(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func (t *BaseTask) SearchROI(
	m gocv.Mat,
	opts ...SearchROIOptionFunc,
) (found bool) {
	opt := SearchROIOption{}
	opt.apply(opts...)

	if opt.isEmpty() {
		return false
	}

	var pt image.Point

	if opt.roi == (image.Rectangle{}) {
		found, pt = t.Im.MatchWithCenterROI(m, domain.MatchOption{
			Path: t.Manager.GetImagePath(opt.path),
		})
	} else {
		found, pt = t.Manager.MatchInROI(m, opt.roi, domain.MatchOption{
			Path:     opt.path,
			PrintVal: opt.debugTemplateMatch,
		})
	}

	if !found {
		return
	}
	if opt.needClick {
		t.Manager.ClickPt(opt.getOffsetPt(pt))
	}
	if opt.wait != nil || *opt.wait > 0 {
		time.Sleep(*opt.wait)
	}
	if opt.nextState != nil {
		t.SetState(*opt.nextState)
	}
	return
}

func (t *BaseTask) GetPtWithROI(roi image.Rectangle, pt image.Point) image.Point {
	return image.Pt(pt.X+roi.Min.X, pt.Y+roi.Min.Y)
}

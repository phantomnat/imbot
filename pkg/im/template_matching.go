package im

import (
	"image"
	"strconv"
	"strings"

	"gocv.io/x/gocv"

	"github.com/phantomnat/imbot/pkg/domain"
)

const (
	thOnePercent      = 0.01
	thOneTenthPercent = 0.001
)

func (m *ImageManager) MatchDefault(src gocv.Mat, p ...string) (bool, image.Point) {
	return m.matchWithCenterROI(src, domain.MatchOption{
		Path: strings.Join(p, "."),
		Th:   thOnePercent,
	})
}

// MatchOneP match should below 0.01 threshold value
// func (m *ImageManager) MatchOneP(src gocv.Mat, p string) (bool, image.Point) {
// 	return m.matchWithMask(src, p, thOnePercent)
// }

// func (m *ImageManager) MatchOneTenth(src gocv.Mat, p string) (bool, image.Point) {
// 	return m.matchWithMask(src, p, thOneTenthPercent)
// }

// Match from <p>.1 to <p>.10
// func (m *ImageManager) MatchNDefault(src gocv.Mat, p string) (bool, image.Point) {
// 	return m.matchWithMask(src, p, thOnePercent)
// }

// Match from <p>.1 to <p>.10
// func (m *ImageManager) MatchNOneTenth(src gocv.Mat, p string) (bool, image.Point) {
// 	return m.matchWithMask(src, p, thOneTenthPercent)
// }

// Match from <p>.1 to <p>.10
// func (m *ImageManager) MatchN(src gocv.Mat, p string, th float32) (bool, image.Point) {
// 	return m.matchWithMask(src, domain.MatchOption{
// 		Path: p,
// 		Th:   th,
// 	})
// }

func (m *ImageManager) MatchWithCenterROI(src gocv.Mat, opt domain.MatchOption) (bool, image.Point) {
	return m.matchWithCenterROI(src, opt)
}

func (m *ImageManager) MatchPoint(src gocv.Mat, opt domain.MatchOption) (bool, image.Point) {
	m.applyDefaultMatchOption(&opt)

	if ok, tpl := m.Get(opt.Path); ok {
		if matched, pt := m.templateMatch(&src, tpl, opt); matched {
			return matched, pt
		}
	}
	for i := 1; i <= 10; i++ {
		if ok, tpl := m.Get(opt.Path + "_" + strconv.Itoa(i)); ok {
			if matched, pt := m.templateMatch(&src, tpl, opt); matched {
				return true, pt
			}
		} else {
			// not found template, so we break
			break
		}
	}
	return false, image.Point{}
}

func (m *ImageManager) MatchMultiPoints(src gocv.Mat, opt domain.MatchOption) []image.Point {
	m.applyDefaultMatchOption(&opt)

	var points []image.Point
	if ok, tpl := m.Get(opt.Path); ok {
		pts := m.templateMatches(&src, tpl, opt)
		if len(pts) > 0 {
			points = append(points, pts...)
		}
	}
	for i := 1; i <= 10; i++ {
		if ok, tpl := m.Get(opt.Path + "_" + strconv.Itoa(i)); ok {
			pts := m.templateMatches(&src, tpl, opt)
			if len(pts) > 0 {
				points = append(points, pts...)
			}
		} else {
			// not found template, so we break
			break
		}
	}
	return points
}

func (m *ImageManager) applyDefaultMatchOption(opt *domain.MatchOption) {
	if opt.Th == 0 {
		opt.Th = thOnePercent
	}
	if ok, tpl := m.Get(opt.Path + "_mask"); ok {
		opt.Mask = tpl
		opt.HasMask = true
	}
}

// matchWithCenterROI from <p>.1 to <p>.10
func (m *ImageManager) matchWithCenterROI(src gocv.Mat, opt domain.MatchOption) (bool, image.Point) {

	m.applyDefaultMatchOption(&opt)

	if ok, tpl := m.Get(opt.Path); ok {
		if matched, pt := m.templateMatch(&src, tpl, opt); matched {
			return matched, image.Point{X: pt.X + tpl.Cols()/2, Y: pt.Y + tpl.Rows()/2}
		}
	}
	for i := 1; i <= 10; i++ {
		if ok, tpl := m.Get(opt.Path + "_" + strconv.Itoa(i)); ok {
			if matched, pt := m.templateMatch(&src, tpl, opt); matched {
				return true, image.Point{X: pt.X + tpl.Cols()/2, Y: pt.Y + tpl.Rows()/2}
			}
		} else {
			// not found template, so we break
			break
		}
	}
	return false, image.Point{}
}

func (m *ImageManager) templateMatch(src *gocv.Mat, tpl *gocv.Mat, o domain.MatchOption) (bool, image.Point) {
	res := m.rawTemplateMatch(src, tpl, o.Mask)
	defer res.Close()

	v, _, l, _ := gocv.MinMaxLoc(res)

	if o.PrintVal {
		m.log.Debugf("matching %s: %.4f at %v (mask: %v, expected: %.4f)", o.Path, v, l, o.HasMask, o.Th)
	}
	if v < o.Th {
		return true, l
	}

	return false, image.Point{}
}

func (m *ImageManager) templateMatches(src *gocv.Mat, tpl *gocv.Mat, o domain.MatchOption) []image.Point {
	res := m.rawTemplateMatch(src, tpl, o.Mask)
	defer res.Close()

	var points []image.Point

	for {
		v, _, l, _ := gocv.MinMaxLoc(res)
		if o.PrintVal {
			m.log.Debugf("matching %s: %.4f at %v (res: (%v, %v), mask: %v, expected: %.4f)", o.Path, v, l, res.Cols(), res.Rows(), o.HasMask, o.Th)
		}
		if v > o.Th {
			break
		}

		points = append(points, image.Pt(l.X, l.Y))
		res.SetFloatAt(l.Y, l.X, o.Th+1)
		// rectangle(img, Rect(max_point.x, max_point.y, templ.cols, templ.rows), Scalar(0,255,0), 2);
		// x2 := l.X + tpl.Cols()
		// y2 := l.Y + tpl.Rows()
		// if x2 > res.Cols() {
		// 	x2 = res.Cols()
		// }
		// if y2 > res.Rows() {
		// 	y2 = res.Rows()
		// }
		// gocv.RectangleWithParams(
		// 	&res,
		// 	image.Rect(l.X, l.Y, x2, y2),
		// 	color.RGBA{R: 255, G: 255, B: 255, A: 255},
		// 	1,
		// 	gocv.Filled,
		// 	0,
		// )
	}
	return points
}

func (m *ImageManager) rawTemplateMatch(src *gocv.Mat, tpl, mask *gocv.Mat) (res gocv.Mat) {
	var s = src
	var t = tpl

	if src.Type() != gocv.MatTypeCV32FC1 {
		m := gocv.NewMatWithSize(src.Rows(), src.Cols(), gocv.MatTypeCV32FC1)
		src.ConvertTo(&m, gocv.MatTypeCV32FC1)
		s = &m
		defer s.Close()
	}
	if tpl.Type() != gocv.MatTypeCV32FC1 {
		m := gocv.NewMatWithSize(tpl.Rows(), tpl.Cols(), gocv.MatTypeCV32FC1)
		tpl.ConvertTo(&m, gocv.MatTypeCV32FC1)
		t = &m
		defer t.Close()
	}
	if mask == nil {
		m := gocv.NewMat()
		mask = &m
		defer mask.Close()
	} else if mask.Type() != gocv.MatTypeCV32FC1 {
		m := gocv.NewMatWithSize(mask.Rows(), mask.Cols(), gocv.MatTypeCV32FC1)
		mask.ConvertTo(&m, gocv.MatTypeCV32FC1)
		mask = &m
		defer mask.Close()
	}

	cols := src.Cols() - tpl.Cols() + 1
	rows := src.Rows() - tpl.Rows() + 1
	res = gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV32FC1)

	gocv.MatchTemplate(*s, *t, &res, gocv.TmSqdiffNormed, *mask)

	return res
}

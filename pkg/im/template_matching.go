package im

import (
	"image"
	"strconv"
	"strings"

	"gocv.io/x/gocv"
)

const (
	thOnePercent      = 0.01
	thOneTenthPercent = 0.001
)

func (m *ImageManager) MatchDefault(src gocv.Mat, p ...string) (bool, image.Point) {
	return m.Match(src, strings.Join(p, "."), thOnePercent)
}

// MatchOneP match should below 0.01 threshold value
func (m *ImageManager) MatchOneP(src gocv.Mat, p string) (bool, image.Point) {
	return m.Match(src, p, thOnePercent)
}

func (m *ImageManager) MatchOneTenth(src gocv.Mat, p string) (bool, image.Point) {
	return m.Match(src, p, thOneTenthPercent)
}

// Match from <p>.1 to <p>.10
func (m *ImageManager) MatchNDefault(src gocv.Mat, p string) (bool, image.Point) {
	return m.Match(src, p, thOnePercent)
}

// Match from <p>.1 to <p>.10
func (m *ImageManager) MatchNOneTenth(src gocv.Mat, p string) (bool, image.Point) {
	return m.Match(src, p, thOneTenthPercent)
}

// Match from <p>.1 to <p>.10
func (m *ImageManager) MatchN(src gocv.Mat, p string, th float32) (bool, image.Point) {
	return m.Match(src, p, th)
}

type MatchOption struct {
	Path      string
	Mask      *gocv.Mat
	HasMask   bool
	Th        float32
	PrintVal  bool
	Normalize bool
}

func (o *MatchOption) applyDefault() {
	if o == nil {
		return
	}
	if o.Th == 0 {
		o.Th = thOnePercent
	}
}

// Match from <p>.1 to <p>.10
func (m *ImageManager) Match(src gocv.Mat, p string, th float32, opt ...MatchOption) (bool, image.Point) {
	o := MatchOption{}
	if len(opt) > 0 {
		o = opt[0]
	}
	o.applyDefault()
	if ok, tpl := m.Get(p + "_mask"); ok {
		o.Mask = tpl
		o.HasMask = true
	} else {
		m := gocv.NewMat()
		o.Mask = &m
	}

	if ok, tpl := m.Get(p); ok {
		if matched, pt := m.match(&src, p, tpl, o); matched {
			return matched, pt
		}
	}
	for i := 1; i <= 10; i++ {
		if ok, tpl := m.Get(p + "_" + strconv.Itoa(i)); ok {
			if matched, pt := m.match(&src, p+"_"+strconv.Itoa(i), tpl, o); matched {
				return true, pt
			}
		} else {
			// not found template, so we break
			break
		}
	}
	return false, image.Point{}
}

func (m *ImageManager) match(src *gocv.Mat, txtTpl string, tpl *gocv.Mat, o MatchOption) (bool, image.Point) {
	var s = src
	var t = tpl
	var mask = o.Mask
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
	if mask.Type() != gocv.MatTypeCV32FC1 {
		m := gocv.NewMatWithSize(mask.Rows(), mask.Cols(), gocv.MatTypeCV32FC1)
		mask.ConvertTo(&m, gocv.MatTypeCV32FC1)
		mask = &m
		defer mask.Close()
	}

	cols := src.Cols() - tpl.Cols() + 1
	rows := src.Rows() - tpl.Rows() + 1
	res := gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV32FC1)
	defer func() {
		res.Close()
	}()

	gocv.MatchTemplate(*s, *t, &res, gocv.TmSqdiffNormed, *mask)
	if o.Normalize {
		gocv.Normalize(res, &res, 0, 1, gocv.NormMinMax)
	}

	// var v2 float32
	v, _, l, _ := gocv.MinMaxLoc(res)
	// if o.HasMask && txtTpl == "swc.btn_top_right_menu" {
	// 	sr := s.Region(image.Rect(l.X, l.Y, l.X+tpl.Cols(), l.Y+tpl.Rows()))
	// 	defer sr.Close()

	// 	t2 := gocv.NewMat()
	// 	defer t2.Close()
	// 	s2 := gocv.NewMat()
	// 	defer s2.Close()
	// 	r2 := gocv.NewMat()
	// 	defer r2.Close()

	// 	gocv.BitwiseAnd(sr, *mask, &s2)
	// 	gocv.BitwiseAnd(*t, *mask, &t2)
	// 	gocv.AbsDiff(s2, t2, &r2)

	// 	gocv.IMWrite(filepath.Join("img", "mask.png"), *mask)
	// 	gocv.IMWrite(filepath.Join("img", "t2.png"), t2)
	// 	gocv.IMWrite(filepath.Join("img", "s2.png"), s2)
	// 	gocv.IMWrite(filepath.Join("img", "r2.png"), r2)

	// 	v2 = float32(r2.Sum().Val1) / float32(mask.Sum().Val1)
	// }
	if o.PrintVal {
		m.log.Debugf("matching %s: %.4f at %v (mask: %v, expected: %.4f)", txtTpl, v, l, o.HasMask, o.Th)
	}
	//m.log.With("path", txtTpl).Debugf("match template min loc: %.4f (expected: %.4f)", v, th)
	if v < o.Th {
		return true, image.Point{X: l.X + tpl.Cols()/2, Y: l.Y + tpl.Rows()/2}
	}

	return false, image.Point{}
}

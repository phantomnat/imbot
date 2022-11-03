package im

import (
	"image"
	"strconv"

	"gocv.io/x/gocv"
)

func (m *ImageManager) RGBMatchOneP(src gocv.Mat, p string) (bool, image.Point) {
	return m.RGBMatch(src, p, thOnePercent)
}

func (m *ImageManager) RGBMatchOneTenth(src gocv.Mat, p string) (bool, image.Point) {
	return m.RGBMatch(src, p, thOneTenthPercent)
}

func (m *ImageManager) RGBMatch(src gocv.Mat, p string, th float32) (bool, image.Point) {
	if ok, tpl := m.GetBGR(p); ok {
		if matched, pt := m.RGBMatchTemplate(src, p, tpl, th); matched {
			return true, pt
		}
	}
	for i := 1; i <= 10; i++ {
		if ok, tpl := m.GetBGR(p + "_" + strconv.Itoa(i)); ok {
			if matched, pt := m.RGBMatchTemplate(src, p+"_"+strconv.Itoa(i), tpl, th); matched {
				return true, pt
			}
		} else {
			break
		}
	}
	return false, image.Point{}
}

func (m *ImageManager) RGBMatchTemplate(src gocv.Mat, txtTpl string, tpl *gocv.Mat, th float32) (bool, image.Point) {
	// grayscale match
	res := m.rawMatchTemplate(&src, tpl)
	// rgb match
	results := m.rgbMatchTemplate(&src, tpl)
	defer func() {
		res.Close()
		for _, r := range results {
			r.Close()
		}
	}()

	// find best possible
	n := res.Cols() * res.Rows()
	for i := 0; i < n; i++ {
		v, _, l, _ := gocv.MinMaxLoc(*res)
		//switch {
		//case strings.Contains(txtTpl, "dispatch"):
		//	switch {
		//	case strings.Contains(txtTpl, "ready"), strings.Contains(txtTpl, "completed"):
		//		m.log.With("path", txtTpl).Debugf("match template min loc: %.4f (expected: %.4f)", v, th)
		//	}
		//}
		if v > th {
			return false, image.Point{}
		}

		isRGBOk := results[0].GetFloatAt(l.Y, l.X) < th &&
			results[1].GetFloatAt(l.Y, l.X) < th &&
			results[2].GetFloatAt(l.Y, l.X) < th
		if isRGBOk {
			return true, l
		}
		// clear that position / find next
		res.SetFloatAt(l.Y, l.X, th+1)
	}
	return false, image.Point{}
}

func (m *ImageManager) rgbMatchTemplate(src *gocv.Mat, tpl *gocv.Mat) (results [3]*gocv.Mat) {
	// rgb match
	srcs := gocv.Split(*src)
	dsts := gocv.Split(*tpl)
	defer func() {
		for _, m := range srcs {
			m.Close()
		}
		for _, m := range dsts {
			m.Close()
		}
	}()
	if len(srcs) != 3 || len(dsts) != 3 {
		return
	}
	for i := 0; i < 3; i++ {
		results[i] = m.rawMatchTemplate(&srcs[i], &dsts[i])
	}
	return
}

// rawMatchTemplate
// auto handle mat type
func (m *ImageManager) rawMatchTemplate(src *gocv.Mat, tpl *gocv.Mat) (result *gocv.Mat) {
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
		defer func() {
			t.Close()
		}()
	}

	cols := src.Cols() - tpl.Cols() + 1
	rows := src.Rows() - tpl.Rows() + 1
	res := gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV32FC1)
	mask := gocv.NewMat()
	defer mask.Close()

	gocv.MatchTemplate(*s, *t, &res, gocv.TmSqdiffNormed, mask)
	result = &res
	return
}

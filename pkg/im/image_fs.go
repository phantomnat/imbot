//go:build !release
// +build !release

package im

import (
	"os"
	"path/filepath"
	"strings"

	"gocv.io/x/gocv"
)

var images = map[string][]byte{}

// func (m *ImageManager) GetFyneResource(name, path string) (ok bool, resource *fyne.StaticResource) {
// 	if data, ok := images[path]; ok {
// 		return true, fyne.NewStaticResource(name, data)
// 	}

// 	p := filepath.Join(
// 		"img",
// 		strings.ReplaceAll(path, ".", string(filepath.Separator))+".png",
// 	)
// 	if _, err := os.Stat(p); err == nil {
// 		data, err := ioutil.ReadFile(p)
// 		if err != nil {
// 			m.log.With("path", path).Errorf("read file: %+v", err)
// 			return
// 		}
// 		resource = fyne.NewStaticResource(name, data)
// 		ok = true
// 	}
// 	return
// }

func (m *ImageManager) Get(path string) (ok bool, mat *gocv.Mat) {
	if _, ok := m.images[path]; ok {
		return true, m.images[path]
	}

	// load from disk
	p := filepath.Join(
		"img",
		strings.ReplaceAll(path, ".", string(filepath.Separator))+".png",
	)
	if _, err := os.Stat(p); err == nil {
		mat := gocv.IMRead(p, gocv.IMReadColor)
		gocv.CvtColor(mat, &mat, gocv.ColorRGBAToBGR)
		defer mat.Close()
		gray := gocv.NewMatWithSize(mat.Rows(), mat.Cols(), gocv.MatTypeCV32FC1)
		mat.ConvertTo(&gray, gocv.MatTypeCV32FC1)
		m.images[path] = &gray
		rgb := gocv.NewMatWithSize(mat.Rows(), mat.Cols(), gocv.MatTypeCV32FC3)
		mat.ConvertTo(&rgb, gocv.MatTypeCV32FC3)
		m.bgrImages[path] = &rgb
		return true, m.images[path]
	}

	return false, nil
}

func (m *ImageManager) GetBGR(path string) (ok bool, mat *gocv.Mat) {
	if _, ok := m.bgrImages[path]; ok {
		return true, m.bgrImages[path]
	}

	// load from disk
	p := filepath.Join(
		"img",
		strings.ReplaceAll(path, ".", string(filepath.Separator))+".png",
	)
	if _, err := os.Stat(p); err == nil {
		mat := gocv.IMRead(p, gocv.IMReadColor)
		gocv.CvtColor(mat, &mat, gocv.ColorRGBAToBGR)
		defer mat.Close()
		gray := gocv.NewMatWithSize(mat.Rows(), mat.Cols(), gocv.MatTypeCV32FC1)
		mat.ConvertTo(&gray, gocv.MatTypeCV32FC1)
		m.images[path] = &gray
		rgb := gocv.NewMatWithSize(mat.Rows(), mat.Cols(), gocv.MatTypeCV32FC3)
		mat.ConvertTo(&rgb, gocv.MatTypeCV32FC3)
		m.bgrImages[path] = &rgb
		return true, m.bgrImages[path]
	}

	return false, nil
}

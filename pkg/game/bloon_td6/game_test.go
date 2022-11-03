package bloon_td6

import (
	"testing"

	. "github.com/onsi/gomega"
	"gocv.io/x/gocv"
)

func TestGetMat(t *testing.T) {
	g := NewGomegaWithT(t)

	game, err := New()
	g.Expect(err).NotTo(HaveOccurred())

	mat, err := game.GetMat()
	g.Expect(err).NotTo(HaveOccurred())
	gocv.IMWrite("my_get_mat.png", mat)
	mat.Close()
}

func TestRobotgoGetMat(t *testing.T) {
	g := NewGomegaWithT(t)

	game, err := New()
	g.Expect(err).NotTo(HaveOccurred())

	mat, err := game.GetMatFromRobotgo()
	g.Expect(err).NotTo(HaveOccurred())
	gocv.IMWrite("robotgo_get_mat.png", mat)
	mat.Close()
}

func BenchmarkMyCapture(b *testing.B) {
	g := NewGomegaWithT(b)

	game, err := New()
	g.Expect(err).NotTo(HaveOccurred())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		func() {
			mat, err := game.GetMat()
			g.Expect(err).NotTo(HaveOccurred())
			defer mat.Close()
		}()
	}
}

func BenchmarkRobotgoCapture(b *testing.B) {
	g := NewGomegaWithT(b)

	game, err := New()
	g.Expect(err).NotTo(HaveOccurred())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		func() {
			mat, err := game.GetMatFromRobotgo()
			g.Expect(err).NotTo(HaveOccurred())
			defer mat.Close()
		}()
	}
}

package huecolor_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rschio/huecolor"
)

func TestXY(t *testing.T) {
	x0, y0, bri0 := float32(0.4962), float32(0.4151), uint8(143)
	xy0 := huecolor.XY{X: x0, Y: y0, Bri: bri0}

	r, g, b, _ := xy0.RGBA()
	x1, y1, bri1 := huecolor.RGBToXY(uint8(r/257), uint8(g/257), uint8(b/257))
	xy1 := huecolor.XY{X: x1, Y: y1, Bri: bri1}

	opt := cmp.Comparer(func(x, y uint8) bool {
		delta := x - y
		if delta < 0 {
			delta = -delta
		}
		return delta <= 3
	})

	if diff := cmp.Diff(xy0, xy1, opt, cmpopts.EquateApprox(0, 0.001)); diff != "" {
		t.Errorf("huecolor.RGBToXY: %s", diff)
	}

	rr, gg, bb := huecolor.XYToRGB(xy1.X, xy1.Y, xy1.Bri)
	x2, y2, bri2 := huecolor.RGBToXY(rr, gg, bb)
	xy2 := huecolor.XY{X: x2, Y: y2, Bri: bri2}
	if diff := cmp.Diff(xy0, xy2, opt, cmpopts.EquateApprox(0, 0.001)); diff != "" {
		t.Errorf("huecolor.RGBToXY: %s", diff)
	}
}

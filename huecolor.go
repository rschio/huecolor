package huecolor

import (
	"math"
)

// XY represents a XY color in CIE color space.
// This type is adapted to work Philips Hue bulbs color gamut.
type XY struct {
	X, Y float32
	Bri  uint8
}

func (c XY) RGBA() (uint32, uint32, uint32, uint32) {
	if c.Y == 0 {
		return 0, 0, 0, 0xffff
	}

	// See https://developers.meethue.com/develop/application-design-guidance/color-conversion-formulas-rgb-to-xy-and-back/
	x := c.X
	y := c.Y
	z := 1 - x - y
	Y := float32(c.Bri) / 254
	X := (Y / y) * x
	Z := (Y / y) * z

	r := X*1.6564920 + Y*-0.354851 + Z*-0.255038
	g := X*-0.707196 + Y*1.6553970 + Z*0.0361520
	b := X*0.0517130 + Y*-0.121364 + Z*1.0115300

	r = xyToRGBGamma(r)
	g = xyToRGBGamma(g)
	b = xyToRGBGamma(b)

	return uint32(to16(r)), uint32(to16(g)), uint32(to16(b)), 0xffff
}

// XYToRGB converts a x, y, bri triple to a RGB triple.
func XYToRGB(x, y float32, bri uint8) (uint8, uint8, uint8) {
	if y == 0 {
		return 0, 0, 0
	}

	z := 1 - x - y
	Y := float32(bri) / 254
	X := (Y / y) * x
	Z := (Y / y) * z

	r := X*1.6564920 + Y*-0.354851 + Z*-0.255038
	g := X*-0.707196 + Y*1.6553970 + Z*0.0361520
	b := X*0.0517130 + Y*-0.121364 + Z*1.0115300

	r = xyToRGBGamma(r)
	g = xyToRGBGamma(g)
	b = xyToRGBGamma(b)

	return to8(r), to8(g), to8(b)
}

// RGBToXY converts a r, g, b triple into a x, y, bri triple.
func RGBToXY(r, g, b uint8) (float32, float32, uint8) {
	_r := rgbToXYGamma(float32(r) / 0xff)
	_g := rgbToXYGamma(float32(g) / 0xff)
	_b := rgbToXYGamma(float32(b) / 0xff)

	// This matrix is the inverse of XYToRGB matrix.
	X := _r*0.664511594963000 + _g*0.1543237180440 + _b*0.1620284095390
	Y := _r*0.283881604719000 + _g*0.6684336040750 + _b*0.0476855704228
	Z := _r*8.81031356321e-05 + _g*0.0723095049022 + _b*0.9860393032600

	x := X / (X + Y + Z)
	y := Y / (X + Y + Z)

	return x, y, uint8(Y * 254)
}

func xyToRGBGamma(c float32) float32 {
	if c <= 0.0031308 {
		return 12.92 * c
	}
	return float32(1.055*math.Pow(float64(c), 1/2.4) - 0.055)
}

func to16(v float32) uint16 {
	v = max(v, 0)
	v = min(v, 1)
	return uint16(v*0xffff + 0.5)
}

func to8(v float32) uint8 {
	v = max(v, 0)
	v = min(v, 1)
	return uint8(v*0xff + 0.5)
}

func rgbToXYGamma(c float32) float32 {
	if c > 0.04045 {
		return float32(math.Pow((float64(c)+0.055)/1.055, 2.4))
	}
	return c / 12.92
}

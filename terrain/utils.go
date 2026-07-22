package terrain

import (
	"math"

	camera "github.com/Guidjy/wireframe/camera"
	"github.com/Guidjy/wireframe/config"
	. "github.com/gen2brain/raylib-go/raylib"
)

var cam *camera.Cam = camera.GetCamInstance()

func drawEdge(v1, v2 Vector3, color Color) {
	p1, _ := cam.ProjectPoint(cam.WorldToCameraSpace(v1))
	p2, _ := cam.ProjectPoint(cam.WorldToCameraSpace(v2))
	DrawLineV(p1, p2, color)
}

func AddColors(c1, c2 Color) Color {
	c := Color{
		R: c1.R + c2.R,
		G: c1.G + c2.G,
		B: c1.B + c2.B,
		A: c1.A + c2.A,
	}
	return c
}

func SubtractColors(c1, c2 Color) Color {
	c := Color{
		R: c1.R - c2.R,
		G: c1.G - c2.G,
		B: c1.B - c2.B,
		A: c1.A - c2.A,
	}
	return c
}

func DivideColors(c1, c2 Color) Color {
	c := Color{
		R: c1.R / c2.R,
		G: c1.G / c2.G,
		B: c1.B / c2.B,
		A: c1.A / c2.A,
	}
	return c
}

func MultiplyColorByFLoat(c1 Color, f float32) Color {
	r := uint8(float32(c1.R) * f)
	g := uint8(float32(c1.G) * f)
	b := uint8(float32(c1.B) * f)
	a := uint8(float32(c1.A) * f)

	c := Color{
		R: r,
		G: g,
		B: b,
		A: a,
	}

	return c
}

func DivideColorByFLoat(c1 Color, f float32) Color {
	if f == 0 {
		return Color{}
	}
	r := uint8(float32(c1.R) / f)
	g := uint8(float32(c1.G) / f)
	b := uint8(float32(c1.B) / f)
	a := uint8(float32(c1.A) / f)

	return Color{R: r, G: g, B: b, A: a}
}

func IsPixelOnScreen(x, y float32) bool {
	return 0 <= int(x) && int(x) < config.ScreenWidth && 0 <= int(y) && int(y) < config.ScreenHeight
}

type ColorFloat struct {
	R, G, B, A float32
}

func ToColorFloat(c Color) ColorFloat {
	return ColorFloat{float32(c.R), float32(c.G), float32(c.B), float32(c.A)}
}

// Safely converts back to uint8, clamping values between 0 and 255 so they never wrap around
func (cf ColorFloat) ToColor() Color {
	return Color{
		R: uint8(math.Min(math.Max(float64(cf.R), 0), 255)),
		G: uint8(math.Min(math.Max(float64(cf.G), 0), 255)),
		B: uint8(math.Min(math.Max(float64(cf.B), 0), 255)),
		A: uint8(math.Min(math.Max(float64(cf.A), 0), 255)),
	}
}

func AddColorFloats(c1, c2 ColorFloat) ColorFloat {
	return ColorFloat{c1.R + c2.R, c1.G + c2.G, c1.B + c2.B, c1.A + c2.A}
}

func SubtractColorFloats(c1, c2 ColorFloat) ColorFloat {
	return ColorFloat{c1.R - c2.R, c1.G - c2.G, c1.B - c2.B, c1.A - c2.A}
}

func DivideColorFloat(c1 ColorFloat, f float32) ColorFloat {
	if f == 0 {
		return ColorFloat{}
	}
	return ColorFloat{c1.R / f, c1.G / f, c1.B / f, c1.A / f}
}

package terrain

import (
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

package terrain

import (
	"cmp"
	"math"
	"slices"

	. "github.com/Guidjy/wireframe/lighting"
	. "github.com/gen2brain/raylib-go/raylib"
)

type Vertex struct {
	Pos Vector3

	Normal Vector3

	Color Color
}

func DirectionToLightSource(v Vector3, ls LightSource) Vector3 {
	return ls.Pos.Subtract(v).Normalize()
}

// Helper to prevent float division-by-zero (which causes +Inf and infinite loop freezes)
func safeDiv(num, den float32) float32 {
	if den == 0 {
		return 0
	}
	return num / den
}

func Rasterize(v0, v1, v2 Vertex) {
	vertices := []Vertex{v0, v1, v2}

	// transofmrs vertices to camera space and projects them
	for i := 0; i < len(vertices); i++ {
		camPos := cam.WorldToCameraSpace(vertices[i].Pos)
		proj2D, depthZ := cam.ProjectPoint(camPos)

		if depthZ <= 0 || math.IsNaN(float64(proj2D.X)) || math.IsNaN(float64(proj2D.Y)) {
			return
		}

		vertices[i].Pos = Vector3{X: proj2D.X, Y: proj2D.Y, Z: depthZ}
	}

	// sorts vertices by screen y position
	slices.SortFunc(vertices, func(va, vb Vertex) int {
		return cmp.Compare(va.Pos.Y, vb.Pos.Y)
	})

	stepXShort := safeDiv(vertices[0].Pos.X-vertices[1].Pos.X, vertices[0].Pos.Y-vertices[1].Pos.Y)
	stepXLong := safeDiv(vertices[0].Pos.X-vertices[2].Pos.X, vertices[0].Pos.Y-vertices[2].Pos.Y)
	x1 := vertices[0].Pos.X
	x2 := vertices[0].Pos.X

	verticalStepZShort := safeDiv(vertices[0].Pos.Z-vertices[1].Pos.Z, vertices[0].Pos.Y-vertices[1].Pos.Y)
	verticalStepZSLong := safeDiv(vertices[0].Pos.Z-vertices[2].Pos.Z, vertices[0].Pos.Y-vertices[2].Pos.Y)
	z1 := vertices[0].Pos.Z
	z2 := vertices[0].Pos.Z

	verticalShortColorStep := DivideColorByFLoat(SubtractColors(vertices[0].Color, vertices[1].Color), vertices[0].Pos.Y-vertices[1].Pos.Y)
	verticalLongColorStep := DivideColorByFLoat(SubtractColors(vertices[0].Color, vertices[2].Color), vertices[0].Pos.Y-vertices[2].Pos.Y)
	color1 := vertices[0].Color
	color2 := vertices[0].Color

	// Paints top part of the triangle
	for y := vertices[0].Pos.Y; y < vertices[1].Pos.Y; y++ {
		var startX, endX, startZ, endZ float32
		var startColor, endColor Color

		if x1 < x2 {
			startX, endX = x1, x2
			startZ, endZ = z1, z2
			startColor, endColor = color1, color2
		} else {
			startX, endX = x2, x1
			startZ, endZ = z2, z1
			startColor, endColor = color2, color1
		}

		horizontalStepZ := safeDiv(startZ-endZ, startX-endX)
		currentZ := startZ

		horizontalStepColor := DivideColorByFLoat(SubtractColors(startColor, endColor), startX-endX)
		currentColor := startColor

		for x := startX; x < endX; x++ {
			if IsPixelOnScreen(x, y) {
				if currentZ < cam.ZBuffer.Depth[int(x)][int(y)] {
					cam.ZBuffer.Depth[int(x)][int(y)] = currentZ
					cam.ZBuffer.Color[int(x)][int(y)] = currentColor
				}
			}
			currentZ += horizontalStepZ
			currentColor = AddColors(currentColor, horizontalStepColor)
		}

		x1 += stepXShort
		x2 += stepXLong
		z1 += verticalStepZShort
		z2 += verticalStepZSLong
		color1 = AddColors(color1, verticalShortColorStep)
		color2 = AddColors(color2, verticalLongColorStep)
	}

	// Paints bottom part of the triangle
	stepXShort = safeDiv(vertices[1].Pos.X-vertices[2].Pos.X, vertices[1].Pos.Y-vertices[2].Pos.Y)
	x1 = vertices[1].Pos.X

	verticalStepZShort = safeDiv(vertices[2].Pos.Z-vertices[1].Pos.Z, vertices[2].Pos.Y-vertices[1].Pos.Y)
	z1 = vertices[1].Pos.Z

	verticalShortColorStep = DivideColorByFLoat(SubtractColors(vertices[2].Color, vertices[1].Color), vertices[2].Pos.Y-vertices[1].Pos.Y)
	color1 = vertices[1].Color

	for y := vertices[1].Pos.Y; y < vertices[2].Pos.Y; y++ {
		var startX, endX, startZ, endZ float32
		var startColor, endColor Color

		if x1 < x2 {
			startX, endX = x1, x2
			startZ, endZ = z1, z2
			startColor, endColor = color1, color2
		} else {
			startX, endX = x2, x1
			startZ, endZ = z2, z1
			startColor, endColor = color2, color1
		}

		horizontalStepZ := safeDiv(startZ-endZ, startX-endX)
		currentZ := startZ

		horizontalStepColor := DivideColorByFLoat(SubtractColors(startColor, endColor), startX-endX)
		currentColor := startColor

		for x := startX; x < endX; x++ {
			if IsPixelOnScreen(x, y) {
				if currentZ < cam.ZBuffer.Depth[int(x)][int(y)] {
					cam.ZBuffer.Depth[int(x)][int(y)] = currentZ
					cam.ZBuffer.Color[int(x)][int(y)] = currentColor
				}
			}
			currentZ += horizontalStepZ
			currentColor = AddColors(currentColor, horizontalStepColor)
		}

		x1 += stepXShort
		x2 += stepXLong
		z1 += verticalStepZShort
		z2 += verticalStepZSLong
		color1 = AddColors(color1, verticalShortColorStep)
		color2 = AddColors(color2, verticalLongColorStep)
	}
}

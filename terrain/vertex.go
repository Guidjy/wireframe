package terrain

import (
	. "github.com/Guidjy/wireframe/lighting"
	. "github.com/gen2brain/raylib-go/raylib"
)

type Vertex struct {
	Pos Vector3

	Normal Vector3

	color Color
}

func (v Vertex) DirectionToLightSource(ls LightSource) Vector3 {
	return ls.Pos.Subtract(v.Pos).Normalize()
}

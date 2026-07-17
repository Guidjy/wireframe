package terrain

import (
	"math"

	. "github.com/Guidjy/wireframe/camera"
	. "github.com/gen2brain/raylib-go/raylib"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var rotationOffset float64 = 0.0

func drawEdge(v1 Vector3, v2 Vector3) {
	cam := GetCamInstance()

	p1 := cam.ProjectPoint(cam.AlignPoint(v1))
	p2 := cam.ProjectPoint(cam.AlignPoint(v2))

	rl.DrawLineV(p1, p2, rl.White)
}

// Draws a cube in world origin for testing
func RenderCube(size float32) {
	step := (math.Pi * 2) / 4
	r := (float64(size) * math.Sqrt(2)) / 2
	vertices := make([]Vector3, 8)

	rotationOffset += 0.0001

	// Creates vertices
	for i := 0; i < 4; i++ {
		theta := step * float64(i)

		x := float32(r * math.Cos(theta+float64(rotationOffset)))
		z := float32(r * math.Sin(theta+float64(rotationOffset)))

		v1 := Vector3{X: x, Y: 0, Z: z}
		v2 := Vector3{X: x, Y: size, Z: z}

		vertices[i] = v1
		vertices[i+4] = v2
	}

	// Draws edges
	for i := 0; i < 4; i++ {
		nextEdgeIndex := (i + 1) % 4
		// Bottom face
		drawEdge(vertices[i], vertices[nextEdgeIndex])
		// Top face
		drawEdge(vertices[i+4], vertices[nextEdgeIndex+4])
		// Vertical Edges
		drawEdge(vertices[i], vertices[i+4])
	}
}

package terrain

import (
	"cmp"
	"math"
	"math/rand/v2"
	"slices"

	"github.com/Guidjy/wireframe/camera"
	config "github.com/Guidjy/wireframe/config"
	. "github.com/Guidjy/wireframe/lighting"
	. "github.com/gen2brain/raylib-go/raylib"
)

const minControlPointDeltaY = 100
const maxControlPointDeltaY = 700
const minControlPointCount = 10
const maxControlPointCount = 30

type Terrain struct {
	controlPoints [][]Vector3

	vertices []Vertex

	controlPointDeltaY int

	controlPointCount int

	ball Ball

	sun LightSource
}

// B(t) B-Spline base functions
func b0(t float32) float32 {
	return 1.0 / 6.0 * float32(math.Pow(float64(1-t), 3))
}
func b1(t float32) float32 {
	return 1.0 / 6.0 * float32((3*math.Pow(float64(t), 3) - 6*math.Pow(float64(t), 2) + 4))
}
func b2(t float32) float32 {
	return 1.0 / 6.0 * float32((-3*math.Pow(float64(t), 3) + 3*math.Pow(float64(t), 2) + 3*float64(t) + 1))
}
func b3(t float32) float32 {
	return 1.0 / 6.0 * float32(math.Pow(float64(t), 3))
}

func (terrain *Terrain) Init() {
	terrain.controlPointDeltaY = minControlPointDeltaY
	terrain.controlPointCount = minControlPointCount

	terrain.ball.init()

	terrain.sun.Init(Vector3{X: 0, Y: 200, Z: 0})

	terrain.GenerateTerrain()
}

// Generates a B-Spline terrain with nxn control points
func (terrain *Terrain) GenerateTerrain() {
	n := terrain.controlPointCount

	terrain.controlPoints = make([][]Vector3, n)
	for i := 0; i < n; i++ {
		terrain.controlPoints[i] = make([]Vector3, n)
	}

	hw := float32(config.TerrainWidth) / 2.0
	step := float32(config.TerrainWidth) / float32(n-1)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			var controlPoint Vector3

			controlPoint.X = step*float32(i) - hw
			controlPoint.Y = rand.Float32() * float32(terrain.controlPointDeltaY)
			controlPoint.Z = step*float32(j) - hw

			terrain.controlPoints[i][j] = controlPoint
		}
	}

	terrain.calculateTerrainVertices()
}

// Returns an interpolated point in a patch of the B-Spline surface and it's normal vector
func (terrain Terrain) calculateSplinePointAndNormal(patchX int, patchZ int, s float32, t float32) (Vector3, Vector3) {
	// A point in a patch of the B-Spline surface is defined by this formula:
	// Q(s, t) = ∑[𝑖=0, 3] ∑[𝑖=0, 3] 𝑃i,j * Bi(s) * Bj(t)
	// Each patch of the surface is defined by a 4x4 block of control points

	p0 := Vector3Zero() // Q(s, t)
	p1 := Vector3Zero() // Q(s+0.01, t)
	p2 := Vector3Zero() // Q(s, t+0.01)

	// weights of the base functions for s and t
	bs0 := []float32{b0(s), b1(s), b2(s), b3(s)}
	bt0 := []float32{b0(t), b1(t), b2(t), b3(t)}

	bs1 := []float32{b0(s + 0.01), b1(s + 0.01), b2(s + 0.01), b3(s + 0.01)}
	bt1 := bt0

	bs2 := bs0
	bt2 := []float32{b0(t + 0.01), b1(t + 0.01), b2(t + 0.01), b3(t + 0.01)}

	// Sum of the influences of the 4x4 control points (∑[𝑖=0, 3] ∑[𝑖=0, 3] 𝑃i,j * Bi(s) * Bj(t))
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			controlPoint := terrain.controlPoints[patchX+i][patchZ+j] // Pij

			weight := bs0[i] * bt0[j]               // Bi(s) * Bj(t)
			p0 = p0.Add(controlPoint.Scale(weight)) // Adds weighted influence of controlPoints[i][j]

			weight = bs1[i] * bt1[j]
			p1 = p1.Add(controlPoint.Scale(weight))

			weight = bs2[i] * bt2[j]
			p2 = p2.Add(controlPoint.Scale(weight))
		}
	}

	normal := p1.Subtract(p0).CrossProduct(p2.Subtract(p0)).Normalize()

	return p0, normal
}

// Returns a point in the terrain surface given global (x, z) coordinates
func (terrain Terrain) GetSurfacePoint(x float32, z float32) Vector3 {
	hw := float32(config.TerrainWidth) / 2.0
	step := float32(config.TerrainWidth) / float32(len(terrain.controlPoints)-1)

	// Normalizes global coords to the grid space
	gridX := float32(x+hw)/step - 1.0 // -1 because the B-Spline doesn't actually touch the edges
	gridZ := float32(z+hw)/step - 1.0

	patchX := int(math.Floor(float64(gridX)))
	patchZ := int(math.Floor(float64(gridZ)))

	s := gridX - float32(patchX)
	t := gridZ - float32(patchZ)

	// makes sure that an out of bounds control point won't be accessed by the window
	maxPatch := len(terrain.controlPoints) - 4
	if patchX < 0 {
		patchX = 0
		s = 0.0
	} else if patchX > maxPatch {
		patchX = maxPatch
		s = 1.0
	}
	if patchZ < 0 {
		patchZ = 0
		t = 0.0
	} else if patchZ > maxPatch {
		patchZ = maxPatch
		t = 1.0
	}

	p, _ := terrain.calculateSplinePointAndNormal(patchX, patchZ, s, t)

	return p
}

// Calculates all of the terrain's vertices and their normals, and stores them in terrain.vertices
func (terrain *Terrain) calculateTerrainVertices() {
	terrain.vertices = terrain.vertices[:0]

	for i := 0; i < len(terrain.controlPoints)-3; i++ {
		for j := 0; j < len(terrain.controlPoints)-3; j++ {

			const step float32 = 0.25
			for s := float32(0); s < 1.0; s += step {
				for t := float32(0); t < 1.0; t += step {
					v0, n0 := terrain.calculateSplinePointAndNormal(i, j, s, t)           // bottom-left
					v1, n1 := terrain.calculateSplinePointAndNormal(i, j, s+step, t)      // top-left
					v2, n2 := terrain.calculateSplinePointAndNormal(i, j, s+step, t+step) // top-right
					v3, n3 := terrain.calculateSplinePointAndNormal(i, j, s, t+step)      // bottom-right

					terrain.vertices = append(terrain.vertices, Vertex{v0, n0, Green}, Vertex{v1, n1, Green}, Vertex{v2, n2, Green}, Vertex{v3, n3, Green})
				}
			}

		}
	}

}

// Renders the terrain as a wireframe mesh
func (terrain Terrain) RenderWifreframe() {
	cam := camera.GetCamInstance()

	for i := 0; i < len(terrain.vertices); i += 4 {
		if i+4 > len(terrain.vertices) {
			return
		}

		p0 := terrain.vertices[i].Pos
		p1 := terrain.vertices[i+1].Pos
		p2 := terrain.vertices[i+2].Pos
		p3 := terrain.vertices[i+3].Pos

		p0 = cam.WorldToCameraSpace(p0)
		p1 = cam.WorldToCameraSpace(p1)
		p2 = cam.WorldToCameraSpace(p2)
		p3 = cam.WorldToCameraSpace(p3)

		if config.ShouldCullHiddenFaces {
			edge1 := p1.Subtract(p0)
			edge2 := p2.Subtract(p0)
			normal := edge1.CrossProduct(edge2)

			viewDir := p0
			if viewDir.DotProduct(normal) > 0 {
				continue
			}
		}

		projectedP0, _ := cam.ProjectPoint(p0)
		projectedP1, _ := cam.ProjectPoint(p1)
		projectedP2, _ := cam.ProjectPoint(p2)
		projectedP3, _ := cam.ProjectPoint(p3)

		DrawLineV(projectedP0, projectedP1, White)
		DrawLineV(projectedP1, projectedP2, White)
		DrawLineV(projectedP2, projectedP3, White)
		DrawLineV(projectedP3, projectedP0, White)
		DrawLineV(projectedP0, projectedP2, White)
	}

}

func (terrain Terrain) rasterize(v0, v1, v2 Vertex) {
	//cam := camera.GetCamInstance()

	// sorts vertices by their y value
	vertices := []Vertex{v0, v1, v2}
	slices.SortFunc(vertices, func(va, vb Vertex) int {
		return cmp.Compare(va.Pos.Y, vb.Pos.Y)
	})

	// paints top part of the triangle
	stepXShort := (vertices[0].Pos.X - vertices[1].Pos.X) / (vertices[0].Pos.Y - vertices[1].Pos.Y)
	stepXLong := (vertices[0].Pos.X - vertices[2].Pos.X) / (vertices[0].Pos.Y - vertices[2].Pos.Y)

	x1 := vertices[0].Pos.X
	x2 := vertices[0].Pos.X

	for y := vertices[0].Pos.Y; y < vertices[1].Pos.Y; y++ {
		//startX := min(x1, x2)
		//endX := max(x1, x2)

		for x := x1; x < x2; x++ {
			// calculate (x, y) piel depth and potentially update zbuffer
			// calculate pixel color and paint
		}

		x1 += stepXShort
		x2 += stepXLong
	}

}

func (terrain *Terrain) updateControlPointDeltaY(increase bool) {
	if increase {
		terrain.controlPointDeltaY += 100
	} else {
		terrain.controlPointDeltaY -= 100
	}

	if terrain.controlPointDeltaY > maxControlPointDeltaY {
		terrain.controlPointDeltaY = maxControlPointDeltaY
		return
	} else if terrain.controlPointDeltaY < minControlPointDeltaY {
		terrain.controlPointDeltaY = minControlPointDeltaY
		return
	}

	terrain.GenerateTerrain()
}

func (terrain *Terrain) updateControlPointCount(increase bool) {
	if increase {
		terrain.controlPointCount += 5
	} else {
		terrain.controlPointCount -= 5
	}

	if terrain.controlPointCount > maxControlPointCount {
		terrain.controlPointCount = maxControlPointCount
		return
	} else if terrain.controlPointCount < minControlPointCount {
		terrain.controlPointCount = minControlPointCount
		return
	}

	terrain.GenerateTerrain()
}

func (terrain *Terrain) handleKeyboardInput() {
	// increase/decrease terrain hill size
	if IsKeyPressed(KeyZ) {
		terrain.updateControlPointDeltaY(true)
	}
	if IsKeyPressed(KeyX) {
		terrain.updateControlPointDeltaY(false)
	}
	// increase/decrease terrain resolution
	if IsKeyPressed(KeyC) {
		terrain.updateControlPointCount(true)
	}
	if IsKeyPressed(KeyV) {
		terrain.updateControlPointCount(false)
	}
	// toggle backface culling
	if IsKeyPressed(KeyB) {
		config.ShouldCullHiddenFaces = !config.ShouldCullHiddenFaces
	}
}

func (terrain *Terrain) Update() {
	terrain.handleKeyboardInput()
	terrain.RenderWifreframe()

	terrain.ball.Update()
	terrain.ball.Pos.Y = terrain.GetSurfacePoint(terrain.ball.Pos.X, terrain.ball.Pos.Z).Y + terrain.ball.Radius/2

}

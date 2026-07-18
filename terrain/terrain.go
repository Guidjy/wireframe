package terrain

import (
	"math"
	"math/rand/v2"

	camera "github.com/Guidjy/wireframe/camera"
	config "github.com/Guidjy/wireframe/config"
	. "github.com/gen2brain/raylib-go/raylib"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const minControlPointDeltaY = 100
const maxControlPointDeltaY = 700
const minControlPointCount = 10
const maxControlPointCount = 30

type Terrain struct {
	controlPoints [][]Vector3

	controlPointDeltaY int

	controlPointCount int

	shouldCullHiddenFaces bool
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
	terrain.shouldCullHiddenFaces = true
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
}

// Returns an interpolated point in a patch of the B-Spline surface
func (terrain Terrain) calculateSplinePoint(patchX int, patchZ int, s float32, t float32) Vector3 {
	// A point in a patch of the B-Spline surface is defined by this formula:
	// Q(s, t) = ∑[𝑖=0, 3] ∑[𝑖=0, 3] 𝑃i,j * Bi(s) * Bj(t)
	// Each patch of the surface is defined by a 4x4 block of control points

	p := Vector3Zero() // Q(s, t)

	// weights of the base functions for s and t
	bs := []float32{b0(s), b1(s), b2(s), b3(s)}
	bt := []float32{b0(t), b1(t), b2(t), b3(t)}

	// Sum of the influences of the 4x4 control points (∑[𝑖=0, 3] ∑[𝑖=0, 3] 𝑃i,j * Bi(s) * Bj(t))
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			controlPoint := terrain.controlPoints[patchX+i][patchZ+j] // Pij
			weight := bs[i] * bt[j]                                   // Bi(s) * Bj(t)
			p = p.Add(controlPoint.Scale(weight))                     // Adds weighted influence of controlPoints[i][j]
		}
	}

	return p
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

	return terrain.calculateSplinePoint(patchX, patchZ, s, t)
}

func (terrain Terrain) Render() {
	cam := camera.GetCamInstance()

	for i := 0; i < len(terrain.controlPoints)-3; i++ {
		for j := 0; j < len(terrain.controlPoints)-3; j++ {

			const step float32 = 0.25
			for s := float32(0); s < 1.0; s += step {
				for t := float32(0); t < 1.0; t += step {

					p0 := terrain.calculateSplinePoint(i, j, s, t)           // bottom-left
					p1 := terrain.calculateSplinePoint(i, j, s+step, t)      // top-left
					p2 := terrain.calculateSplinePoint(i, j, s, t+step)      // bottom-right
					p3 := terrain.calculateSplinePoint(i, j, s+step, t+step) // top-right

					// TODO: backface culling

					projectedP0 := cam.ProjectPoint(cam.WorldToCameraSpace(p0))
					projectedP1 := cam.ProjectPoint(cam.WorldToCameraSpace(p1))
					projectedP2 := cam.ProjectPoint(cam.WorldToCameraSpace(p2))
					projectedP3 := cam.ProjectPoint(cam.WorldToCameraSpace(p3))

					rl.DrawLineV(projectedP0, projectedP1, rl.White)
					rl.DrawLineV(projectedP1, projectedP3, rl.White)
					rl.DrawLineV(projectedP3, projectedP2, rl.White)
					rl.DrawLineV(projectedP2, projectedP0, rl.White)
					rl.DrawLineV(projectedP0, projectedP3, rl.White)
				}
			}

		}
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
	if rl.IsKeyPressed(rl.KeyZ) {
		terrain.updateControlPointDeltaY(true)
	}
	if rl.IsKeyPressed(rl.KeyX) {
		terrain.updateControlPointDeltaY(false)
	}
	// increase/decrease terrain resolution
	if rl.IsKeyPressed(rl.KeyC) {
		terrain.updateControlPointCount(true)
	}
	if rl.IsKeyPressed(rl.KeyV) {
		terrain.updateControlPointCount(false)
	}
}

func (terrain *Terrain) Update() {
	terrain.handleKeyboardInput()

	terrain.Render()
}

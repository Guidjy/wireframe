package camera

import (
	"math"
	"sync"

	config "github.com/Guidjy/wireframe/config"
	. "github.com/gen2brain/raylib-go/raylib"
)

var cameraInstance *Cam
var once sync.Once

type Cam struct {
	eye    Vector3 // Camera position
	target Vector3 // Point to which the camera is looking at
	up     Vector3 // Defines which direction is up so that the camera isn't on its side

	yaw   float32 // Horizontal camera rotation
	pitch float32 // Vertical camera rotation

	n Vector3
	u Vector3
	v Vector3

	focalLength float32 // Distance from the projection plane
}

func GetCamInstance() *Cam {
	// once.Do guarantees that the function inside of it runs only once. Important for thread-safe initialization (probably not going to matter tho ¯\_(ツ)_/¯)
	once.Do(func() {
		cameraInstance = &Cam{
			eye:         Vector3{X: 0, Y: 5, Z: -10},
			target:      Vector3Zero(),
			up:          Vector3{X: 0, Y: 1, Z: 0},
			yaw:         0,
			pitch:       0,
			n:           Vector3Zero(),
			u:           Vector3Zero(),
			v:           Vector3Zero(),
			focalLength: config.CameraFolcaLength,
		}
	})

	return cameraInstance
}

// TODO: review n and forward

// Aligns camera and world space axes
func (cam *Cam) ChangeBasis() {
	forward := cam.target.Subtract(cam.eye)

	cam.n = forward.Normalize()
	cam.u = cam.up.CrossProduct(forward).Normalize()
	cam.v = cam.n.CrossProduct(cam.u)
}

// Aligns a point in 3D space with the camera. Should be done after chagin bases
func (cam Cam) AlignPoint(p Vector3) Vector3 {
	relativeP := p.Subtract(cam.eye)
	newPoint := relativeP

	newPoint.X = relativeP.DotProduct(cam.u)
	newPoint.Y = relativeP.DotProduct(cam.v)
	newPoint.Z = relativeP.DotProduct(cam.n)

	// clipping
	if newPoint.Z <= 0 {
		newPoint.X = float32(math.Inf(1))
		newPoint.Y = float32(math.Inf(1))
		newPoint.Z = float32(math.Inf(1))
	}

	return newPoint
}

func (cam Cam) ProjectPoint(p3d Vector3) Vector2 {
	var p2d Vector2

	projectedX := p3d.X * cam.focalLength / p3d.Z
	projectedY := p3d.Y * cam.focalLength / p3d.Z

	// Centers points on screen.
	p2d.X = projectedX + float32(config.ScreenWidth)/2.0
	p2d.Y = -projectedY + float32(config.ScreenHeight)/2.0
	return p2d
}

func (cam *Cam) Update() {
	cam.ChangeBasis()
}

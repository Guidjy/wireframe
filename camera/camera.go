package camera

import (
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
			eye:         Vector3One(),
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

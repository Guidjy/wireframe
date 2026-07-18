package camera

import (
	"math"
	"sync"

	config "github.com/Guidjy/wireframe/config"
	. "github.com/gen2brain/raylib-go/raylib"
)

var cameraInstance *Cam
var once sync.Once

const movementSpeed = 100
const mouseSensitivity = 0.003
const maxPitch = float32(math.Pi/2) - 0.01 // used to prevent the camera fro flipping upside down

type Cam struct {
	eye     Vector3 // Camera position
	target  Vector3 // Point to which the camera is looking at
	worldUp Vector3 // Defines which direction is up so that the camera isn't on its side

	yaw   float32 // Horizontal camera rotation
	pitch float32 // Vertical camera rotation

	forward  Vector3 // n
	right    Vector3 // u
	cameraUp Vector3 // v

	focalLength float32 // Distance from the projection plane

}

func GetCamInstance() *Cam {
	// once.Do guarantees that the function inside of it runs only once. Important for thread-safe initialization (probably not going to matter tho ¯\_(ツ)_/¯)
	once.Do(func() {
		cameraInstance = &Cam{
			eye:         Vector3{X: 0, Y: 5, Z: -10},
			target:      Vector3Zero(),
			worldUp:     Vector3{X: 0, Y: 1, Z: 0},
			yaw:         0,
			pitch:       0,
			forward:     Vector3Zero(),
			right:       Vector3Zero(),
			cameraUp:    Vector3Zero(),
			focalLength: config.CameraFolcaLength,
		}
	})

	return cameraInstance
}

// Updates the camera's local nuv basis vectors (forward, right, up)
func (cam *Cam) UpdateBasis() {
	forward := cam.target.Subtract(cam.eye)

	cam.forward = forward.Normalize()
	cam.right = forward.CrossProduct(cam.worldUp).Normalize()
	cam.cameraUp = cam.right.CrossProduct(cam.forward)
}

// Transforms a 3D point from world space into camera space
func (cam Cam) WorldToCameraSpace(p Vector3) Vector3 {
	relativeP := p.Subtract(cam.eye)
	newPoint := relativeP

	// rotates the point in world space so that the camera acts as the origin
	newPoint.X = relativeP.DotProduct(cam.right)
	newPoint.Y = relativeP.DotProduct(cam.cameraUp)
	newPoint.Z = relativeP.DotProduct(cam.forward)

	// clipping
	if newPoint.Z <= 0 {
		newPoint.X = float32(math.Inf(1))
		newPoint.Y = float32(math.Inf(1))
		newPoint.Z = float32(math.Inf(1))
	}

	return newPoint
}

// Projects a point in 3D space onto the projection plane
func (cam Cam) ProjectPoint(p3d Vector3) Vector2 {
	var p2d Vector2

	projectedX := p3d.X * cam.focalLength / p3d.Z
	projectedY := p3d.Y * cam.focalLength / p3d.Z

	// Centers points on screen.
	p2d.X = projectedX + float32(config.ScreenWidth)/2.0
	p2d.Y = -projectedY + float32(config.ScreenHeight)/2.0
	return p2d
}

func (cam *Cam) handleKeyboardInput() {
	velocity := Vector3Zero()
	step := movementSpeed * GetFrameTime()

	if IsKeyDown(KeyW) {
		velocity = velocity.Add(cam.forward.Scale(step))
	}
	if IsKeyDown(KeyS) {
		velocity = velocity.Subtract(cam.forward.Scale(step))
	}
	if IsKeyDown(KeyD) {
		velocity = velocity.Add(cam.right.Scale(step))
	}
	if IsKeyDown(KeyA) {
		velocity = velocity.Subtract(cam.right.Scale(step))
	}
	if IsKeyDown(KeyE) {
		velocity = velocity.Add(cam.cameraUp.Scale(step))
	}
	if IsKeyDown(KeyQ) {
		velocity = velocity.Subtract(cam.cameraUp.Scale(step))
	}

	cam.eye = cam.eye.Add(velocity)
	cam.target = cam.target.Add(velocity)
}

func (cam *Cam) handleMouseMovement() {
	delta := GetMouseDelta()

	cam.yaw -= delta.X * mouseSensitivity
	cam.pitch -= delta.Y * mouseSensitivity

	if cam.pitch > maxPitch {
		cam.pitch = maxPitch
	} else if cam.pitch < -maxPitch {
		cam.pitch = -maxPitch
	}

	var forward Vector3
	forward.X = float32(math.Cos(float64(cam.pitch)) * math.Sin(float64(cam.yaw)))
	forward.Y = float32(math.Sin(float64(cam.pitch)))
	forward.Z = float32(math.Cos(float64(cam.pitch)) * math.Cos(float64(cam.yaw)))

	cam.target = cam.eye.Add(forward.Normalize())
}

func (cam *Cam) Update() {
	cam.handleMouseMovement()
	cam.handleKeyboardInput()

	cam.UpdateBasis()

}

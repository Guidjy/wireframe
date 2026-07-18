package terrain

import (
	camera "github.com/Guidjy/wireframe/camera"
	. "github.com/gen2brain/raylib-go/raylib"
)

var cam *camera.Cam = camera.GetCamInstance()

func drawEdge(v1, v2 Vector3, color Color) {
	p1, _ := cam.ProjectPoint(cam.WorldToCameraSpace(v1))
	p2, _ := cam.ProjectPoint(cam.WorldToCameraSpace(v2))
	DrawLineV(p1, p2, color)
}

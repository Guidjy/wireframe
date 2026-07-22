package main

import (
	. "github.com/Guidjy/wireframe/camera"
	"github.com/Guidjy/wireframe/config"
	. "github.com/Guidjy/wireframe/hud"
	. "github.com/Guidjy/wireframe/terrain"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	rl.InitWindow(config.ScreenWidth, config.ScreenHeight, "Guidjy's Wireframe")

	rl.DisableCursor()

	cam := GetCamInstance()

	var terrain Terrain
	terrain.Init()

	config.Init()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.DarkGray)

		cam.Update()
		terrain.Update()

		DisplayHUD()

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

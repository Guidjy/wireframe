package main

import (
	"fmt"

	. "github.com/Guidjy/wireframe/camera"
	. "github.com/Guidjy/wireframe/config"
	. "github.com/Guidjy/wireframe/terrain"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	rl.InitWindow(ScreenWidth, ScreenHeight, "raylib [core] example - basic window")

	if ShouldLimitFPS {
		rl.SetTargetFPS(60)
	}

	rl.DisableCursor()

	cam := GetCamInstance()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.DarkGray)

		rl.DrawText(fmt.Sprintf("Current FPS: %d", rl.GetFPS()), 0, 0, 20, rl.LightGray)

		cam.Update()

		RenderCube(2)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

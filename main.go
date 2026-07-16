package main

import (
	"fmt"

	. "github.com/Guidjy/wireframe/config"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	rl.InitWindow(ScreenWidth, ScreenHeight, "raylib [core] example - basic window")

	if ShouldLimitFPS {
		rl.SetTargetFPS(60)
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.DarkGray)

		rl.DrawText(fmt.Sprintf("Current FPS: %d", rl.GetFPS()), 0, 0, 20, rl.LightGray)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

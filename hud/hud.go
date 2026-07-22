package hud

import (
	"fmt"

	. "github.com/Guidjy/wireframe/config"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func DisplayHUD() {
	rl.DrawText(fmt.Sprintf("Current FPS: %d", rl.GetFPS()), ScreenWidth-200, 0, 20, rl.LightGray)

	rl.DrawText("Move camera: WASD - QE", 2, 0, 20, rl.LightGray)
	rl.DrawText("Move ball: Arrow Keys", 2, 22, 20, rl.LightGray)
	rl.DrawText("Increase/Decrease Hill Size: Z/X", 2, 44, 20, rl.LightGray)
	rl.DrawText("Increase/Decrease Terrain Resolution: C/V", 2, 66, 20, rl.LightGray)
	rl.DrawText("Toggle Backface Culling: B", 2, 88, 20, rl.LightGray)
	rl.DrawText("Toggle Terrain Rasterization: R", 2, 110, 20, rl.LightGray)
}

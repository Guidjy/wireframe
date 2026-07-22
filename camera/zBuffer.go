package camera

import (
	"math"

	config "github.com/Guidjy/wireframe/config"
	. "github.com/gen2brain/raylib-go/raylib"
)

type ZBuffer struct {
	Depth [][]float32

	Color [][]Color
}

func (zb *ZBuffer) Init() {
	zb.Depth = make([][]float32, config.ScreenWidth)
	zb.Color = make([][]Color, config.ScreenWidth)

	for i := 0; i < config.ScreenWidth; i++ {
		zb.Depth[i] = make([]float32, config.ScreenHeight)
		zb.Color[i] = make([]Color, config.ScreenHeight)

		for j := 0; j < config.ScreenHeight; j++ {
			zb.Depth[i][j] = float32(math.Inf(1))
			zb.Color[i][j] = DarkGray
		}
	}
}

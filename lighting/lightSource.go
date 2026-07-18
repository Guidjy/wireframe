package lighting

import . "github.com/gen2brain/raylib-go/raylib"

type LightSource struct {
	Pos Vector3
}

func (ls *LightSource) Init(pos Vector3) {
	ls.Pos = pos
}

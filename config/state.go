package config

var ShouldCullHiddenFaces bool

var ShouldApplyShading bool

var ShouldRasterizeTerrain bool

func Init() {
	ShouldCullHiddenFaces = true
	ShouldApplyShading = false
	ShouldRasterizeTerrain = false
}

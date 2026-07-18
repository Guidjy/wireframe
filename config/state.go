package config

var ShouldCullHiddenFaces bool

var ShouldApplyShading bool

func Init() {
	ShouldCullHiddenFaces = true
	ShouldApplyShading = false
}

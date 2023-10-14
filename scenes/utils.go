package scenes

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/vcgj7-game/assets"
	"github.com/quasilyte/vcgj7-game/eui"
)

func initUI(scene *ge.Scene, root *widget.Container) {
	bg := scene.NewSprite(assets.ImageMenuBg)
	bg.Centered = false
	scene.AddGraphics(bg)

	uiObject := eui.NewSceneObject(root)
	scene.AddGraphics(uiObject)
	scene.AddObject(uiObject)
}

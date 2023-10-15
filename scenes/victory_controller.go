package scenes

import (
	"os"
	"runtime"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/vcgj7-game/assets"
	"github.com/quasilyte/vcgj7-game/eui"
	"github.com/quasilyte/vcgj7-game/session"
	"github.com/quasilyte/vcgj7-game/styles"
)

type VictoryController struct {
	scene *ge.Scene
	state *session.State
}

func NewVictoryController(state *session.State) *VictoryController {
	return &VictoryController{state: state}
}

func (c *VictoryController) Init(scene *ge.Scene) {
	c.scene = scene

	scene.Audio().PauseCurrentMusic()

	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	rowContainer := eui.NewRowLayoutContainerWithMinWidth(320, 8, nil)
	root.AddChild(rowContainer)

	bigFont := assets.BitmapFont3

	rowContainer.AddChild(eui.NewCenteredLabel("Victory!", bigFont))

	rowContainer.AddChild(eui.NewSeparator(nil, styles.TransparentColor))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, "BACK TO MENU", func() {
		scene.Context().ChangeScene(NewMainMenuController(c.state))
	}))

	if runtime.GOARCH != "wasm" {
		rowContainer.AddChild(eui.NewButton(c.state.UIResources, "EXIT", func() {
			os.Exit(0)
		}))
	}

	initUI(scene, root)
}

func (c *VictoryController) Update(delta float64) {}

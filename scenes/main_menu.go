package scenes

import (
	"fmt"
	"os"
	"runtime"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/vcgj7-game/assets"
	"github.com/quasilyte/vcgj7-game/eui"
	"github.com/quasilyte/vcgj7-game/gamedata"
	"github.com/quasilyte/vcgj7-game/session"
	"github.com/quasilyte/vcgj7-game/styles"
)

type MainMenuController struct {
	state *session.State
}

func NewMainMenuController(state *session.State) *MainMenuController {
	return &MainMenuController{state: state}
}

func (c *MainMenuController) Init(scene *ge.Scene) {
	scene.Audio().SetGroupVolume(assets.SoundGroupEffect, assets.VolumeMultiplier(c.state.Settings.SoundLevel))
	scene.Audio().SetGroupVolume(assets.SoundGroupMusic, assets.VolumeMultiplier(c.state.Settings.MusicLevel))

	scene.Audio().PauseCurrentMusic()

	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	rowContainer := eui.NewRowLayoutContainerWithMinWidth(320, 8, nil)
	root.AddChild(rowContainer)

	bigFont := assets.BitmapFont3
	tinyFont := assets.BitmapFont1

	rowContainer.AddChild(eui.NewCenteredLabel("Pixelspace Rangers", bigFont))

	rowContainer.AddChild(eui.NewSeparator(nil, styles.TransparentColor))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, "PLAY", func() {
		c.state.World = gamedata.NewWorld(scene.Rand())
		scene.Context().ChangeScene(NewChoiceController(c.state))
	}))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, "SETTINGS", func() {
		scene.Context().ChangeScene(NewSettingsController(c.state))
	}))

	b := eui.NewButton(c.state.UIResources, "CREDITS", func() {
		// TODO
	})
	b.GetWidget().Disabled = true
	rowContainer.AddChild(b)

	if runtime.GOARCH != "wasm" {
		rowContainer.AddChild(eui.NewButton(c.state.UIResources, "EXIT", func() {
			os.Exit(0)
		}))
	}

	rowContainer.AddChild(eui.NewSeparator(nil, styles.TransparentColor))
	rowContainer.AddChild(eui.NewCenteredLabel(fmt.Sprintf("#vas3kclubjam build %d", currentBuild), tinyFont))

	initUI(scene, root)
}

func (c *MainMenuController) Update(delta float64) {}

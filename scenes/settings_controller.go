package scenes

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/vcgj7-game/assets"
	"github.com/quasilyte/vcgj7-game/controls"
	"github.com/quasilyte/vcgj7-game/eui"
	"github.com/quasilyte/vcgj7-game/session"
	"github.com/quasilyte/vcgj7-game/styles"
)

type settingsController struct {
	state *session.State
	scene *ge.Scene
}

func NewSettingsController(state *session.State) *settingsController {
	return &settingsController{
		state: state,
	}
}

func (c *settingsController) Init(scene *ge.Scene) {
	c.scene = scene
	c.initUI(scene)
}

func (c *settingsController) initUI(scene *ge.Scene) {
	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	rowContainer := eui.NewRowLayoutContainerWithMinWidth(400, 16, nil)
	root.AddChild(rowContainer)

	rowContainer.AddChild(eui.NewCenteredLabel("Settings", assets.BitmapFont2))

	rowContainer.AddChild(eui.NewSelectButton(eui.SelectButtonConfig{
		Resources:  c.state.UIResources,
		Input:      c.state.Input,
		Value:      &c.state.Settings.SoundLevel,
		Label:      "Effects volume",
		ValueNames: []string{"off", "1", "2", "3", "4", "5"},
		OnPressed: func() {
			if c.state.Settings.SoundLevel != 0 {
				scene.Audio().SetGroupVolume(assets.SoundGroupEffect, assets.VolumeMultiplier(c.state.Settings.SoundLevel))
				scene.Audio().PlaySound(assets.AudioPhotonCannon1)
			}
		},
	}))

	rowContainer.AddChild(eui.NewSelectButton(eui.SelectButtonConfig{
		Resources:  c.state.UIResources,
		Input:      c.state.Input,
		Value:      &c.state.Settings.MusicLevel,
		Label:      "Music volume",
		ValueNames: []string{"off", "1", "2", "3", "4", "5"},
		OnPressed: func() {
			c.scene.Audio().SetGroupVolume(assets.SoundGroupMusic, assets.VolumeMultiplier(c.state.Settings.MusicLevel))
		},
	}))

	rowContainer.AddChild(eui.NewSeparator(nil, styles.TransparentColor))
	rowContainer.AddChild(eui.NewButton(c.state.UIResources, "OK", func() {
		c.leave()
	}))

	initUI(scene, root)
}

func (c *settingsController) Update(delta float64) {
	if c.state.Input.ActionIsJustPressed(controls.ActionBack) {
		c.leave()
	}
}

func (c *settingsController) leave() {
	c.scene.Context().SaveGameData("save", c.state.Settings)
	c.scene.Context().ChangeScene(NewMainMenuController(c.state))
}

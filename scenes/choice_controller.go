package scenes

import (
	"fmt"
	"strings"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/vcgj7-game/assets"
	"github.com/quasilyte/vcgj7-game/controls"
	"github.com/quasilyte/vcgj7-game/eui"
	"github.com/quasilyte/vcgj7-game/gamedata"
	"github.com/quasilyte/vcgj7-game/session"
	"github.com/quasilyte/vcgj7-game/styles"
	"github.com/quasilyte/vcgj7-game/worldsim"
)

type ChoiceController struct {
	scene *ge.Scene
	state *session.State

	mapPosMarkerRotation gmath.Rad
	mapPosMarkerBase     gmath.Vec
	mapPosMarker         *ge.Sprite

	planetSectorLabels []*ge.Sprite
	planetSectorTitles []*ge.Label

	statusPanelText *widget.Text
	textPanelText   *widget.Text

	selectedChoice *worldsim.Choice
	runner         *worldsim.Runner

	choiceButtons []*choiceButton
}

type choiceButton struct {
	choice *worldsim.Choice
	widget *widget.Button
}

func NewChoiceController(state *session.State) *ChoiceController {
	return &ChoiceController{state: state}
}

func (c *ChoiceController) Init(scene *ge.Scene) {
	c.scene = scene
	c.initUI()

	scene.Audio().PauseCurrentMusic()
	scene.Audio().PlayMusic(assets.AudioMusicGlobal)

	{
		c.mapPosMarker = scene.NewSprite(assets.ImageMapLocation)
		c.mapPosMarker.Pos.Base = &c.mapPosMarkerBase
		c.mapPosMarker.Rotation = &c.mapPosMarkerRotation
		scene.AddGraphics(c.mapPosMarker)
	}

	{
		c.planetSectorTitles = make([]*ge.Label, len(c.state.World.Planets))
		c.planetSectorLabels = make([]*ge.Sprite, len(c.state.World.Planets))
		mapBase := &gmath.Vec{X: 752, Y: 76 - 19}
		for i, p := range c.state.World.Planets {
			s := ge.NewSprite(scene.Context())
			s.Pos.Base = mapBase
			s.Visible = false
			s.Pos.Offset = p.Info.MapOffset
			c.planetSectorLabels[i] = s
			scene.AddGraphics(s)

			l := ge.NewLabel(assets.BitmapFont1)
			l.AlignHorizontal = ge.AlignHorizontalCenter
			l.AlignVertical = ge.AlignVerticalCenter
			l.Text = strings.TrimPrefix(p.Info.Name, "Planet ")
			l.Width = 48
			l.Height = 20
			l.Pos.Base = mapBase
			l.Pos.Offset = p.Info.MapOffset.Add(gmath.Vec{X: -23, Y: 12})
			scene.AddGraphics(l)
		}
	}

	c.runner = worldsim.NewRunner(c.state.World)
	c.runner.Init(scene)
	c.runner.EventStartBattle.Connect(nil, c.onBattleStart)
	c.runner.EventGameOver.Connect(nil, c.onGameOver)

	c.replaceChoices()
	c.updateUI()
}

func (c *ChoiceController) onGameOver(victory bool) {
	if !victory {
		c.scene.Context().ChangeScene(NewMainMenuController(c.state))
	} else {
		c.scene.Context().ChangeScene(NewVictoryController(c.state))
	}
}

func (c *ChoiceController) replaceChoices() {
	result := c.runner.GenerateChoices()
	choices := result.Choices

	for i, c := range c.choiceButtons {
		if i > len(choices)-1 {
			c.widget.Text().Label = ""
			c.widget.GetWidget().Disabled = true
			continue
		}
		c.choice = &choices[i]
		c.widget.GetWidget().Disabled = false
		c.widget.Text().Label = fmt.Sprintf("%d. %s", i+1, c.choice.Text)
	}

	c.textPanelText.Label = result.Text
}

func (c *ChoiceController) selectChoice(i int) {
	if c.choiceButtons[i].choice == nil {
		return
	}
	c.selectedChoice = c.choiceButtons[i].choice

	if c.selectedChoice.Mode != gamedata.ModeUnknown {
		c.state.World.Player.Mode = c.selectedChoice.Mode
	}

	if c.selectedChoice.Time > 0 {
		if !c.runner.AdvanceTime(c.selectedChoice.Time) {
			c.replaceChoices()
			c.updateUI()
			return
		}
	}

	postMode := c.choiceButtons[i].choice.OnResolved()
	c.state.World.Player.Mode = postMode
	c.replaceChoices()
	c.updateUI()
}

func (c *ChoiceController) onBattleStart(info worldsim.BattleInfo) {
	c.scene.Context().ChangeScene(NewBattleController(c.state, info.Enemy))
}

func (c *ChoiceController) initUI() {
	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	rowContainer := eui.NewRowLayoutContainerWithMinWidth(900, 8, nil)
	root.AddChild(rowContainer)

	upperGrid := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(3),
			widget.GridLayoutOpts.Stretch([]bool{true, false, false}, nil),
			widget.GridLayoutOpts.Spacing(8, 8))))
	rowContainer.AddChild(upperGrid)

	textPanel := eui.NewPanelWithPadding(c.state.UIResources, 100, 100, widget.NewInsetsSimple(16))
	upperGrid.AddChild(textPanel)

	textPanelText := widget.NewText(
		widget.TextOpts.Text("", assets.BitmapFont1, styles.ButtonTextColor),
		widget.TextOpts.ProcessBBCode(true),
		widget.TextOpts.MaxWidth(450),
	)
	c.textPanelText = textPanelText
	textPanel.AddChild(textPanelText)

	picPanel := eui.NewPanelWithPadding(c.state.UIResources, 196, 196, widget.NewInsetsSimple(8))
	upperGrid.AddChild(picPanel)

	mapPanel := eui.NewPanelWithPadding(c.state.UIResources, 196, 196, widget.NewInsetsSimple(8))
	upperGrid.AddChild(mapPanel)

	mapBg := eui.NewGraphic(c.state.UIResources, assets.ImageSystemMap)
	mapPanel.AddChild(mapBg)

	lowerGrid := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{true, false}, nil),
			widget.GridLayoutOpts.Spacing(8, 8))))
	rowContainer.AddChild(lowerGrid)

	optionsList := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, nil),
			widget.GridLayoutOpts.Spacing(8, 8))))
	lowerGrid.AddChild(optionsList)

	c.choiceButtons = make([]*choiceButton, worldsim.MaxChoices)
	for i := range c.choiceButtons {
		id := i
		b := eui.NewButtonWithConfig(c.state.UIResources, eui.ButtonConfig{
			AlignLeft: true,
			Text:      fmt.Sprintf("[%d] button", i+1),
			OnClick: func() {
				c.selectChoice(id)
			},
			Font: assets.BitmapFont1,
		})
		b.GetWidget().Disabled = true
		optionsList.AddChild(b)
		c.choiceButtons[i] = &choiceButton{
			widget: b,
		}
	}

	statusPanel := eui.NewPanelWithPadding(c.state.UIResources, 196*2+8, 100, widget.NewInsetsSimple(16))
	lowerGrid.AddChild(statusPanel)

	c.statusPanelText = widget.NewText(
		widget.TextOpts.Text("", assets.BitmapFont1, styles.ButtonTextColor),
		widget.TextOpts.MaxWidth(360),
	)
	statusPanel.AddChild(c.statusPanelText)

	initUI(c.scene, root)
}

func (c *ChoiceController) Update(delta float64) {
	c.mapPosMarkerRotation += gmath.Rad(2 * delta)

	c.handleInput()
}

func (c *ChoiceController) handleInput() {
	for i := 0; i < worldsim.MaxChoices; i++ {
		if c.choiceButtons[i].choice == nil {
			continue
		}
		a := controls.ActionChoice1 + input.Action(i)
		if c.state.Input.ActionIsJustPressed(a) {
			c.selectChoice(i)
			break
		}
	}
}

func (c *ChoiceController) updateUI() {
	c.mapPosMarkerBase = (gmath.Vec{X: 752, Y: 75 - 19}).Add(c.state.World.Player.Planet.Info.MapOffset)

	p := c.state.World.Player
	{
		day := (c.state.World.GameTime / 24) + 1
		hours := c.state.World.GameTime % 24
		salary := gamedata.GetSalary(p.Experience) + p.ExtraSalary
		lines := []string{
			fmt.Sprintf("Day %d, %02d:00", day, hours),
			"",
			fmt.Sprintf("Combat experience: %d (salary is %d credits/day)", p.Experience, salary),
			fmt.Sprintf("Credits: %d", p.Credits),
			fmt.Sprintf("Vessel structure: %d%%", gmath.Clamp(int(100*p.VesselHP), 0, 100)),
			fmt.Sprintf("Fuel: %d/%d", p.Fuel, p.MaxFuel),
			fmt.Sprintf("Cargo: %d/%d", p.Cargo, p.MaxCargo),
		}
		lines = append(lines, "")
		q := c.state.World.CurrentQuest
		if q != nil && q.Active {
			lines = append(lines, fmt.Sprintf("Delivery quest destination: %s", q.Receiver.Info.Name))
		}
		if len(p.Artifacts) != 0 {
			lines = append(lines, fmt.Sprintf("Artifacts: %s", strings.Join(p.Artifacts, ", ")))
		} else {
			lines = append(lines, "Artifacts: <none>")
		}
		c.statusPanelText.Label = strings.Join(lines, "\n")
	}

	for i, s := range c.planetSectorLabels {
		p := c.state.World.Planets[i]
		switch p.Faction {
		case gamedata.FactionA:
			s.SetImage(c.scene.Context().Loader.LoadImage(assets.ImageAlliedPlanet))
			s.Visible = true
		case gamedata.FactionB, gamedata.FactionC:
			s.SetImage(c.scene.Context().Loader.LoadImage(assets.ImageHostilePlanet))
			s.Visible = true
		default:
			s.Visible = false
		}
	}
}

func (c *ChoiceController) formatChoiceTime(h int) string {
	days := (h / 24)
	hours := h % 24
	if days == 0 {
		return fmt.Sprintf("%d hours", hours)
	}
	if hours == 0 {
		return fmt.Sprintf("%d days", days)
	}
	return fmt.Sprintf("%d days and %d hours", days, hours)
}

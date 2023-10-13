package scenes

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/vcgj7-game/assets"
	"github.com/quasilyte/vcgj7-game/eui"
	"github.com/quasilyte/vcgj7-game/session"
)

type ChoiceController struct {
	scene *ge.Scene
	state *session.State

	choiceButtons []*widget.Button
}

func NewChoiceController(state *session.State) *ChoiceController {
	return &ChoiceController{state: state}
}

func (c *ChoiceController) Init(scene *ge.Scene) {
	c.scene = scene
	c.initUI()
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

	textPanel := eui.NewPanelWithPadding(c.state.UIResources, 100, 100, widget.NewInsetsSimple(8))
	upperGrid.AddChild(textPanel)

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

	c.choiceButtons = make([]*widget.Button, 6)
	for i := range c.choiceButtons {
		b := eui.NewButtonWithConfig(c.state.UIResources, eui.ButtonConfig{
			AlignLeft: true,
			Text:      fmt.Sprintf("[%d] button", i+1),
			OnClick:   func() {},
			Font:      assets.BitmapFont1,
		})
		b.GetWidget().Disabled = true
		optionsList.AddChild(b)
	}

	statusPanel := eui.NewPanelWithPadding(c.state.UIResources, 196*2+8, 100, widget.NewInsetsSimple(8))
	lowerGrid.AddChild(statusPanel)

	initUI(c.scene, root)
}

func (c *ChoiceController) Update(delta float64) {}

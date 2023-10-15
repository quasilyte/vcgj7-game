package worldsim

import (
	"fmt"
	"strings"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/vcgj7-game/gamedata"
)

var colorReplacer = strings.NewReplacer(
	"</>", "[/color]",
	"<g>", "[color=7AE168]",
	"<p>", "[color=B392FF]",
	"<r>", "[color=FF6363]",
	"<y>", "[color=FFF163]",
)

func cfmt(format string, args ...any) string {
	format = colorReplacer.Replace(format)
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}

type textColor int

const (
	colorDefault textColor = iota
	colorGreen
	colorPurple
	colorRed
	colorYellow
)

func colorizeText(s string, clr textColor) string {
	switch clr {
	case colorGreen:
		return "[color=7AE168]" + s + "[/color]"
	case colorPurple:
		return "[color=B392FF]" + s + "[/color]"
	case colorRed:
		return "[color=FF6363]" + s + "[/color]"
	case colorYellow:
		return "[color=FFF163]" + s + "[/color]"
	default:
		return s
	}
}

func genModeText(scene *ge.Scene, world *gamedata.World) string {
	switch world.Player.Mode {
	case gamedata.ModeOrbiting:
		return genOrbitingModeText(scene, world)
	case gamedata.ModeJustEntered:
		return genJustEnteredModeText(scene, world)
	case gamedata.ModeDocked:
		return genDockedModeText(scene, world)
	default:
		return "?"
	}
}

func genDockedModeText(scene *ge.Scene, world *gamedata.World) string {
	picker := gmath.NewRandPicker[string](scene.Rand())

	player := world.Player
	planet := player.Planet

	if planet.Info.GasGiant {
		picker.AddOption("Exploring the station hallways.", 1.2)
		picker.AddOption("Admiring the gas giant from the station windows.", 1.0)
		picker.AddOption("Spent some time doing nothing on this station.", 0.7)
	} else {
		picker.AddOption("Spending some time on the streets.", 1.1)
		picker.AddOption("Exploring the space decks.", 1.0)
		picker.AddOption("Walking through the local market.", 0.9)
		picker.AddOption(fmt.Sprintf("Spent some time doing nothing on %s.", planet.Info.Name), 0.7)
	}

	return picker.Pick()
}

func genJustEnteredModeText(scene *ge.Scene, world *gamedata.World) string {
	picker := gmath.NewRandPicker[string](scene.Rand())

	player := world.Player
	planet := player.Planet

	if planet.Faction == player.Faction {
		if planet.Info.GasGiant {
			picker.AddOption(fmt.Sprintf("Entering the allied %s station orbit.", planet.Info.Name), 1.5)
		} else {
			picker.AddOption(fmt.Sprintf("Entering the allied %s orbit.", planet.Info.Name), 1.5)
		}
	} else if planet.Faction == gamedata.FactionNone {
		if planet.Info.GasGiant {
			picker.AddOption(fmt.Sprintf("Entering the %s gas giant vicinity.", planet.Info.Name), 1.5)
		} else {
			picker.AddOption(fmt.Sprintf("Entering the %s orbit.", planet.Info.Name), 1.5)
		}
	} else {
		picker.AddOption(fmt.Sprintf("Approaching the %s. Danger: enemies detected.", planet.Info.Name), 1.5)
		picker.AddOption(fmt.Sprintf("Entering the hostile %s orbit.", planet.Info.Name), 1.4)
		picker.AddOption(fmt.Sprintf("Getting in range of a hostile %s.", planet.Info.Name), 1.2)
	}

	return picker.Pick()
}

func genOrbitingModeText(scene *ge.Scene, world *gamedata.World) string {
	picker := gmath.NewRandPicker[string](scene.Rand())

	player := world.Player
	planet := player.Planet

	if planet.Faction == player.Faction {
		if planet.Info.GasGiant {
			picker.AddOption("Flying around the allied station near "+planet.Info.Name+".", 1.5)
		} else {
			picker.AddOption(fmt.Sprintf("Orbiting around allied %s.", planet.Info.Name), 1.5)
		}
	} else if planet.Faction == gamedata.FactionNone {
		if planet.Info.GasGiant {
			picker.AddOption("Navigating around the "+planet.Info.Name+" gas giant.", 1.5)
		} else {
			picker.AddOption(fmt.Sprintf("Orbiting around neutral %s.", planet.Info.Name), 1.5)
		}
	} else {
		picker.AddOption("Hiding from the hostile fleet of "+planet.Info.Name+".", 1.5)
		picker.AddOption("Observing the hostile "+planet.Info.Name+".", 1.2)
		picker.AddOption("Spying on "+planet.Info.Name+".", 1.1)
	}

	return picker.Pick()
}

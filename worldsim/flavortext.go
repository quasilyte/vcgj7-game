package worldsim

import (
	"fmt"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/vcgj7-game/gamedata"
)

func genModeText(scene *ge.Scene, world *gamedata.World) string {
	switch world.Player.Mode {
	case gamedata.ModeOrbiting:
		return genOrbitingModeText(scene, world)
	case gamedata.ModeJustEntered:
		return genJustEnteredModeText(scene, world)
	default:
		return "?"
	}
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

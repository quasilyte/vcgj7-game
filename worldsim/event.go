package worldsim

import (
	"fmt"
	"math"

	"github.com/quasilyte/gmath"
	"github.com/quasilyte/vcgj7-game/gamedata"
)

type eventKind int

const (
	eventUnknown eventKind = iota

	eventFuelScavenge
	eventMineralsHunt

	eventBuyFuel
	eventSellMinerals
)

func (r *Runner) generateEventChoices(event eventInfo) string {
	player := r.world.Player
	planet := player.Planet

	switch event.kind {
	case eventFuelScavenge:
		fuelScavenged := r.scene.Rand().IntRange(3, 12)
		r.choices = append(r.choices, Choice{
			Text: "Done",
			OnSelected: func() {
				player.Fuel = gmath.ClampMax(player.Fuel+fuelScavenged, player.MaxFuel)
				r.commitChoice(gamedata.ModeOrbiting)
			},
		})
		if r.scene.Rand().Bool() {
			return fmt.Sprintf("%d fuel units acquired.", fuelScavenged)
		}
		return fmt.Sprintf("Scavenged %d fuel units.", fuelScavenged)

	case eventMineralsHunt:
		mineralsFound := r.scene.Rand().IntRange(20, 40)
		if r.scene.Rand().Chance(0.3) {
			mineralsFound *= 2
		}
		loaded := mineralsFound
		freeCargo := player.FreeCargoSpace()
		if freeCargo < loaded {
			loaded = freeCargo
		}
		r.choices = append(r.choices, Choice{
			Text: "Done",
			OnSelected: func() {
				if r.scene.Rand().Chance(0.7) {
					planet.MineralsDelay = r.scene.Rand().FloatRange(3, 14)
				}
				player.LoadCargo(mineralsFound)
				r.commitChoice(gamedata.ModeOrbiting)
			},
		})
		if loaded < mineralsFound {
			return fmt.Sprintf("Found %d minerals, but could only collect %d.", mineralsFound, loaded)
		}
		return fmt.Sprintf("Collected %d minerals.", loaded)

	case eventSellMinerals:
		mineralsDemand := 1.0
		s := "The minerals are in demand here."
		switch {
		case planet.MineralDeposit < 50:
			mineralsDemand = 1.3
			s = fmt.Sprintf("The minerals are in %s demand here.", colorizeText("high", colorYellow))
		case planet.MineralDeposit > 100:
			mineralsDemand = 0.8
			s = "The minerals have normal price here."
		case planet.MineralDeposit > 250:
			mineralsDemand = 0.5
			s = "The minerals have are not in demand here."
		case planet.MineralDeposit > 500:
			mineralsDemand = 0.3
			s = "The minerals have low price here."
		case planet.MineralDeposit > 900:
			mineralsDemand = 0.1
			s = "The minerals have very low price here."
		}
		price := r.scene.Rand().FloatRange(0.8, 1.6)
		totalCost := int(math.Ceil(float64(player.Cargo) * price * mineralsDemand))
		r.choices = append(r.choices, Choice{
			Text: "Accept deal",
			OnSelected: func() {
				planet.MineralDeposit += player.Cargo
				player.Cargo = 0
				player.Credits += totalCost
				r.commitChoice(gamedata.ModeDocked)
			},
		})
		r.choices = append(r.choices, Choice{
			Text: "Decline deal",
			OnSelected: func() {
				r.commitChoice(gamedata.ModeDocked)
			},
		})
		return fmt.Sprintf("%s\n\nSell %d minerals for %d credits?", s, player.Cargo, totalCost)

	case eventBuyFuel:
		fuelPrice := 3
		if r.scene.Rand().Chance(0.3) {
			fuelPrice = 2
		}
		maxSpent := 90
		if maxSpent > player.Credits {
			maxSpent = player.Credits
		}
		bought := maxSpent / fuelPrice
		if player.Fuel+bought > player.MaxFuel {
			bought = player.MaxFuel - player.Fuel
		}
		spent := bought * fuelPrice
		r.choices = append(r.choices, Choice{
			Text: "Done",
			OnSelected: func() {
				player.Credits -= spent
				player.Fuel += bought
				r.commitChoice(gamedata.ModeDocked)
			},
		})
		return fmt.Sprintf("Bought %d fuel units for %d credits.", bought, spent)

	default:
		panic(fmt.Sprintf("unexpected event kind: %d", event.kind))
	}
}

package worldsim

import (
	"fmt"
	"math"
	"strings"

	"github.com/quasilyte/gmath"
	"github.com/quasilyte/vcgj7-game/gamedata"
)

type eventKind int

const (
	eventUnknown eventKind = iota

	eventBattle
	eventBattleInterrupt

	eventFuelScavenge
	eventMineralsHunt
	eventScanArea

	eventBuyFuel
	eventUpgradeLab
	eventSellMinerals
)

func (r *Runner) afterBattleChoices() string {
	player := r.world.Player

	reward := player.BattleRewards
	player.BattleRewards = gamedata.BattleRewards{}

	if !reward.Victory {
		r.choices = append(r.choices, Choice{
			Text: "The great ranger's life has come to an end",
			OnResolved: func() gamedata.Mode {
				r.EventGameOver.Emit(false)
				return 0
			},
		})
		return "Your vessel was destroyed in battle."
	}

	r.choices = append(r.choices, Choice{
		Text: "Done",
		OnResolved: func() gamedata.Mode {
			player.Experience += reward.Experience
			player.Credits += reward.Credits
			player.LoadCargo(reward.Cargo)
			player.Fuel = gmath.ClampMax(player.Fuel+reward.Fuel, player.MaxFuel)
			return gamedata.ModeOrbiting
		},
	})

	lines := make([]string, 0, 5)

	lines = append(lines, "You are victorious!")
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Earned %d combat experience.", reward.Experience))
	if reward.Credits != 0 {
		lines = append(lines, fmt.Sprintf("Found %d credits equivalent.", reward.Credits))
	}
	if reward.Cargo != 0 {
		lines = append(lines, fmt.Sprintf("Scavenged %d resource units.", reward.Cargo))
	}
	if reward.Fuel != 0 {
		lines = append(lines, fmt.Sprintf("Recovered %d fuel units.", reward.Fuel))
	}

	return strings.Join(lines, "\n")
}

func (r *Runner) generateEventChoices(event eventInfo) string {
	player := r.world.Player
	planet := player.Planet

	switch event.kind {
	case eventScanArea:
		lines := make([]string, 0, 6)
		lines = append(lines, "Scanning area...")
		foundAnyone := false
		for i, num := range planet.VesselsByFaction {
			if num == 0 {
				continue
			}
			if !foundAnyone {
				lines = append(lines, "")
			}
			foundAnyone = true
			f := gamedata.Faction(i)
			lines = append(lines, fmt.Sprintf("%s vessels: %d", f.Name(), num))
		}
		if !foundAnyone {
			lines = append(lines, "")
			lines = append(lines, "No vessels detected.")
		}
		r.choices = append(r.choices, Choice{
			Text: "Done",
			OnResolved: func() gamedata.Mode {
				return gamedata.ModeOrbiting
			},
		})
		return strings.Join(lines, "\n")

	case eventUpgradeLab:
		var s string
		price := 0
		switch r.world.UpgradeAvailable {
		case gamedata.UpgradeJumpMaxDistance:
			s = fmt.Sprintf("A jump engine booster that increases its %s.", colorizeText("max jump distance", colorGreen))
			price = 20
		case gamedata.UpgradeMaxFuel:
			s = fmt.Sprintf("A special fuel tank extender to increase its %s.", colorizeText("max capacity", colorGreen))
			price = 45
		case gamedata.UpgradeMaxCargo:
			s = fmt.Sprintf("A better storage compactor, it will %s of your vessel.", colorizeText("increase max cargo", colorGreen))
			price = 70
		case gamedata.UpgradeJumpSpeed:
			s = fmt.Sprintf("A jump engine cooling system that allows you to %s.", colorizeText("travel between the planets faster", colorGreen))
			price = 15
		}

		if player.Credits >= price {
			r.choices = append(r.choices, Choice{
				Text: "Buy this upgrade",
				OnResolved: func() gamedata.Mode {
					r.world.UpgradeRerollDelay = 0
					r.world.NextUpgradeDelay = r.scene.Rand().FloatRange(30, 45)
					player.Credits -= price
					switch r.world.UpgradeAvailable {
					case gamedata.UpgradeJumpMaxDistance:
						player.MaxJumpDist += float64(r.scene.Rand().IntRange(3, 6))
					case gamedata.UpgradeMaxFuel:
						player.MaxFuel += r.scene.Rand().IntRange(5, 15)
					case gamedata.UpgradeMaxCargo:
						player.MaxCargo += r.scene.Rand().IntRange(5, 20)
					case gamedata.UpgradeJumpSpeed:
						player.JumpSpeed += float64(r.scene.Rand().IntRange(15, 30))
					}
					return gamedata.ModeDocked
				},
			})
		}

		r.choices = append(r.choices, Choice{
			Text: "Leave lab",
			OnResolved: func() gamedata.Mode {
				r.world.NextUpgradeDelay = r.scene.Rand().FloatRange(2, 5)
				return gamedata.ModeDocked
			},
		})
		lines := []string{
			"You visited an experimental research lab. A person in white coat approaches you.",
			"",
			"After a quick discission, one particular upgrade catched your attention... " + s,
			"",
			fmt.Sprintf("It will cost you %d credits.", price),
		}
		return strings.Join(lines, "\n")

	case eventBattle, eventBattleInterrupt:
		r.choices = append(r.choices, Choice{
			Text: "Fight!",
			Mode: gamedata.ModeCombat,
			OnResolved: func() gamedata.Mode {
				planet.VesselsByFaction[event.enemy.Faction]--
				r.EventStartBattle.Emit(BattleInfo{
					Enemy:          event.enemy,
					ChallengeLevel: 1,
				})
				return gamedata.ModeAfterCombat
			},
		})
		if event.kind == eventBattleInterrupt {
			return fmt.Sprintf("Your actions were interrupted by a %s. Prepare for battle.", colorizeText("hostile vessel", colorRed))
		}
		return "This is a battle test."

	case eventFuelScavenge:
		fuelScavenged := r.scene.Rand().IntRange(3, 12)
		r.choices = append(r.choices, Choice{
			Text: "Done",
			OnResolved: func() gamedata.Mode {
				player.Fuel = gmath.ClampMax(player.Fuel+fuelScavenged, player.MaxFuel)
				return gamedata.ModeOrbiting
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
		foundShipwreck := r.scene.Rand().Chance(0.2)
		fuelGained := 0
		if foundShipwreck {
			fuelGained = r.scene.Rand().IntRange(4, 8)
		}
		damaged := r.scene.Rand().Chance(0.4)
		r.choices = append(r.choices, Choice{
			Text: "Done",
			OnResolved: func() gamedata.Mode {
				if r.scene.Rand().Chance(0.85) {
					planet.MineralsDelay = r.scene.Rand().FloatRange(5, 50)
					if r.scene.Rand().Chance(0.3) {
						planet.MineralsDelay *= 2
					}
				}
				if damaged {
					player.VesselHP -= r.scene.Rand().FloatRange(0.1, 0.2)
				}
				player.Fuel = gmath.ClampMax(player.Fuel+fuelGained, player.MaxFuel)
				player.LoadCargo(mineralsFound)
				return gamedata.ModeOrbiting
			},
		})
		lines := make([]string, 0, 3)
		if loaded < mineralsFound {
			lines = append(lines, fmt.Sprintf("Found %d minerals, but could only collect %d.", mineralsFound, loaded))
		} else {
			lines = append(lines, fmt.Sprintf("Collected %d minerals.", loaded))
		}
		if foundShipwreck {
			lines = append(lines, "")
			lines = append(lines, fmt.Sprintf("While flying near asteroids, you discovered a shipwreck site. You found recyclable objects worth %d fuel units.", fuelGained))
		}
		if damaged {
			lines = append(lines, "")
			lines = append(lines, "Your vessel hull was damaged during the act.")
		}
		return strings.Join(lines, "\n")

	case eventSellMinerals:
		mineralsDemand := 1.0
		s := "The minerals are in demand here."
		switch {
		case planet.MineralDeposit < 50:
			mineralsDemand = 1.4
			s = fmt.Sprintf("The minerals are in %s demand here.", colorizeText("high", colorYellow))
		case planet.MineralDeposit > 100:
			mineralsDemand = 1.1
			s = "The minerals have normal price here."
		case planet.MineralDeposit > 200:
			mineralsDemand = 0.9
			s = "The minerals have are not in demand here."
		case planet.MineralDeposit > 300:
			mineralsDemand = 0.4
			s = "The minerals have low price here."
		case planet.MineralDeposit > 500:
			mineralsDemand = 0.1
			s = "The minerals have very low price here."
		}
		price := r.scene.Rand().FloatRange(0.8, 1.6)
		totalCost := int(math.Ceil(float64(player.Cargo) * price * mineralsDemand))
		r.choices = append(r.choices, Choice{
			Text: "Accept deal",
			OnResolved: func() gamedata.Mode {
				planet.MineralDeposit += player.Cargo
				player.Cargo = 0
				player.Credits += totalCost
				return gamedata.ModeDocked
			},
		})
		r.choices = append(r.choices, Choice{
			Text: "Decline deal",
			OnResolved: func() gamedata.Mode {
				return gamedata.ModeDocked
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
			OnResolved: func() gamedata.Mode {
				player.Credits -= spent
				player.Fuel += bought
				return gamedata.ModeDocked
			},
		})
		return fmt.Sprintf("Bought %d fuel units for %d credits.", bought, spent)

	default:
		panic(fmt.Sprintf("unexpected event kind: %d", event.kind))
	}
}

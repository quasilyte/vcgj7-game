package worldsim

import (
	"fmt"
	"math"
	"strings"

	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/vcgj7-game/assets"
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

	eventTakeQuest
	eventCompleteQuest
	eventNews
	eventBuyFuel
	eventUpgradeLab
	eventWeaponShop
	eventShipyard
	eventWorkshop
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
			if reward.Artifact != "" {
				player.Artifacts = append(player.Artifacts, reward.Artifact)
			}
			if reward.SystemLiberated {
				if player.ExtraSalary < 20 {
					player.ExtraSalary += 3
				} else {
					player.Credits += 30
				}
			}
			return gamedata.ModeOrbiting
		},
	})

	lines := make([]string, 0, 5)

	lines = append(lines, "You are victorious!")
	lines = append(lines, "")
	lines = append(lines, cfmt("Earned <y>%d</> combat experience.", reward.Experience))
	if reward.Credits != 0 {
		lines = append(lines, cfmt("Found <y>%d</> credits equivalent.", reward.Credits))
	}
	if reward.Cargo != 0 {
		lines = append(lines, cfmt("Scavenged <y>%d</> resource units.", reward.Cargo))
	}
	if reward.Fuel != 0 {
		lines = append(lines, cfmt("Recovered <y>%d</> fuel units.", reward.Fuel))
	}
	if reward.Artifact != "" {
		desc := ""
		switch reward.Artifact {
		case "Fuel Generator":
			desc = "Acquired <g>Fuel Generator</> artifact."
		case "Repair Bots":
			desc = "Acquired <g>Repair Bots</> artifact."
		case "Scantide":
			desc = "Acquired <g>Scantide artifact</> (makes scanning faster)."
		case "Lucky Charm":
			desc = "Acquired <g>Lucky Charm</> artifact."
		case "Jumper":
			desc = "Acquired <g>Jumper</> artifact (makes jumps cost less fuel)."
		}
		lines = append(lines, cfmt(desc))
	}

	if reward.SystemLiberated {
		lines = append(lines, "")
		if player.ExtraSalary < 20 {
			lines = append(lines, cfmt("System liberation bonus: <y>+3</> salary."))
		} else {
			lines = append(lines, cfmt("System liberation bonus: <y>30</> credits."))
		}
	}

	return strings.Join(lines, "\n")
}

func (r *Runner) generateEventChoices(event eventInfo) string {
	player := r.world.Player
	planet := player.Planet

	switch event.kind {
	case eventCompleteQuest:
		q := r.world.CurrentQuest
		r.world.CurrentQuest = nil
		r.choices = append(r.choices, Choice{
			Text: "Done",
			OnResolved: func() gamedata.Mode {
				r.world.QuestRerollDelay = float64(r.scene.Rand().IntRange(60, 90))
				player.Credits += q.CreditsReward
				player.Experience += q.ExpReward
				return gamedata.ModeDocked
			},
		})
		return cfmt("Quest completed!\n\nReceived <y>%d</> credits and <y>%d</> experience points.", q.CreditsReward, q.ExpReward)

	case eventTakeQuest:
		q := r.world.CurrentQuest
		lines := make([]string, 0, 8)
		lines = append(lines, cfmt("This quest requires you to deliver this very important object to <p>%s</>.", q.Receiver.Info.Name))
		lines = append(lines, "")
		lines = append(lines, cfmt("Reward: <y>%d</> credits and <y>%d</> experience points.", q.CreditsReward, q.ExpReward))
		r.choices = append(r.choices, Choice{
			Text: "Accept quest",
			OnResolved: func() gamedata.Mode {
				q.Active = true
				return gamedata.ModeDocked
			},
		})
		r.choices = append(r.choices, Choice{
			Text: "Decline quest",
			OnResolved: func() gamedata.Mode {
				return gamedata.ModeDocked
			},
		})
		return strings.Join(lines, "\n")

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
			if f == player.Faction {
				lines = append(lines, cfmt("<g>%s</> vessels: <y>%d</>", f.Name(), num))
			} else {
				lines = append(lines, cfmt("<r>%s</> vessels: <y>%d</>", f.Name(), num))
			}
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

	case eventNews:
		lines := make([]string, 0, 8)
		lines = append(lines, "The latest system-wide news:")
		lines = append(lines, "")
		for _, e := range r.world.RecentEvents {
			day := (e.Time / 24) + 1
			hours := e.Time % 24
			dateString := cfmt("Day %d, %02d:00", day, hours)
			s := fmt.Sprintf("* [%s] %s", dateString, e.Text)
			lines = append(lines, s)
		}
		r.choices = append(r.choices, Choice{
			Text: "Done",
			OnResolved: func() gamedata.Mode {
				return gamedata.ModeDocked
			},
		})
		return strings.Join(lines, "\n")

	case eventShipyard:
		lines := make([]string, 0, 8)
		lines = append(lines, "A new, improved vessel is available for the veterans.")
		lines = append(lines, "")
		lines = append(lines, "In comparison with your current vessel:")
		lines = append(lines, cfmt("<g>+50</> health"))
		lines = append(lines, cfmt("<g>+20</> max energy"))
		lines = append(lines, cfmt("<g>+20</> cargo space"))
		lines = append(lines, cfmt("<r>-15%</> rotation speed"))
		lines = append(lines, "")
		lines = append(lines, cfmt("It costs <y>350</> credits"))
		if player.Credits >= 350 {
			r.choices = append(r.choices, Choice{
				Text: "Buy new vessel",
				OnResolved: func() gamedata.Mode {
					player.ImprovedHull = true
					player.MaxCargo += 20
					player.VesselDesign.Image = assets.ImageVesselPlayerElite
					player.VesselDesign.MaxEnergy += 20
					player.VesselDesign.MaxHP += 50
					player.VesselDesign.RotationSpeed -= 0.8
					player.Credits -= 350
					return gamedata.ModeDocked
				},
			})
		}
		r.choices = append(r.choices, Choice{
			Text: "Leave shipyard",
			OnResolved: func() gamedata.Mode {
				return gamedata.ModeDocked
			},
		})
		return strings.Join(lines, "\n")

	case eventWorkshop:
		lines := make([]string, 0, 8)
		lines = append(lines, "You can improve your vessel combat stats here.")
		lines = append(lines, "")
		lines = append(lines, "Your vessel stats:")

		armorUpgradeCost := 30 + (10 * (player.ArmorLevel - 1))
		lines = append(lines, cfmt("* Armor (level <g>%d</>) - <y>%d</> credits to increase", player.ArmorLevel, armorUpgradeCost))
		if player.Credits >= armorUpgradeCost {
			r.choices = append(r.choices, Choice{
				Text: "Increase armor level",
				Time: 20,
				OnResolved: func() gamedata.Mode {
					player.ArmorLevel++
					player.VesselDesign.MaxHP += float64(r.scene.Rand().IntRange(10, 20) + (3 * player.ArmorLevel))
					player.Credits -= armorUpgradeCost
					return gamedata.ModeDocked
				},
			})
		}

		energyUpgradeCost := 25 + (10 * (player.EnergyLevel - 1))
		lines = append(lines, cfmt("* Energy capacity (level <g>%d</>) - <y>%d</> credits to increase", player.EnergyLevel, energyUpgradeCost))
		if player.Credits >= energyUpgradeCost {
			r.choices = append(r.choices, Choice{
				Text: "Increase energy level",
				Time: 15,
				OnResolved: func() gamedata.Mode {
					player.EnergyLevel++
					player.VesselDesign.MaxEnergy += float64(r.scene.Rand().IntRange(10, 20) + (2 * player.EnergyLevel))
					player.VesselDesign.EnergyRegen += r.scene.Rand().FloatRange(0.1, 0.2)
					player.Credits -= energyUpgradeCost
					return gamedata.ModeDocked
				},
			})
		}

		speedUpgradeCost := 15 + (5 * (player.SpeedLevel - 1))
		lines = append(lines, cfmt("* Max speed (level <g>%d</>) - <y>%d</> credits to increase", player.SpeedLevel, speedUpgradeCost))
		if player.Credits >= speedUpgradeCost {
			r.choices = append(r.choices, Choice{
				Text: "Increase max speed level",
				Time: 10,
				OnResolved: func() gamedata.Mode {
					player.SpeedLevel++
					player.VesselDesign.MaxSpeed += float64(r.scene.Rand().IntRange(20, 35))
					player.Credits -= speedUpgradeCost
					return gamedata.ModeDocked
				},
			})
		}

		accelerationUpgradeCost := 15 + (5 * (player.AccelerationLevel - 1))
		lines = append(lines, cfmt("* Acceleration (level <g>%d</>) - <y>%d</> credits to increase", player.AccelerationLevel, accelerationUpgradeCost))
		if player.Credits >= accelerationUpgradeCost {
			r.choices = append(r.choices, Choice{
				Text: "Increase acceleration level",
				Time: 5,
				OnResolved: func() gamedata.Mode {
					player.AccelerationLevel++
					player.VesselDesign.Acceleration += float64(r.scene.Rand().IntRange(30, 40))
					player.Credits -= accelerationUpgradeCost
					return gamedata.ModeDocked
				},
			})
		}

		rotationUpgradeCost := 25 + (10 * (player.RotationLevel - 1))
		lines = append(lines, cfmt("* Rotation speed (level <g>%d</>) - <y>%d</> credits to increase", player.RotationLevel, rotationUpgradeCost))
		if player.Credits >= rotationUpgradeCost {
			r.choices = append(r.choices, Choice{
				Text: "Increase rotation speed level",
				Time: 20,
				OnResolved: func() gamedata.Mode {
					player.RotationLevel++
					player.VesselDesign.RotationSpeed += gmath.Rad(r.scene.Rand().FloatRange(0.25, 0.4))
					player.Credits -= rotationUpgradeCost
					return gamedata.ModeDocked
				},
			})
		}

		r.choices = append(r.choices, Choice{
			Text: "Leave workshop",
			OnResolved: func() gamedata.Mode {
				return gamedata.ModeDocked
			},
		})

		return strings.Join(lines, "\n")

	case eventWeaponShop:
		formatWeapon := func(w *gamedata.WeaponDesign) string {
			if w.Primary {
				return cfmt("%s (<g>primary</>)", w.Name)
			}
			return cfmt("%s (<p>secondary</>)", w.Name)
		}
		lines := make([]string, 0, 6)
		if len(planet.WeaponsAvailable) > 0 {
			lines = append(lines, "The weapon selection include:")
			for _, weaponName := range planet.WeaponsAvailable {
				w := gamedata.FindWeaponDesign(weaponName)
				cost := cfmt(" - <y>%d</> credits", w.Cost)
				lines = append(lines, "* "+formatWeapon(w)+cost)
				if player.Credits >= w.Cost && !player.HasWeapon(w) {
					r.choices = append(r.choices, Choice{
						Text: "Buy " + w.Name,
						OnResolved: func() gamedata.Mode {
							planet.WeaponsAvailable = xslices.Remove(planet.WeaponsAvailable, weaponName)
							player.Credits -= w.Cost
							if w.Primary {
								player.VesselDesign.MainWeapon = w
							} else {
								player.VesselDesign.SecondaryWeapon = w
							}
							return gamedata.ModeDocked
						},
					})
				}
			}
		} else {
			lines = append(lines, "This weapon shop is empty at the moment. Come again later.")
		}
		lines = append(lines, "")
		lines = append(lines, "Your current weapons:")
		if player.VesselDesign.MainWeapon != nil {
			lines = append(lines, "* "+formatWeapon(player.VesselDesign.MainWeapon))
		}
		if player.VesselDesign.SecondaryWeapon != nil {
			lines = append(lines, "* "+formatWeapon(player.VesselDesign.SecondaryWeapon))
		}
		r.choices = append(r.choices, Choice{
			Text: "Leave weapon shop",
			OnResolved: func() gamedata.Mode {
				return gamedata.ModeDocked
			},
		})
		return strings.Join(lines, "\n")

	case eventUpgradeLab:
		var s string
		price := 0
		switch r.world.UpgradeAvailable {
		case gamedata.UpgradeJumpMaxDistance:
			s = fmt.Sprintf("A jump engine booster that increases its %s.", colorizeText("max jump distance", colorGreen))
			price = 30
		case gamedata.UpgradeMaxFuel:
			s = fmt.Sprintf("A special fuel tank extender to increase its %s.", colorizeText("max capacity", colorGreen))
			price = 60
		case gamedata.UpgradeMaxCargo:
			s = fmt.Sprintf("A better storage compactor, it will %s of your vessel.", colorizeText("increase max cargo", colorGreen))
			price = 70
		case gamedata.UpgradeJumpSpeed:
			s = fmt.Sprintf("A jump engine cooling system that allows you to %s.", colorizeText("travel between the planets faster", colorGreen))
			price = 25
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
			"After a quick discussion, one particular upgrade catched your attention... " + s,
			"",
			cfmt("It will cost you <y>%d</> credits.", price),
		}
		return strings.Join(lines, "\n")

	case eventBattle, eventBattleInterrupt:
		lastDefender := planet.Faction == event.enemy.Faction && planet.VesselsByFaction[event.enemy.Faction] == 1
		event.enemy.LastDefender = lastDefender
		pirateAttack := event.enemy.Image == assets.ImageVesselPirate
		r.choices = append(r.choices, Choice{
			Text: "Fight!",
			Mode: gamedata.ModeCombat,
			OnResolved: func() gamedata.Mode {
				if pirateAttack {
					r.world.PirateSeq++
				}
				if event.enemy.Faction != gamedata.FactionNone {
					planet.VesselsByFaction[event.enemy.Faction]--
				}
				if planet.Faction == event.enemy.Faction && planet.VesselsByFaction[event.enemy.Faction] == 0 {
					planet.Faction = gamedata.FactionNone
					r.world.PushEvent(fmt.Sprintf("%s lost control over %s", event.enemy.Faction.Name(), planet.Info.Name))
				}
				r.EventStartBattle.Emit(BattleInfo{
					Enemy: event.enemy,
				})
				return gamedata.ModeAfterCombat
			},
		})
		lines := make([]string, 0, 4)
		if player.Mode == gamedata.ModeAttack {
			lines = append(lines, "Enemy spotted!")
		} else if event.kind == eventBattleInterrupt {
			if pirateAttack {
				lines = append(lines, cfmt("An <r>unidentified vessel</> opens fire at you."))
				if player.Fuel >= 5 {
					r.choices = append(r.choices, Choice{
						Text: "Retreat [5 fuel]",
						OnResolved: func() gamedata.Mode {
							player.Fuel -= 5
							return gamedata.ModeOrbiting
						},
					})
				}
			} else {
				lines = append(lines, fmt.Sprintf("Your actions were interrupted by a %s. Prepare for battle.", colorizeText("hostile vessel", colorRed)))
			}
		}
		if r.world.Player.Battles < 5 {
			lines = append(lines, "")
			lines = append(lines, cfmt("--- <p>Tutorial</> ---"))
			lines = append(lines, cfmt("Your shield blocks <y>75%</> primary weapon damage and <p>consumes energy</>."))
			lines = append(lines, "Blocking is the most efficient way to recover energy mid-battle.")
			lines = append(lines, "")
			lines = append(lines, "Controls:")
			lines = append(lines, "* Style 1: WASD for movement, mouse buttons to fire.")
			lines = append(lines, "* Style 2: WASD for movement, [O] and [P] to fire.")
			lines = append(lines, "* Style 3: arrows for movement, [Z] and [X] to fire.")

		}
		return strings.Join(lines, "\n")

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
			return cfmt("<y>%d</> fuel units acquired.", fuelScavenged)
		}
		return cfmt("Scavenged <y>%d</> fuel units.", fuelScavenged)

	case eventMineralsHunt:
		mineralsFound := r.scene.Rand().IntRange(20, 40)
		if r.scene.Rand().Chance(0.3) {
			mineralsFound *= 2
		}
		if !player.HasArtifact("Lucky Charm") {
			if r.scene.Rand().Chance(0.06) {
				mineralsFound = 0
			}
		} else {
			mineralsFound += r.scene.Rand().IntRange(4, 14)
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
				if r.scene.Rand().Chance(0.9) {
					planet.MineralsDelay = r.scene.Rand().FloatRange(15, 55)
					if r.scene.Rand().Chance(0.35) {
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
		if mineralsFound == 0 {
			lines = append(lines, "No valuable minerals found.")
		} else {
			if loaded < mineralsFound {
				lines = append(lines, cfmt("Found <y>%d</> minerals, but could only collect <y>%d</>.", mineralsFound, loaded))
			} else {
				lines = append(lines, cfmt("Collected <y>%d</> minerals.", loaded))
			}
		}

		if foundShipwreck {
			lines = append(lines, "")
			lines = append(lines, cfmt("While flying near asteroids, you discovered a shipwreck site. You found recyclable objects worth <y>%d</> fuel units.", fuelGained))
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
			mineralsDemand = 1.8
			s = "The minerals are in high demand here."
		case planet.MineralDeposit > 100:
			mineralsDemand = 1.5
			s = "The minerals have normal price here."
		case planet.MineralDeposit > 200:
			mineralsDemand = 1.0
			s = "The minerals have are not in demand here."
		case planet.MineralDeposit > 300:
			mineralsDemand = 0.5
			s = "The minerals have low price here."
		case planet.MineralDeposit > 500:
			mineralsDemand = 0.2
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
		return cfmt("%s\n\nSell <y>%d</> minerals for <y>%d</> credits?", s, player.Cargo, totalCost)

	case eventBuyFuel:
		fuelPrice := 0.5
		maxSpent := 90.0
		if maxSpent > float64(player.Credits) {
			maxSpent = float64(player.Credits)
		}
		bought := maxSpent / fuelPrice
		if float64(player.Fuel)+bought > float64(player.MaxFuel) {
			bought = float64(player.MaxFuel - player.Fuel)
		}
		spent := bought * fuelPrice
		r.choices = append(r.choices, Choice{
			Text: "Done",
			OnResolved: func() gamedata.Mode {
				player.Credits -= int(math.Ceil(spent))
				player.Fuel += int(math.Ceil(bought))
				return gamedata.ModeDocked
			},
		})
		return cfmt("Bought <y>%d</> fuel units for <y>%d</> credits.", int(math.Ceil(bought)), int(math.Ceil(spent)))

	default:
		panic(fmt.Sprintf("unexpected event kind: %d", event.kind))
	}
}

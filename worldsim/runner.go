package worldsim

import (
	"fmt"
	"math"
	"strings"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/vcgj7-game/assets"
	"github.com/quasilyte/vcgj7-game/gamedata"
)

type Runner struct {
	world       *gamedata.World
	scene       *ge.Scene
	choices     []Choice
	jumpOptions []jumpOption
	textLines   []string

	eventInfo eventInfo

	EventChoiceSelected gsignal.Event[gamedata.Mode]
	EventStartBattle    gsignal.Event[BattleInfo]
	EventGameOver       gsignal.Event[bool]
}

type BattleInfo struct {
	ChallengeLevel int
	Enemy          *gamedata.VesselDesign
}

type eventInfo struct {
	kind eventKind

	enemy *gamedata.VesselDesign
}

type jumpOption struct {
	planet   *gamedata.Planet
	fuelCost int
	time     int
}

type GeneratedChoices struct {
	Text    string
	Choices []Choice
}

func NewRunner(w *gamedata.World) *Runner {
	return &Runner{
		world:       w,
		choices:     make([]Choice, 0, 8),
		jumpOptions: make([]jumpOption, 0, 8),
		textLines:   make([]string, 0, 20),
	}
}

func (r *Runner) Init(scene *ge.Scene) {
	r.scene = scene
}

func (r *Runner) GenerateChoices() GeneratedChoices {
	r.textLines = r.textLines[:0]
	r.choices = r.choices[:0]
	r.jumpOptions = r.jumpOptions[:0]

	player := r.world.Player
	planet := r.world.Player.Planet

	if player.Mode == gamedata.ModeAfterCombat {
		player.Mode = gamedata.ModeOrbiting
		s := r.afterBattleChoices()
		return GeneratedChoices{
			Choices: r.choices,
			Text:    s,
		}
	}

	event := r.eventInfo
	r.eventInfo = eventInfo{}
	if event.kind != eventUnknown {
		s := r.generateEventChoices(event)
		return GeneratedChoices{
			Choices: r.choices,
			Text:    s,
		}
	}

	r.textLines = append(r.textLines, genModeText(r.scene, r.world))

	canJump := true

	if len(r.choices) < MaxChoices && planet.Faction == player.Faction {
		switch player.Mode {
		case gamedata.ModeJustEntered, gamedata.ModeOrbiting:
			s := "Enter the planetary docks"
			h := 3
			if planet.Info.GasGiant {
				s = "Dock the station"
				h = 1
			}
			r.choices = append(r.choices, Choice{
				Time: h,
				Text: s,
				OnSelected: func() {
					player.Planet.AreasVisited = gamedata.PlanetVisitStatus{}
					r.commitChoiceExtra(gamedata.ModeOrbiting, gamedata.ModeDocked)
				},
			})
		}
	}

	if player.Mode == gamedata.ModeDocked {
		canJump = false
	}

	if len(r.choices) < MaxChoices && player.Mode == gamedata.ModeDocked {
		if player.VesselHP < 1.0 {
			price := r.scene.Rand().FloatRange(0.3, 0.5)
			repairAmount := 1.0 - player.VesselHP
			fullPrice := int(math.Ceil((100 * repairAmount) * price))
			if player.Credits > fullPrice {
				// Every 5% is 1 hour.
				// Repair of 100% is 20 hours.
				repairTime := int(math.Ceil(player.VesselHP * 20))
				r.choices = append(r.choices, Choice{
					Time: repairTime,
					Text: "Repair vessel",
					OnSelected: func() {
						player.Credits -= fullPrice
						player.VesselHP = 1.0
						r.commitChoice(gamedata.ModeDocked)
					},
				})
			}
		}
	}

	if len(r.choices) < MaxChoices && player.Mode == gamedata.ModeDocked {
		if player.Credits > 0 && player.Fuel < player.MaxFuel {
			r.choices = append(r.choices, Choice{
				Time: 2,
				Text: "Buy fuel",
				OnSelected: func() {
					r.eventInfo = eventInfo{kind: eventBuyFuel}
					r.commitChoice(gamedata.ModeDocked)
				},
			})
		}
	}

	if len(r.choices) < MaxChoices && player.Mode == gamedata.ModeDocked {
		if r.world.NextUpgradeDelay == 0 {
			r.choices = append(r.choices, Choice{
				Time: 1,
				Text: "Visit upgrade lab",
				OnSelected: func() {
					r.eventInfo = eventInfo{kind: eventUpgradeLab}
					r.commitChoice(gamedata.ModeDocked)
				},
			})
		}
	}

	if len(r.choices) < MaxChoices && player.Mode == gamedata.ModeDocked && !planet.AreasVisited.VisitedMineralsMarket {
		if player.Cargo > 0 && r.scene.Rand().Chance(0.9) {
			r.choices = append(r.choices, Choice{
				Time: 2,
				Text: "Sell minerals",
				OnSelected: func() {
					planet.AreasVisited.VisitedMineralsMarket = true
					r.eventInfo = eventInfo{kind: eventSellMinerals}
					r.commitChoice(gamedata.ModeDocked)
				},
			})
		}
	}

	hasFuel := false
	if canJump {
		// Find all possible routes first.
		for _, p := range r.world.Planets {
			if p == player.Planet {
				continue
			}
			dist := player.Planet.Info.MapOffset.DistanceTo(p.Info.MapOffset)
			fuelNeeded := gmath.ClampMin(int(dist*player.FuelUsage), 1)
			if player.Fuel < fuelNeeded {
				continue
			}
			hasFuel = true
			if dist > player.MaxJumpDist {
				continue
			}
			hours := int(math.Ceil(dist / player.JumpSpeed))
			r.jumpOptions = append(r.jumpOptions, jumpOption{
				planet:   p,
				fuelCost: fuelNeeded,
				time:     hours,
			})
		}
		gmath.Shuffle(r.scene.Rand(), r.jumpOptions)
		// Add as many travel options as possible.
		for len(r.jumpOptions) > 0 && len(r.choices) < MaxChoices {
			j := r.jumpOptions[len(r.jumpOptions)-1]
			r.jumpOptions = r.jumpOptions[:len(r.jumpOptions)-1]
			r.choices = append(r.choices, Choice{
				Time: j.time,
				Text: fmt.Sprintf("Jump to %s [%d fuel]", j.planet.Info.Name, j.fuelCost),
				OnSelected: func() {
					player.Planet = j.planet
					player.Fuel -= j.fuelCost
					r.commitChoiceExtra(gamedata.ModeJump, gamedata.ModeJustEntered)
				},
			})
		}
	}

	if len(r.choices) < MaxChoices && player.Mode == gamedata.ModeDocked {
		r.choices = append(r.choices, Choice{
			Time: 4,
			Text: "Take off",
			OnSelected: func() {
				r.commitChoice(gamedata.ModeOrbiting)
			},
		})
	}

	if len(r.choices) < MaxChoices && player.Mode != gamedata.ModeDocked {
		r.choices = append(r.choices, Choice{
			Time: 1,
			Text: "Combat test",
			OnSelected: func() {
				r.eventInfo = eventInfo{
					kind: eventBattle,
					enemy: &gamedata.VesselDesign{
						Image:           assets.ImageVesselMarauder,
						MaxHP:           150,
						MaxEnergy:       120,
						EnergyRegen:     3.0,
						MaxSpeed:        180,
						Acceleration:    90,
						RotationSpeed:   2.5,
						MainWeapon:      gamedata.FindWeaponDesign("Pulse Laser"),
						SecondaryWeapon: gamedata.FindWeaponDesign("Homing Missile Launcher"),
					},
				}
				r.commitChoice(gamedata.ModeCombat)
			},
		})
	}

	if len(r.choices) < MaxChoices && player.Cargo < player.MaxCargo && player.VesselHP > 0.3 {
		switch r.world.Player.Mode {
		case gamedata.ModeJustEntered, gamedata.ModeOrbiting:
			if planet.MineralsDelay == 0 && r.scene.Rand().Chance(0.7) {
				r.choices = append(r.choices, Choice{
					Time: 7,
					Text: "Hunt asteroids for minerals",
					OnSelected: func() {
						r.eventInfo = eventInfo{kind: eventMineralsHunt}
						r.commitChoice(gamedata.ModeScavenging)
					},
				})
			}
		}
	}

	if len(r.choices) < MaxChoices {
		switch r.world.Player.Mode {
		case gamedata.ModeJustEntered, gamedata.ModeOrbiting:
			canScavenge := !hasFuel || (player.Fuel < 100 && r.scene.Rand().Chance(0.4))
			if canScavenge {
				r.choices = append(r.choices, Choice{
					Time: 8,
					Text: "Scavenge for fuel",
					OnSelected: func() {
						r.eventInfo = eventInfo{kind: eventFuelScavenge}
						r.commitChoice(gamedata.ModeScavenging)
					},
				})
			}
		}
	}

	return GeneratedChoices{
		Text:    strings.Join(r.textLines, "\n"),
		Choices: r.choices,
	}
}

func (r *Runner) commitChoice(m gamedata.Mode) {
	r.world.Player.Mode = m
	r.EventChoiceSelected.Emit(m)
}

func (r *Runner) commitChoiceExtra(m, postMode gamedata.Mode) {
	r.world.Player.Mode = m
	r.EventChoiceSelected.Emit(postMode)
}
